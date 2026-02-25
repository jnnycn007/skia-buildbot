package main

import (
	"context"
	"crypto/sha256"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"strings"
	"unicode/utf8"

	"cloud.google.com/go/storage"
	"github.com/rs/cors"
	"go.skia.org/infra/go/common"
	"go.skia.org/infra/go/httputils"
	"go.skia.org/infra/go/skerr"
	"go.skia.org/infra/go/sklog"
	"go.skia.org/infra/go/util"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

var (
	// Flags.
	host     = flag.String("host", "localhost", "HTTP service host")
	port     = flag.String("port", ":8000", "HTTP service port (e.g., ':8000')")
	promPort = flag.String("prom_port", ":20000", "Metrics service address (e.g., ':10110')")
	local    = flag.Bool("local", false, "Running locally if true. As opposed to in production.")
	bucket   = flag.String("bucket", "", "GCS bucket containing the files to serve.")

	// Global GCS client.
	gcsClient *storage.Client

	// These URLs get backed up to GCS on startup, if they aren't already present.
	backupUrlToSha256 = map[string]string{
		"https://snapshot-cloudflare.debian.org/archive/debian/20260201T022025Z/pool/main/c/ca-certificates/ca-certificates_20250419_all.deb": "ef590f89563aa4b46c8260d49d1cea0fc1b181d19e8df3782694706adf05c184",
	}
)

func validatePath(path string) error {
	if path == "" {
		return errors.New("path is empty")
	}
	if !utf8.ValidString(path) {
		return errors.New("path is not valid UTF-8")
	}
	if path == "." || path == ".." {
		return errors.New("'.' and '..' are not allowed in GCS paths")
	}
	return nil
}

func storageHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/")
	if err := validatePath(path); err != nil {
		sklog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	gcsPath := fmt.Sprintf("gs://%s/%s", *bucket, path)

	obj := gcsClient.Bucket(*bucket).Object(path)
	reader, err := obj.NewReader(r.Context())
	if err == storage.ErrObjectNotExist {
		http.Error(w, "object not found", http.StatusNotFound)
		return
	} else if err != nil {
		sklog.Errorf("Error creating object reader for %s: %s", gcsPath, err)
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}
	defer util.Close(reader)

	w.Header().Set("Content-Type", reader.Attrs.ContentType)
	w.Header().Set("Content-Encoding", reader.Attrs.ContentEncoding)
	size, err := io.Copy(w, reader)
	if err != nil {
		sklog.Errorf("Error reading object %s: %s", gcsPath, err)
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}
	if size != reader.Attrs.Size {
		errMsg := fmt.Sprintf("Read incorrect number of bytes for %s. Expected %d but read %d", gcsPath, reader.Attrs.Size, size)
		sklog.Error(errMsg)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}
}

func main() {
	common.InitWithMust(
		"autoroll-fe",
		common.PrometheusOpt(promPort),
	)
	defer common.Defer()

	if *bucket == "" {
		sklog.Fatal("--bucket is required.")
	}
	*bucket = strings.TrimPrefix(*bucket, "gs://")

	ctx := context.Background()
	// Note: storage.NewClient() will use Application Default Credentials by
	// default if option.WithHTTPClient is not provided, but doing to results in
	// the caller needing to have the serviceusage.services.use permission for
	// the project in question, which our developer accounts do not seem to
	// have by default.
	ts, err := google.DefaultTokenSource(ctx, storage.ScopeReadOnly)
	if err != nil {
		sklog.Fatal(err)
	}
	httpClient := httputils.DefaultClientConfig().WithTokenSource(ts).Client()
	gcsClient, err = storage.NewClient(ctx, option.WithScopes(storage.ScopeReadOnly), storage.WithJSONReads(), option.WithHTTPClient(httpClient))
	if err != nil {
		sklog.Fatal(err)
	}

	// Run file backup in a separate goroutine to avoid downtime if it fails.
	go func() {
		if err := ensureBackups(ctx, httpClient, backupUrlToSha256); err != nil {
			sklog.Errorf("Failed backups: %s", err)
		}
	}()

	serverURL := "https://" + *host
	if *local {
		serverURL = "http://" + *host + *port
	}

	h := httputils.LoggingRequestResponse(http.HandlerFunc(storageHandler))
	h = httputils.XFrameOptionsDeny(h)
	if !*local {
		h = cors.New(cors.Options{
			AllowedOrigins: []string{"*"},
			Debug:          true,
		}).Handler(h)
		h = httputils.HealthzAndHTTPS(h)
	}
	sklog.Infof("Ready to serve on %s", serverURL)
	sklog.Fatal(http.ListenAndServe(*port, h))
}

func ensureBackups(ctx context.Context, httpClient *http.Client, urlToSha256 map[string]string) error {
	sklog.Infof("Backing up specified URLs...")
	for url, sha256 := range urlToSha256 {
		if err := maybeBackupURL(ctx, httpClient, url, sha256); err != nil {
			return skerr.Wrapf(err, "failed to back up %q", url)
		}
	}
	sklog.Infof("Done backing up specified URLs.")
	return nil
}

func maybeBackupURL(ctx context.Context, httpClient *http.Client, url, expectSha256 string) error {
	// Per storage API docs, if the context is cancelled before Writer.Close()
	// is called, the writes are not saved. We'll use this to ensure that the
	// object doesn't end up existing if the digest doesn't match the
	// expectation.
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	obj := gcsClient.Bucket(*bucket).Object(expectSha256)
	if _, err := obj.Attrs(ctx); err == nil {
		// The object already exists; nothing to do.
		return nil
	}

	// We'll write the object while computing the digest.
	const bufSize = 1024 * 1024
	objWriter := obj.NewWriter(ctx)
	objWriter.ChunkSize = bufSize
	hash := sha256.New()
	w := io.MultiWriter(objWriter, hash)

	// Read the URL.
	resp, err := httpClient.Get(url)
	if err != nil {
		return skerr.Wrap(err)
	}
	defer resp.Body.Close()

	// Perform the copy.
	buf := make([]byte, bufSize)
	if _, err := io.CopyBuffer(w, resp.Body, buf); err != nil {
		return skerr.Wrap(err)
	}

	// Ensure that the digest matches the expectation.
	actualSha256 := fmt.Sprintf("%x", hash.Sum(nil))
	if actualSha256 != expectSha256 {
		return skerr.Fmt("incorrect sha256 digest for %q: expected %s but got %s", url, expectSha256, actualSha256)
	}
	sklog.Infof("Backed up %s with digest %s", url, expectSha256)

	return skerr.Wrap(objWriter.Close())
}
