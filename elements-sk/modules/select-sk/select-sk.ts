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

/** @module elements-sk/select-sk
 *
 * @description <h2><code>select-sk</code></h2>
 *
 * <p>
 *   Clicking on the children will cause them to be selected.
 * </p>
 *
 * <p>
 *   The select-sk elements monitors for the addition and removal of child
 *   elements and will update the 'selected' property as needed. Note that it
 *   does not monitor the 'selected' attribute of child elements, and will not
 *   update the 'selected' property if they are changed directly.
 * </p>
 *
 * @example
 *
 *   <select-sk>
 *     <div></div>
 *     <div></div>
 *     <div selected></div>
 *     <div></div>
 *   </select-sk>
 *
 * @attr disabled - Indicates whether the element is disabled.
 *
 * @evt selection-changed - Sent when an item is clicked and the selection is changed.
 *   The detail of the event contains the child element index:
 *
 *   <pre>
 *     detail: {
 *       selection: 1,
 *     }
 *   </pre>
 *
 */

import { LitElement } from 'lit';
import { customElement, property } from 'lit/decorators.js';

export interface SelectSkSelectionChangedEventDetail {
  readonly selection: number;
}

@customElement('select-sk')
export class SelectSk extends LitElement {
  private _selection: number = -1;

  @property({ type: Number, noAccessor: true })
  get selection(): number | string | null | undefined {
    return this._selection;
  }

  set selection(val: number | string | null | undefined) {
    if (this.disabled) {
      return;
    }

    // Sanitize
    if (val === undefined || val === null) {
      val = -1;
    }
    let numVal = +val;
    if (isNaN(numVal)) {
      numVal = -1;
    }

    // Validate bounds. Reset to -1 if invalid or out of bounds.
    if (numVal >= this.children.length || numVal < 0) {
      numVal = -1;
    }

    const oldSelection = this._selection;
    this._selection = numVal;

    // Synchronously update attributes to match legacy behavior
    this._rationalize();

    this.requestUpdate('selection', oldSelection);
  }

  private _disabled: boolean = false;

  @property({ type: Boolean, noAccessor: true })
  get disabled(): boolean {
    return this._disabled;
  }

  set disabled(val: boolean) {
    const oldVal = this._disabled;
    this._disabled = val;

    // Synchronous side effects
    if (this._disabled) {
      this.setAttribute('disabled', '');
      this.setAttribute('aria-disabled', 'true');
      // Do NOT clear selection when disabled to match 'stays fixed' test expectations.
      this.removeAttribute('tabindex');
      this.blur();
    } else {
      this.removeAttribute('disabled');
      this.setAttribute('aria-disabled', 'false');
      this.setAttribute('tabindex', '0');
      this._bubbleUp(); // Find selected child and update selection
    }

    this.requestUpdate('disabled', oldVal);
  }

  private _obs: MutationObserver;

  constructor() {
    super();
    // Keep _selection up to date by monitoring DOM changes.
    this._obs = new MutationObserver(() => this._bubbleUp());
  }

  createRenderRoot() {
    return this;
  }

  connectedCallback(): void {
    super.connectedCallback();
    this.addEventListener('click', this._click);
    this.addEventListener('keydown', this._onKeyDown);
    this.observerConnect();
    // Use a small timeout to allow children to be upgraded/rendered if they are custom elements
    // and to ensure attributes are settled (like disabled). 0ms is sufficient to let the
    // parser finish populating the DOM.
    setTimeout(() => this._bubbleUp(), 0);
  }

  disconnectedCallback(): void {
    super.disconnectedCallback();
    this.removeEventListener('click', this._click);
    this.removeEventListener('keydown', this._onKeyDown);
    this.observerDisconnect();
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

  private _click(e: MouseEvent): void {
    if (this.disabled) {
      return;
    }
    const oldIndex = this.selection;
    // Look up the DOM path until we find an element that is a child of
    // 'this', and set _selection based on that.
    let target: Element | null = e.target as Element;
    while (target && target.parentElement !== this) {
      target = target.parentElement;
    }
    if (target?.parentElement === this) {
      for (let i = 0; i < this.children.length; i++) {
        if (this.children[i] === target) {
          this.selection = i;
          break;
        }
      }
    }
    // _rationalize is called by setter if selection changes
    // But if selection didn't change (clicked same item), ensure DOM is correct
    if (oldIndex === this.selection) {
      this._rationalize();
    }

    if (oldIndex !== this.selection) {
      this._emitEvent();
    }
  }

  private _emitEvent(): void {
    this.dispatchEvent(
      new CustomEvent<SelectSkSelectionChangedEventDetail>('selection-changed', {
        detail: {
          selection: this.selection as number,
        },
        bubbles: true,
      })
    );
  }

  // Loop over all immediate child elements and make sure at most only one is selected.
  private _rationalize(): void {
    if (this.disabled) {
      return;
    }
    this.observerDisconnect();
    if (!this.hasAttribute('role')) {
      this.setAttribute('role', 'listbox');
    }
    if (!this.hasAttribute('tabindex') && !this.disabled) {
      this.setAttribute('tabindex', '0');
    }
    for (let i = 0; i < this.children.length; i++) {
      const child = this.children[i];
      if (!child.hasAttribute('role')) {
        child.setAttribute('role', 'option');
      }
      if (this.selection === i) {
        child.setAttribute('selected', '');
        child.setAttribute('aria-selected', 'true');
      } else {
        child.removeAttribute('selected');
        child.setAttribute('aria-selected', 'false');
      }
    }
    this.observerConnect();
  }

  // Loop over all immediate child elements and find the first one selected.
  private _bubbleUp(): void {
    const oldSelection = this.selection;
    let newSelection = -1;
    if (!this.disabled) {
      for (let i = 0; i < this.children.length; i++) {
        if (this.children[i].hasAttribute('selected')) {
          newSelection = i;
          break;
        }
      }
    }

    if (newSelection !== oldSelection) {
      this.selection = newSelection;
    } else {
      this._rationalize();
    }
  }

  private _onKeyDown(e: KeyboardEvent): void {
    if (e.altKey || this.disabled) return;
    const oldIndex = (this.selection as number) ?? -1;
    let newIndex = oldIndex;

    switch (e.key) {
      case 'ArrowDown':
        if (newIndex < this.children.length - 1) {
          newIndex += 1;
        }
        e.preventDefault();
        break;
      case 'ArrowUp':
        if (newIndex > 0) {
          newIndex -= 1;
        }
        e.preventDefault();
        break;
      case 'Home':
        newIndex = 0;
        e.preventDefault();
        break;
      case 'End':
        newIndex = this.children.length - 1;
        e.preventDefault();
        break;
      default:
        break;
    }

    if (newIndex !== oldIndex) {
      this.selection = newIndex;
      this._emitEvent();
    }
  }
}
