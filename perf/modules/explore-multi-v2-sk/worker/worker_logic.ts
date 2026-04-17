import { TraceData, Query } from './worker-types';

export function appendTraces(traceData: TraceData, newTraces: ArrayBuffer): void {
  const newTracesView = new Uint16Array(newTraces);
  traceData.paramSets.set(newTracesView, traceData.numTraces * traceData.stride);
  traceData.numTraces += newTracesView.length / traceData.stride;
}

export async function filterTraces(
  _queries: Query[],
  _traceData: TraceData,
  _fetchTraces: (query: string) => Promise<ArrayBuffer>
): Promise<void> {}
