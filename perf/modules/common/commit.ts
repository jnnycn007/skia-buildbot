/**
 * @module modules/common/commit
 * @description Commit utility functions.
 *
 */

export const TrimHash = (hash: string): string => {
  return hash.substring(0, 9);
};
