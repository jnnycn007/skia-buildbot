// Copyright 2018 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

/**
 * @module elements-sk/tabs-sk
 * @description <h2><code>tabs-sk</code></h2>
 *
 * <p>
 * The tabs-sk custom element declaration, used in conjunction with button and
 * the [tabs-panel-sk]{@link module:elements-sk/tabs-panel-sk} element
 * allows you to create tabbed interfaces. The association between the buttons
 * and the tabs displayed in [tabs-panel-sk]{@link module:elements-sk/tabs-panel-sk}
 * is document order, i.e. the first button shows the first panel, second
 * button shows second panel, etc.
 * </p>
 *
 * @example
 *
 * <tabs-sk>
 *   <button class=selected>Query</button>
 *   <button>Results</button>
 * </tabs-sk>
 * <tabs-panel-sk>
 *   <div>
 *     This is the query tab.
 *   </div>
 *   <div>
 *     This is the results tab.
 *   </div>
 * </tabs-panel-sk>
 *
 * @attr selected - The index of the selected tab.
 *
 * @evt tab-selected-sk - Event sent when the user clicks on a tab. The events
 *        value of detail.index is the index of the selected tab.
 *
 */
import { LitElement, PropertyValues, noChange } from 'lit';
import { customElement } from 'lit/decorators.js';
import { property } from 'lit/decorators.js';

export interface TabSelectedSkEventDetail {
  readonly index: number;
}

@customElement('tabs-sk')
export class TabsSk extends LitElement {
  @property({ type: Number, reflect: true })
  selected: number = 0;

  createRenderRoot() {
    return this;
  }

  render() {
    return noChange;
  }

  connectedCallback(): void {
    super.connectedCallback();
    this.addEventListener('click', this);
    // Ensure initial state is consistent
    this._updateChildren();
    this._updatePanel();
  }

  disconnectedCallback(): void {
    super.disconnectedCallback();
    this.removeEventListener('click', this);
  }

  handleEvent(e: Event): void {
    e.stopPropagation();
    this.querySelectorAll('button').forEach((ele, i) => {
      if (ele === e.target) {
        this.select(i, true);
      }
    });
  }

  updated(changedProperties: PropertyValues) {
    if (changedProperties.has('selected')) {
      this._updateChildren();
      this._updatePanel();
    }
  }

  private _updateChildren() {
    this.querySelectorAll('button').forEach((ele, i) => {
      ele.classList.toggle('selected', i === this.selected);
    });
  }

  private _updatePanel() {
    if (this.nextElementSibling?.tagName === 'TABS-PANEL-SK') {
      this.nextElementSibling.setAttribute('selected', String(this.selected));
    }
  }

  /**
   * Force the selection of a tab
   *
   * @param index The index of the tab to select.
   * @param trigger If true then trigger the 'tab-selected-sk' event.
   */
  select(index: number, trigger = false): void {
    this.selected = index;
    if (trigger) {
      this.dispatchEvent(
        new CustomEvent<TabSelectedSkEventDetail>('tab-selected-sk', {
          bubbles: true,
          detail: { index: index },
        })
      );
    }
  }
}
