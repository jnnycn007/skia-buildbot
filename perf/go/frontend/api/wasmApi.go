package api

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"go.skia.org/infra/go/httputils"
	"go.skia.org/infra/go/paramtools"
	"go.skia.org/infra/go/query"
	"go.skia.org/infra/go/skerr"
	"go.skia.org/infra/go/sklog"
	"go.skia.org/infra/perf/go/config"
	"go.skia.org/infra/perf/go/psrefresh"
	"go.skia.org/infra/perf/go/tracestore"
	"go.skia.org/infra/perf/go/tracestore/sqltracestore"
	"go.skia.org/infra/perf/go/types"
)

const (
	defaultCacheTTL = 5 * time.Minute
	fileCacheTTL    = 14 * 24 * time.Hour
)

type Param struct {
	Id    uint16 `json:"id"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

// WasmSupport defines the methods needed by WasmApi that are not in the main TraceStore interface.
type WasmSupport interface {
	GetNonEmptyTraceIDs(ctx context.Context, startCommit, endCommit types.CommitNumber) ([][]byte, error)
}

type wasmApi struct {
	traceStore  tracestore.TraceStore
	psRefresher psrefresh.ParamSetRefresher
	cacheDir    string
	filterQuery string

	mutex sync.Mutex
	cache *wasmCache
}

type wasmCache struct {
	tileNumber types.TileNumber
	version    string
	meta       []byte
	params     []byte
	traces     []byte
	createdAt  time.Time
}

func inferredFilterQuery(cfg *config.QueryConfig) string {
	var parts []string
	if cfg.DefaultParamSelections != nil {
		for k, v := range cfg.DefaultParamSelections {
			if len(v) > 0 {
				if len(v) == 1 {
					parts = append(parts, fmt.Sprintf("%s=%s", k, v[0]))
				} else {
					parts = append(parts, fmt.Sprintf("%s=~(%s)", k, strings.Join(v, "|")))
				}
			}
		}
	}
	if cfg.ConditionalDefaults != nil {
		var metrics []string
		for _, rule := range cfg.ConditionalDefaults {
			if rule.Trigger.Param == "metric" {
				metrics = append(metrics, rule.Trigger.Values...)
			}
		}
		if len(metrics) > 0 {
			parts = append(parts, fmt.Sprintf("metric=~(%s)", strings.Join(metrics, "|")))
		}
	}
	return strings.Join(parts, "&")
}

func NewWasmApi(traceStore tracestore.TraceStore, psRefresher psrefresh.ParamSetRefresher, cacheDir string, cfg *config.InstanceConfig) *wasmApi {
	return &wasmApi{
		traceStore:  traceStore,
		psRefresher: psRefresher,
		cacheDir:    cacheDir,
		filterQuery: inferredFilterQuery(&cfg.QueryConfig),
	}
}

func (api *wasmApi) Start(ctx context.Context) {
	if api.traceStore == nil {
		sklog.Warningf("TraceStore is nil, not starting background Wasm cache generator")
		return
	}
	sklog.Infof("Starting background Wasm cache generator")
	go func() {
		if err := api.ensureCache(ctx); err != nil {
			sklog.Errorf("Failed to generate initial Wasm cache: %v", err)
		}

		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := api.ensureCache(ctx); err != nil {
					sklog.Errorf("Failed to refresh Wasm cache: %v", err)
				}
			}
		}
	}()
}

func (api *wasmApi) RegisterHandlers(router *chi.Mux) {
	router.Get("/_/wasm/meta.json", api.metaHandler)
	router.Get("/_/wasm/params.json", api.paramsHandler)
	router.Get("/_/wasm/traces.bin", api.tracesHandler)
}

func (api *wasmApi) ensureCache(ctx context.Context) error {
	api.mutex.Lock()
	defer api.mutex.Unlock()

	tileCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	tile, err := api.traceStore.GetLatestTile(tileCtx)
	if err != nil {
		return skerr.Wrap(err)
	}

	if api.cache != nil && api.cache.tileNumber == tile {
		// The latest tile is actively updated, so we enforce a TTL.
		if time.Since(api.cache.createdAt) < defaultCacheTTL {
			return nil
		}
		sklog.Infof("Refreshing Wasm cache for tile %d (TTL expired)", tile)
	}

	if err := os.MkdirAll(api.cacheDir, 0755); err != nil {
		return skerr.Wrapf(err, "Failed to create cache dir %q", api.cacheDir)
	}

	tracesFile := filepath.Join(api.cacheDir, fmt.Sprintf("traces_%d.bin", tile))
	metaFile := filepath.Join(api.cacheDir, fmt.Sprintf("meta_%d.json", tile))
	paramsFile := filepath.Join(api.cacheDir, fmt.Sprintf("params_%d.json", tile))

	stat, err := os.Stat(tracesFile)
	if err == nil && time.Since(stat.ModTime()) < fileCacheTTL {

		_, errMeta := os.Stat(metaFile)
		_, errParams := os.Stat(paramsFile)
		if errMeta == nil && errParams == nil {
			sklog.Infof("Loading Wasm cache from files for tile %d", tile)
			tracesBuf, err1 := os.ReadFile(tracesFile)
			metaBuf, err2 := os.ReadFile(metaFile)
			paramsBuf, err3 := os.ReadFile(paramsFile)
			if err1 == nil && err2 == nil && err3 == nil {
				api.cache = &wasmCache{
					tileNumber: tile,
					traces:     tracesBuf,
					meta:       metaBuf,
					params:     paramsBuf,
					createdAt:  stat.ModTime(),
				}
				return nil
			}
			sklog.Errorf("Failed to read cache files: %v, %v, %v", err1, err2, err3)
		}
	}

	sklog.Infof("Generating Wasm memory cache for tile %d", tile)

	lookup, stride, params := api.buildLookup()

	// Fetch traces to build traces.bin.
	fmt.Println("filterQuery: ", api.filterQuery)
	q, err := query.NewFromString(api.filterQuery)
	if err != nil {
		return skerr.Wrap(err)
	}
	queryCtx, cancelQuery := context.WithTimeout(ctx, 60*time.Second)
	defer cancelQuery()
	queryCtx = context.WithValue(queryCtx, sqltracestore.UseInvertedIndex, true)
	outParams, err := api.traceStore.QueryTracesIDOnly(queryCtx, tile, q)
	if err != nil {
		return skerr.Wrap(err)
	}

	var allParams []paramtools.Params
	var allKeys []string
	for p := range outParams {
		allParams = append(allParams, p)
		key, err := query.MakeKeyFast(p)
		if err != nil {
			continue
		}
		allKeys = append(allKeys, key)
	}

	// Filter by non-empty traces
	wasmSupport, ok := api.traceStore.(WasmSupport)
	if !ok {
		return skerr.Fmt("TraceStore does not support Wasm operations")
	}
	tileSize := api.traceStore.TileSize()
	beginCommit, endCommit := types.TileCommitRangeForTileNumber(tile, tileSize)

	// Use a timeout for this potentially heavy operation
	ctxTimeout, cancel := context.WithTimeout(ctx, 20*time.Minute)
	defer cancel()
	nonEmptyTraceIDs, err := wasmSupport.GetNonEmptyTraceIDs(ctxTimeout, beginCommit, endCommit)
	if err != nil {
		return skerr.Wrap(err)
	}

	nonEmptyKeys := filterNonEmptyKeys(allKeys, nonEmptyTraceIDs)
	sklog.Infof("Filtered to %d non-empty traces out of %d", len(nonEmptyKeys), len(allKeys))

	tracesBinary, traceCount := encodeTraces(allParams, allKeys, nonEmptyKeys, lookup, stride)

	version := fmt.Sprintf("%d", time.Now().Unix())

	meta := struct {
		Stride  int    `json:"stride"`
		Count   int    `json:"count"`
		Version string `json:"version"`
	}{
		Stride:  stride,
		Count:   traceCount,
		Version: version,
	}

	metaBytes, err := json.Marshal(meta)
	if err != nil {
		return skerr.Wrap(err)
	}

	paramsBytes, err := json.Marshal(params)
	if err != nil {
		return skerr.Wrap(err)
	}

	api.cache = &wasmCache{
		tileNumber: tile,
		version:    version,
		meta:       metaBytes,
		params:     paramsBytes,
		traces:     tracesBinary,
		createdAt:  time.Now(),
	}

	sklog.Infof("Generated Wasm cache: traces=%d stride=%d", traceCount, stride)

	// Save to cache files
	if err := os.WriteFile(tracesFile, tracesBinary, 0644); err != nil {
		sklog.Errorf("Failed to save traces cache: %v", err)
	}
	if err := os.WriteFile(metaFile, metaBytes, 0644); err != nil {
		sklog.Errorf("Failed to save meta cache: %v", err)
	}
	if err := os.WriteFile(paramsFile, paramsBytes, 0644); err != nil {
		sklog.Errorf("Failed to save params cache: %v", err)
	}

	return nil
}

func (api *wasmApi) metaHandler(w http.ResponseWriter, r *http.Request) {
	if err := api.ensureCache(r.Context()); err != nil {
		httputils.ReportError(w, err, "Failed to ensure cache", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(api.cache.meta); err != nil {
		sklog.Errorf("Failed to write meta response: %v", err)
	}
}

func (api *wasmApi) paramsHandler(w http.ResponseWriter, r *http.Request) {
	if err := api.ensureCache(r.Context()); err != nil {
		httputils.ReportError(w, err, "Failed to ensure cache", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(api.cache.params); err != nil {
		sklog.Errorf("Failed to write params response: %v", err)
	}
}

func (api *wasmApi) tracesHandler(w http.ResponseWriter, r *http.Request) {
	if err := api.ensureCache(r.Context()); err != nil {
		httputils.ReportError(w, err, "Failed to ensure cache", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	if _, err := w.Write(api.cache.traces); err != nil {
		sklog.Errorf("Failed to write traces response: %v", err)
	}
}

func (api *wasmApi) buildLookup() (map[string]map[string]uint16, int, []Param) {
	ps := api.psRefresher.GetAll()
	lookup := map[string]map[string]uint16{}
	var idCounter uint16 = 1
	var params []Param

	for key, values := range ps {
		lookup[key] = map[string]uint16{}
		for _, val := range values {
			id := idCounter
			idCounter++
			params = append(params, Param{Id: id, Key: key, Value: val})
			lookup[key][val] = id
		}
	}

	stride := len(ps)
	if stride%8 != 0 {
		stride = (stride/8 + 1) * 8
	}
	return lookup, stride, params
}

func filterNonEmptyKeys(allKeys []string, nonEmptyTraceIDs [][]byte) map[string]bool {
	traceNameMap := map[[16]byte]string{}
	for _, key := range allKeys {
		hash := md5.Sum([]byte(key))
		traceNameMap[hash] = key
	}

	nonEmptyKeys := map[string]bool{}
	for _, id := range nonEmptyTraceIDs {
		var idArray [16]byte
		copy(idArray[:], id)
		if key, ok := traceNameMap[idArray]; ok {
			nonEmptyKeys[key] = true
		}
	}
	return nonEmptyKeys
}

func encodeTraces(allParams []paramtools.Params, allKeys []string, nonEmptyKeys map[string]bool, lookup map[string]map[string]uint16, stride int) ([]byte, int) {
	var tracesBinary []byte
	traceCount := 0
	for idx, p := range allParams {
		key := allKeys[idx]
		if !nonEmptyKeys[key] {
			continue
		}

		row := make([]uint16, stride)
		i := 0
		for key, val := range p {
			if l, ok := lookup[key]; ok {
				if id, ok := l[val]; ok {
					row[i] = id
					i++
				}
			}
		}

		for _, v := range row {
			tracesBinary = append(tracesBinary, byte(v), byte(v>>8))
		}
		traceCount++
	}
	return tracesBinary, traceCount
}
