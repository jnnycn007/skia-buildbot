import { expect } from 'chai';
import { loadCachedTestBed, takeScreenshot, TestBed } from '../../../puppeteer-tests/util';
import { Page } from 'puppeteer';
import { NewBugDialogSkPO } from './new-bug-dialog-sk_po';

describe('new-bug-dialog-sk', () => {
  let testBed: TestBed;
  let newBugDialogSkPO: NewBugDialogSkPO;

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
    newBugDialogSkPO = new NewBugDialogSkPO((await testBed.page.$('new-bug-dialog-sk'))!);
  });

  afterEach(async () => {
    testBed.page.removeAllListeners('request');
    await testBed.page.setRequestInterception(false);
  });

  it('should render the demo page', async () => {
    expect(await testBed.page.$$('new-bug-dialog-sk')).to.have.length(1);
    await takeScreenshot(testBed.page, 'perf', 'new-bug-dialog-sk');
  });

  it('should open and close the dialog', async () => {
    // The dialog is closed by default.
    expect(await newBugDialogSkPO.isDialogOpen()).to.be.false;

    // Click the button in the demo page to open the dialog.
    await testBed.page.click('#open-dialog');
    expect(await newBugDialogSkPO.isDialogOpen()).to.be.true;

    // Click the close icon inside the dialog.
    await newBugDialogSkPO.closeIcon.click();
    expect(await newBugDialogSkPO.isDialogOpen()).to.be.false;
  });

  it('should toggle label checkboxes', async () => {
    // Click the button in the demo page to open the dialog.
    await testBed.page.click('#open-dialog');
    expect(await newBugDialogSkPO.isDialogOpen()).to.be.true;

    // Find the first label checkbox.
    const checkbox = await testBed.page.$('input.buglabel');
    expect(checkbox).to.not.be.null;

    // It should be checked by default.
    let isChecked = await checkbox!.getProperty('checked').then((p) => p.jsonValue());
    expect(isChecked).to.be.true;

    // Click to uncheck.
    await checkbox!.click();
    isChecked = await checkbox!.getProperty('checked').then((p) => p.jsonValue());
    expect(isChecked).to.be.false;

    // Click to check again.
    await checkbox!.click();
    isChecked = await checkbox!.getProperty('checked').then((p) => p.jsonValue());
    expect(isChecked).to.be.true;
  });

  it('should switch between component radio buttons', async () => {
    // Click the button in the demo page to open the dialog.
    await testBed.page.click('#open-dialog');
    expect(await newBugDialogSkPO.isDialogOpen()).to.be.true;

    // Find the component radio buttons.
    const radios = await testBed.page.$$('input[name="component"]');
    expect(radios).to.have.length(2);

    const [radio1, radio2] = radios;

    // The first radio button should be checked by default.
    let radio1Checked = await radio1.getProperty('checked').then((p) => p.jsonValue());
    let radio2Checked = await radio2.getProperty('checked').then((p) => p.jsonValue());
    expect(radio1Checked).to.be.true;
    expect(radio2Checked).to.be.false;

    // Click the second radio button.
    await radio2.click();
    radio1Checked = await radio1.getProperty('checked').then((p) => p.jsonValue());
    radio2Checked = await radio2.getProperty('checked').then((p) => p.jsonValue());
    expect(radio1Checked).to.be.false;
    expect(radio2Checked).to.be.true;
  });

  it('should submit the form, fire an event, and open a new window', async () => {
    // Click the button in the demo page to open the dialog.
    await testBed.page.click('#open-dialog');
    expect(await newBugDialogSkPO.isDialogOpen()).to.be.true;

    const newPagePromise = new Promise<Page>((resolve) =>
      testBed.page.browser().once('targetcreated', (target) => resolve(target.page()!))
    );

    // Listen for the 'anomaly-changed' event on the dialog element.
    // We need to do this before submitting the form.
    const eventFiredPromise = testBed.page.evaluate(
      () =>
        new Promise((resolve) => {
          const dialog = document.querySelector('new-bug-dialog-sk')!;
          dialog.addEventListener(
            'anomaly-changed',
            (e: any) => {
              resolve(e.detail);
            },
            { once: true }
          );
        })
    );

    // Submit the form.
    await testBed.page.evaluate(() => {
      const dialog = document.querySelector('new-bug-dialog-sk')!;
      const form = dialog.querySelector<HTMLFormElement>('#new-bug-form')!;
      form.requestSubmit();
    });

    // Wait for the event and the new page.
    const eventDetail: any = await eventFiredPromise;
    const newPage = await newPagePromise;

    // Check that the event was fired with the correct detail.
    expect(eventDetail.bugId).to.equal(358011161);

    // Check that a new page was opened with the correct URL.
    expect(newPage.url()).to.equal('https://issues.chromium.org/issues/358011161');
    await newPage.close();

    // Also check that the dialog is now closed.
    expect(await newBugDialogSkPO.isDialogOpen()).to.be.false;
  });
});
