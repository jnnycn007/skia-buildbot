import { expect } from 'chai';
import { loadCachedTestBed, takeScreenshot, TestBed } from '../../../puppeteer-tests/util';
import { AlertsPageSkPO } from './alerts-page-sk_po';

describe('alerts-page-sk', () => {
  let testBed: TestBed;
  let alertsPageSkPO: AlertsPageSkPO;

  before(async () => {
    testBed = await loadCachedTestBed();
  });

  beforeEach(async () => {
    await testBed.page.goto(testBed.baseUrl);
    alertsPageSkPO = new AlertsPageSkPO((await testBed.page.$('alerts-page-sk'))!);
    await testBed.page.setViewport({ width: 1200, height: 800 });
  });

  afterEach(async () => {
    // Ensure interception is off after every test to prevent side effects
    await testBed.page.setRequestInterception(false);
    // Remove listeners to prevent memory leaks and duplicate triggers
    testBed.page.removeAllListeners('request');
  });

  it('should render the demo page', async () => {
    // Smoke test.
    expect(await testBed.page.$$('alerts-page-sk')).to.have.length(2);
  });

  it('opens the dialog when Edit button is clicked', async () => {
    await alertsPageSkPO.editButton.click();
    // Wait for dialog to be visible and assert it is open.
    await testBed.page.waitForSelector(alertsPageSkPO.dialogSelector);
    expect(await alertsPageSkPO.isDialogOpen()).to.be.true;
  });

  it('opens the dialog when New button is clicked', async () => {
    await alertsPageSkPO.newButton.click();
    // Wait for dialog to be visible and assert it is open.
    await testBed.page.waitForSelector(alertsPageSkPO.dialogSelector);
    expect(await alertsPageSkPO.isDialogOpen()).to.be.true;
  });

  it('closes the dialog when cancel is clicked', async () => {
    await alertsPageSkPO.newButton.click();
    // Wait for dialog to be visible.
    await testBed.page.waitForSelector(alertsPageSkPO.dialogSelector);
    expect(await alertsPageSkPO.isDialogOpen()).to.be.true;
    // Click the Cancel button.
    await alertsPageSkPO.cancelButton.click();
    // The dialog should be closed.
    expect(await alertsPageSkPO.isDialogOpen()).to.be.false;
  });

  it('updates the list when Show deleted configs is clicked', async () => {
    await alertsPageSkPO.showDeletedCheckbox.click();

    // Give the table a chance to redraw.
    await new Promise((resolve) => setTimeout(resolve, 1000));

    // Verify the deleted Foo row is displayed.
    const tableContent = await alertsPageSkPO.getTableContent();
    expect(tableContent).to.contain('Foo');
  });

  it('saves the alert when accept is clicked', async () => {
    await alertsPageSkPO.editButton.click();
    // Wait for dialog to be visible.
    await testBed.page.waitForSelector(alertsPageSkPO.dialogSelector);
    expect(await alertsPageSkPO.isDialogOpen()).to.be.true;

    await alertsPageSkPO.acceptButton.click();
    expect(await alertsPageSkPO.isDialogOpen()).to.be.false;
  });

  it('deletes an alert', async () => {
    const initialRowCount = await alertsPageSkPO.getRowCount();
    expect(initialRowCount).to.equal(2);

    await alertsPageSkPO.deleteButton.click();

    // Give the table a chance to redraw.
    await new Promise((resolve) => setTimeout(resolve, 1000));

    const afterDeleteRowCount = await alertsPageSkPO.getRowCount();
    expect(afterDeleteRowCount).to.equal(1);

    const warning = await alertsPageSkPO.noAlertsWarning;
    expect(await warning.innerText).to.equal('No alerts have been configured.');
  });

  describe('screenshots', () => {
    it('shows the default view', async () => {
      await takeScreenshot(testBed.page, 'perf', 'alerts-page-sk');
    });

    it('clicks on "New"', async () => {
      await alertsPageSkPO.newButton.click();
      await takeScreenshot(testBed.page, 'perf', 'alerts-page-sk_new_dialog');
    });

    it('clicks on "Show deleted configs"', async () => {
      await alertsPageSkPO.showDeletedCheckbox.click();
      await takeScreenshot(testBed.page, 'perf', 'alerts-page-sk_show_deleted');
    });

    it('clicks on "Edit"', async () => {
      await alertsPageSkPO.editButton.click();
      await takeScreenshot(testBed.page, 'perf', 'alerts-page-sk_edit_dialog');
    });
  });
});
