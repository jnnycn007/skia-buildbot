import { expect } from 'chai';
import { loadCachedTestBed, takeScreenshot, TestBed } from '../../../puppeteer-tests/util';
import { Page } from 'puppeteer';
import { assert } from 'chai';

describe('new-bug-dialog-sk', () => {
  let testBed: TestBed;

  const mockResponses = {
    '/_/login/status': {
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({ email: 'test@google.com', roles: ['editor'] }),
    },
    '/_/triage/file_bug': {
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({ bug_id: 358011161 }),
    },
  };

  before(async () => {
    testBed = await loadCachedTestBed();
  });

  beforeEach(async () => {
    await testBed.page.setRequestInterception(true);
    testBed.page.on('request', (request) => {
      const matchingPath = Object.keys(mockResponses).find((path) => request.url().endsWith(path));
      if (matchingPath) {
        request.respond(mockResponses[matchingPath as keyof typeof mockResponses]);
      } else {
        request.continue();
      }
    });
    await testBed.page.goto(testBed.baseUrl);
  });

  afterEach(async () => {
    testBed.page.removeAllListeners('request');
    await testBed.page.setRequestInterception(false);
  });

  it('should render the demo page', async () => {
    expect(await testBed.page.$$('new-bug-dialog-sk')).to.have.length(1);
    await takeScreenshot(testBed.page, 'perf', 'new-bug-dialog-sk');
  });

  it('should redirect to Buganizer when New Bug button is clicked', async () => {
    const fileBug = async () => await testBed.page.click('#file-bug');

    // Start waiting for the popup and click the link at the same time.
    const [buganizerPage] = await Promise.all([
      new Promise<Page>((resolve) => testBed.page.once('popup', resolve)),
      fileBug(),
    ]);

    assert.isNotNull(buganizerPage);
    expect(buganizerPage.url()).to.include('https://issues.chromium.org/issues/358011161');
    await buganizerPage.close();
  });
});
