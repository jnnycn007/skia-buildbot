package regrshortcutstore

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/json"
	"slices"
	"strings"

	"go.skia.org/infra/go/cache"
	"go.skia.org/infra/go/skerr"
	"go.skia.org/infra/perf/go/types"
)

// cacheRegressionsShortcutStore provides an implementation of regrshortcut.Store
// which stores these shortcuts in a cache instead of the database.
// The primary use case is when we connect to prod database on a local instance,
// using multigraph needs write access to the database in order to write the
// regrshortcuts. This store prevents the need to elevate to breakglass by
// using a local cache to store this data.
type cacheRegressionsShortcutStore struct {
	cacheClient cache.Cache
}

// NewCacheRegressionsShortcutStore returns a new instance of cacheRegressionsShortcutStore.
func NewCacheRegressionsShortcutStore(cacheClient cache.Cache) *cacheRegressionsShortcutStore {
	return &cacheRegressionsShortcutStore{
		cacheClient: cacheClient,
	}
}

// Create implements the regrshortcut.Store interface.
func (c *cacheRegressionsShortcutStore) Create(ctx context.Context, regrIdList []string) (string, error) {
	if len(regrIdList) == 0 {
		return "", skerr.Fmt("regression id list cannot be empty")
	}

	slices.Sort(regrIdList)
	shortcut := c.calcHash(regrIdList)

	var buff bytes.Buffer
	err := json.NewEncoder(&buff).Encode(regrIdList)
	if err != nil {
		return "", skerr.Wrapf(err, "Failed to encode regression id list")
	}

	err = c.cacheClient.SetValue(ctx, shortcut, buff.String())
	if err != nil {
		return "", skerr.Wrapf(err, "Failed to set value in cache")
	}

	return shortcut, nil
}

// Get implements the regrshortcut.Store interface.
func (c *cacheRegressionsShortcutStore) Get(ctx context.Context, shortcut string) ([]string, error) {
	if !strings.HasPrefix(shortcut, "\\x") {
		shortcut = "\\x" + shortcut
	}
	value, err := c.cacheClient.GetValue(ctx, shortcut)
	if err != nil {
		return nil, skerr.Wrapf(err, "Failed to get value from cache")
	}

	var regrIdList []string
	if err := json.Unmarshal([]byte(value), &regrIdList); err != nil {
		return nil, skerr.Wrapf(err, "Failed to decode regression id list")
	}

	return regrIdList, nil
}

func (c *cacheRegressionsShortcutStore) calcHash(regrIdList []string) string {
	hash := md5.Sum([]byte(strings.Join(regrIdList, ",")))
	return string(types.TraceIDForSQLFromTraceIDAsBytes(hash[:]))
}
