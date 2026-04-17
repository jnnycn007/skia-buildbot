import { expect } from 'chai';
import { appendTraces, filterTraces } from './worker_logic';
import { TraceData, Query } from './worker-types';

describe('worker_logic', () => {
  it('should append traces to TraceData', () => {
    const memory = new WebAssembly.Memory({ initial: 1 });
    const stride = 8;
    const dataPtr = 0;
    const paramSets = new Uint16Array(memory.buffer, dataPtr, 10 * stride);

    const traceData: TraceData = {
      memory,
      paramSets,
      matchingParams: new Int32Array(10),
      filteredTraceIndices: new Int32Array(10),
      stride,
      numTraces: 0,
      maxParamId: 10,
      dataPtr,
      matchingParamsPtr: 100,
      outPtr: 200,
      bitsetSize: 10,
    };

    const newTraces = new Uint16Array([1, 2, 3, 4, 5, 6, 7, 8]); // 1 trace

    appendTraces(traceData, newTraces.buffer);

    expect(traceData.numTraces).to.equal(1);
    expect(traceData.paramSets[0]).to.equal(1);
    expect(traceData.paramSets[7]).to.equal(8);
  });

  it('should NOT request traces when empty', async () => {
    const memory = new WebAssembly.Memory({ initial: 1 });
    const stride = 8;
    const dataPtr = 0;
    const paramSets = new Uint16Array(memory.buffer, dataPtr, 10 * stride);

    const traceData: TraceData = {
      memory,
      paramSets,
      matchingParams: new Int32Array(10),
      filteredTraceIndices: new Int32Array(10),
      stride,
      numTraces: 0, // Empty
      maxParamId: 10,
      dataPtr,
      matchingParamsPtr: 100,
      outPtr: 200,
      bitsetSize: 10,
    };

    const queries: Query[] = [{ benchmark: ['motionmark'] }];

    let fetchCalled = false;
    const fetchTraces = async (_query: string) => {
      fetchCalled = true;
      return new Uint16Array([1, 2, 3, 4, 5, 6, 7, 8]).buffer;
    };

    await filterTraces(queries, traceData, fetchTraces);

    expect(fetchCalled).to.be.false;
    expect(traceData.numTraces).to.equal(0);
  });
});
