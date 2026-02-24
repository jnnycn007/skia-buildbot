import { expect } from 'chai';
import { calculateNudgeList } from './nudge-util';
import { MISSING_DATA_SENTINEL } from '../const/const';
import { ColumnHeader, Anomaly, Trace } from '../json';
import { AnomalyData } from '../common/anomaly-data';

describe('calculateNudgeList', () => {
  const mockHeader = (offsets: number[]): (ColumnHeader | null)[] => {
    return offsets.map(
      (offset) =>
        ({
          offset,
          timestamp: offset,
          hash: '',
          author: '',
          message: '',
          url: '',
        }) as unknown as ColumnHeader
    );
  };

  const mockAnomalyData: AnomalyData = {
    anomaly: {} as Anomaly,
    x: 0,
    y: 0,
    highlight: false,
  };

  it('correctly calculates nudge list for contiguous data', () => {
    // Trace: [10, 20, 30, 40, 50]
    // Header: [1, 2, 3, 4, 5]
    // Current Index: 2 (Value 30, Offset 3)
    const trace = [10, 20, 30, 40, 50] as unknown as Trace;
    const header = mockHeader([1, 2, 3, 4, 5]);
    const currentIndex = 2;

    const nudgeList = calculateNudgeList(trace, header, currentIndex, mockAnomalyData, 1);

    // Expect 3 entries: -1, 0, 1
    expect(nudgeList).to.have.length(3);

    // Entry 0 (Current)
    const currentEntry = nudgeList.find((n) => n.display_index === 0);
    expect(currentEntry).to.exist;
    expect(currentEntry!.x).to.equal(2);
    expect(currentEntry!.selected).to.be.true;

    // Entry -1 (Previous)
    const prevEntry = nudgeList.find((n) => n.display_index === -1);
    expect(prevEntry).to.exist;
    expect(prevEntry!.x).to.equal(1);
    expect(prevEntry!.end_revision).to.equal(2);

    // Entry 1 (Next)
    const nextEntry = nudgeList.find((n) => n.display_index === 1);
    expect(nextEntry).to.exist;
    expect(nextEntry!.x).to.equal(3);
    expect(nextEntry!.end_revision).to.equal(4);
  });

  it('skips missing data points (sparse data)', () => {
    // Trace: [10, MISSING, MISSING, 40, 50]
    // Header: [1, 2, 3, 4, 5]
    // Current Index: 0 (Value 10, Offset 1)
    const trace = [10, MISSING_DATA_SENTINEL, MISSING_DATA_SENTINEL, 40, 50] as unknown as Trace;
    const header = mockHeader([1, 2, 3, 4, 5]);
    const currentIndex = 0;

    const nudgeList = calculateNudgeList(trace, header, currentIndex, mockAnomalyData, 1);

    // Entry 0 (Current)
    const currentEntry = nudgeList.find((n) => n.display_index === 0);
    expect(currentEntry).to.exist;
    expect(currentEntry!.x).to.equal(0);
    expect(currentEntry!.start_revision).to.equal(1);
    expect(currentEntry!.end_revision).to.equal(1);

    // Entry 1 (Next) - Should skip indices 1 and 2, and find index 3 (Value 40)
    const nextEntry = nudgeList.find((n) => n.display_index === 1);
    expect(nextEntry).to.exist;
    expect(nextEntry!.x).to.equal(3);
    // start_revision should be previous VALID point's offset + 1 (header[0].offset + 1 = 2)
    expect(nextEntry!.start_revision).to.equal(2);
    expect(nextEntry!.end_revision).to.equal(4);
    expect(nextEntry!.y).to.equal(40);

    // Entry -1 (Previous) - Out of bounds, should not exist
    const prevEntry = nudgeList.find((n) => n.display_index === -1);
    expect(prevEntry).to.not.exist;
  });

  it('handles the first data point correctly', () => {
    // Trace: [10, 20, 30]
    // Header: [1, 2, 3]
    // Current Index: 0
    const trace = [10, 20, 30] as unknown as Trace;
    const header = mockHeader([1, 2, 3]);
    const currentIndex = 0;

    const nudgeList = calculateNudgeList(trace, header, currentIndex, mockAnomalyData, 1);

    // Entry 0 (Current)
    const currentEntry = nudgeList.find((n) => n.display_index === 0);
    expect(currentEntry).to.exist;
    expect(currentEntry!.start_revision).to.equal(1);
    expect(currentEntry!.end_revision).to.equal(1);
  });

  it('skips missing data points backwards', () => {
    // Trace: [10, 20, MISSING, MISSING, 50]
    // Header: [1, 2, 3, 4, 5]
    // Current Index: 4 (Value 50, Offset 5)
    const trace = [10, 20, MISSING_DATA_SENTINEL, MISSING_DATA_SENTINEL, 50] as unknown as Trace;
    const header = mockHeader([1, 2, 3, 4, 5]);
    const currentIndex = 4;

    const nudgeList = calculateNudgeList(trace, header, currentIndex, mockAnomalyData, 2);

    // Entry -1 (Previous) - Should skip indices 3 and 2, find index 1 (Value 20)
    const prevEntry1 = nudgeList.find((n) => n.display_index === -1);
    expect(prevEntry1).to.exist;
    expect(prevEntry1!.x).to.equal(1);
    // start_revision should be previous VALID point's offset + 1 (header[0].offset + 1 = 2)
    expect(prevEntry1!.start_revision).to.equal(2);
    expect(prevEntry1!.end_revision).to.equal(2);

    // Entry -2 (Previous 2) - Should find index 0 (Value 10)
    const prevEntry2 = nudgeList.find((n) => n.display_index === -2);
    expect(prevEntry2).to.exist;
    expect(prevEntry2!.x).to.equal(0);
    expect(prevEntry2!.start_revision).to.equal(1); // No prior valid point
    expect(prevEntry2!.end_revision).to.equal(1);
  });

  it('handles disjoint traces (bug reproduction)', () => {
    // Scenario from user:
    // Trace A: [1, 3, 5, 7] -> indices in global header: 0, 2, 4, 6
    // Trace B: [2, 4, 6, 8] -> indices in global header: 1, 3, 5, 7
    // Global Header: [1, 2, 3, 4, 5, 6, 7, 8]
    // Trace B Data in Global Frame: [MISSING, 2, MISSING, 4, MISSING, 6, MISSING, 8]

    const header = mockHeader([1, 2, 3, 4, 5, 6, 7, 8]);
    const traceB = [
      MISSING_DATA_SENTINEL, // 1 (Trace A)
      2, // 2 (Trace B)
      MISSING_DATA_SENTINEL, // 3 (Trace A)
      4, // 4 (Trace B)
      MISSING_DATA_SENTINEL, // 5 (Trace A)
      6, // 6 (Trace B) -> Index 5
      MISSING_DATA_SENTINEL, // 7 (Trace A)
      8, // 8 (Trace B) -> Anomaly Here (Index 7)
    ] as unknown as Trace;

    const currentIndex = 7; // Anomaly at 8

    // Test nudge -1
    const nudgeList = calculateNudgeList(traceB, header, currentIndex, mockAnomalyData, 1);

    // Expect nudge -1 to land on index 5 (Value 6)
    const prevEntry = nudgeList.find((n) => n.display_index === -1);
    expect(prevEntry).to.exist;
    expect(prevEntry!.x).to.equal(5);

    // start_revision should be previous VALID point's offset + 1 (header[3].offset + 1 = 4 + 1 = 5)
    expect(prevEntry!.start_revision).to.equal(5);
    expect(prevEntry!.end_revision).to.equal(6);
    expect(prevEntry!.y).to.equal(6);
  });
});
