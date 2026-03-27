package regrshortcutstore

import (
	"context"
	"crypto/md5"
	"errors"
	"slices"
	"strings"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"go.skia.org/infra/go/skerr"
	"go.skia.org/infra/go/sql/pool"
	"go.skia.org/infra/perf/go/types"
)

// RegressionsShortcutStore implements the regrshortcut.Store interface.
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

// Create implements the regrshortcut.Store interface.
func (rss *RegressionsShortcutStore) Create(ctx context.Context, regrIdList []string) (string, error) {
	slices.Sort(regrIdList)
	shortcut := rss.calcHash(regrIdList)

	if _, err := rss.db.Exec(ctx, `INSERT INTO RegressionsShortcuts(sid, anomaly_ids) VALUES ($1, $2)`, shortcut, regrIdList); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			// Shortcut is already present, we continue gracefully.
			// We don't guard against md5 collisions.
			return shortcut, nil
		}
		return "", skerr.Fmt("failed to write new regressions shortcut: %s", err)
	}
	return shortcut, nil
}

func (rss *RegressionsShortcutStore) calcHash(regrIdList []string) string {
	hash := md5.Sum([]byte(strings.Join(regrIdList, ",")))
	return string(types.TraceIDForSQLFromTraceIDAsBytes(hash[:]))
}
