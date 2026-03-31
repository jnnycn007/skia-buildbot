import { expect } from 'chai';
import { loadCachedTestBed, takeScreenshot, TestBed } from '../../../puppeteer-tests/util';
import { TestPickerSkPO } from './test-picker-sk_po';
import { DEFAULT_VIEWPORT } from '../common/puppeteer-test-util';
import {
  BENCHMARK,
  BOT,
  SUBTEST_1,
  SUBTEST_1_NEW,
  SUBTEST_2,
  SUBTEST_2_NEW,
  TEST,
  TEST_NEW,
} from './test_data';

describe('test-picker-sk', () => {
  let testBed: TestBed;
  before(async () => {
    testBed = await loadCachedTestBed();
  });

  beforeEach(async () => {
    await testBed.page.goto(testBed.baseUrl);
    await testBed.page.setViewport(DEFAULT_VIEWPORT);
  });

  it('should render the component', async () => {
    await testBed.page.waitForSelector('test-picker-sk');
  });

  it('selects items one by one and verifies the query', async () => {
    const testPickerPO = new TestPickerSkPO((await testBed.page.$('test-picker-sk'))!);

    // Wait for the first field to be available.
    await testPickerPO.waitForPickerField(0);
    const benchmarkField = await testPickerPO.getPickerField(0);
    // 'blink_perf.css' is a valid option.
    await benchmarkField.selectExact(BENCHMARK);
    await testPickerPO.waitForSpinnerInactive();
    await takeScreenshot(testBed.page, 'perf', 'test-picker-sk-benchmark-selected');

    // Wait for the next field to appear (Bot).
    await testPickerPO.waitForPickerField(1);
    const botField = await testPickerPO.getPickerField(1);
    await botField.selectExact(BOT);
    await testPickerPO.waitForSpinnerInactive();

    // Wait for the next field (Test).
    await testPickerPO.waitForPickerField(2);
    const testField = await testPickerPO.getPickerField(2);
    await testField.clear();
    await testField.selectExact(TEST);
    await testPickerPO.waitForSpinnerInactive();

    // Wait for the next field (Subtest1).
    await testPickerPO.waitForPickerField(3);
    const subtest1Field = await testPickerPO.getPickerField(3);
    await subtest1Field.selectExact(SUBTEST_1);
    await testPickerPO.waitForSpinnerInactive();

    // Wait for the next field (Subtest2).
    await testPickerPO.waitForPickerField(4);
    const subtest2Field = await testPickerPO.getPickerField(4);
    await subtest2Field.selectExact(SUBTEST_2);
    await testPickerPO.waitForSpinnerInactive();

    // Click the plot button.
    await testPickerPO.clickPlotButton();

    // Verify the query event.
    // In the demo page, the event detail is dumped into <pre id="events">.
    const eventsPre = (await testBed.page.$('#events'))!;
    await testBed.page.waitForFunction(
      (el) => el.textContent && el.textContent.length > 0,
      {},
      eventsPre
    );
    const query = await eventsPre.evaluate((el) => el.textContent);

    const expectedQuery = [
      `benchmark=${BENCHMARK}`,
      `&bot=${BOT}`,
      `&subtest1=${SUBTEST_1}`,
      `&subtest2=${SUBTEST_2}`,
      `&test=${encodeURIComponent(TEST)}`,
    ].join('');

    expect(query).to.equal(expectedQuery);
  });

  it('selects all, deletes middle, and refills with another path', async () => {
    const testPickerPO = new TestPickerSkPO((await testBed.page.$('test-picker-sk'))!);

    // 1. Fill all selectors
    // Benchmark
    await testPickerPO.waitForPickerField(0);
    const benchmarkField = await testPickerPO.getPickerField(0);
    await benchmarkField.selectExact(BENCHMARK);
    await testPickerPO.waitForSpinnerInactive();

    // Bot
    await testPickerPO.waitForPickerField(1);
    const botField = await testPickerPO.getPickerField(1);
    await botField.selectExact(BOT);
    await testPickerPO.waitForSpinnerInactive();

    // Test
    await testPickerPO.waitForPickerField(2);
    const testField = await testPickerPO.getPickerField(2);
    await testField.clear();
    await testField.selectExact(TEST);
    await testPickerPO.waitForSpinnerInactive();

    // Subtest1
    await testPickerPO.waitForPickerField(3);
    const subtest1Field = await testPickerPO.getPickerField(3);
    await subtest1Field.selectExact(SUBTEST_1);
    await testPickerPO.waitForSpinnerInactive();

    // Subtest2
    await testPickerPO.waitForPickerField(4);
    const subtest2Field = await testPickerPO.getPickerField(4);
    await subtest2Field.selectExact(SUBTEST_2);
    await testPickerPO.waitForSpinnerInactive();

    // 2. Delete in the middle (Test field)
    await testField.clear();
    await testPickerPO.waitForSpinnerInactive();

    // 3. Refill with another path
    // Refill Test
    await testField.selectExact(TEST_NEW);
    await testPickerPO.waitForSpinnerInactive();

    // Refill Subtest1
    await testPickerPO.waitForPickerField(3);
    const subtest1FieldNew = await testPickerPO.getPickerField(3);
    await subtest1FieldNew.selectExact(SUBTEST_1_NEW);
    await testPickerPO.waitForSpinnerInactive();

    // Refill Subtest2
    await testPickerPO.waitForPickerField(4);
    const subtest2FieldNew = await testPickerPO.getPickerField(4);
    await subtest2FieldNew.selectExact(SUBTEST_2_NEW);
    await testPickerPO.waitForSpinnerInactive();

    // Click plot
    await testPickerPO.clickPlotButton();

    // Verify
    const eventsPre = (await testBed.page.$('#events'))!;
    await testBed.page.waitForFunction(
      (el) => el.textContent && el.textContent.length > 0,
      {},
      eventsPre
    );
    const query = await eventsPre.evaluate((el) => el.textContent);
    const expectedQuery = [
      `benchmark=${BENCHMARK}`,
      `&bot=${BOT}`,
      `&subtest1=${SUBTEST_1_NEW}`,
      `&subtest2=${SUBTEST_2_NEW}`,
      `&test=${encodeURIComponent(TEST_NEW)}`,
    ].join('');
    expect(query).to.equal(expectedQuery);
  });

  it('does not overlap elements on small screens', async () => {
    // Set to a small viewport to force wrapping.
    await testBed.page.setViewport({ width: 300, height: 800 });

    const testPickerPO = new TestPickerSkPO((await testBed.page.$('test-picker-sk'))!);

    // Populate multiple fields so they wrap.
    await testPickerPO.waitForPickerField(0);
    const benchmarkField = await testPickerPO.getPickerField(0);
    await benchmarkField.selectExact(BENCHMARK);
    await testPickerPO.waitForSpinnerInactive();

    await testPickerPO.waitForPickerField(1);
    const botField = await testPickerPO.getPickerField(1);
    await botField.selectExact(BOT);
    await testPickerPO.waitForSpinnerInactive();

    await testPickerPO.waitForPickerField(2);
    const testField = await testPickerPO.getPickerField(2);
    await testField.clear();
    await testField.selectExact(TEST);
    await testPickerPO.waitForSpinnerInactive();

    // Give it a moment to render fully.
    await new Promise((resolve) => setTimeout(resolve, 100));

    // Retrieve the bounding boxes of the fields.
    const fields = await testBed.page.$$('test-picker-sk picker-field-sk');
    expect(fields.length).to.be.greaterThan(1);

    const boxes = [];
    for (const field of fields) {
      const box = await field.boundingBox();
      if (box) {
        boxes.push(box);
      }
    }

    // Check that no bounding boxes overlap.
    for (let i = 0; i < boxes.length; i++) {
      for (let j = i + 1; j < boxes.length; j++) {
        const box1 = boxes[i];
        const box2 = boxes[j];

        const xOverlap = box1.x < box2.x + box2.width && box1.x + box1.width > box2.x;
        const yOverlap = box1.y < box2.y + box2.height && box1.y + box1.height > box2.y;

        // If they overlap, this will fail.
        expect(xOverlap && yOverlap).to.be.false;
      }
    }

    // Take a screenshot to capture the layout.
    // https://screenshot.googleplex.com/BDxE2JysXAT5yzw
    await takeScreenshot(testBed.page, 'perf', 'test-picker-sk-small-viewport');
  });

  it('expands the next field and NOT the current one when a selection is made', async () => {
    const testPickerPO = new TestPickerSkPO((await testBed.page.$('test-picker-sk'))!);

    // Wait for the first field to be available.
    await testPickerPO.waitForPickerField(0);
    const benchmarkField = await testPickerPO.getPickerField(0);

    // Make a selection in the first field.
    await benchmarkField.selectExact(BENCHMARK);
    await testPickerPO.waitForSpinnerInactive();

    // Verify the first field is CLOSED.
    expect(await benchmarkField.isOpened()).to.be.false;

    // Wait for the next field to appear (Bot).
    await testPickerPO.waitForPickerField(1);
    const botField = await testPickerPO.getPickerField(1);

    // Give it a moment to expand using robust polling.
    await testPickerPO.waitForFieldOpened(1);

    // Verify the second field is OPENED.
    expect(await botField.isOpened()).to.be.true;
  });
});
