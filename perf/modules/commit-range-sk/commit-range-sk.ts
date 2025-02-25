/**
 * @module modules/commit-range-sk
 * @description <h2><code>commit-range-sk</code></h2>
 *
 * Displays a link that describes a range of commits in a repo. This element
 * uses the global `window.perf.commit_range_url`, which can be set on Perf via
 * the command line.
 */
import { html } from 'lit/html.js';
import { define } from '../../../elements-sk/modules/define';
import { ElementSk } from '../../../infra-sk/modules/ElementSk';
import { lookupCids } from '../cid/cid';
import { MISSING_DATA_SENTINEL } from '../const/const';
import { ColumnHeader, CommitNumber } from '../json';
import '../window/window';

// Converts CommitNumbers to Git hashes.
type commitNumberToHashes = (commitNumbers: CommitNumber[]) => Promise<string[]>;

/** The default implementation for commitNumberToHashes run the commit numbers
 *  through cid lookup to get the hashes by making a request to the server.
 */
const defaultcommitNumberToHashes = async (cids: CommitNumber[]): Promise<string[]> => {
  const json = await lookupCids(cids);
  return [json.commitSlice![0].hash, json.commitSlice![1].hash];
};

export class CommitRangeSk extends ElementSk {
  private _trace: number[] = [];

  private _commitIndex: number = -1;

  private _header: (ColumnHeader | null)[] | null = null;

  /** The calculated URL. */
  private _url: string = '';

  /** The link text to display. */
  private _text: string = '';

  private _htmlTemplate = html``;

  /** Determines if text contains links. */
  showLinks: boolean = false;

  // commitNumberToHashes can be replaced to make testing easier.
  private commitNumberToHashes: commitNumberToHashes = defaultcommitNumberToHashes;

  constructor() {
    super(CommitRangeSk.template);
  }

  private static template = (ele: CommitRangeSk) => html`${ele.htmlTemplate}`;

  connectedCallback(): void {
    super.connectedCallback();
    this._upgradeProperty('trace');
    this._upgradeProperty('commitIndex');
    this._upgradeProperty('header');
    this._render();
  }

  reset(): void {
    this._commitIndex = -1;
    this._trace = [];
    this._header = null;
    this._htmlTemplate = html``;
    this._render();
  }

  clear(): void {
    this._url = '';
    this._text = '';
    this._render();
  }

  async recalcLink(): Promise<void> {
    const commitIndex = this._commitIndex;
    if (
      window.perf.commit_range_url === '' ||
      commitIndex === -1 ||
      this._trace.length === 0 ||
      this._header === null
    ) {
      this.clear();
      return;
    }
    // First the previous commit that has data.
    let prevCommit = commitIndex - 1;

    while (prevCommit > 0 && this._trace[prevCommit] === MISSING_DATA_SENTINEL) {
      prevCommit -= 1;
    }

    // If we don't find a second commit then we can't present the information.
    if (prevCommit < 0) {
      this.clear();
      return;
    }

    const startOffset = this._header[prevCommit]?.offset ?? null;
    const endOffset = this._header[commitIndex]?.offset ?? null;
    if (startOffset === null || endOffset === null) {
      this.clear();
      return;
    }
    const cids: CommitNumber[] = [startOffset, endOffset];

    try {
      // Run the commit numbers through cid lookup to get the hashes.
      const hashes = await this.commitNumberToHashes(cids);
      // Create the URL.
      let url = window.perf.commit_range_url;
      url = url.replace('{begin}', hashes[0]);
      url = url.replace('{end}', hashes[1]);

      // Now populate link, including text and url.
      this._url = url;

      // Add +1 to the previous commit to only show commits after previous.
      this._text = `${startOffset + 1} - ${endOffset}`;
      // Check if there are no points between start and end.
      if (startOffset + 1 === endOffset) {
        this._text = `${endOffset}`;
      }
      // Show only text of commit if hover, else show full text with links.
      this._htmlTemplate = html`${this._text}`;
      if (this.showLinks) {
        this._htmlTemplate = html`<a href="${this._url}" target="_blank">${this._text}</a>`;
      }
      this._render();
    } catch (error) {
      console.log(error);
      this.clear();
    }
  }

  /** A single trace. */
  get trace(): number[] {
    return this._trace;
  }

  set trace(val: number[]) {
    this._trace = val;
    this.recalcLink();
  }

  /** An index into trace, the location of the commit being referenced. */
  get commitIndex(): number {
    return this._commitIndex;
  }

  set commitIndex(val: number) {
    this._commitIndex = val;
    this.recalcLink();
  }

  /** The ColumnHeader of the DataFrame that contains the trace. */
  get header(): (ColumnHeader | null)[] | null {
    return this._header;
  }

  set header(val: (ColumnHeader | null)[] | null) {
    this._header = val;
    this.recalcLink();
  }

  get url(): string {
    return this._url;
  }

  get text(): string {
    return this._text;
  }

  get htmlTemplate() {
    return this._htmlTemplate;
  }
}

define('commit-range-sk', CommitRangeSk);
