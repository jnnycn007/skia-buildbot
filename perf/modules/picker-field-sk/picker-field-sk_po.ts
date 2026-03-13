import { PageObject } from '../../../infra-sk/modules/page_object/page_object';
import { PageObjectElement } from '../../../infra-sk/modules/page_object/page_object_element';
import { CheckOrRadio } from '../../../elements-sk/modules/checkbox-sk/checkbox-sk';

export class PickerFieldSkPO extends PageObject {
  get comboBox(): PageObjectElement {
    return this.bySelector('vaadin-multi-select-combo-box');
  }

  get splitByCheckbox(): PageObjectElement {
    return this.bySelector('checkbox-sk#split-by');
  }

  async isSplitChecked(): Promise<boolean> {
    return await this.splitByCheckbox.applyFnToDOMNode((el) => (el as CheckOrRadio).checked);
  }

  async checkSplit(): Promise<void> {
    if (!(await this.isSplitChecked())) {
      await this.splitByCheckbox.applyFnToDOMNode((el) => {
        const checkbox = el as any;
        checkbox.checked = true;
        el.dispatchEvent(new Event('change', { bubbles: true, composed: true }));
      });
    }
  }

  async uncheckSplit(): Promise<void> {
    if (await this.isSplitChecked()) {
      await this.splitByCheckbox.applyFnToDOMNode((el) => {
        const checkbox = el as any;
        checkbox.checked = false;
        el.dispatchEvent(new Event('change', { bubbles: true, composed: true }));
      });
    }
  }

  get selectPrimaryCheckbox(): PageObjectElement {
    return this.bySelector('checkbox-sk#select-primary');
  }

  get selectAllCheckbox(): PageObjectElement {
    return this.bySelector('checkbox-sk#select-all');
  }

  async isAllChecked(): Promise<boolean> {
    return await this.selectAllCheckbox.applyFnToDOMNode((el) => (el as CheckOrRadio).checked);
  }

  async checkAll(): Promise<void> {
    if (!(await this.isAllChecked())) {
      await this.selectAllCheckbox.applyFnToDOMNode((el) => {
        const checkbox = el as any;
        checkbox.checked = true;
        el.dispatchEvent(new Event('change', { bubbles: true, composed: true }));
      });
    }
  }

  async getLabel(): Promise<string> {
    const label = await this.comboBox.getAttribute('label');
    return label || '';
  }

  async getSelectedItems(): Promise<string[]> {
    return await this.comboBox.applyFnToDOMNode((n) => (n as any).selectedItems);
  }

  async openOverlay(): Promise<void> {
    await this.comboBox.click();
  }

  async select(value: string): Promise<void> {
    await this.openOverlay();
    await this.comboBox.type(value);
    await this.comboBox.press('Enter');
  }

  async selectExact(value: string): Promise<void> {
    await this.openOverlay();
    // Use evaluate to avoid serialization issues with CustomEvent detail on older puppeteer versions when possible,
    // or just dispatch the event that the component listens to.
    await this.comboBox.applyFnToDOMNode((el: any, val: unknown) => {
      const v = val as string;
      if (!el.selectedItems) {
        el.selectedItems = [];
      }
      if (!el.selectedItems.includes(v)) {
        // Create new array to ensure change detection
        el.selectedItems = [...el.selectedItems, v];
      }

      // The vaadin-multi-select-combo-box component dispatches 'selected-items-changed'
      // when its selection changes. We dispatch it here so `picker-field-sk`'s
      // @selected-items-changed listener catches it just like a user interaction.
      el.dispatchEvent(
        new CustomEvent('selected-items-changed', {
          detail: { value: el.selectedItems },
          bubbles: true,
          composed: true,
        })
      );
      // Close overlay to prevent obscuring DOM elements
      el.opened = false;
    }, value);
  }

  async search(value: string): Promise<void> {
    await this.select(value);
  }

  async clear(): Promise<void> {
    await this.comboBox.applyFnToDOMNode((el: any) => {
      el.selectedItems = [];
      el.dispatchEvent(
        new CustomEvent('selected-items-changed', {
          detail: { value: [] },
          bubbles: true,
          composed: true,
        })
      );
    });
  }

  /**
   * Removes a selected option from the combo box.
   *
   * It searches for a chip with a matching title (or inner text) and clicks its remove button.
   *
   * @param label The label of the option to remove (e.g., "Android").
   */
  async removeSelectedOption(label: string): Promise<void> {
    const chips = this.comboBox.bySelectorAll('vaadin-multi-select-combo-box-chip');
    const chipToRemove = await chips.find(async (chip) => {
      // Check title attribute first (observed behavior).
      const title = await chip.getAttribute('title');
      if (title && title.trim() === label) {
        return true;
      }
      // Fallback to innerText just in case.
      const text = await chip.innerText;
      return text.trim() === label;
    });

    if (!chipToRemove) {
      throw new Error(`Could not find chip with label '${label}'`);
    }

    await chipToRemove.applyFnToDOMNode((el) => {
      const root = el.shadowRoot || el;

      // Try to find the remove button by part name (standard Vaadin).
      const removePart = root.querySelector('[part~="remove-button"]');
      if (removePart) {
        (removePart as HTMLElement).click();
        return;
      }

      throw new Error('Could not find remove button in chip');
    });
  }

  async isDisabled(): Promise<boolean> {
    return await this.comboBox.hasAttribute('readonly');
  }
}
