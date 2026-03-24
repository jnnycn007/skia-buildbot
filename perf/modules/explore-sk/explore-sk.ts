/**
 * @module module/explore-sk
 * @description <h2><code>explore-sk</code></h2>
 *
 * Main page of Perf, for exploring data.
 */
import { html, LitElement } from 'lit';
import { customElement, query, state } from 'lit/decorators.js';
import {
  ExploreSimpleSk,
  State as ExploreSimpleSkState,
} from '../explore-simple-sk/explore-simple-sk';
import { stateReflector } from '../../../infra-sk/modules/stateReflector';
import { HintableObject } from '../../../infra-sk/modules/hintable';
import { QueryConfig } from '../json';
import { jsonOrThrow } from '../../../infra-sk/modules/jsonOrThrow';

import '../explore-simple-sk';
import '../favorites-dialog-sk';
import '../test-picker-sk';

import { LoggedIn } from '../../../infra-sk/modules/alogin-sk/alogin-sk';
import { Status as LoginStatus } from '../../../infra-sk/modules/json';
import { errorMessage } from '../errorMessage';
import { TestPickerSk } from '../test-picker-sk/test-picker-sk';
import { queryFromKey } from '../paramtools';
import { PlotSelectionEventDetails } from '../plot-google-chart-sk/plot-google-chart-sk';

import '@material/web/button/outlined-button.js';

@customElement('explore-sk')
export class ExploreSk extends LitElement {
  @query('explore-simple-sk') private exploreSimpleSk!: ExploreSimpleSk;

  @query('#test-picker') private testPicker!: TestPickerSk;

  private stateHasChanged: (() => void) | null = null;

  @state() private defaults: QueryConfig | null = null;

  @state() private _useTestPicker = false;

  createRenderRoot() {
    return this;
  }

  connectedCallback() {
    super.connectedCallback();

    void this.init();

    document.addEventListener('keydown', (e) => {
      if (this.exploreSimpleSk) {
        this.exploreSimpleSk.keyDown(e);
      }
    });

    this.addEventListener('remove-explore', () => {
      this.exploreSimpleSk?.reset();
    });

    this.addEventListener('selection-range-changed', (e) => {
      const detail = (e as CustomEvent<PlotSelectionEventDetails>).detail;
      const state = this.exploreSimpleSk!.state;
      if (!detail.value) {
        return;
      }

      let newBegin = detail.value.begin;
      let newEnd = detail.value.end;

      if (detail.domain === 'commit') {
        const header = this.exploreSimpleSk!.getHeader();
        if (header) {
          if (detail.start !== undefined && header[detail.start]) {
            newBegin = header[detail.start]!.timestamp;
          }
          if (detail.end !== undefined && header[detail.end]) {
            newEnd = header[detail.end]!.timestamp;
          }
        }
      }

      state.begin = newBegin;
      state.end = newEnd;
      if (this.stateHasChanged) {
        this.stateHasChanged();
      }
    });

    LoggedIn()
      .then((status: LoginStatus) => {
        if (this.exploreSimpleSk) {
          this.exploreSimpleSk.state.enable_favorites =
            status.email !== null && status.email !== '';
        }
      })
      .catch(errorMessage);
  }

  private async init() {
    await this.updateComplete;
    await this.initializeDefaults();
    this.stateHasChanged = stateReflector(
      () => this.exploreSimpleSk.state as unknown as HintableObject,
      async (hintableState) => {
        const newState = hintableState as unknown as ExploreSimpleSkState;
        this.exploreSimpleSk.state = newState;
        const isV2 = window.perf.enable_v2_ui && localStorage.getItem('v2_ui') !== 'false';
        this._useTestPicker = !!newState.use_test_picker_query || isV2;
        if (this._useTestPicker) {
          await this.initializeTestPicker();
        }
      }
    );
  }

  render() {
    return html`
      <test-picker-sk
        id="test-picker"
        class="test-picker ${this._useTestPicker ? '' : 'hidden'}"
        @plot-button-clicked=${this.onPlotButtonClicked}></test-picker-sk>
      <explore-simple-sk
        .defaults=${this.defaults}
        .navOpen=${true}
        .openQueryByDefault=${true}
        .useTestPicker=${this._useTestPicker}
        @state_changed=${this.onStateChanged}
        @rendered_traces=${() => this.requestUpdate()}
        @populate-query=${this.onPopulateQuery}></explore-simple-sk>
    `;
  }

  private onStateChanged() {
    if (this.stateHasChanged) {
      this.stateHasChanged();
    }
  }

  private async onPlotButtonClicked() {
    const explore = this.exploreSimpleSk;
    if (explore) {
      const queryStr = this.testPicker.createQueryFromFieldData();
      await explore.addFromQueryOrFormula(true, 'query', queryStr, '');
    }
  }

  private async onPopulateQuery(e: Event) {
    const trace_key = (e as CustomEvent).detail.key;
    const testPickerParams = this.defaults?.include_params ?? null;
    await this.testPicker.populateFieldDataFromQuery(
      queryFromKey(trace_key),
      testPickerParams!,
      {}
    );
    this.testPicker.scrollIntoView();
  }

  /**
   * Fetches defaults from backend and passes them down to the
   * ExploreSimpleSk element.
   */
  private async initializeDefaults() {
    await fetch(`/_/defaults/`, {
      method: 'GET',
    })
      .then(jsonOrThrow)
      .then((json) => {
        this.defaults = json;
      })
      .catch(errorMessage);
  }

  // Initialize TestPickerSk
  private async initializeTestPicker() {
    const testPickerParams = this.defaults?.include_params ?? null;

    if (this.exploreSimpleSk.state.queries && this.exploreSimpleSk.state.queries.length > 0) {
      await this.testPicker.populateFieldDataFromQuery(
        this.exploreSimpleSk.state.queries.join('&'),
        testPickerParams!,
        {}
      );
    } else {
      await this.testPicker.initializeTestPicker(
        testPickerParams!,
        this.defaults?.default_param_selections ?? {},
        false
      );
    }
  }
}
