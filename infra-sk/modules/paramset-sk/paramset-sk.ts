/**
 * @module module/paramset-sk
 * @description <h2><code>paramset-sk</code></h2>
 *
 * The paramset-sk element displays a paramset and generates events as the
 * params and labels are clicked.
 *
 * @evt paramset-key-click - Generated when the key for a paramset is clicked.
 *     The name of the key will be sent in e.detail.key. The value of
 *     e.detail.ctrl is true if the control key was pressed when clicking.
 *
 *      {
 *        key: "arch",
 *        ctrl: false,
 *      }
 *
 * @evt paramset-key-value-click - Generated when one value for a paramset is
 *     clicked. The name of the key will be sent in e.detail.key, the value in
 *     e.detail.value. The value of e.detail.ctrl is true if the control key was
 *     pressed when clicking.
 *
 *      {
 *        key: "arch",
 *        value: "x86",
 *        ctrl: false,
 *      }
 *
 * @evt paramset-checkbox-click - Generated when the checkbox for a paramset value
 *     is clicked.
 *
 *      {
 *        key: "arch",
 *        value: "x86",
 *        selected: true
 *      }
 *
 *  * @evt paramset-key-checkbox-click - Generated when the checkbox for a key in the paramset
 *     is clicked.
 *
 *      {
 *        key: "arch",
 *        values: ["x86", "x64"],
 *        selected: true
 *      }
 *
 * @evt plus-click - Generated when the plus sign is clicked. The element must
 *     have the 'clickable_plus' attribute set. The details of the event
 *     contains both the key and the values for the row, for example:
 *
 *      {
 *        key: "arch",
 *        values" ["x86", "risc-v"],
 *      }
 *
 * @evt paramset-value-remove-click - Generated when one value for a paramset is
 *     removed. The name of the key will be sent in e.detail.key, the value in
 *     e.detail.value.
 *
 *      {
 *        key: "arch",
 *        value: "x86",
 *      }
 *
 * @attr {string} clickable - If true then keys and values look like they are
 *     clickable i.e. via color, text-decoration, and cursor. If clickable is
 *     false then this element won't generate the events listed below, and the
 *     keys and values are not styled to look clickable. Setting both clickable
 *     and clickable_values is unsupported.
 *
 * @attr {string} clickable_values - If true then only the values are clickable.
 *     Setting both clickable and clickable_values is unsupported.
 *
 * @attr {string} clickable_plus - If true then a plus sign is added to every
 * row in the right hand column, that when pressed emits the plus-click event
 * that contains the key and values for that row.
 *
 * @attr {string} removable_values - If true then the cancel icon is displayed
 * next to each value in the paramset to remove the values from the set
 *
 * @attr {string} checkbox_values - If true, then the values displayed will have
 * a checkbox to let the user select/unselect the specific value.
 *
 */
import { html, TemplateResult, LitElement, PropertyValues } from 'lit';
import { customElement, property, state } from 'lit/decorators.js';
import { ParamSet } from '../query';
import '../../../elements-sk/modules/icons/add-icon-sk';
import { ToastSk } from '../../../elements-sk/modules/toast-sk/toast-sk';
import '../../../elements-sk/modules/icons/cancel-icon-sk';
import '../../../elements-sk/modules/checkbox-sk';
import '../../../elements-sk/modules/toast-sk';
import { CheckOrRadio } from '../../../elements-sk/modules/checkbox-sk/checkbox-sk';

export interface ParamSetSkClickEventDetail {
  readonly key: string;
  readonly value?: string;
  readonly ctrl: boolean;
}

export interface ParamSetSkPlusClickEventDetail {
  readonly key: string;
  readonly values: string[];
}

export interface ParamSetSkRemoveClickEventDetail {
  readonly key: string;
  readonly value: string;
}

export interface ParamSetSkCheckboxClickEventDetail {
  readonly key: string;
  readonly value: string;
  readonly selected: boolean;
}

export interface ParamSetSkKeyCheckboxClickEventDetail {
  readonly key: string;
  readonly values: string[];
  readonly selected: boolean;
}

