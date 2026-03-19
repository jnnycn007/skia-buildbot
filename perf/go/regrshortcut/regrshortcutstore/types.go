package regressionsshortcutstore

import "context"

// Store persists regression shortcuts
type Store interface {
	// Write creates a new shortcut for a list of regressions
	Write(ctx context.Context, regressionIdList []string) error
}
