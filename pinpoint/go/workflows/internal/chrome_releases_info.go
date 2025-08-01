package internal

import (
	"context"
	"encoding/json"
	"fmt"

	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"go.skia.org/infra/go/exec"
	"go.skia.org/infra/go/httputils"
	"go.skia.org/infra/go/skerr"
	"go.skia.org/infra/go/sklog"
	"go.skia.org/infra/go/util"
)

// BuildInfo fromt the ChromiumDash response.
type BuildInfo struct {
	// Browser contains 'Chrome', 'Edge', 'Safari'
	Browser string `json:"browser"`

	// Channel contains 'Canary', 'Dev', 'Beta', 'Stable'
	Channel string `json:"channel"`

	// Platform contains the build platform. e.g. 'Windows'
	Platform string `json:"platform"`

	// Version contains the latest Chrome build version e.g. `136.0.7103.153`
	Version string `json:"version"`
}

const (
	// chromiumDashUrl response contains the latest Chrome build versions.
	chromiumDashUrl = "https://chromiumdash.appspot.com/fetch_releases?num=1"
	// chromeInternalBucket is the bucket to save the build info JSON files.
	chromeInternalBucket = "chrome-perf-non-public"
	// cbbRefInfoPath is the root of the build info files in the bucket.
	cbbRefInfoPath = "cbb_ref_info/chrome/%s/%s.json"
	// cbbRefInfoRepo is the root of the build info files in the chromium/src.
	cbbRefInfoRepo = "testing/perf/cbb_ref_info/chrome/%s/%s.json"
	// cbbBranchName provides a default name to create a new branch.
	cbbBranchName = "cbb-autoroll"
	// cbbCommitMessage provides a default commit message.
	cbbCommitMessage = "Update CBB autorolll for the builds refs"
	// clNumberStatus to get CL# and status from `git cl status` output.
	// e.g. "  * cbb-autoroll : https://crrev.com/c/12345 (closed)"
	// match[1] == "12345", match[2] == '(closed)'
	clNumberStatus = "%s.*:.*https://crrev.com/c/(\\d+) (.+)"
	// clCommitNumber to get CL commit number from `git cl status` output.
	// e.g. "  Cr-Commit-Position: refs/heads/main@{#99999}"
	// match[1] == "99999"
	clCommitNumber = ".*Cr-Commit-Position: refs/heads/main@{#(\\d+)}"
	// crrevUrl to get git hash from a commit position from the crrev.com
	crrevUrl = "https://crrev.com/%s"
	// crrevCommitHash to commit hash from a redirect crrev URL.
	// e.g. "https://chromium.googlesource.com/chromium/src/+/12345abcdef"
	// match[1] == "12345abcdef'
	crrevCommitHash = "https://chromium.googlesource.com/chromium/src/\\+/(.*)"
)

var (
	// Keys match the ChromiumDash and Values match the subfolders in the GCS.
	cbbChannels = map[string]string{
		"Dev":    "dev",
		"Stable": "stable",
	}
	cbbPlatforms = map[string]string{
		"Android": "Android",
		"Mac":     "macOS",
		"Windows": "Windows",
	}
	// httpClient shares the http client object.
	httpClient *http.Client
)

// getChromiumDashInfo detects new Chrome releases, submits their info to the
// main branch, and returns a commit position.
func GetChromeReleasesInfoActivity(ctx context.Context) (*ChromeReleaseInfo, error) {
	// TODO(b/388894957): Create HTTP Client in the Orchestrator to share.
	httpClient = httputils.NewTimeoutClient()
	resp, err := httputils.GetWithContext(ctx, httpClient, chromiumDashUrl)
	if err != nil {
		sklog.Fatalf("Failed to get ChromiumDash response: %s", err)
	}
	var builds []BuildInfo
	if err := json.NewDecoder(resp.Body).Decode(&builds); err != nil {
		sklog.Fatalf("Invalid ChromiumDash response:%s, err: %s", resp.Body, err)
	}

	newBuilds, err := filterBuilds(ctx, builds)
	if err != nil {
		return nil, skerr.Wrap(err)
	}

	return commitBuildsInfo(ctx, newBuilds)
}

// filterBuilds removes supported builds if their version hasn't changed.
func filterBuilds(ctx context.Context, builds []BuildInfo) ([]BuildInfo, error) {
	var store, err = NewStore(ctx, chromeInternalBucket, true)
	if err != nil {
		return nil, skerr.Wrap(err)
	}
	var newBuilds []BuildInfo
	for _, build := range builds {
		if _, found := cbbChannels[build.Channel]; !found {
			continue
		}
		if _, found := cbbPlatforms[build.Platform]; !found {
			continue
		}
		filePath := fmt.Sprintf(cbbRefInfoPath, cbbChannels[build.Channel], cbbPlatforms[build.Platform])
		if store.Exists(filePath) {
			var content, err = store.GetFileContent(filePath)
			if err != nil {
				return nil, skerr.Wrap(err)
			}
			var gcsBuild BuildInfo
			if err := json.Unmarshal(content, &gcsBuild); err != nil {
				return nil, skerr.Wrap(err)
			}
			if build.Version == gcsBuild.Version {
				sklog.Infof("Version did not change. store: %v, repo: %v", gcsBuild, build)
				continue
			}
		} else {
			sklog.Infof("No history found for %s", filePath)
		}
		build.Browser = "Chrome"

		// TODO(b/388894957): We may need to update the GCS after committing.
		jsonData, err := json.MarshalIndent(build, "", "  ")
		if err != nil {
			return nil, skerr.Wrap(err)
		}
		if err := store.WriteFile(filePath, string(jsonData)); err != nil {
			return nil, skerr.Wrap(err)
		}

		newBuilds = append(newBuilds, build)
	}
	return newBuilds, nil
}