@customElement('paramset-sk')
export class ParamSetSk extends LitElement {
  @property({ type: Array })
  titles: string[] = [];

  @property({ type: Array })
  paramsets: ParamSet[] = [];

  @property({ type: Object })
  highlight: { [key: string]: string } = {};

  @property({ type: Boolean, reflect: true })
  clickable: boolean = false;

  @property({ type: Boolean, reflect: true, attribute: 'clickable_values' })
  clickable_values: boolean = false;

  @property({ type: Boolean, reflect: true, attribute: 'clickable_plus' })
  clickable_plus: boolean = false;

  @property({ type: Boolean, reflect: true, attribute: 'checkbox_values' })
  checkbox_values: boolean = false;

  @property({ type: Boolean, reflect: true, attribute: 'removable_values' })
  removable_values: boolean = false;

  @property({ type: Boolean, reflect: true, attribute: 'copy-content' })
  copy_content: boolean = false;

  @state()
  private _sortedKeys: string[] = [];

  @state()
  private _unchecked: Map<string, Set<string>> = new Map();

  private _toast: ToastSk | null = null;

  createRenderRoot() {
    return this;
  }

  protected firstUpdated() {
    this._toast = this.querySelector('toast-sk');
  }

  protected willUpdate(changedProperties: PropertyValues) {
    if (changedProperties.has('paramsets')) {
      // Compute a rolled up set of all parameter keys across all paramsets.
      const allKeys = new Set<string>();
      this.paramsets.forEach((p) => {
        Object.keys(p).forEach((key) => {
          allKeys.add(key);
        });
      });
      this._sortedKeys = Array.from(allKeys).sort();
      this._unchecked = new Map();
    }
  }

  render() {
    return html`
      <table @click=${this._click} class=${this.computeClass()}>
        <tbody>
          <tr>
            <th></th>
            ${this.titlesTemplate()}
          </tr>
          ${this.rowsTemplate()}
        </tbody>
      </table>
      <toast-sk duration="2000">Copied</toast-sk>
    `;
  }

  private titlesTemplate() {
    return this.normalizedTitles().map((t) => html`<th>${t}</th>`);
  }

  private rowsTemplate() {
    return this._sortedKeys.map((key) => this.rowTemplate(key));
  }

  private rowTemplate(key: string) {
    if (this.checkbox_values) {
      // If this row contains only one value, let's disable the key checkbox.
      let disabled = false;
      this.paramsets.forEach((ps) => {
        const vals = ps[key];
        if (vals && vals.length <= 1) {
          disabled = true;
        }
      });
      return html`<tr>
        <th data-key=${key}>
          <checkbox-sk
            id="selectAll-${key}"
            name=""
            @change=${(e: MouseEvent) => this.paramsetKeySelectAllHandler(e, key)}
            label=""
            checked
            ?disabled=${disabled}
            title="Select/Unselect this value from the graph.">
          </checkbox-sk>
          ${key}
        </th>
        ${this.paramsetValuesTemplate(key)}
      </tr>`;
    } else {
      return html`<tr>
        <th data-key=${key}>${key}</th>
        ${this.paramsetValuesTemplate(key)}
      </tr>`;
    }
  }

  private paramsetValuesTemplate(key: string) {
    const ret: TemplateResult[] = [];
    this.paramsets.forEach((p) =>
      ret.push(
        html`<td>${this.paramsetValueTemplate(key, p[key] || [])}</td>`,
        this.optionalPlusSign(key, p),
        this.optionalCopyContent(key, p)
      )
    );
    return ret;
  }

  private optionalPlusSign(key: string, p: ParamSet): TemplateResult {
    if (!this.clickable_plus) {
      return html``;
    }
    return html` <td>
      <add-icon-sk data-key=${key} data-values=${JSON.stringify(p[key])}></add-icon-sk>
    </td>`;
  }

