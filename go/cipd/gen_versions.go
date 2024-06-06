//go:build ignore
// +build ignore

package main

/*
	Generate asset_versions_gen.go.
*/

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"go.skia.org/infra/go/cipd"
	"go.skia.org/infra/go/depot_tools/deps_parser"
	"go.skia.org/infra/go/exec"
	"go.skia.org/infra/go/skerr"
	"go.skia.org/infra/go/sklog"
	"go.skia.org/infra/go/util"
)

const (
	TARGET_FILE = "asset_versions_gen.go"
	HEADER      = `// Code generated by "go run gen_versions.go"; DO NOT EDIT

package cipd

var PACKAGES = map[string]*Package{
`
	FOOTER = `}`
)

func genAssetVersionsGo(depsEntries deps_parser.DepsEntries, pkgDir, rootDir string) error {
	// Pull out only the CIPD dependencies.
	var pkgs []*cipd.Package
	for _, entry := range depsEntries {
		if entry.Type != deps_parser.DepType_Cipd {
			continue
		}
		pkg := &cipd.Package{
			Path:    entry.Path,
			Name:    entry.Id,
			Version: entry.Version,
		}
		if pkg.Path == "" {
			pkg.Path = "."
		}
		pkgs = append(pkgs, pkg)
	}

	// List the assets.
	assetsDir := path.Join(rootDir, "infra", "bots", "assets")
	entries, err := os.ReadDir(assetsDir)
	if err != nil {
		return skerr.Wrap(err)
	}
	for _, e := range entries {
		if e.IsDir() {
			contents, err := os.ReadFile(path.Join(assetsDir, e.Name(), "VERSION"))
			if err == nil {
				name := e.Name()
				fullName := fmt.Sprintf("skia/bots/%s", name)
				pkgs = append(pkgs, &cipd.Package{
					Path:    name,
					Name:    fullName,
					Version: cipd.VersionTag(strings.TrimSpace(string(contents))),
				})
			} else if !os.IsNotExist(err) {
				return skerr.Wrap(err)
			}
		}
	}

	// Write the file.
	sort.Sort(cipd.PackageSlice(pkgs))
	targetFile := path.Join(pkgDir, TARGET_FILE)
	if err := util.WithWriteFile(targetFile, func(w io.Writer) error {
		_, err := w.Write([]byte(HEADER))
		if err != nil {
			return err
		}
		for _, pkg := range pkgs {
			_, err := fmt.Fprintf(w, fmt.Sprintf(`	"%s": &Package{
		Path: "%s",
		Name: "%s",
		Version: "%s",
	},
`, pkg.Name, pkg.Path, pkg.Name, pkg.Version))
			if err != nil {
				return err
			}
		}
		_, err = w.Write([]byte(FOOTER))
		return err
	}); err != nil {
		return skerr.Wrap(err)
	}
	if _, err := exec.RunCwd(context.Background(), ".", "gofmt", "-s", "-w", targetFile); err != nil {
		return skerr.Wrap(err)
	}
	return nil
}

func genEnsureFile(entries deps_parser.DepsEntries, rootDir string) error {
	// We only use the ensure file on Linux.
	const platform = cipd.PlatformLinuxAmd64

	// Extract the dependency on CIPD itself.
	cipdEntry, ok := entries[cipd.PkgNameCIPD]
	if !ok {
		sklog.Fatal("Unable to find critical package %s", cipd.PkgNameCIPD)
	}
	delete(entries, cipdEntry.Id)

	// Sort package names for consistency.
	entryIDs := make([]string, 0, len(entries))
	for id := range entries {
		entryIDs = append(entryIDs, id)
	}
	sort.Strings(entryIDs)

	// Organize packages by destination path.
	byPath := map[string][]*deps_parser.DepsEntry{}
	longestIDLength := 0
	for _, entryID := range entryIDs {
		entry := entries[entryID]
		if entry.Type != deps_parser.DepType_Cipd {
			continue
		}
		pkgSplit := strings.Split(entry.Id, "/")
		pkgPlatform := pkgSplit[len(pkgSplit)-1]
		if util.In(pkgPlatform, cipd.Platforms) && pkgPlatform != platform {
			continue
		}
		byPath[entry.Path] = append(byPath[entry.Path], entry)
		if len(entry.Id) > longestIDLength {
			longestIDLength = len(entry.Id)
		}
	}

	// Sort path names for consistency.
	paths := make([]string, 0, len(byPath))
	for path := range byPath {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	// Construct the ensure file.
	lines := strings.Split(fmt.Sprintf(`# The CIPD server to use.
$ServiceURL %s

# This is the CIPD client itself.
%s %s
`, cipd.DefaultServiceURL, cipdEntry.Id, cipdEntry.Version), "\n")

	for _, path := range paths {
		entriesInPath := byPath[path]
		if path == "" {
			lines = append(lines, "@Subdir")
		} else {
			lines = append(lines, fmt.Sprintf("@Subdir %s", path))
		}
		for _, entry := range entriesInPath {
			spacer := strings.Repeat(" ", longestIDLength-len(entry.Id))
			lines = append(lines, fmt.Sprintf("%s %s%s", entry.Id, spacer, entry.Version))
		}
		lines = append(lines, "")
	}
	content := []byte(strings.Join(lines, "\n"))

	// Write the file.
	if err := os.WriteFile(filepath.Join(rootDir, "cipd.ensure"), content, os.ModePerm); err != nil {
		sklog.Fatal(err)
	}
	return nil
}

func main() {
	// Read and parse the DEPS file.
	_, filename, _, _ := runtime.Caller(0)
	pkgDir := path.Dir(filename)
	rootDir := path.Join(pkgDir, "..", "..")

	depsContents, err := os.ReadFile(filepath.Join(rootDir, "DEPS"))
	if err != nil {
		sklog.Fatal(err)
	}
	depsEntries, err := deps_parser.ParseDepsNoNormalize(string(depsContents))
	if err != nil {
		sklog.Fatal(err)
	}

	if err := genAssetVersionsGo(depsEntries, pkgDir, rootDir); err != nil {
		sklog.Fatal(err)
	}

	if err := genEnsureFile(depsEntries, rootDir); err != nil {
		sklog.Fatal(err)
	}
}
