export interface Param {
  id: number;
  key: string;
  value: string;
}

export interface TraceData {
  memory: WebAssembly.Memory;
  paramSets: Uint16Array;
  matchingParams: Int32Array;
  filteredTraceIndices: Int32Array;
  stride: number;
  numTraces: number;
  maxParamId: number;
  dataPtr: number;
  matchingParamsPtr: number;
  outPtr: number;
  bitsetSize: number;
}

export interface Query {
  [key: string]: string[];
}

export interface WasmExports {
  heap_base: { value: number };
  memory: WebAssembly.Memory;
  filterTraces(
    batchSize: number,
    startTraceIndex: number,
    stride: number,
    queryPtr: number,
    queryLen: number,
    dataPtr: number,
    outPtr: number,
    outWriteOffset: number,
    limit: number,
    matchingParamsPtr: number,
    bitsetSize: number
  ): number;
  scanTraces(
    batchSize: number,
    startTraceIndex: number,
    stride: number,
    queryPtr: number,
    queryLen: number,
    dataPtr: number,
    outPtr: number,
    outWriteOffset: number,
    limit: number
  ): number;
  filterTracesBitmask(
    batchSize: number,
    startTraceIndex: number,
    stride: number,
    queryPtr: number,
    queryLen: number,
    dataPtr: number,
    outPtr: number,
    outWriteOffset: number,
    limit: number,
    matchingParamsPtr: number,
    bitsetSize: number
  ): number;
  scanTracesBitmask(
    batchSize: number,
    startTraceIndex: number,
    stride: number,
    queryPtr: number,
    queryLen: number,
    dataPtr: number,
    outPtr: number,
    outWriteOffset: number,
    limit: number,
    bitsetSize: number
  ): number;
}

export interface SearchCache {
  query: string;
  contextStr: string;
  indices: Int32Array | null;
}
