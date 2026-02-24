/**
 * @module modules/domain-picker-sk
 * @description <h2><code>domain-picker-sk</code></h2>
 *
 * Allows picking either a date range for commits, or for
 * picking a number of commits to show before a selected
 * date.
 *
 * @attr {string} force_request_type - A value of 'dense' or 'range' will
 *   force the corresponding request_type to be always set.
 *
 */
import { html, LitElement } from 'lit';
import { customElement, property, state } from 'lit/decorators.js';
import { RequestType } from '../json';

import '../../../elements-sk/modules/radio-sk';
import '../calendar-input-sk';

// Types of domain ranges we can choose.
// TODO(jcgregorio) Make the underlying dataframe.RequestType a string.
const RANGE = 0; // Specify a begin and end time.
const DENSE = 1; // Specify an end time and the number of commits with data.

type ForceRequestType = 'range' | 'dense' | '';

const toDate = (seconds: number) => new Date(seconds * 1000);

const toForceRequestType = (s: string | null): ForceRequestType => {
  if (s === 'range') {
    return 'range';
  }
  if (s === 'dense') {
    return 'dense';
  }
  return '';
};

/** The state of the DomainPickerSk control. */
export interface DomainPickerState {
  /**  Beginning of time range in Unix timestamp seconds. */
  begin: number;
  /**  End of time range in Unix timestamp seconds. */
  end: number;
  /**
   * If RequestType is REQUEST_COMPACT (1), then this value is the number of
   * commits to show before End, and the value of Begin is ignored.
   */
  num_commits: number;
  request_type: RequestType;
}

@customElement('domain-picker-sk')
export class DomainPickerSk extends LitElement {
  @state()
  private begin: number;

  @state()
  private end: number;

  @state()
  private num_commits: number;

  @state()
  private selected_request_type: RequestType;

  @property({
    attribute: 'force_request_type',
    reflect: true,
    converter: {
      fromAttribute: toForceRequestType,
    },
  })
  force_request_type: ForceRequestType = '';

  constructor() {
    super();
    const now = Date.now();
    this.begin = Math.floor(now / 1000 - 24 * 60 * 60);
    this.end = Math.floor(now / 1000);
    this.num_commits = 50;
    this.selected_request_type = RANGE;
  }

  createRenderRoot() {
    return this;
  }

  get request_type(): RequestType {
    if (this.force_request_type === 'dense') {
      return DENSE;
    }
    if (this.force_request_type === 'range') {
      return RANGE;
    }
    return this.selected_request_type;
  }

  render() {
    return html`
      ${this.showRadio()}
      <div class="ranges">
        ${this.renderRequestType()}
        <label>
          <span class="prefix">End:</span>
          <calendar-input-sk
            @input=${this.endChange}
            .displayDate=${toDate(this.end)}></calendar-input-sk>
        </label>
      </div>
    `;
  }

  private showRadio() {
    if (!this.force_request_type) {
      return html`
        <radio-sk
          @change=${this.typeRange}
          ?checked=${this.request_type === RANGE}
          label="Date Range"
          name="daterange"></radio-sk>
        <radio-sk
          @change=${this.typeDense}
          ?checked=${this.request_type === DENSE}
          label="Dense"
          name="daterange"></radio-sk>
      `;
    }
    return html``;
  }

  private renderRequestType() {
    if (this.request_type === RANGE) {
      return html`
        <p>Display all points in the date range.</p>
        <label>
          <span class="prefix">Begin:</span>
          <calendar-input-sk
            @input=${this.beginChange}
            .displayDate=${toDate(this.begin)}></calendar-input-sk>
        </label>
      `;
    }
    return html`
      <p>Display only the points that have data before the date.</p>
      <label>
        <span class="prefix">Points</span>
        <input
          @change=${this.numChanged}
          type="number"
          .value="${this.num_commits.toString()}"
          min="1"
          max="5000"
          list="defaultNumbers"
          title="The number of points." />
      </label>
      <datalist id="defaultNumbers">
        <option value="50"></option>
        <option value="100"></option>
        <option value="250"></option>
        <option value="500"></option>
      </datalist>
    `;
  }

  private typeRange() {
    this.selected_request_type = RANGE;
  }

  private typeDense() {
    this.selected_request_type = DENSE;
  }

  private beginChange(e: CustomEvent<Date>) {
    this.begin = Math.floor(e.detail.valueOf() / 1000);
  }

  private endChange(e: CustomEvent<Date>) {
    this.end = Math.floor(e.detail.valueOf() / 1000);
  }

  private numChanged(e: MouseEvent) {
    this.num_commits = (e.target! as HTMLInputElement).valueAsNumber;
  }

  get state(): DomainPickerState {
    return {
      begin: this.begin,
      end: this.end,
      num_commits: this.num_commits,
      request_type: this.request_type,
    };
  }

  set state(val: DomainPickerState) {
    if (!val) {
      return;
    }
    this.begin = val.begin;
    this.end = val.end;
    this.num_commits = val.num_commits;
    this.selected_request_type = val.request_type;
  }
}
