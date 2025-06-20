// Workflow to run all CBB benchmarks on a particular browser / device

package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"go.skia.org/infra/go/skerr"
	"go.skia.org/infra/go/sklog"
	"go.skia.org/infra/perf/go/ingest/format"
	"go.skia.org/infra/perf/go/perfresults"
	"go.skia.org/infra/pinpoint/go/common"
	"go.skia.org/infra/pinpoint/go/workflows"

	"go.temporal.io/sdk/workflow"
)

// CbbRunnerParams defines the parameters for CbbRunnerWorkflow.
type CbbRunnerParams struct {
	// Name of the device configuration, e.g. "mac-m3-pro-perf-cbb".
	BotConfig string

	// The commit to run the benchmarks on. For CBB, this should be a commit in
	// the Chromium main branch, with updated CBB marker file.
	Commit *common.CombinedCommit

	// Name of the Browser to test. Supports "chrome", "safari", and "edge".
	Browser string

	// Browser channel to test. All browsers support "Stable". Chrome and Edge
	// also support "Dev", while Safari also supports "tp" (Technology Preview).
	Channel string

	// The set of benchmarks to run, and the number of iterations for each.
	// If nil, use default.
	Benchmarks []BenchmarkRunConfig
}

// Configuration for a particular benchmark.
type BenchmarkRunConfig struct {
	Benchmark  string
	Iterations int32
}

func setupBenchmarks(cbb *CbbRunnerParams) []BenchmarkRunConfig {
	if cbb.Benchmarks != nil {
		return cbb.Benchmarks
	}

	botConfig := cbb.BotConfig
	var benchmarks []BenchmarkRunConfig
	if strings.HasPrefix(botConfig, "mac-") {
		benchmarks = append(benchmarks, BenchmarkRunConfig{"speedometer3", 150})
	} else {
		benchmarks = append(benchmarks, BenchmarkRunConfig{"speedometer3", 22})
	}

	benchmarks = append(benchmarks, BenchmarkRunConfig{"jetstream2", 22})
	benchmarks = append(benchmarks, BenchmarkRunConfig{"motionmark1.3", 22})

	return benchmarks
}

// Info about the browser we are testing. Retrieved from CBB info file in chromium repo.
type browserInfo struct {
	Browser                    string `json:"browser"`
	Channel                    string `json:"channel"`
	Platform                   string `json:"platform"`
	Version                    string `json:"version"`
	ChromiumMainBranchPosition int    `json:"chromium_main_branch_position"`
}

func getBrowserInfo(ctx workflow.Context, cbb *CbbRunnerParams) (*browserInfo, error) {
	var platformName string
	if strings.HasPrefix(cbb.BotConfig, "mac-") {
		platformName = "macOS"
	} else if strings.HasPrefix(cbb.BotConfig, "win-") {
		platformName = "Windows"
	} else if strings.HasPrefix(cbb.BotConfig, "android-") {
		platformName = "Android"
	} else {
		return nil, errors.New(fmt.Sprintf("Unable to determine platform for bot %s", cbb.BotConfig))
	}
	gitPath := fmt.Sprintf("testing/perf/cbb_ref_info/%s/%s/%s.json", cbb.Browser, cbb.Channel, platformName)

	var content []byte
	if err := workflow.ExecuteActivity(ctx, ReadGitFileActivity, cbb.Commit, gitPath).Get(ctx, &content); err != nil {
		return nil, skerr.Wrapf(err, "Unable to fetch CBB info file %s", gitPath)
	}

	var bi browserInfo
	if err := json.Unmarshal(content, &bi); err != nil {
		return nil, skerr.Wrapf(err, "Unable to parse contents of CBB info file %s", gitPath)
	}
	sklog.Infof("CBB browser info: %v", bi)

	return &bi, nil
}

func setupBrowser(bi *browserInfo) string {
	browser := fmt.Sprintf("--official-browser=%s-%s", strings.ToLower(bi.Browser), strings.ToLower(bi.Channel))
	if bi.Browser == "Chrome" {
		browser = fmt.Sprintf("%s-%s", browser, bi.Version)
	}
	return browser
}

// Workflow to run all CBB benchmarks on a particular browser / bot config.
// TODO(b/388894957): Upload the results.
func CbbRunnerWorkflow(ctx workflow.Context, cbb *CbbRunnerParams) (*map[string]*format.Format, error) {
	ctx = workflow.WithActivityOptions(ctx, regularActivityOptions)
	ctx = workflow.WithChildOptions(ctx, runBenchmarkWorkflowOptions)

	bi, err := getBrowserInfo(ctx, cbb)
	if err != nil {
		return nil, skerr.Wrap(err)
	}

	benchmarks := setupBenchmarks(cbb)
	browser := setupBrowser(bi)

	results := map[string]*format.Format{}

	for _, b := range benchmarks {
		p := &SingleCommitRunnerParams{
			BotConfig:      cbb.BotConfig,
			Benchmark:      b.Benchmark + ".crossbench",
			Story:          "default",
			CombinedCommit: cbb.Commit,
			Iterations:     b.Iterations,
			ExtraArgs:      []string{browser},
		}

		var cr *CommitRun
		if err := workflow.ExecuteChildWorkflow(ctx, workflows.SingleCommitRunner, p).Get(ctx, &cr); err != nil {
			return nil, skerr.Wrap(err)
		}

		results[b.Benchmark] = formatResult(cr, cbb.BotConfig, p.Benchmark)
	}

	return &results, nil
}

// Taking all swarming task results for one benchmark on one bot config,
// and convert the results into the format required by the perf dashboard.
func formatResult(cr *CommitRun, bot string, benchmark string) *format.Format {
	data := format.Format{
		Version: 1,
		GitHash: cr.Build.Commit.Main.GitHash,
		Key: map[string]string{
			"bot":       bot,
			"benchmark": benchmark,
		},
	}

	values := map[string][]float64{}
	units := map[string]string{}
	for _, run := range cr.Runs {
		for c, v := range run.Values {
			values[c] = append(values[c], v...)
		}
		for c, u := range run.Units {
			units[c] = u
		}
	}

	for c, v := range values {
		h := perfresults.Histogram{
			SampleValues: v,
		}
		u := units[c]
		r := format.Result{
			Key: map[string]string{
				"test": c,
				"unit": u,
			},
			Measurements: map[string][]format.SingleMeasurement{
				// Following the convention used in CBB v2 (bench-o-matic) and v3 (swarming-in-google3),
				// we use median (instead of mean) as the data value. We also return standard deviation
				// as a measurement of noise. Other statistics can be added in the future when needed.
				"stat": {
					{
						Value:       "value",
						Measurement: float32(h.Median()),
					},
					{
						Value:       "error",
						Measurement: float32(h.Stddev()),
					},
				},
			},
		}
		if strings.HasSuffix(units[c], "_smallerIsBetter") {
			r.Key["improvement_direction"] = "down"
		} else if strings.HasSuffix(units[c], "_biggerIsBetter") {
			r.Key["improvement_direction"] = "up"
		}
		data.Results = append(data.Results, r)
	}

	return &data
}
