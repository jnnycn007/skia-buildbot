package regrshortcut

import "context"

// Store persists regression shortcuts
type Store interface {
	// Create creates and saves a new shortcut for a list of regressions.
	// Returns gracefully in case when an entry for the list already exists.
	Create(ctx context.Context, regressionIdList []string) (string, error)
}
