import { expect } from 'chai';
import { calculateFetchRequests } from './fetch-logic';

describe('calculateFetchRequests', () => {
  it('should return a box fetch for missing IDs', () => {
    const visibleIds = ['t1', 't2'];
    const loadedIds = new Set<string>();
    const viewRange = { min: 10, max: 20 };
    const requests = calculateFetchRequests(visibleIds, loadedIds, viewRange, null, null);
    expect(requests).to.deep.equal([
      {
        ids: ['t1', 't2'],
        min: 10,
        max: 20,
      },
    ]);
  });

  it('should return a left directional fetch when view range is below loaded bounds', () => {
    const visibleIds = ['t1'];
    const loadedIds = new Set(['t1']);
    const viewRange = { min: 5, max: 15 };
    const loadedBounds = { t1: { min: 10, max: 20 } };
    const requests = calculateFetchRequests(visibleIds, loadedIds, viewRange, loadedBounds, null);
    expect(requests.length).to.equal(1);
    expect(requests[0].order).to.equal('DESC');
    expect(requests[0].ids).to.deep.equal(['t1']);
    expect(requests[0].max).to.equal(9); // loadedBounds.min - 1
    expect(requests[0].min).to.equal(5 - 200); // viewRange.min - prefetch
  });

  it('should return a right directional fetch when view range is above loaded bounds', () => {
    const visibleIds = ['t1'];
    const loadedIds = new Set(['t1']);
    const viewRange = { min: 15, max: 25 };
    const loadedBounds = { t1: { min: 10, max: 20 } };
    const requests = calculateFetchRequests(visibleIds, loadedIds, viewRange, loadedBounds, null);
    expect(requests.length).to.equal(1);
    expect(requests[0].order).to.equal('ASC');
    expect(requests[0].ids).to.deep.equal(['t1']);
    expect(requests[0].min).to.equal(21); // loadedBounds.max + 1
    expect(requests[0].max).to.equal(25 + 200); // viewRange.max + prefetch
  });

  it('should return no requests when all data is loaded', () => {
    const visibleIds = ['t1'];
    const loadedIds = new Set(['t1']);
    const viewRange = { min: 12, max: 18 };
    const loadedBounds = { t1: { min: 10, max: 20 } };
    const requests = calculateFetchRequests(visibleIds, loadedIds, viewRange, loadedBounds, null);
    expect(requests).to.deep.equal([]);
  });
});