  private optionalCopyContent(key: string, p: ParamSet): TemplateResult {
    if (!this.copy_content) {
      return html``;
    }
    return html` <td>
      <div class="icon-sk copy-content" @click=${() => this.copyContent(`${key}=${p[key]}`)}>
        content_copy
      </div>
    </td>`;
  }

  private paramsetValueTemplate(key: string, params: string[]) {
    // Figure out if we are down to just one checkbox being checked. If so we'll
    // want to disable that checkbox so that it can't be unchecked, otherwise
    // all the data will disappear from the display.
    let downToJustOneCheckedCheckboxForThisKey = false;

    // Count the number of unchecked values for this key.
    let numUnchecked = 0;
    const uncheckedSet = this._unchecked.get(key);
    if (uncheckedSet !== undefined) {
      numUnchecked = uncheckedSet.size;
    }

    if (params.length - numUnchecked <= 1) {
      downToJustOneCheckedCheckboxForThisKey = true;
    }
    return params.map((value) => {
      if (this.checkbox_values) {
        let disabled = false;
        const currentCheckboxChecked = uncheckedSet === undefined || !uncheckedSet.has(value);
        if (downToJustOneCheckedCheckboxForThisKey && currentCheckboxChecked) {
          disabled = true;
        }

        return html`
          <div class=${this.highlighted(key, value)} data-key=${key} data-value=${value}>
            <checkbox-sk
              id="checkbox-${key}-${value}"
              name=""
              @change=${(e: MouseEvent) => this.checkboxValueClickHandler(e, key, value)}
              label=""
              checked
              ?disabled=${disabled}
              title="Select/Unselect this value from the graph.">
            </checkbox-sk>
            ${value}
          </div>
        `;
      }
      return html`<div class=${this.highlighted(key, value)} data-key=${key} data-value=${value}>
        <span>${value}</span> ${this.cancelIconTemplate(key, value)}
      </div> `;
    });
  }

  private cancelIconTemplate(key: string, value: string): TemplateResult {
    if (this.removable_values) {
      return html`<cancel-icon-sk
        id="${key}-${value}-remove"
        data-key=${key}
        data-value=${value}
        title="Negative"></cancel-icon-sk>`;
    }
    return html``;
  }

  private computeClass() {
    if (this.clickable_values) {
      return 'clickable_values';
    }
    if (this.clickable) {
      return 'clickable';
    }
    return '';
  }

  private highlighted(key: string, value: string) {
    return this.highlight[key] === value ? 'highlight' : '';
  }

  private async copyContent(body: string) {
    await navigator.clipboard.writeText(body);
    this._toast?.show();
  }

  private fixUpDisabledStateOnRemainingCheckboxes(isChecked: boolean, key: string, value: string) {
    // Update the unchecked status and then re-render.
    const set = this._unchecked.get(key) || new Set();
    if (isChecked) {
      set.delete(value);
    } else {
      set.add(value);
    }
    this._unchecked.set(key, set);
  }

  private paramsetKeySelectAllHandler(e: MouseEvent, key: string) {
    const selectAll = (e.target! as HTMLInputElement).checked;
    this.paramsets.forEach((p) => {
      const vals = p[key];
      if (!vals) return;
      let keepCheckedIndex = -1;
      const valuesToUpdate: string[] = [];
      for (let i = 0; i < vals.length; i++) {
        const checkbox_id = `checkbox-${key}-${vals[i]}`;
        const checkbox = this.querySelector(`#${CSS.escape(checkbox_id)}`) as CheckOrRadio;
        if (!checkbox) continue;

        // If we are unselecting, we want to keep the first checked item
        // from selected values around. This is because we would need at least
        // one trace in the graph.
        if (keepCheckedIndex === -1 && checkbox.checked && !selectAll) {
          keepCheckedIndex = i;
        } else {
          valuesToUpdate.push(vals[i]);
          this.fixUpDisabledStateOnRemainingCheckboxes(selectAll, key, vals[i]);
        }
      }
      this.requestUpdate();

      const detail: ParamSetSkKeyCheckboxClickEventDetail = {
        selected: selectAll,
        key: key,
        values: valuesToUpdate,
      };
      // Send the event so that the graph can be updated to match the selection.
      this.dispatchEvent(
        new CustomEvent<ParamSetSkKeyCheckboxClickEventDetail>('paramset-key-checkbox-click', {
          detail,
          bubbles: true,
        })
      );
    });
  }

