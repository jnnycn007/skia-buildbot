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

/** @module elements-sk/multi-select-sk
 *
 * @description <h2><code>multi-select-sk</code></h2>
 *
 * <p>
 *   Clicking on the children will cause them to be selected.
 * </p>
 *
 * <p>
 *   The multi-select-sk elements monitors for the addition and removal of child
 *   elements and will update the 'selected' property as needed. Note that it
 *   does not monitor the 'selected' attribute of child elements, and will not
 *   update the 'selected' property if they are changed directly.
 * </p>
 *
 * @example
 *
 *   <multi-select-sk>
 *     <div></div>
 *     <div></div>
 *     <div selected></div>
 *     <div></div>
 *     <div selected></div>
 *   </multi-select-sk>
 *
 * @evt selection-changed - Sent when an item is clicked and the selection is changed.
 *   The detail of the event contains the indices of the children elements:
 *
 *   <pre>
 *     detail: {
 *       selection: [2,4],
 *     }
 *   </pre>
 *
 */
import { LitElement, PropertyValues } from 'lit';
import { customElement, property } from 'lit/decorators.js';

@customElement('multi-select-sk')
export class MultiSelectSk extends LitElement {
  private _obs: MutationObserver;

  private _selection: number[] = [];

  constructor() {
    super();
    // Keep _selection up to date by monitoring DOM changes.
    this._obs = new MutationObserver(() => this._bubbleUp());
  }

  connectedCallback(): void {
    super.connectedCallback();
    this.addEventListener('click', this._click);
    this.observerConnect();
    this._bubbleUp();
  }

  disconnectedCallback(): void {
    super.disconnectedCallback();
    this.removeEventListener('click', this._click);
    this.observerDisconnect();
  }

  createRenderRoot(): this {
    return this;
  }

  observerDisconnect() {
    this._obs.disconnect();
  }

  observerConnect() {
    this._obs.observe(this, {
      subtree: true,
      childList: true,
      attributes: true,
      attributeFilter: ['selected'],
    });
  }

  /** Whether this element should respond to input. */
  @property({ type: Boolean, reflect: true })
  disabled: boolean = false;

  /**
   * A sorted array of indices that are selected or [] if nothing is selected.
   * If selection is set to a not sorted array, it will be sorted anyway.
   */
  @property({ type: Array, attribute: false })
  get selection(): number[] {
    return this._selection;
  }

  set selection(val: number[]) {
    if (this.disabled) {
      return;
    }
    if (!val || !val.sort) {
      val = [];
    }
    val.sort();
    const oldVal = this._selection;
    this._selection = val;
    this.requestUpdate('selection', oldVal);
  }

  updated(changedProperties: PropertyValues) {
    if (changedProperties.has('disabled')) {
      if (this.disabled) {
        // When disabled, we do nothing to _selection. It is frozen.
        // We don't call _rationalize either.
      } else {
        // When re-enabled, we rely on the DOM attributes to restore selection.
        this._bubbleUp();
      }
    }

    if (changedProperties.has('selection')) {
      // Only sync changes to DOM if enabled.
      // If disabled, we might have cleared _selection, but we don't want to clear DOM.
      if (!this.disabled) {
        this._rationalize();
      }
    }
  }

  private _click(e: MouseEvent): void {
    if (this.disabled) {
      return;
    }
    // Look up the DOM path until we find an element that is a child of
    // 'this', and set _selection based on that.
    let target: Element | null = e.target as Element;
    while (target && target.parentElement !== this) {
      target = target.parentElement;
    }
    if (!target || target.parentElement !== this) {
      return; // not a click we care about
    }
    if (target.hasAttribute('selected')) {
      target.removeAttribute('selected');
    } else {
      target.setAttribute('selected', '');
    }
    this._bubbleUp();
    this.dispatchEvent(
      new CustomEvent<MultiSelectSkSelectionChangedEventDetail>('selection-changed', {
        detail: {
          selection: this._selection,
        },
        bubbles: true,
      })
    );
  }

  // Loop over all immediate child elements update the selected attributes
  // based on the selected property of this element.
  private _rationalize(): void {
    this.observerDisconnect();
    // assume this.selection is sorted when this is called.
    let s = 0;
    for (let i = 0; i < this.children.length; i++) {
      if (this.children[i].getAttribute('tabindex') === null) {
        this.children[i].setAttribute('tabindex', '0');
      }
      if (this._selection[s] === i) {
        this.children[i].setAttribute('selected', '');
        s++;
      } else {
        this.children[i].removeAttribute('selected');
      }
    }
    this.observerConnect();
  }

  // Loop over all immediate child elements and find all with the selected
  // attribute.
  private _bubbleUp(): void {
    if (this.disabled) {
      return;
    }
    this._selection = [];
    for (let i = 0; i < this.children.length; i++) {
      if (this.children[i].hasAttribute('selected')) {
        this._selection.push(i);
      }
    }
    // Since _bubbleUp changes _selection, we should notify Lit if we want strict reactivity,
    // though _rationalize already synced the DOM.
    this.requestUpdate('selection');
  }
}

export interface MultiSelectSkSelectionChangedEventDetail {
  readonly selection: number[];
}
