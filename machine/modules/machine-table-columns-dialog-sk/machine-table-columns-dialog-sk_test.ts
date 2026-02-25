import './index';
import { assert } from 'chai';
import { $$ } from '../../../infra-sk/modules/dom';
import { CheckOrRadio } from '../../../elements-sk/modules/checkbox-sk/checkbox-sk';
import {
  ColumnOrder,
  ColumnTitles,
  MachineTableColumnsDialogSk,
} from './machine-table-columns-dialog-sk';

import { setUpElementUnderTest } from '../../../infra-sk/modules/test_util';

describe('machine-table-columns-dialog-sk', () => {
  const newInstance = setUpElementUnderTest<MachineTableColumnsDialogSk>(
    'machine-table-columns-dialog-sk'
  );

  let element: MachineTableColumnsDialogSk;
  beforeEach(() => {
    element = newInstance();
  });

  it('returns undefined on cancel', async () => {
    const promise = element.edit(['Annotation']);
    $$<HTMLButtonElement>('#cancel', element)!.click();
    assert.isUndefined(await promise);
  });

  it('returns original hidden list on OK if no checkboxes are clicked', async () => {
    const startValue: ColumnTitles[] = ['Annotation'];
    const promise = element.edit(startValue);
    $$<HTMLButtonElement>('#ok', element)!.click();
    const result = await promise;
    assert.deepEqual(result, startValue);
  });

  it('returns modified hidden list on OK if a checkbox is clicked', async () => {
    const startValue: ColumnTitles[] = ['Annotation'];
    const promise = element.edit(startValue);

    // Click first checkbox.
    const checkbox = $$<CheckOrRadio>('checkbox-sk', element)!;
    await checkbox.updateComplete;
    checkbox.click();
    $$<HTMLButtonElement>('#ok', element)!.click();
    const result = await promise;
    const expected = [ColumnOrder[0], ...startValue];
    assert.deepEqual(result, expected);
  });
});
