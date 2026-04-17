// assembly/filter.ts

// Export heap base so JS knows where to put data
// @ts-expect-error: __heap_base is provided by AssemblyScript
export const heap_base = __heap_base;

/**
 * Filter traces and collect faceted search data using SIMD on 16-bit integers.
 */
export function filterTraces(
  batchSize: i32,
  startTraceIndex: i32,
  stride: i32,
  queryPtr: i32,
  _queryLen: i32, // Size of query buffer in Int32 elements
  dataPtr: i32,
  outPtr: i32,
  outWriteOffset: i32, // Where to start writing in outPtr
  limit: i32, // Global limit
  matchingParamsPtr: i32,
  bitsetSize: i32
): i32 {
  let matchCount = 0;
  let collectedCount = 0;

  const numKeys = load<i32>(queryPtr);
  const startPtr = queryPtr + 4;

  const isSingleKey = numKeys === 1;

  for (let i = 0; i < batchSize; ++i) {
    const absIndex = startTraceIndex + i;
    // Offset for Uint16 elements is (start + i) * stride * 2
    const tOffset = dataPtr + absIndex * stride * 2;

    let misses = 0;
    let missedKeyIndex = -1;
    let currPtr = startPtr;

    if (numKeys > 0) {
      for (let k = 0; k < numKeys; ++k) {
        const keyLen = load<i32>(currPtr);
        currPtr += 4;

        let groupMatch = false;

        for (let v = 0; v < keyLen; ++v) {
          const qVal = load<i32>(currPtr);
          currPtr += 4;

          if (groupMatch) continue;

          const qVec = i16x8.splat(<i16>qVal);

          let j = 0;
          for (; j <= stride - 8; j += 8) {
            const tVec = v128.load(tOffset + (j << 1));
            if (v128.any_true(i16x8.eq(tVec, qVec))) {
              groupMatch = true;
              break;
            }
          }

          if (!groupMatch) {
            for (; j < stride; ++j) {
              if (<i32>load<u16>(tOffset + (j << 1)) === qVal) {
                groupMatch = true;
                break;
              }
            }
          }
        }

        if (!groupMatch) {
          misses++;
          missedKeyIndex = k;
          if (misses > 1) break;
        }
      }
    }

    if (misses === 0) {
      if (outWriteOffset + collectedCount < limit) {
        store<i32>(outPtr + (outWriteOffset + collectedCount) * 4, absIndex);
        collectedCount++;
      }
      matchCount++;

      for (let j = 0; j < stride; ++j) {
        const val = <i32>load<u16>(tOffset + j * 2);
        if (val === 0) break;

        // Main: Accumulate counts (i32)
        const mainPtr = matchingParamsPtr + val * 4;
        store<i32>(mainPtr, load<i32>(mainPtr) + 1);

        if (!isSingleKey) {
          for (let k = 0; k < numKeys; ++k) {
            const keyPtr = matchingParamsPtr + ((k + 1) * bitsetSize + val) * 4;
            store<i32>(keyPtr, load<i32>(keyPtr) + 1);
          }
        }
      }
    } else if (misses === 1) {
      if (!isSingleKey) {
        const bitsetOffset = (missedKeyIndex + 1) * bitsetSize;

        for (let j = 0; j < stride; ++j) {
          const val = <i32>load<u16>(tOffset + j * 2);
          if (val === 0) break;

          const ptr = matchingParamsPtr + (bitsetOffset + val) * 4;
          store<i32>(ptr, load<i32>(ptr) + 1);
        }
      }
    }
  }

  return matchCount;
}

/**
 * Read-only scan of traces. Returns matching indices but DOES NOT update bitsets.
 */
export function scanTraces(
  batchSize: i32,
  startTraceIndex: i32,
  stride: i32,
  queryPtr: i32,
  _queryLen: i32,
  dataPtr: i32,
  outPtr: i32,
  outWriteOffset: i32,
  limit: i32
): i32 {
  let matchCount = 0;
  let collectedCount = 0;

  const numKeys = load<i32>(queryPtr);
  const startPtr = queryPtr + 4;

  for (let i = 0; i < batchSize; ++i) {
    const absIndex = startTraceIndex + i;
    const tOffset = dataPtr + absIndex * stride * 2;

    let misses = 0;
    let currPtr = startPtr;

    if (numKeys > 0) {
      for (let k = 0; k < numKeys; ++k) {
        const keyLen = load<i32>(currPtr);
        currPtr += 4;

        let groupMatch = false;

        for (let v = 0; v < keyLen; ++v) {
          const qVal = load<i32>(currPtr);
          currPtr += 4;

          if (groupMatch) continue;

          const qVec = i16x8.splat(<i16>qVal);

          let j = 0;
          for (; j <= stride - 8; j += 8) {
            const tVec = v128.load(tOffset + (j << 1));
            if (v128.any_true(i16x8.eq(tVec, qVec))) {
              groupMatch = true;
              break;
            }
          }

          if (!groupMatch) {
            for (; j < stride; ++j) {
              if (<i32>load<u16>(tOffset + (j << 1)) === qVal) {
                groupMatch = true;
                break;
              }
            }
          }
        }

        if (!groupMatch) {
          misses++;
          break;
        }
      }
    }

    if (misses === 0) {
      if (outWriteOffset + collectedCount < limit) {
        store<i32>(outPtr + (outWriteOffset + collectedCount) * 4, absIndex);
        collectedCount++;
      }
      matchCount++;
    }
  }

  return matchCount;
}

/**
 * Filter traces using a bitmask array. Faster for large queries.
 */
