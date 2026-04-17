export interface FetchRequest {
  ids: string[];
  min?: number;
  max?: number;
  order?: 'ASC' | 'DESC';
}

export function calculateFetchRequests(
  visibleIds: string[],
  loadedIds: Set<string>,
  viewRange: { min: number; max: number } | null,
  loadedBounds: Record<string, { min: number; max: number }> | null,
  globalBounds: Record<string, { min: number; max: number }> | null,
  getPrimaryKey: (id: string) => string = (id) => id,
  isDateMode: boolean = false,
  pendingRequests: FetchRequest[] = []
): FetchRequest[] {
  const requests: FetchRequest[] = [];

  // 1. Missing IDs (Box Fetch)
  const missingIds = visibleIds.filter((id) => !loadedIds.has(getPrimaryKey(id)));
  const uniqueMissing = missingIds.filter((id) => {
    return !pendingRequests.some((req) => req.ids.includes(id));
  });

  if (uniqueMissing.length > 0) {
    console.log(
      `[FetchLogic] Queueing box fetch for ${uniqueMissing.length} missing IDs.`,
      uniqueMissing
    );
    return [
      {
        ids: uniqueMissing,
        min: viewRange?.min,
        max: viewRange?.max,
      },
    ];
  }

  // 2. Gaps for Existing IDs (Directional Fetch)
  if (viewRange && loadedBounds) {
    const presentIds = visibleIds.filter((id) => loadedIds.has(getPrimaryKey(id)));

    const leftIds = new Set<string>();
    let leftMax = -Infinity;

    const rightIds = new Set<string>();
    let rightMin = Infinity;

    for (const id of presentIds) {
      const primaryId = getPrimaryKey(id);
      const lBounds = loadedBounds[primaryId];
      const gBounds = globalBounds ? globalBounds[primaryId] : null;

      if (lBounds) {
        // Check Left Gap (History)
        if (viewRange.min < lBounds.min) {
          if (!gBounds || lBounds.min > gBounds.min) {
            const boundary = lBounds.min - 1;
            leftIds.add(id);
            if (boundary > leftMax) leftMax = boundary;
          }
        }

        // Check Right Gap (Future)
        if (viewRange.max > lBounds.max) {
          if (!gBounds || lBounds.max < gBounds.max) {
            const boundary = lBounds.max + 1;
            rightIds.add(id);
            if (boundary < rightMin) rightMin = boundary;
          }
        }
      }
    }

    if (leftIds.size > 0) {
      const filteredIds = Array.from(leftIds).filter((id) => {
        return !pendingRequests.some((req) => req.order === 'DESC' && req.ids.includes(id));
      });
      if (filteredIds.length > 0) {
        const prefetch = isDateMode
          ? 7 * 24 * 3600
          : Math.max(200, (viewRange.max - viewRange.min) * 2);
        console.log(
          `[FetchLogic] Queueing left (DESC) directional fetch for ${
            filteredIds.length
          } IDs from max ${leftMax} with min bound ${viewRange.min - prefetch}.`
        );
        requests.push({
          ids: filteredIds,
          max: leftMax,
          min: viewRange.min - prefetch,
          order: 'DESC',
        });
      }
    }

    if (rightIds.size > 0) {
      const filteredIds = Array.from(rightIds).filter((id) => {
        return !pendingRequests.some((req) => req.order === 'ASC' && req.ids.includes(id));
      });
      if (filteredIds.length > 0) {
        const prefetch = isDateMode
          ? 7 * 24 * 3600
          : Math.max(200, (viewRange.max - viewRange.min) * 2);
        console.log(
          `[FetchLogic] Queueing right (ASC) directional fetch for ${
            filteredIds.length
          } IDs from min ${rightMin} with max bound ${viewRange.max + prefetch}.`
        );
        requests.push({
          ids: filteredIds,
          min: rightMin,
          max: viewRange.max + prefetch,
          order: 'ASC',
        });
      }
    }
  }

  return requests;
}