  private checkboxValueClickHandler(e: MouseEvent, key: string, value: string) {
    const isChecked = (e.target! as HTMLInputElement).checked;
    const detail: ParamSetSkCheckboxClickEventDetail = {
      selected: isChecked,
      key: key,
      value: value,
    };
    this.dispatchEvent(
      new CustomEvent<ParamSetSkCheckboxClickEventDetail>('paramset-checkbox-click', {
        detail,
        bubbles: true,
      })
    );

    this.fixUpDisabledStateOnRemainingCheckboxes(isChecked, key, value);
    this.requestUpdate();
  }

  private _click(e: MouseEvent) {
    if (
      !this.clickable &&
      !this.clickable_values &&
      !this.clickable_plus &&
      !this.removable_values
    ) {
      return;
    }

    const t = e.target as HTMLElement;
    // Check if we clicked "content_copy" which is a div
    if (t.classList.contains('copy-content')) {
      return;
    }

    const target = t.closest('[data-key]') as HTMLElement | null;
    if (!target) {
      return;
    }

    if (target.nodeName === 'TH') {
      if (!this.clickable) {
        return;
      }
      const detail: ParamSetSkClickEventDetail = {
        key: target.dataset.key!,
        ctrl: e.ctrlKey,
      };
      this.dispatchEvent(
        new CustomEvent<ParamSetSkClickEventDetail>('paramset-key-click', {
          detail,
          bubbles: true,
        })
      );
    } else if (target.nodeName === 'DIV') {
      // It could be a value click
      const detail: ParamSetSkClickEventDetail = {
        key: target.dataset.key!,
        value: target.dataset.value,
        ctrl: e.ctrlKey,
      };
      this.dispatchEvent(
        new CustomEvent<ParamSetSkClickEventDetail>('paramset-key-value-click', {
          detail,
          bubbles: true,
        })
      );
    } else {
      // Check if target is ADD-ICON-SK or CANCEL-ICON-SK
      // Since closest('[data-key]') might find them.
      if (target.nodeName === 'ADD-ICON-SK') {
        const detail: ParamSetSkPlusClickEventDetail = {
          key: target.dataset.key!,
          values: JSON.parse(target.dataset.values!) as string[],
        };
        this.dispatchEvent(
          new CustomEvent<ParamSetSkPlusClickEventDetail>('plus-click', {
            detail,
            bubbles: true,
          })
        );
      } else if (target.nodeName === 'CANCEL-ICON-SK') {
        this.removeParam(target.dataset.key!, target.dataset.value!);
      }
    }
  }

  // Returns the titles specified by the user, or an empty title for each paramset
  // if the number of specified titles and the number of paramsets don't match.
  private normalizedTitles() {
    if (this.titles.length === this.paramsets.length) {
      return this.titles;
    }
    return new Array<string>(this.paramsets.length).fill('');
  }

  removeParam(key: string, value: string) {
    const paramsets: ParamSet[] = [];
    this.paramsets.forEach((paramset) => {
      const values = paramset[key];
      if (!values) {
        paramsets.push(paramset);
        return;
      }
      const valIndex = values.indexOf(value);
      if (valIndex > -1) {
        values.splice(valIndex, 1);
        if (values.length === 0) {
          delete paramset[key];
        } else {
          paramset[key] = values;
        }
      }
      paramsets.push(paramset);
    });

    // Set the current paramsets to the updated value
    this.paramsets = paramsets;
    const detail: ParamSetSkRemoveClickEventDetail = {
      key: key,
      value: value,
    };
    this.dispatchEvent(
      new CustomEvent<ParamSetSkRemoveClickEventDetail>('paramset-value-remove-click', {
        detail,
        bubbles: true,
      })
    );
  }
}
