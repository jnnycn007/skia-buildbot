import { ColumnHeader, Trace } from '../json';
import { NudgeEntry } from '../triage-menu-sk/triage-menu-sk';
import { MISSING_DATA_SENTINEL } from '../const/const';
import { AnomalyData } from '../common/anomaly-data';

/**
 * Scans the trace in a given direction to find indices containing valid data.
 */
function getValidIndices(
  trace: Trace,
  headerLength: number,
  startIndex: number,
  direction: -1 | 1,
  maxCount: number
): number[] {
  const indices: number[] = [];
  let probe = startIndex + direction;

  while (probe >= 0 && probe < headerLength && indices.length < maxCount) {
    if (trace[probe] !== MISSING_DATA_SENTINEL) {
      indices.push(probe);
    }
    probe += direction;
  }
  return indices;
}

/**
 * Constructs a NudgeEntry based on the target index and its preceding valid index.
 */
function createNudgeEntry(
  header: (ColumnHeader | null)[],
  trace: Trace,
  anomalyData: AnomalyData,
  displayIndex: number,
  targetIndex: number,
  prevValidIndex: number | null,
  xOffset: number
): NudgeEntry {
  // If there's a previous valid point in the trace, start_revision is the commit after it.
  // Otherwise, fallback to the target index's own offset (e.g., the very first point in the trace).
  const start_revision =
    prevValidIndex !== null ? header[prevValidIndex]!.offset + 1 : header[targetIndex]!.offset;

  return {
    display_index: displayIndex,
    anomaly_data: anomalyData,
    selected: displayIndex === 0,
    start_revision: start_revision,
    end_revision: header[targetIndex]!.offset,
    x: targetIndex - xOffset,
    y: trace[targetIndex],
  };
}

/**
 * Calculates the nudge list for an anomaly, ensuring that nudges only target data points
 * belonging to this specific trace (skipping missing data points in the global timeline).
 *
 * The start_revision of a nudged point is strictly calculated as the commit immediately
 * following the *previous valid point* in this trace, properly handling sparse data.
 *
 * @param trace The data trace (array of numbers).
 * @param header The dataframe header (array of ColumnHeader).
 * @param currentIndex The index of the anomaly in the trace/header.
 * @param anomalyData The anomaly data object.
 * @param nudgeRange The range of nudge steps (default 2).
 * @param xOffset Optional offset to subtract from the absolute index for NudgeEntry.x.
 * @returns A list of NudgeEntry objects.
 */
export function calculateNudgeList(
  trace: Trace,
  header: (ColumnHeader | null)[],
  currentIndex: number,
  anomalyData: AnomalyData,
  nudgeRange: number = 2,
  xOffset: number = 0
): NudgeEntry[] {
  const nudgeList: NudgeEntry[] = [];
  const headerLength = header.length;

  if (currentIndex < 0 || currentIndex >= headerLength) {
    return nudgeList;
  }

  // Gather valid indices to the left (need nudgeRange + 1 to calculate start_revision)
  const leftIndices = getValidIndices(trace, headerLength, currentIndex, -1, nudgeRange + 1);
  // Gather valid indices to the right
  const rightIndices = getValidIndices(trace, headerLength, currentIndex, 1, nudgeRange);

  for (let i = -nudgeRange; i <= nudgeRange; i++) {
    if (i < 0) {
      const leftPos = Math.abs(i) - 1; // i = -1 -> index 0
      if (leftPos < leftIndices.length) {
        const targetIndex = leftIndices[leftPos];
        const prevValidIndex = leftPos + 1 < leftIndices.length ? leftIndices[leftPos + 1] : null;
        nudgeList.push(
          createNudgeEntry(header, trace, anomalyData, i, targetIndex, prevValidIndex, xOffset)
        );
      }
    } else if (i === 0) {
      const prevValidIndex = leftIndices.length > 0 ? leftIndices[0] : null;
      nudgeList.push(
        createNudgeEntry(header, trace, anomalyData, 0, currentIndex, prevValidIndex, xOffset)
      );
    } else if (i > 0) {
      const rightPos = i - 1; // i = 1 -> index 0
      if (rightPos < rightIndices.length) {
        const targetIndex = rightIndices[rightPos];
        let prevValidIndex: number | null;

        if (i === 1) {
          // The previous point is either the current index (if valid), or the first valid point to the left.
          prevValidIndex =
            trace[currentIndex] !== MISSING_DATA_SENTINEL
              ? currentIndex
              : leftIndices.length > 0
                ? leftIndices[0]
                : null;
        } else {
          prevValidIndex = rightIndices[rightPos - 1];
        }
        nudgeList.push(
          createNudgeEntry(header, trace, anomalyData, i, targetIndex, prevValidIndex, xOffset)
        );
      }
    }
  }

  return nudgeList;
}
