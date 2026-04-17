import { WasmExports, TraceData } from './worker-types';

export const BATCH_SIZE = 50000;
export const YIELD_INTERVAL_MS = 12;

export async function yieldToMain() {
  if (typeof globalThis !== 'undefined' && (globalThis as any).scheduler?.yield) {
    await (globalThis as any).scheduler.yield();
  } else {
    await new Promise((resolve) => setTimeout(resolve, 0));
  }
}

export async function runWasmBatch(
  wasmFilter: WasmExports,
  traceData: TraceData,
  queryPtr: number,
  queryLen: number,
  outputLimit: number,
  totalQueryValues: number,
  checkInterrupt: () => boolean,
  yieldFn: () => Promise<void> = yieldToMain
): Promise<number> {
  let processed = 0;
  let matchCount = 0;
  let lastYield = performance.now();

  // Choose the best filter implementation based on the number of query values.
  // If there are many values (e.g., wildcards expanded), the bitmask approach is much faster.
  const useBitmask = totalQueryValues > 1;

  while (processed < traceData.numTraces) {
    if (checkInterrupt()) return -1;

    const batch = Math.min(BATCH_SIZE, traceData.numTraces - processed);

    let batchMatches = 0;
    if (useBitmask) {
      batchMatches = wasmFilter.filterTracesBitmask(
        batch,
        processed,
        traceData.stride,
        queryPtr,
        queryLen,
        traceData.dataPtr,
        traceData.outPtr,
        matchCount,
        outputLimit,
        traceData.matchingParamsPtr,
        traceData.bitsetSize
      );
    } else {
      batchMatches = wasmFilter.filterTraces(
        batch,
        processed,
        traceData.stride,
        queryPtr,
        queryLen,
        traceData.dataPtr,
        traceData.outPtr,
        matchCount,
        outputLimit,
        traceData.matchingParamsPtr,
        traceData.bitsetSize
      );
    }

    matchCount += batchMatches;
    processed += batch;

    if (processed < traceData.numTraces && performance.now() - lastYield > YIELD_INTERVAL_MS) {
      await yieldFn();
      lastYield = performance.now();
    }
  }
  return matchCount;
}

export async function scanWasmBatch(
  wasmFilter: WasmExports,
  traceData: TraceData,
  queryPtr: number,
  queryLen: number,
  outputLimit: number,
  totalQueryValues: number,
  checkInterrupt: () => boolean,
  yieldFn: () => Promise<void> = yieldToMain
): Promise<number> {
  let processed = 0;
  let matchCount = 0;
  let lastYield = performance.now();

  const useBitmask = totalQueryValues > 1;

  while (processed < traceData.numTraces) {
    if (checkInterrupt()) return -1;

    const batch = Math.min(BATCH_SIZE, traceData.numTraces - processed);

    let batchMatches = 0;
    if (useBitmask) {
      batchMatches = wasmFilter.scanTracesBitmask(
        batch,
        processed,
        traceData.stride,
        queryPtr,
        queryLen,
        traceData.dataPtr,
        traceData.outPtr,
        matchCount,
        outputLimit,
        traceData.bitsetSize
      );
    } else {
      batchMatches = wasmFilter.scanTraces(
        batch,
        processed,
        traceData.stride,
        queryPtr,
        queryLen,
        traceData.dataPtr,
        traceData.outPtr,
        matchCount,
        outputLimit
      );
    }

    matchCount += batchMatches;
    processed += batch;

    if (processed < traceData.numTraces && performance.now() - lastYield > YIELD_INTERVAL_MS) {
      await yieldFn();
      lastYield = performance.now();
    }
  }
  return matchCount;
}