// commitBuildsInfo creates JSON files and uploads the associated commit.
func commitBuildsInfo(ctx context.Context, builds []BuildInfo) (*ChromeReleaseInfo, error) {
	sklog.Infof("Builds to commit thier info: %v", builds)
	if len(builds) == 0 {
		sklog.Infof("No new build was detected.")
		return nil, nil
	}
	client, err := NewGitChromium(ctx)
	if err != nil {
		return nil, skerr.Wrap(err)
	}
	if err := client.ShallowClone(cbbBranchName); err != nil {
		return nil, skerr.Wrap(err)
	}

	for _, build := range builds {
		filename := fmt.Sprintf(cbbRefInfoRepo, cbbChannels[build.Channel], cbbPlatforms[build.Platform])
		path := filepath.Join(client.repoDir, filename)
		jsonData, err := json.MarshalIndent(build, "", "  ")
		if err != nil {
			return nil, skerr.Wrapf(err, "Failed to convert %v to JSON", build)
		}
		if err := os.WriteFile(path, []byte(jsonData), 0644); err != nil {
			return nil, skerr.Wrapf(err, "Failed to write: %s", path)
		}
		if _, err := exec.RunCwd(client.ctx, client.repoDir, client.gitExec, "add", filename); err != nil {
			return nil, skerr.Wrapf(err, "Failed to add %s in Git", filename)
		}
		sklog.Infof("Git added %s", filename)
	}

	if _, err := exec.RunCwd(client.ctx, client.repoDir, client.gitExec, "commit", "-m", cbbCommitMessage); err != nil {
		return nil, skerr.Wrapf(err, "Failed to commit")
	}
	sklog.Infof("Git committed successfully! Start uploading it.")

	stdout, err := exec.RunCwd(
		client.ctx, client.repoDir, client.gitExec,
		"push", "origin", "HEAD:refs/for/main",
		"-o", "r=rubber-stamper@appspot.gserviceaccount.com",
		"-o", "l=Auto-Submit+1")
	if err != nil {
		return nil, skerr.Wrapf(err, "Failed to upload the change.")
	}
	sklog.Infof("Git uploaded successfully! stdput=%s", stdout)

	return waitForSubmitCl(client)
}

// waitForSubmitCl waits for the 'rubber-stamper' to submit uploaded CLs, then
// returns the commit position.
func waitForSubmitCl(client *gitClient) (*ChromeReleaseInfo, error) {
	var commitPosition string
	statusPattern := fmt.Sprintf(clNumberStatus, cbbBranchName)
	sklog.Infof("Waiting for CL to be submitted.")
	start := time.Now()
	for {
		if time.Now().Sub(start) > ClSubmissionTimeout {
			return nil, fmt.Errorf("waitForSubmitCl timeout!")
		}
		// TODO(b/433796566): Re-implement without using "cl" command.
		stdout, err := exec.RunCwd(
			client.ctx, client.repoDir, client.gitExec, "cl", "status")
		if err != nil {
			return nil, skerr.Wrapf(err, "Failed to run `git cl status`.")
		}
		re := regexp.MustCompile(statusPattern)
		match := re.FindStringSubmatch(stdout)
		if len(match) == 3 && match[2] == "(closed)" {
			re = regexp.MustCompile(clCommitNumber)
			match = re.FindStringSubmatch(stdout)
			if len(match) != 2 {
				return nil, fmt.Errorf("Failed to detect Commit Number: %s", stdout)
			}
			commitPosition = match[1]
			sklog.Infof("Detected commit number=%s", commitPosition)
			return findCommitHash(client.ctx, commitPosition)
		} else {
			sklog.Infof("CL status: stdout=%s\nmatch=%v", stdout, match)
		}

		time.Sleep(10 * time.Second)
	}
}

// findCommitHash finds the commit hash by hitting the crrev.com with the
// commit position. The redirected url includes the commit hash value.
func findCommitHash(ctx context.Context, commitPosition string) (*ChromeReleaseInfo, error) {
	url := fmt.Sprintf(crrevUrl, commitPosition)
	resp, err := httputils.GetWithContext(ctx, httpClient, url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		util.Close(resp.Body)
		return nil, skerr.Fmt("findCommitHash got status %q", resp.Status)
	}
	redirectUrl := resp.Request.URL.String()
	sklog.Infof("crrev.com redirect URL=%s", redirectUrl)
	re := regexp.MustCompile(crrevCommitHash)
	match := re.FindStringSubmatch(redirectUrl)
	if len(match) != 2 {
		return nil, fmt.Errorf("Failed to detect Commit Hash: %v", resp.Request.URL)
	}
	commitHash := match[1]
	commitInfo := &ChromeReleaseInfo{
		CommitPosition: commitPosition,
		CommitHash:     commitHash,
	}
	return commitInfo, nil
}
