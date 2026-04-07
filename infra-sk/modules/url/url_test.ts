import { assert } from 'chai';
import { rootDomain } from './index';

describe('rootDomain', () => {
  it('returns skia.org for perf.skia.org', () => {
    assert.equal(rootDomain('perf.skia.org'), 'skia.org');
  });

  it('returns chrome-perf.corp.goog for chrome-perf.corp.goog', () => {
    assert.equal(rootDomain('chrome-perf.corp.goog'), 'chrome-perf.corp.goog');
  });

  it('returns example.com for example.com', () => {
    assert.equal(rootDomain('example.com'), 'example.com');
  });
});
