package regressionsshortcutstore

import (
	"cmp"
	"context"
	"crypto/md5"
	"slices"
	"strings"

	"go.skia.org/infra/go/skerr"
	"go.skia.org/infra/go/sql/pool"
	"go.skia.org/infra/perf/go/types"
)

// RegressionsShortcutStore implements the regressionsshortcut.Store interface.
type RegressionsShortcutStore struct {
	// db is the underlying database.
	db pool.Pool
}

// New returns a *RegressionsShortcutStore
func New(db pool.Pool) *RegressionsShortcutStore {
	return &RegressionsShortcutStore{
		db: db,
	}
}

// Write implements the RegressionsShortcutStore interface
func (rss *RegressionsShortcutStore) Write(ctx context.Context, regrIdList []string) error {
	slices.SortFunc(regrIdList, func(a, b string) int {
		return cmp.Compare(a, b)
	})
	shortcut := rss.calcHash(regrIdList)
	_, err := rss.db.Exec(ctx, `INSERT INTO RegressionsShortcuts(sid, anomaly_ids) VALUES ($1, $2)`, shortcut, regrIdList)
	if err != nil {
		return skerr.Fmt("failed to write new regressions shortcut: %s", err)
	}
	return nil
}

func (rss *RegressionsShortcutStore) calcHash(regrIdList []string) string {
	hash := md5.Sum([]byte(strings.Join(regrIdList, ",")))
	return string(types.TraceIDForSQLFromTraceIDAsBytes(hash[:]))
}
