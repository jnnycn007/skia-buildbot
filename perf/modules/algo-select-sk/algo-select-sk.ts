/**
 * @module modules/algo-select-sk
 * @description <h2><code>algo-select-sk</code></h2>
 *
 * Displays and allows changing the clustering algorithm.
 *
 * @evt algo-change - Sent when the algo has changed. The value is stored
 *    in e.detail.algo.
 *
 * @attr {string} algo - The algorithm name.
 */
import '../../../elements-sk/modules/select-sk';
import { html, LitElement } from 'lit';
import { customElement, property } from 'lit/decorators.js';
import { $ } from '../../../infra-sk/modules/dom';
import { SelectSkSelectionChangedEventDetail } from '../../../elements-sk/modules/select-sk/select-sk';
import { ClusterAlgo } from '../json';

function toClusterAlgo(s: string): ClusterAlgo {
  const allowed = ['kmeans', 'stepfit'];
  if (allowed.indexOf(s) !== -1) {
    return s as ClusterAlgo;
  }
  return 'kmeans';
}

export interface AlgoSelectAlgoChangeEventDetail {
  algo: ClusterAlgo;
}

@customElement('algo-select-sk')
export class AlgoSelectSk extends LitElement {
  private _algo: ClusterAlgo = 'kmeans';

  @property({ type: String, noAccessor: true })
  get algo(): ClusterAlgo {
    return this._algo;
  }

  set algo(val: ClusterAlgo) {
    const oldVal = this._algo;
    this._algo = toClusterAlgo(val);
    if (this._algo !== oldVal) {
      this.setAttribute('algo', this._algo);
      this.requestUpdate('algo', oldVal);
    }
  }

  createRenderRoot() {
    return this;
  }

  render() {
    return html`
      <select-sk @selection-changed=${this._selectionChanged}>
        <div
          class="kmeans"
          value="kmeans"
          ?selected=${this.algo === 'kmeans'}
          title="Use k-means clustering on the trace shapes and look for a step on the cluster centroid.">
          K-Means
        </div>
        <div
          class="stepfit"
          value="stepfit"
          ?selected=${this.algo === 'stepfit'}
          title="Look for a step in each individual trace.">
          Individual
        </div>
      </select-sk>
    `;
  }

  private _selectionChanged(e: CustomEvent<SelectSkSelectionChangedEventDetail>) {
    let index = e.detail.selection;
    if (index < 0) {
      index = 0;
    }
    this.algo = toClusterAlgo($('div', this)[index].getAttribute('value') || '');
    const detail = {
      algo: this.algo,
    };
    this.dispatchEvent(
      new CustomEvent<AlgoSelectAlgoChangeEventDetail>('algo-change', {
        detail,
        bubbles: true,
      })
    );
  }
}
