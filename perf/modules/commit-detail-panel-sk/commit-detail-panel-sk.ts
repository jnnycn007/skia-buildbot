/**
 * @module modules/commit-detail-panel-sk
 * @description <h2><code>commit-detail-panel-sk</code></h2>
 *
 * @evt commit-selected - Event produced when a commit is selected. The
 *     the event detail contains the serialized cid.CommitDetail and
 *     a simplified description of the commit:
 *
 *     <pre>
 *     {
 *       selected: 2,
 *       description: "foo (foo@example.org) 62W Commit from foo.",
 *       commit: {
 *         author: "foo (foo@example.org)",
 *         url: "skia.googlesource.com/bar",
 *         message: "Commit from foo.",
 *         ts: 1439649751,
 *       },
 *     }
 *     </pre>
 *
 * @attr {Boolean} selectable - A boolean attribute that if true means
 *     that the commits are selectable, and when selected
 *     the 'commit-selected' event is generated.
 *
 * @attr {Number} selected - The index of the selected commit.
 */
import { html, LitElement, TemplateResult } from 'lit';
import { property, customElement } from 'lit/decorators.js';
import { findParent } from '../../../infra-sk/modules/dom';
import '../commit-detail-sk';
import { Commit } from '../json';

export interface CommitDetailPanelSkCommitSelectedDetails {
  selected: number;
  description: string;
  commit: Commit;
}

@customElement('commit-detail-panel-sk')
export class CommitDetailPanelSk extends LitElement {
  @property({ type: Array })
  details: Commit[] = [];

  @property({ type: Boolean })
  hide: boolean = false;

  @property({ type: String })
  trace_id: string = '';

  @property({ type: Boolean, reflect: true })
  selectable: boolean = false;

  @property({ type: Number, reflect: true })
  selected: number = -1;

  protected createRenderRoot() {
    return this;
  }

  private renderRows(): TemplateResult[] {
    if (this.hide) {
      return [html``];
    }
    return this.details.map(
      (item, index) => html`
        <tr data-id="${index}" ?selected="${this._isSelected(index)}">
          <td>
            <commit-detail-sk .cid=${item} .trace_id=${this.trace_id}></commit-detail-sk>
          </td>
        </tr>
      `
    );
  }

  render() {
    return html`
      <table @click=${this._click}>
        ${this.renderRows()}
      </table>
    `;
  }

  private _isSelected(index: number) {
    return this.selectable && index === this.selected;
  }

  private _click(e: MouseEvent) {
    if (!this.selectable) {
      return;
    }
    const ele = findParent(e.target as HTMLElement, 'TR');
    if (!ele) {
      return;
    }
    this.selected = +(ele.dataset.id || '0');
    if (this.selected > this.details.length - 1) {
      return;
    }
    const commit = this.details[this.selected];
    const detail = {
      selected: this.selected,
      description: `${commit.author} -  ${commit.message}`,
      commit,
    };
    this.dispatchEvent(
      new CustomEvent<CommitDetailPanelSkCommitSelectedDetails>('commit-selected', {
        detail,
        bubbles: true,
      })
    );
  }
}
