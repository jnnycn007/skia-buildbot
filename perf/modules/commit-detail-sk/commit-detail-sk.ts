/**
 * @module modules/commit-detail-sk
 * @description <h2><code>commit-detail-sk</code></h2>
 *
 * An element to display information around a single commit.
 *
 */
import { html, LitElement } from 'lit';
import { property, customElement } from 'lit/decorators.js';
import { diffDate } from '../../../infra-sk/modules/human';
import { fromObject } from '../../../infra-sk/modules/query';
import { Commit, CommitNumber } from '../json';

import '@material/web/button/outlined-button.js';

// The range of time (in seconds) to display around a specific commit in the Explore view.
// +/- 4 days provides a reasonable window to see context.
const FOUR_DAYS_IN_SECONDS = 4 * 24 * 60 * 60;

@customElement('commit-detail-sk')
export class CommitDetailSk extends LitElement {
  @property({ type: Object })
  cid: Commit = {
    author: '',
    message: '',
    url: '',
    ts: 0,
    hash: '',
    offset: CommitNumber(0),
    body: '',
  };

  @property({ type: String })
  trace_id: string = '';

  protected createRenderRoot() {
    return this;
  }

  render() {
    return html`
      <div class="linkish">
        <pre>
${this.cid.hash.slice(0, 8)} - ${this.cid.author} - ${diffDate(this.cid.ts * 1000)} - ${this.cid
            .message}</pre
        >
        <div class="tip">
          <md-outlined-button
            @click=${() => {
              if (this.trace_id) {
                const query = {
                  begin: this.cid.ts - FOUR_DAYS_IN_SECONDS,
                  end: this.cid.ts + FOUR_DAYS_IN_SECONDS,
                  keys: this.trace_id,
                  num_commits: 50,
                  request_type: 1,
                  xbaroffset: this.cid.offset,
                };
                this.openLink(`/e/?${fromObject(query)}`);
              } else {
                this.openLink(`/g/e/${this.cid.hash}`);
              }
            }}>
            Explore
          </md-outlined-button>
          <md-outlined-button @click=${() => this.openLink(`/g/c/${this.cid.hash}`)}>
            Cluster
          </md-outlined-button>
          <md-outlined-button @click=${() => this.openLink(`/g/t/${this.cid.hash}`)}>
            Triage
          </md-outlined-button>
          <md-outlined-button @click=${() => this.openLink(this.cid.url)}>
            Commit
          </md-outlined-button>
        </div>
      </div>
    `;
  }

  private openLink(link: string): void {
    window.open(link, '_blank');
  }
}
