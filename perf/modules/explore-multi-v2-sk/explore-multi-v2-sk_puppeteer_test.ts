import { expect } from 'chai';
import { loadCachedTestBed, TestBed } from '../../../puppeteer-tests/util';

describe('explore-multi-v2-sk', () => {
  let testBed: TestBed;

  before(async () => {
    testBed = await loadCachedTestBed();
  });

  beforeEach(async () => {
    const page = testBed.page;
    page.on('console', (msg) => console.log('PAGE LOG:', msg.text()));

    await page.evaluateOnNewDocument(() => {
      (window as any).WORKER_URL =
        'data:application/javascript,self.postMessage({ type: "LOADED" }); self.onmessage = (e) => { if (e.data.type === "INIT") { self.postMessage({ type: "READY" }); } else if (e.data.type === "SUGGEST") { self.postMessage({ type: "SUGGEST_RESULT", idx: e.data.idx, payload: [{ params: [{ key: "test", value: "Score" }], count: 42, score: 100 }] }); } };';

      // Trap fetchMock when it is assigned to window
      let fm: any;
      Object.defineProperty(window, 'fetchMock', {
        get() {
          return fm;
        },
        set(v) {
          fm = v;
          console.log('MOCK LOG: fetchMock trapped!');
          // Allow falling back to network or original fetch
          if (fm.config) {
            fm.config.fallbackToNetwork = true;
          }
          // Register wasm mocks immediately
          fm.get(
            'glob:*/_/wasm/meta.json*',
            { version: 'test-version' },
            { overwriteRoutes: true }
          );
          fm.get('glob:*/_/wasm/params.json*', {}, { overwriteRoutes: true });
          fm.get('glob:*/dist/explore-multi-v2-sk/filter.wasm*', new ArrayBuffer(0), {
            overwriteRoutes: true,
          });
          fm.get('glob:*/_/wasm/traces.bin*', new ArrayBuffer(0), { overwriteRoutes: true });
        },
        configurable: true,
      });
    });

    await page.goto(testBed.baseUrl);
    await page.waitForSelector('explore-multi-v2-sk');
  });

  it('should trigger fetch when panning in Date Mode', async () => {
    const page = testBed.page;

    // Mock /_/trace_values
    await page.evaluate(() => {
      (window as any).fetchMock.post(
        '/_/trace_values',
        { results: { 'test-trace': [{ x: 1000, y: 10 }] } },
        { overwriteRoutes: true }
      );
    });

    // Toggle Date Mode
    await page.evaluate(() => {
      const explore = document.querySelector('explore-multi-v2-sk');
      const exploreEl = explore as any;
      exploreEl._matchingTraceIds = ['t1', 't2'];
      exploreEl._pageSize = 10;
      exploreEl._tracePage = 0;
      exploreEl._dateMode = true;
      exploreEl._globalBounds = {};
      exploreEl._loadedBounds = {};
      exploreEl._viewportMinX = null;
      exploreEl._viewportMaxX = null;
    });

    // Simulate panning by calling _handleViewportChanged directly
    await page.evaluate(() => {
      const explore = document.querySelector('explore-multi-v2-sk') as any;
      if (explore) {
        explore._handleViewportChanged({
          detail: { minCommit: 500, maxCommit: 1500 },
        });
      }
    });

    // Wait for fetch to be called
    await page.waitForFunction(() => {
      const calls = (window as any).fetchMock.calls('/_/trace_values');
      return calls.length > 0;
    });

    // Verify request body
    const calls = await page.evaluate(() => {
      return (window as any).fetchMock.calls('/_/trace_values').map((c: any) => c[1].body);
    });

    expect(calls.length).to.be.greaterThan(0);
    const body = JSON.parse(calls[0]);
    expect(body.begin).to.be.a('number');
    expect(body.end).to.be.a('number');
    expect(body.begin).to.equal(500);
    expect(body.end).to.equal(1500);
  });

  it('should update suggestion counts when typing', async () => {
    const page = testBed.page;

    // Set availableParams on query-bar-sk
    await page.evaluate(() => {
      const explore = document.querySelector('explore-multi-v2-sk') as any;
      const queryBar = explore.shadowRoot.querySelector('query-bar-sk') as any;
      queryBar.availableParams = [{ id: 1, key: 'test', value: 'Score', count: 0 }];
      queryBar.query = {};
    });

    // Type in the query bar and trigger input event
    await page.evaluate(async () => {
      const explore = document.querySelector('explore-multi-v2-sk') as any;
      const queryBar = explore.shadowRoot.querySelector('query-bar-sk') as any;
      const input = queryBar.shadowRoot.querySelector('md-outlined-text-field');
      input.value = 'Sc';
      input.dispatchEvent(new Event('input'));
      await queryBar.updateComplete;
    });

    // Wait for the count to be updated to (42)
    await page.waitForFunction(() => {
      const explore = document.querySelector('explore-multi-v2-sk') as any;
      const queryBar = explore.shadowRoot.querySelector('query-bar-sk') as any;
      const countEl = queryBar.shadowRoot.querySelector('.s-count.right');
      return countEl && countEl.textContent === '(42)';
    });

    const countText = await page.evaluate(() => {
      const explore = document.querySelector('explore-multi-v2-sk') as any;
      const queryBar = explore.shadowRoot.querySelector('query-bar-sk') as any;
      const countEl = queryBar.shadowRoot.querySelector('.s-count.right');
      return countEl ? countEl.textContent : '';
    });

    expect(countText).to.equal('(42)');
  });

  it('should load worker and become ready', async () => {
    const page = testBed.page;

    const workerReady = await page.evaluate(async () => {
      const explore = document.querySelector('explore-multi-v2-sk') as any;
      // Wait up to 1 second for worker to initialize
      for (let i = 0; i < 20; i++) {
        if (explore._workerController && explore._workerController.isReady()) return true;
        await new Promise((resolve) => setTimeout(resolve, 50));
      }
      return false;
    });

    expect(workerReady).to.be.true;
  });
});
