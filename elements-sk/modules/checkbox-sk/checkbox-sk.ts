/* eslint-disable no-self-assign */
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
 * @module elements-sk/checkbox-sk
 * @description <h2><code>checkbox-sk</code></h2>
 *
 * <p>
 *   The checkbox-sk and element contains a native 'input'
 *   element in light DOM so that it can participate in a form element.
 * </p>
 *
 * <p>
 *    Each element also supports the following attributes exactly as the
 *    native checkbox element:
 *    <ul>
 *      <li>checked</li>
 *      <li>disabled</li>
 *      <li>name</li>
 *     </ul>
 * </p>
 *
 * <p>
 *    All the normal events of a native checkbox are supported.
 * </p>
 *
 * @attr label - A string, with no markup, that is to be used as the label for
 *            the checkbox. If you wish to have a label with markup then set
 *            'label' to the empty string and create your own
 *            <code>label</code> element in the DOM with the 'for' attribute
 *            set to match the name of the checkbox-sk.
 *
 * @prop checked This mirrors the checked attribute.
 * @prop disabled This mirrors the disabled attribute.
 * @prop name This mirrors the name attribute.
 * @prop label This mirrors the label attribute.
 *
 */
import { html, LitElement, TemplateResult } from 'lit';
import { property } from 'lit/decorators.js';
import { define } from '../define';

export class CheckOrRadio extends LitElement {
  private static nextUniqueId = 0;

  protected readonly uniqueId = `${CheckOrRadio.nextUniqueId++}`;

  @property({ type: Boolean, reflect: true })
  checked: boolean = false;

  @property({ type: Boolean, reflect: true })
  disabled: boolean = false;

  @property({ type: String, reflect: true })
  name: string = '';

  @property({ type: String, reflect: true })
  label: string = '';

  protected get _role(): string {
    return 'checkbox';
  }

  // Allow subclasses to override icons
  protected get checkedIcon(): string {
    return 'check_box';
  }

  protected get uncheckedIcon(): string {
    return 'check_box_outline_blank';
  }

  connectedCallback() {
    super.connectedCallback();
    this.addEventListener('click', this.handleHostClick);
  }

  disconnectedCallback() {
    super.disconnectedCallback();
    this.removeEventListener('click', this.handleHostClick);
  }

  private handleHostClick = (e: MouseEvent) => {
    if (e.target === this) {
      if (this.checked && this._role === 'radio') {
        return;
      }
      const input = this.querySelector('input');
      input?.click();
    }
  };

  createRenderRoot() {
    return this;
  }

  render(): TemplateResult {
    return html`
      <label for="${this._role}-${this.uniqueId}">
        <input
          type=${this._role}
          name=${this.name}
          id="${this._role}-${this.uniqueId}"
          ?disabled=${this.disabled}
          .checked=${this.checked}
          @change=${this.handleChange} />
        <span class="icons">
          <span class="icon-sk unchecked">${this.uncheckedIcon}</span>
          <span class="icon-sk checked">${this.checkedIcon}</span>
        </span>
        <span class="label">${this.label}</span>
      </label>
    `;
  }

  private handleChange(e: Event) {
    // The native input event 'change' bubbles, but we want to ensure
    // the property is updated.
    this.checked = (e.target as HTMLInputElement).checked;
  }
}

define('checkbox-sk', CheckOrRadio);
