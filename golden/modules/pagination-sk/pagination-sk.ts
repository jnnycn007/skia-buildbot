/**
 * @module module/pagination-sk
 * @description <h2><code>pagination-sk</code></h2>
 *
 * Widget to let user page forward and backward. Page buttons will be
 * disabled/enabled depending on the offset/total/page_size values.
 *
 * Compatible with the server code httputils.PaginationParams.
 *
 * @attr offset {int} indicates the offset into the list of the items we are paged to.
 *
 * @attr page_size {int} the number of items we go forward/backward on a single page.
 *
 * @attr total {int} the total number of items that can be paged through or MANY if the
 * server doesn't know.
 *
 * @evt page-changed - Sent when user pages forward or backward. Check
 * e.detail.delta for how many pages changed and which direction.
 */

import { html, LitElement } from 'lit';
import { customElement, property } from 'lit/decorators.js';

// MANY (2^31-1, aka math.MaxInt32) is a special value meaning the
// server doesn't know how many items there are, only that it's more
// than are currently being displayed.
const MANY = 2147483647;

export interface PaginationSkPageChangedEventDetail {
  readonly delta: number;
}

@customElement('pagination-sk')
export class PaginationSk extends LitElement {
  @property({ type: Number, reflect: true })
  offset: number = 0;

  @property({ type: Number, attribute: 'page_size', reflect: true })
  pageSize: number = 0;

  @property({ type: Number, reflect: true })
  total: number = 0;

  createRenderRoot() {
    return this;
  }

  render() {
    if (this.pageSize >= this.total && this.total !== MANY) {
      // Hide buttons if only one page is needed.
      return html``;
    }

    return html`
      <button
        ?disabled=${this._currPage() <= 1}
        title="Go to previous page of results."
        @click=${() => this._page(-1)}
        class="prev">
        Prev
      </button>
      <div class="counter">page ${this._currPage()}</div>
      <button
        ?disabled=${!this._canGoNext(this.offset + this.pageSize)}
        title="Go to next page of results."
        @click=${() => this._page(1)}
        class="next">
        Next
      </button>
      <button
        ?disabled=${!this._canGoNext(this.offset + 5 * this.pageSize)}
        title="Skip forward 5 pages of results."
        @click=${() => this._page(5)}
        class="skip">
        +5
      </button>
    `;
  }

  private _currPage() {
    return Math.round(this.offset / this.pageSize) + 1;
  }

  private _canGoNext(next: number) {
    return this.total === MANY ? true : next < this.total;
  }

  private _page(n: number) {
    this.dispatchEvent(
      new CustomEvent<PaginationSkPageChangedEventDetail>('page-changed', {
        detail: {
          delta: n,
        },
        bubbles: true,
        composed: true,
      })
    );
  }
}
