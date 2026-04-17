export interface ScreenPoint {
  px: number;
  py: number;
  rawPy: number;
  rawX: number;
}

/**
 * Smooths an array of screen points using a bidirectional adaptive Bilateral EMA.
 * Preserves edges (step changes) while aggressively rejecting transient outliers.
 */
export function smoothPoints(
  points: ScreenPoint[],
  smoothingRadius: number,
  edgeDetectionFactor: number,
  edgeLookahead: number
): { smoothed: number[]; std: number[] } {
  const n = points.length;
  if (n === 0) return { smoothed: [], std: [] };
  if (n === 1) return { smoothed: [points[0].py], std: [0] };

  const MIN_PIXEL_NOISE = 2;
  const sf = new Float64Array(n);
  const sb = new Float64Array(n);
  const vf = new Float64Array(n);
  const vb = new Float64Array(n);

  // Forward Pass
  let prevMean = points[0].rawPy;
  let prevVar = MIN_PIXEL_NOISE * MIN_PIXEL_NOISE;
  sf[0] = prevMean;
  vf[0] = prevVar;

  for (let i = 1; i < n; i++) {
    const raw = points[i].rawPy;
    const dx = Math.abs(points[i].px - points[i - 1].px);
    const alphaX = 1 - Math.exp(-Math.max(0, dx) / smoothingRadius);

    const diffCurr = raw - prevMean;
    const stdDev = Math.max(MIN_PIXEL_NOISE, Math.sqrt(prevVar));
    const devCurr = Math.abs(diffCurr) / stdDev;

    let edgeScore = devCurr;
    let allSameSign = true;

    for (let step = 1; step <= edgeLookahead; step++) {
      const targetIdx = Math.min(i + step, n - 1);
      const nextRaw = points[targetIdx].rawPy;
      const diffNext = nextRaw - prevMean;
      const devNext = Math.abs(diffNext) / stdDev;

      if (diffCurr * diffNext <= 0) {
        allSameSign = false;
        break;
      }
      edgeScore = Math.min(edgeScore, devNext);
    }

    if (!allSameSign) edgeScore = 0;

    const edgeFactor = 1 - Math.exp(-Math.pow(edgeScore / 2.0, 2));
    const transientScore = Math.max(0, devCurr - edgeScore);
    const outlierFactor = 1 - Math.exp(-Math.pow(transientScore / 2.0, 2));

    let adaptiveAlphaMean = alphaX * (1 - outlierFactor * 0.95);
    let adaptiveAlphaVar = alphaX * (1 - outlierFactor);

    adaptiveAlphaMean =
      adaptiveAlphaMean * (1 - edgeFactor) +
      Math.max(adaptiveAlphaMean, edgeDetectionFactor) * edgeFactor;
    adaptiveAlphaVar =
      adaptiveAlphaVar * (1 - edgeFactor) +
      Math.max(adaptiveAlphaVar, edgeDetectionFactor) * edgeFactor;

    prevMean = adaptiveAlphaMean * raw + (1 - adaptiveAlphaMean) * prevMean;
    sf[i] = prevMean;

    const newDiffCurr = raw - prevMean;
    prevVar = adaptiveAlphaVar * (newDiffCurr * newDiffCurr) + (1 - adaptiveAlphaVar) * prevVar;
    vf[i] = prevVar;
  }

  // Backward Pass
  prevMean = points[n - 1].rawPy;
  prevVar = MIN_PIXEL_NOISE * MIN_PIXEL_NOISE;
  sb[n - 1] = prevMean;
  vb[n - 1] = prevVar;

  for (let i = n - 2; i >= 0; i--) {
    const raw = points[i].rawPy;
    const dx = Math.abs(points[i + 1].px - points[i].px);
    const alphaX = 1 - Math.exp(-Math.max(0, dx) / smoothingRadius);

    const diffCurr = raw - prevMean;
    const stdDev = Math.max(MIN_PIXEL_NOISE, Math.sqrt(prevVar));
    const devCurr = Math.abs(diffCurr) / stdDev;

    let edgeScore = devCurr;
    let allSameSign = true;

    for (let step = 1; step <= edgeLookahead; step++) {
      const targetIdx = Math.max(i - step, 0);
      const nextRaw = points[targetIdx].rawPy;
      const diffNext = nextRaw - prevMean;
      const devNext = Math.abs(diffNext) / stdDev;

      if (diffCurr * diffNext <= 0) {
        allSameSign = false;
        break;
      }
      edgeScore = Math.min(edgeScore, devNext);
    }

    if (!allSameSign) edgeScore = 0;

    const edgeFactor = 1 - Math.exp(-Math.pow(edgeScore / 2.0, 2));
    const transientScore = Math.max(0, devCurr - edgeScore);
    const outlierFactor = 1 - Math.exp(-Math.pow(transientScore / 2.0, 2));

    let adaptiveAlphaMean = alphaX * (1 - outlierFactor * 0.95);
    let adaptiveAlphaVar = alphaX * (1 - outlierFactor);

    adaptiveAlphaMean =
      adaptiveAlphaMean * (1 - edgeFactor) +
      Math.max(adaptiveAlphaMean, edgeDetectionFactor) * edgeFactor;
    adaptiveAlphaVar =
      adaptiveAlphaVar * (1 - edgeFactor) +
      Math.max(adaptiveAlphaVar, edgeDetectionFactor) * edgeFactor;

    prevMean = adaptiveAlphaMean * raw + (1 - adaptiveAlphaMean) * prevMean;
    sb[i] = prevMean;

    const newDiffCurr = raw - prevMean;
    prevVar = adaptiveAlphaVar * (newDiffCurr * newDiffCurr) + (1 - adaptiveAlphaVar) * prevVar;
    vb[i] = prevVar;
  }

  const smoothed: number[] = [];
  const std: number[] = [];
  for (let i = 0; i < n; i++) {
    smoothed.push((sf[i] + sb[i]) / 2);
    // Take the max variance between forward and backward passes to be conservative, then sqrt for std
    std.push(Math.sqrt(Math.max(vf[i], vb[i])));
  }

  return { smoothed, std };
}