export function filterTracesBitmask(
  batchSize: i32,
  startTraceIndex: i32,
  stride: i32,
  queryPtr: i32,
  queryLen: i32, // Size of query buffer in Int32 elements
  dataPtr: i32,
  outPtr: i32,
  outWriteOffset: i32, // Where to start writing in outPtr
  limit: i32, // Global limit
  matchingParamsPtr: i32,
  bitsetSize: i32
): i32 {
  let matchCount = 0;
  let collectedCount = 0;

  const numKeys = load<i32>(queryPtr);
  const startPtr = queryPtr + 4;

  const isSingleKey = numKeys === 1;

  let queryMapPtr = queryPtr + queryLen * 4;
  queryMapPtr = (queryMapPtr + 7) & ~7; // align to 8 bytes

  // Initialize the bitmask map (only the needed size)
  memory.fill(queryMapPtr, 0, bitsetSize * 8);

  let currPtr = startPtr;
  for (let k = 0; k < numKeys; ++k) {
    const keyLen = load<i32>(currPtr);
    currPtr += 4;
    for (let v = 0; v < keyLen; ++v) {
      const qVal = load<i32>(currPtr);
      currPtr += 4;
      if (qVal < bitsetSize) {
        store<u64>(queryMapPtr + qVal * 8, load<u64>(queryMapPtr + qVal * 8) | ((1 as u64) << k));
      }
    }
  }

  let targetMask: u64 = ((1 as u64) << numKeys) - 1;
  if (numKeys === 0) targetMask = 0;

  for (let i = 0; i < batchSize; ++i) {
    const absIndex = startTraceIndex + i;
    const tOffset = dataPtr + absIndex * stride * 2;

    let matchedKeysMask: u64 = 0;

    for (let j = 0; j < stride; ++j) {
      const val = <i32>load<u16>(tOffset + j * 2);
      if (val === 0) break;
      if (val < bitsetSize) {
        matchedKeysMask |= load<u64>(queryMapPtr + val * 8);
      }
    }

    const misses = numKeys - <i32>popcnt(targetMask & matchedKeysMask);

    if (misses === 0) {
      if (outWriteOffset + collectedCount < limit) {
        store<i32>(outPtr + (outWriteOffset + collectedCount) * 4, absIndex);
        collectedCount++;
      }
      matchCount++;

      for (let j = 0; j < stride; ++j) {
        const val = <i32>load<u16>(tOffset + j * 2);
        if (val === 0) break;

        const mainPtr = matchingParamsPtr + val * 4;
        store<i32>(mainPtr, load<i32>(mainPtr) + 1);

        if (!isSingleKey) {
          for (let k = 0; k < numKeys; ++k) {
            const keyPtr = matchingParamsPtr + ((k + 1) * bitsetSize + val) * 4;
            store<i32>(keyPtr, load<i32>(keyPtr) + 1);
          }
        }
      }
    } else if (misses === 1) {
      if (!isSingleKey) {
        const missedBit = targetMask & ~matchedKeysMask;
        const missedKeyIndex = <i32>ctz(missedBit);
        const bitsetOffset = (missedKeyIndex + 1) * bitsetSize;

        for (let j = 0; j < stride; ++j) {
          const val = <i32>load<u16>(tOffset + j * 2);
          if (val === 0) break;

          const ptr = matchingParamsPtr + (bitsetOffset + val) * 4;
          store<i32>(ptr, load<i32>(ptr) + 1);
        }
      }
    }
  }

  return matchCount;
}

export function scanTracesBitmask(
  batchSize: i32,
  startTraceIndex: i32,
  stride: i32,
  queryPtr: i32,
  queryLen: i32,
  dataPtr: i32,
  outPtr: i32,
  outWriteOffset: i32,
  limit: i32,
  bitsetSize: i32
): i32 {
  let matchCount = 0;
  let collectedCount = 0;

  const numKeys = load<i32>(queryPtr);
  const startPtr = queryPtr + 4;

  let queryMapPtr = queryPtr + queryLen * 4;
  queryMapPtr = (queryMapPtr + 7) & ~7;
  memory.fill(queryMapPtr, 0, bitsetSize * 8);

  let currPtr = startPtr;
  for (let k = 0; k < numKeys; ++k) {
    const keyLen = load<i32>(currPtr);
    currPtr += 4;
    for (let v = 0; v < keyLen; ++v) {
      const qVal = load<i32>(currPtr);
      currPtr += 4;
      if (qVal < bitsetSize) {
        store<u64>(queryMapPtr + qVal * 8, load<u64>(queryMapPtr + qVal * 8) | ((1 as u64) << k));
      }
    }
  }

  let targetMask: u64 = ((1 as u64) << numKeys) - 1;
  if (numKeys === 0) targetMask = 0;

  for (let i = 0; i < batchSize; ++i) {
    const absIndex = startTraceIndex + i;
    const tOffset = dataPtr + absIndex * stride * 2;

    let matchedKeysMask: u64 = 0;

    for (let j = 0; j < stride; ++j) {
      const val = <i32>load<u16>(tOffset + j * 2);
      if (val === 0) break;
      if (val < bitsetSize) {
        matchedKeysMask |= load<u64>(queryMapPtr + val * 8);
      }
    }

    const misses = numKeys - <i32>popcnt(targetMask & matchedKeysMask);

    if (misses === 0) {
      if (outWriteOffset + collectedCount < limit) {
        store<i32>(outPtr + (outWriteOffset + collectedCount) * 4, absIndex);
        collectedCount++;
      }
      matchCount++;
    }
  }

  return matchCount;
}
