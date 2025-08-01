// Workflow to run all CBB benchmarks on a particular browser / device

package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"path"
	"strings"
	"time"

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

	// Name of the Google cloud storage bucket to upload results to.
	Bucket string
}

// Configuration for a particular benchmark.
type BenchmarkRunConfig struct {
	Benchmark  string
	Iterations int32
}

func validateParameters(cbb *CbbRunnerParams) error {
	switch cbb.Browser {
	case "chrome", "edge":
		if cbb.Channel != "stable" && cbb.Channel != "dev" {
			return fmt.Errorf(
				"Unrecognized browser channel %s, %s only supports stable and dev",
				cbb.Channel, cbb.Browser)
		}
	case "safari":
		if cbb.Channel != "stable" && cbb.Channel != "technology-preview" {
			return fmt.Errorf(
				"Unrecognized browser channel %s, %s only supports stable and technology-preview",
				cbb.Channel, cbb.Browser)
		}
	default:
		return fmt.Errorf(
			"Unrecognized browser %s, only chrome, safari, and edge are supported",
			cbb.Browser)
	}
	return nil
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

	// TODO(b/433537961) Double number of iterations on Android until we
	// figure out why benchmarks fail frequently on it.
	if strings.HasPrefix(botConfig, "android-") {
		for i := range benchmarks {
			benchmarks[i].Iterations *= 2
		}
	}

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
		return nil, fmt.Errorf("Unable to determine platform for bot %s", cbb.BotConfig)
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

func setupBrowser(cbb *CbbRunnerParams, bi *browserInfo) string {
	browser := fmt.Sprintf("--official-browser=%s-%s", cbb.Browser, cbb.Channel)
	if cbb.Browser == "chrome" {
		browser = fmt.Sprintf("%s-%s", browser, bi.Version)
	}
	return browser
}

// Generate Pinpoint job ID. Uses abbreviations to keep the job ID short.
func genJobId(bi *browserInfo, cbb *CbbRunnerParams, benchmark string) string {
	// Shorten browser/channel name.
	// * Browser name is shortened to 3 characters.
	// * Channel name is omitted if it is "Stable".
	// * Special case: Safari Technology Preview is abbreviated as "STP".
	browser := bi.Browser[:3]
	if bi.Channel != "Stable" {
		if bi.Browser == "Safari" && bi.Channel == "Technology Preview" {
			browser = "STP"
		} else {
			browser += " " + bi.Channel
		}
	}

	// Chrome and Edge versions are used as-is, but Safari versions are too long.
	// * Safari Stable version looks like "1.2.3 (12345.6.7.8)". Truncate at the space.
	// * STP version looks like "1.2.3 (Release 123, 12345.6.7.8)". Keep only the release number.
	version := bi.Version
	if bi.Browser == "Safari" {
		if bi.Channel == "Stable" {
			space := strings.Index(version, " ")
			if space != -1 {
				version = version[:space]
			}
		} else {
			release := strings.Index(version, "Release ")
			if release != -1 {
				version = version[release+8:]
			}
			comma := strings.Index(version, ",")
			if comma != -1 {
				version = version[:comma]
			}
		}
	}

	// Shorten bot name.
	bot := cbb.BotConfig
	switch bot {
	case "android-pixel-tangor-perf-cbb":
		bot = "tang"
	case "mac-m3-pro-perf-cbb":
		bot = "m3"
	case "win-victus-perf-cbb":
		bot = "vic"
	case "win-arm64-snapdragon-plus-perf-cbb":
		bot = "plus"
	case "win-arm64-snapdragon-elite-perf-cbb":
		bot = "elite"
	}

	// Shorten benchmark name.
	if strings.HasPrefix(benchmark, "speedometer") {
		benchmark = "SP"
	} else if strings.HasPrefix(benchmark, "jetstream") {
		benchmark = "JS"
	} else if strings.HasPrefix(benchmark, "motionmark") {
		benchmark = "MM"
	}

	return fmt.Sprintf("CBB %s %s %s %s", browser, version, bot, benchmark)
}

// Workflow to run all CBB benchmarks on a particular browser / bot config and upload the results.
func CbbRunnerWorkflow(ctx workflow.Context, cbb *CbbRunnerParams) (*map[string]*format.Format, error) {
	startTime := time.Now()

	ctx = workflow.WithActivityOptions(ctx, regularActivityOptions)
	ctx = workflow.WithChildOptions(ctx, runBenchmarkWorkflowOptions)

	err := validateParameters(cbb)
	if err != nil {
		return nil, skerr.Wrap(err)
	}

	bi, err := getBrowserInfo(ctx, cbb)
	if err != nil {
		return nil, skerr.Wrap(err)
	}

	benchmarks := setupBenchmarks(cbb)
	browser := setupBrowser(cbb, bi)

	results := map[string]*format.Format{}

	for _, b := range benchmarks {
		p := &SingleCommitRunnerParams{
			PinpointJobID:  genJobId(bi, cbb, b.Benchmark),
			BotConfig:      cbb.BotConfig,
			Benchmark:      b.Benchmark,
			Story:          "default",
			CombinedCommit: cbb.Commit,
			Iterations:     b.Iterations,
			ExtraArgs:      []string{browser},
		}

		if !strings.HasSuffix(p.Benchmark, ".crossbench") {
			p.Benchmark += ".crossbench"
		}

		var cr *CommitRun
		if err := workflow.ExecuteChildWorkflow(ctx, workflows.SingleCommitRunner, p).Get(ctx, &cr); err != nil {
			return nil, skerr.Wrap(err)
		}

		r := formatResult(cr, cbb.BotConfig, p.Benchmark, bi)
		results[b.Benchmark] = r

		if cbb.Bucket != "" {
			var swc StringWriterCloser = StringWriterCloser{
				builder: new(strings.Builder),
			}
			err = r.Write(swc)
			if err != nil {
				return nil, skerr.Wrap(err)
			}

			var gsPath string
			err = workflow.ExecuteActivity(
				ctx, UploadCbbResultsActivity, startTime, cbb.BotConfig, p.Benchmark, swc.builder.String(), cbb.Bucket).Get(ctx, &gsPath)
			if err != nil {
				return nil, skerr.Wrap(err)
			}
			sklog.Infof("Uploaded results to %s", gsPath)
		} else {
			sklog.Warning("No GS bucket provided to upload results")
		}
	}

	return &results, nil
}

// Provide a StringWriterCloser interface wrapper around strings.Builder,
// required for serializing the results to a string.
type StringWriterCloser struct {
	builder *strings.Builder
}

func (swc StringWriterCloser) Write(p []byte) (int, error) {
	return swc.builder.Write(p)
}

func (StringWriterCloser) Close() error {
	return nil
}

// Taking all swarming task results for one benchmark on one bot config,
// and convert the results into the format required by the perf dashboard.
func formatResult(cr *CommitRun, bot string, benchmark string, bi *browserInfo) *format.Format {
	data := format.Format{
		Version: 1,
		GitHash: fmt.Sprintf("CP:%d", cr.Build.Commit.Main.CommitPosition),
		Key: map[string]string{
			"master":    "ChromiumPerf",
			"bot":       bot,
			"benchmark": benchmark,
		},
		Links: map[string]string{
			"Browser Version": bi.Version,
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
				"test":      c,
				"unit":      u,
				"subtest_1": fmt.Sprintf("%s %s", bi.Browser, bi.Channel),
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

// Activity to upload CBB results to cloud storage, for import into perf dashboard.
func UploadCbbResultsActivity(ctx context.Context, t time.Time, bot string, benchmark string, results string, bucket string) (string, error) {
	if bucket == "" {
		return "", nil
	}
	store, err := NewStore(ctx, bucket, false)
	if err != nil {
		return "", skerr.Wrap(err)
	}

	// The file path inside the GS bucket uses a pattern similar to the one used by perf waterfall.
	datePath := fmt.Sprintf("%04d/%02d/%02d", t.Year(), t.Month(), t.Day())
	timestamp := fmt.Sprintf("%02d%02d%02d", t.Hour(), t.Minute(), t.Second())
	filename := fmt.Sprintf(
		"skia_results_%s_%s_%04d_%02d_%02d_%02d_%02d_%02d.json",
		benchmark, bot, t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(),
	)
	storeFilePath := path.Join("ingest", datePath, "ChromiumPerf", bot, timestamp, benchmark, filename)

	err = store.WriteFile(storeFilePath, results)
	if err != nil {
		return "", skerr.Wrap(err)
	}

	gsPath := fmt.Sprintf("gs://%s/%s", bucket, storeFilePath)
	return gsPath, nil
}
