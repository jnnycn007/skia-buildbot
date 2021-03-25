// Package tjstore defines an interface for storing TryJob-related data
// as needed for operating Gold.
package tjstore

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.skia.org/infra/go/paramtools"
	ci "go.skia.org/infra/golden/go/continuous_integration"
	"go.skia.org/infra/golden/go/types"
)

// Store (sometimes called TryJobStore) is an interface around a database
// for storing TryJobs and TryJobResults. Of note, we will only store data for
// TryJobs which uploaded data to Gold (e.g. via ingestion); the purpose of
// this interface is not to store data about every TryJob.
type Store interface {
	// GetTryJob returns the TryJob corresponding to the given id.
	// Returns NotFound if it doesn't exist.
	GetTryJob(ctx context.Context, id, cisName string) (ci.TryJob, error)

	// GetTryJobs returns all TryJobs associated with a given Changelist and Patchset.
	// The returned slice could be empty if the CL or PS don't exist.
	// TryJobs should be sorted by DisplayName.
	GetTryJobs(ctx context.Context, psID CombinedPSID) ([]ci.TryJob, error)

	// GetResults returns any TryJobResults for a given Changelist and Patchset.
	// The returned slice could be empty and is not sorted. If updatedAfter is not
	// a zero time, it will be used to return the subset of results created on or after
	// the given time.
	GetResults(ctx context.Context, psID CombinedPSID, updatedAfter time.Time) ([]TryJobResult, error)
}

var ErrNotFound = errors.New("not found")

type TryJobResult struct {
	// GroupParams describe the general configuration that produced
	// the digest/image. This includes things like the model of device
	// that drew the image. GroupParams are likely to be shared among
	// many, if not all, the TryJobResults for a single TryJob, and
	// by making them a separate parameter, the map can be shared rather
	// than copied. Clients should treat this as read-only and not modify
	// it, as it could be shared by multiple different TryJobResults.
	GroupParams paramtools.Params

	// ResultParams describe the specialized configuration that
	// produced the digest/image. This includes the test name and corpus,
	// things that change for each result. This map is safe to be written
	// to by the client.
	// In the event of conflict, ResultParams should override Options which
	// override GroupParams.
	ResultParams paramtools.Params

	// Options give extra details about this result. This includes things
	// like the file format. Skia uses this for things like gamma_correctness.
	// Clients should treat this as read-only and not modify it, as it could
	// be shared by multiple different TryJobResults.
	Options paramtools.Params

	// Digest references the image that was generated by the test.
	Digest types.Digest

	// TryjobID represents the tryjob that produced this result.
	TryjobID string
	// System represents the CI System that the Tryjob belonged to.
	System string
}

// CombinedPSID represents an identifier that uniquely refers to a Patchset.
type CombinedPSID struct {
	CL  string
	CRS string
	PS  string
}

// Key creates a probably unique id for a given
// Patchset using the id of the Changelist it belongs to and the
// ChangeReviewSystem it is a part of. We say "probably unique" because
// a malicious person could try to control the clID and the psID to make
// two different inputs make the same result, but that is unlikely for
// ids that are valid (i.e. exist on a system like Gerrit).
func (c CombinedPSID) Key() string {
	return fmt.Sprintf("%s__%s__%s", c.CL, c.CRS, c.PS)
}

// Equal returns true if the IDs are identical, false otherwise.
func (c CombinedPSID) Equal(other CombinedPSID) bool {
	return c.CL == other.CL && c.PS == other.PS && c.CRS == other.CRS
}
