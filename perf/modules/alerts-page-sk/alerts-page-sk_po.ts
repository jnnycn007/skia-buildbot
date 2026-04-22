import { PageObject } from '../../../infra-sk/modules/page_object/page_object';
import {
  PageObjectElement,
  PageObjectElementList,
} from '../../../infra-sk/modules/page_object/page_object_element';

/** A page object for the alerts-page-sk component. */
export class AlertsPageSkPO extends PageObject {
  get newButton(): PageObjectElement {
    return this.bySelector('button.action');
  }

  get showDeletedCheckbox(): PageObjectElement {
    return this.bySelector('#showDeletedConfigs');
  }

  get editButton(): PageObjectElement {
    return this.bySelector('create-icon-sk');
  }

  get dialog(): PageObjectElement {
    return this.bySelector('dialog');
  }

  get cancelButton(): PageObjectElement {
    return this.bySelector('button.cancel');
  }

  get acceptButton(): PageObjectElement {
    return this.bySelector('button.accept');
  }

  get deleteButton(): PageObjectElement {
    return this.bySelector('delete-icon-sk');
  }

  get dialogSelector(): string {
    return 'dialog[open]';
  }

  async isDialogOpen(): Promise<boolean> {
    return this.dialog.applyFnToDOMNode((el) => (el as HTMLDialogElement).open);
  }

  private get rows(): PageObjectElementList {
    return this.bySelectorAll('section#light-section table#alerts-table > tbody > tr');
  }

  async getRowCount(): Promise<number> {
    return await this.rows.length;
  }

  async getTableContent(): Promise<string> {
    const table = await this.bySelector('section#light-section table#alerts-table');
    return table.innerText;
  }

  get noAlertsWarning(): PageObjectElement {
    return this.bySelector('div.warning');
  }
}
