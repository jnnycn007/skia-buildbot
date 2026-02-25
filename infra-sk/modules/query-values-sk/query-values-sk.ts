/**
 * @module modules/query-values-sk
 * @description <h2><code>query-values-sk</code></h2>
 *
 * The right-hand side of the query-sk element, the values for a single key
 * in a query/paramset.
 *
 * @evt query-values-changed - Triggered only when the selections have actually
 *     changed. The selection is available in e.detail.
 *
 * @attr {boolean} hide_invert - If the option to invert a query should be made available to
 *       the user.
 * @attr {boolean} hide_regex - If the option to include regex in the query should be made
 *       available to the user.
 */
import { html, LitElement } from 'lit';
import { customElement, property, state } from 'lit/decorators.js';
import { MultiSelectSkSelectionChangedEventDetail } from '../../../elements-sk/modules/multi-select-sk/multi-select-sk';
import { CheckOrRadio } from '../../../elements-sk/modules/checkbox-sk/checkbox-sk';
import '../../../elements-sk/modules/checkbox-sk';
import '../../../elements-sk/modules/multi-select-sk';

export interface QueryValuesSkQueryValuesChangedEventDetail {
  invert: boolean;
  regex: boolean;
  values: string[];
}

@customElement('query-values-sk')
export class QueryValuesSk extends LitElement {
  private static nextUniqueId = 0;

  readonly uniqueId = `${QueryValuesSk.nextUniqueId++}`;

  private _options: string[] = [];

  private _selected: string[] = [];

  // Custom setter for options to handle legacy filtering logic
  @property({ attribute: false, noAccessor: true })
  get options(): string[] {
    return this._options;
  }

  set options(val: string[]) {
    const oldVal = this._options;
    this._options = val;
    this._selected = []; // Legacy options setter resets selected
    this._fastFilter();
    this.requestUpdate('options', oldVal);
  }

  // Custom setter for selected to handle legacy cleaning and state synchronization
  @property({ attribute: false, noAccessor: true })
  get selected(): string[] {
    return this._selected;
  }

  set selected(val: string[]) {
    const oldVal = this._selected;
    this._selected = val;
    this._syncStateFromSelected();
    this.requestUpdate('selected', oldVal);
  }

  @property({ type: Boolean, reflect: true })
  hide_invert: boolean = false;

  @property({ type: Boolean, reflect: true })
  hide_regex: boolean = false;

  @state()
  private _filtering: boolean = false;

  @state()
  private _filteredOptions: string[] = [];

  @state()
  private _invertChecked: boolean = false;

  @state()
  private _regexChecked: boolean = false;

  @state()
  private _regexValue: string = '';

  // We keep track of filter string to re-filter when options change
  private _filterString: string = '';

  createRenderRoot() {
    return this;
  }

  private _syncStateFromSelected() {
    const val = this._selected;
    this._invertChecked = !!(val.length >= 1 && val[0].startsWith('!'));
    this._regexChecked = !!(val.length === 1 && val[0].startsWith('~'));

    // Clean _selected in place
    this._selected = val.map((v) => {
      if (v.startsWith('!') || v.startsWith('~')) {
        return v.slice(1);
      }
      return v;
    });

    if (this._selected.length && this._regexChecked) {
      this._regexValue = this._selected[0];
    }
  }

  /**
   * Filter the options displayed based on text entered in the
   * filter text box
   */
  private _fastFilter() {
    const filterString = this._filterString ? this._filterString.trim() : '';

    if (!filterString) {
      this._filtering = false;
      this._filteredOptions = [];
      return;
    }

    this._filtering = true;
    const filters = filterString.toLowerCase().split(/\s+/);
    const matches = (s: string): boolean => {
      s = s.toLowerCase();
      return filters.some((f) => s.includes(f));
    };
    this._filteredOptions = this._options.filter(matches);
  }

  protected render() {
    // Determine which options to display (filtered or all)
    const displayOptions = this._filtering ? this._filteredOptions : this._options;

    return html`
      <checkbox-sk
        id="invert-${this.uniqueId}"
        .checked=${this._invertChecked}
        @change=${this._invertChange}
        title="Match items not selected below."
        label="Invert"
        ?hidden=${this.hide_invert}></checkbox-sk>
      <checkbox-sk
        id="regex-${this.uniqueId}"
        class="regex"
        .checked=${this._regexChecked}
        @change=${this._regexChange}
        title="Match items via regular expression."
        label="Regex"
        ?hidden=${this.hide_regex}></checkbox-sk>
      <input
        type="text"
        id="regexValue-${this.uniqueId}"
        class="regexValue"
        .value=${this._regexValue}
        @input=${this._regexInputChange} />
      <div class="filtering">
        <input
          id="filter-${this.uniqueId}"
          .value=${this._filterString}
          @input=${this._onFilterInput}
          placeholder="Filter Values"
          name="query-value-sk-filter-val"
          autocomplete="off" />
        ${this._filtering
          ? html`<button @click=${this.clearFilter} class="clear_filters" title="Clear filter">
              &cross;
            </button>`
          : ''}
      </div>
      <multi-select-sk
        id="values-${this.uniqueId}"
        class="values"
        @selection-changed=${this._selectionChange}>
        ${displayOptions.map(
          (v) => html` <div value=${v} ?selected=${this._selected.includes(v)}>${v}</div> `
        )}
      </multi-select-sk>
    `;
  }

  private _invertChange(e: Event) {
    const checkbox = e.target as CheckOrRadio;
    this._invertChecked = checkbox.checked;
    if (this._invertChecked && this._regexChecked) {
      this._regexChecked = false;
    }
    this._fireEvent();
  }

  private _regexChange(e: Event) {
    const checkbox = e.target as CheckOrRadio;
    this._regexChecked = checkbox.checked;
    if (this._regexChecked && this._invertChecked) {
      this._invertChecked = false;
    }
    this._fireEvent();
  }

  private _onFilterInput(e: Event) {
    this._filterString = (e.target as HTMLInputElement).value;
    this._fastFilter();
  }

  private _regexInputChange(e: Event) {
    this._regexValue = (e.target as HTMLInputElement).value;
    this._fireEvent();
  }

  private _selectionChange(e: CustomEvent<MultiSelectSkSelectionChangedEventDetail>) {
    const currentOptions = this._filtering ? this._filteredOptions : this._options;

    // Convert indices to values
    const newSelected = e.detail.selection.map((i) => currentOptions[i]);

    this._selected = newSelected;
    this._fireEvent();
    this.requestUpdate();
  }

  private _fireEvent() {
    const prefix = this._invertChecked ? '!' : '';
    let selectedValues: string[] = [];

    if (this._regexChecked) {
      selectedValues = [`~${this._regexValue}`];
    } else {
      // Only include selected values that are currently visible options
      // to match legacy behavior where multi-select-sk indices are used.
      const currentOptions = this._filtering ? this._filteredOptions : this._options;
      selectedValues = this._selected
        .filter((v) => currentOptions.includes(v))
        .map((v) => prefix + v);
    }

    this.dispatchEvent(
      new CustomEvent<QueryValuesSkQueryValuesChangedEventDetail>('query-values-changed', {
        detail: {
          invert: this._invertChecked,
          regex: this._regexChecked,
          values: selectedValues,
        },
        bubbles: true,
      })
    );
  }

  public clearFilter(): void {
    this._filterString = '';
    this._fastFilter();
  }
}
