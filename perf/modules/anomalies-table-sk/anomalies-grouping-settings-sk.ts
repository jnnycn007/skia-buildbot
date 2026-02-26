import { html, LitElement } from 'lit';
import { customElement, property } from 'lit/decorators.js';
import '../../../elements-sk/modules/icons/help-icon-sk';
import { AnomalyGroupingConfig, RevisionGroupingMode, GroupingCriteria } from './grouping';

@customElement('anomalies-grouping-settings-sk')
export class AnomaliesGroupingSettingsSk extends LitElement {
  @property({ attribute: false })
  config!: AnomalyGroupingConfig;

  @property({ type: String })
  uniqueId: string = '';

  createRenderRoot() {
    return this;
  }

  render() {
    if (!this.config) {
      return html``;
    }
    const safeId = this.uniqueId || 'default';
    return html`
      <details class="grouping-settings">
        <summary>Grouping Settings</summary>
        <div class="grouping-settings-panel">
          <div class="grouping-setting-group">
            <div class="grouping-setting-label">
              <label for="revision-mode-select-${safeId}">Revision Grouping</label>
              <help-icon-sk
                title="Determines how anomalies are grouped based on their commit range.

* Range Overlap: Groups anomalies with overlapping ranges.
* Exact Range: Groups anomalies with exact same ranges.
* Ignore Range: Groups all anomalies into a single group."></help-icon-sk>
            </div>
            <select
              id="revision-mode-select-${safeId}"
              @change=${(e: Event) => this.onRevisionModeChange(e)}>
              <option value="OVERLAPPING" ?selected=${this.config.revisionMode === 'OVERLAPPING'}>
                Range Overlap
              </option>
              <option value="EXACT" ?selected=${this.config.revisionMode === 'EXACT'}>
                Exact Range
              </option>
              <option value="ANY" ?selected=${this.config.revisionMode === 'ANY'}>
                Ignore Range
              </option>
            </select>
          </div>

          <div class="grouping-setting-group">
            <div class="grouping-setting-label">
              <span>Split Groups By</span>
              <help-icon-sk
                title="Determines more granular group split *after* grouping by revision."></help-icon-sk>
            </div>
            <div class="checkbox-container">
              <label>
                <input
                  type="checkbox"
                  value="BENCHMARK"
                  ?checked=${this.config.groupBy.has('BENCHMARK')}
                  @change=${(e: Event) => this.onGroupByChange(e, 'BENCHMARK')} />
                Benchmark
              </label>
              <label>
                <input
                  type="checkbox"
                  value="BOT"
                  ?checked=${this.config.groupBy.has('BOT')}
                  @change=${(e: Event) => this.onGroupByChange(e, 'BOT')} />
                Bot
              </label>
              <label>
                <input
                  type="checkbox"
                  value="TEST"
                  ?checked=${this.config.groupBy.has('TEST')}
                  @change=${(e: Event) => this.onGroupByChange(e, 'TEST')} />
                Test (without subtests)
              </label>
            </div>
          </div>
        </div>
      </details>
    `;
  }

  private onRevisionModeChange(e: Event) {
    const select = e.target as HTMLSelectElement;
    this.dispatchEvent(
      new CustomEvent<RevisionGroupingMode>('revision-mode-change', {
        detail: select.value as RevisionGroupingMode,
        bubbles: true,
      })
    );
  }

  private onGroupByChange(e: Event, criteria: GroupingCriteria) {
    const checkbox = e.target as HTMLInputElement;
    this.dispatchEvent(
      new CustomEvent<{ criteria: GroupingCriteria; enabled: boolean }>('group-by-change', {
        detail: {
          criteria,
          enabled: checkbox.checked,
        },
        bubbles: true,
      })
    );
  }
}
