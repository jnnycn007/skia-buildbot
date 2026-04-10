package formatter

import (
	"context"
	"fmt"
	"strings"

	"go.skia.org/infra/perf/go/config"
	"go.skia.org/infra/perf/go/git"
	"go.skia.org/infra/perf/go/git/provider"
	"go.skia.org/infra/perf/go/notify"
	"go.skia.org/infra/perf/go/types"
)

// NewCommitRangeFormatter returns a standard CommitRangeFormatter that builds Git log URLs using the instance GitRepoConfig.
// Ideally, the behavior should be exactly the same as in tooltip.
// TODO(b/485178559) Make it identical.
func NewCommitRangeFormatter(perfGit git.Git) types.CommitRangeFormatter {
	return func(ctx context.Context, startCommit, endCommit int64) string {
		startHash, err := perfGit.GitHashFromCommitNumber(ctx, types.CommitNumber(startCommit))
		if err != nil {
			return fmt.Sprintf("%d -> %d", startCommit, endCommit)
		}
		endHash, err := perfGit.GitHashFromCommitNumber(ctx, types.CommitNumber(endCommit))
		if err != nil {
			return fmt.Sprintf("%d -> %d", startCommit, endCommit)
		}

		startDisplayed := startHash[:min(len(startHash), 8)]
		endDisplayed := endHash[:min(len(endHash), 8)]

		basePath := config.Config.GitRepoConfig.URL
		var urlTemplate string
		// TODO(b/485178559) rework this to be config.git.source based
		// most of instances use gitiles, but flutter uses just github.
		// Currently, we don't have to support any other VCS.
		// The todo above will make it more clear.
		if strings.Contains(basePath, "googlesource.com") {
			urlTemplate = basePath + "/+log/{begin}..{end}"
		} else {
			urlTemplate = basePath + "/compare/{begin}...{end}"
		}
		commitUrl := notify.URLFromCommitRange(
			provider.Commit{GitHash: endHash},
			provider.Commit{GitHash: startHash},
			urlTemplate,
		)
		return fmt.Sprintf("[%s..%s](%s)", startDisplayed, endDisplayed, commitUrl)
	}
}
