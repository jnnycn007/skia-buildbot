import { Anomaly } from '../json';
import { AnomalyGroup } from './grouping';
import { getPercentChange } from '../common/anomaly';

/**
 * Shape of the processed anomaly data for display.
 */
export interface ProcessedAnomaly {
  bugId: number;
  revision: number;
  bot: string;
  testsuite: string;
  test: string;
  delta: number;
  isImprovement: boolean;
}

export class AnomalyTransformer {
  /**
   * Transforms a raw Anomaly into a display-ready ProcessedAnomaly object.
   */
  static getProcessedAnomaly(anomaly: Anomaly): ProcessedAnomaly {
    const bugId = anomaly.bug_id;
    const testPathPieces = anomaly.test_path.split('/');
    const bot = testPathPieces[1];
    const testsuite = testPathPieces[2];
    const test = testPathPieces.slice(3, testPathPieces.length).join('/');
    const revision = anomaly.start_revision;
    const delta = getPercentChange(anomaly.median_before_anomaly, anomaly.median_after_anomaly);
    return {
      bugId,
      revision,
      bot,
      testsuite,
      test,
      delta,
      isImprovement: anomaly.is_improvement,
    };
  }

  /**
   * Computes a readable revision range string.
   */
  static computeRevisionRange(start: number | null, end: number | null): string {
    if (start === null || end === null) {
      return '';
    }
    if (start === end) {
      return '' + end;
    }
    return start + ' - ' + end;
  }

  /**
   * Finds the longest common test path prefix for a list of anomalies.
   * Used for group summary rows.
   */
  static findLongestSubTestPath(anomalyList: Anomaly[]): string {
    if (anomalyList.length === 0) {
      return '';
    }

    const getTestPartTokens = (testPath: string) => {
      return testPath.split('/').slice(3);
    };

    let commonTokens = getTestPartTokens(anomalyList[0]!.test_path);

    if (commonTokens.length === 0 || (commonTokens.length === 1 && commonTokens[0] === '')) {
      return '*';
    }

    for (let i = 1; i < anomalyList.length; i++) {
      const currentTokens = getTestPartTokens(anomalyList[i].test_path);
      let matchCount = 0;
      for (let j = 0; j < commonTokens.length && j < currentTokens.length; j++) {
        if (commonTokens[j] === currentTokens[j]) {
          matchCount++;
        } else {
          break;
        }
      }
      commonTokens = commonTokens.slice(0, matchCount);

      if (commonTokens.length === 0) {
        return '*';
      }
    }

    const originalTokens = getTestPartTokens(anomalyList[0]!.test_path);
    if (commonTokens.length !== originalTokens.length) {
      return commonTokens.join('/') + '/*';
    }
    return commonTokens.join('/');
  }

  /**
   * Determines the summary delta for a group of anomalies.
   * Returns [deltaValue, isRegression].
   */
  static determineSummaryDelta(anomalyGroup: AnomalyGroup): [number, boolean] {
    const regressions = anomalyGroup.anomalies.filter((a) => !a.is_improvement);
    let targetAnomalies = anomalyGroup.anomalies;
    if (regressions.length > 0) {
      // If there are regressions, find the one with the largest magnitude.
      targetAnomalies = regressions;
    }

    if (targetAnomalies.length === 0) {
      return [0, false];
    }

    const biggestChangeAnomaly = targetAnomalies.reduce((prev, current) => {
      const prevDelta = Math.abs(
        getPercentChange(prev.median_before_anomaly, prev.median_after_anomaly)
      );
      const currentDelta = Math.abs(
        getPercentChange(current.median_before_anomaly, current.median_after_anomaly)
      );
      return prevDelta > currentDelta ? prev : current;
    });

    return [
      getPercentChange(
        biggestChangeAnomaly.median_before_anomaly,
        biggestChangeAnomaly.median_after_anomaly
      ),
      regressions.length > 0,
    ];
  }
}
