import { PageObject } from '../../../infra-sk/modules/page_object/page_object';
import { PageObjectElement } from '../../../infra-sk/modules/page_object/page_object_element';

export class RegressionPageSkPO extends PageObject {
  get sheriffSelect(): PageObjectElement {
    return this.bySelector('select[id^="filter-"]');
  }

  get triagedButton(): PageObjectElement {
    return this.bySelector('#btnTriaged');
  }

  get improvementsButton(): PageObjectElement {
    return this.bySelector('#btnImprovements');
  }

  get showMoreButton(): PageObjectElement {
    return this.bySelector('#showMoreAnomalies');
  }

  async selectSheriff(sheriff: string) {
    await this.sheriffSelect.enterValue(sheriff);
  }
}
