// ==UserScript==
// @name         Super Monkey Test
// @namespace    http://tampermonkey.net/
// @version      1.30
// @description  Automates UI interactions to stress-test the MultiGraph viewer.
// @author       seawardt@google.com
// @match        *://chrome-perf.corp.goog/m/*
// @grant        none
// ==/UserScript==

/**
 * SUPER MONKEY TEST PLAN - v1.30
 *
 * DESCRIPTION:
 * This script automates a stress test sequence for the MultiGraph UI. It navigates
 * through filter fields, generates graphs, tests split/unsplit logic, and verifies
 * state consistency (URL, Sidebar, Pagination).
 *
 * HOW TO RUN:
 * Option 1: Chrome DevTools Snippet (Manual)
 *   1. Open DevTools (F12) -> Sources -> Snippets.
 *   2. Create a new snippet, paste this code.
 *   3. Right-click and "Run" or press Ctrl+Enter.
 *
 * Option 2: Tampermonkey (Automatic)
 *   1. Install Tampermonkey extension.
 *   2. Create a new script and paste this code.
 *   3. Navigate to a MultiGraph page (/m/).
 *
 * FEATURES:
 * - Dynamic path discovery (finds valid datasets automatically).
 * - Robust error handling (retries, timeouts, toast detection).
 * - State verification (URL parameters, Graph counts, Sidebar rows).
 * - Pagination iteration and validation.
 */

(async function SuperMonkeyTest() {
  const VERSION = '1.30';
  console.clear();
  console.log(`SUPER MONKEY v${VERSION} - DATA AGNOSTIC LOADED`);
  const START_TIME = Date.now();

  // --- GLOBAL ERROR TRAP ---
  window.addEventListener('unhandledrejection', function (event) {
    console.error('UNHANDLED PROMISE REJECTION:', event.reason);
    if (typeof TEST_HISTORY !== 'undefined') {
      TEST_HISTORY.push(`[ERROR] Unhandled Rejection: ${event.reason}`);
    }
  });

  // --- 1. PRE-FLIGHT CHECK ---
  if (!window.location.pathname.includes('/m')) {
    console.error('WRONG PAGE! This script requires the MultiGraph page (/m).');
    alert('WRONG PAGE!\nPlease navigate to the MultiGraph page (/m) before running this script.');
    return;
  }

  // --- 2. SELF-CLEANUP ---
  if (window.__TOAST_OBSERVER__) {
    console.log('Disconnecting previous Toast Observer...');
    window.__TOAST_OBSERVER__.disconnect();
  }

  // --- GLOBAL STATE ---
  let IS_PAUSED = false;
  let FATAL_ERROR = null;
  let BACKEND_BUSY_UNTIL = 0; // Timestamp to wait until
  let TEST_COUNTER = 0;
  const TEST_HISTORY = [];

  const CONTEXT = {
    path: [],
    highVolumePath: [],
    targetSplitField: null,
    splitOptions: [],
    initialTraceCount: 0,
    graphsRendered: false,
  };

  // --- CONFIGURATION ---
  const CONFIG = {
    mainAppTag: 'explore-multi-sk',
    pickerTag: 'test-picker-sk',
    timeouts: {
      short: 3000,
      medium: 8000,
      long: 15000,
      xl: 45000,
    },
    thresholds: {
      minTraces: 15,
      maxTraces: 80,
      highVolume: 50,
    },
  };

  // --- HELPER FUNCTIONS ---
  const checkFailure = () => {
    if (FATAL_ERROR) throw new Error(`FATAL ERROR: ${FATAL_ERROR}`);
  };

  const rawSleep = (ms) => new Promise((r) => setTimeout(r, ms));

  const togglePause = () => {
    IS_PAUSED = !IS_PAUSED;
    TestHUD.update(null, IS_PAUSED ? 'PAUSED' : 'RESUMED', 'WARN');
    if (IS_PAUSED) console.warn('SCRIPT PAUSED');
    else console.info('SCRIPT RESUMED');
  };

  const checkPauseState = async () => {
    if (IS_PAUSED) {
      while (IS_PAUSED) await rawSleep(200);
    }
    // Check if backend is flagged as busy
    if (Date.now() < BACKEND_BUSY_UNTIL) {
      const wait = BACKEND_BUSY_UNTIL - Date.now();
      console.log(`Backend busy (Toast detected). Waiting ${wait}ms...`);
      await rawSleep(wait);
    }
  };

  const sleep = async (ms) => {
    checkFailure();
    await checkPauseState();
    return new Promise((r) => {
      setTimeout(() => {
        checkFailure();
        r();
      }, ms);
    });
  };

  const shuffleArray = (array) => {
    for (let i = array.length - 1; i > 0; i--) {
      const j = Math.floor(Math.random() * (i + 1));
      [array[i], array[j]] = [array[j], array[i]];
    }
    return array;
  };

  const getTimestamp = () => {
    const now = new Date();
    return now.toISOString().split('T')[1].slice(0, -1);
  };

  // --- DOM UTILS ---
  const DomUtils = {
    getExploreApp: () => document.querySelector(CONFIG.mainAppTag),
    getTestPicker: () => {
      let app = DomUtils.getExploreApp();
      return app ? (app.shadowRoot || app).querySelector(CONFIG.pickerTag) : null;
    },
    getPickerFields: () => {
      const picker = DomUtils.getTestPicker();
      if (!picker) return [];
      const root = picker.shadowRoot || picker;
      return Array.from(root.querySelectorAll('picker-field-sk'));
    },
    getComboBox: (field) => {
      if (!field) return null;
      return field.shadowRoot
        ? field.shadowRoot.querySelector('vaadin-multi-select-combo-box')
        : field.querySelector('vaadin-multi-select-combo-box');
    },
    getGraphs: (silent = false) => {
      const app = DomUtils.getExploreApp();
      if (!app) return [];
      const root = app.shadowRoot || app;
      const container = root.querySelector('#graphContainer');
      if (!container) return [];
      const all = Array.from(container.querySelectorAll('explore-simple-sk'));

      const visible = all.filter((el, i) => {
        const style = getComputedStyle(el);
        const rect = el.getBoundingClientRect();
        const isVisible =
          style.display !== 'none' &&
          style.visibility !== 'hidden' &&
          rect.width > 50 &&
          rect.height > 50;
        return isVisible;
      });

      if (visible.length === 0 && all.length > 0) {
        if (!silent) {
          Logger.warn(
            `getGraphs: Found ${all.length} graphs in DOM, ` +
              `but NONE are visible. Returning all.`
          );
        }
        return all;
      }
      return visible;
    },
    getPlotButton: () => {
      const picker = DomUtils.getTestPicker();
      return picker ? (picker.shadowRoot || picker).querySelector('#plot-button') : null;
    },
    getTraceCount: () => {
      const picker = DomUtils.getTestPicker();
      if (!picker) return 0;
      const root = picker.shadowRoot || picker;
      const el = root.querySelector('.test-picker-sk-matches-container');

      if (el && el.querySelector('spinner-sk[active]')) {
        return -1;
      }
      if (el && el.innerText.trim() === 'Traces:') return -1;

      const match = el ? el.innerText.match(/Traces:\s*([\d,]+)/) : null;
      return match ? parseInt(match[1].replace(/,/g, ''), 10) : -1;
    },
    getSidebarCount: () => {
      const graphs = DomUtils.getGraphs();
      let total = 0;
      graphs.forEach((g) => {
        const root = g.shadowRoot || g;
        const findRows = (el) => {
          // Note: This recursive search is fragile (Tracked in b/485902011).
          // Ideally, add data-testid="sidebar-rows" to the <div id="rows">
          // in explore-simple-sk.ts for direct selection.
          if (el.id === 'rows' && el.querySelector('li')) return el;
          if (el.shadowRoot) {
            const res = findRows(el.shadowRoot);
            if (res) return res;
          }
          if (el.children) {
            for (let i = 0; i < el.children.length; i++) {
              const res = findRows(el.children[i]);
              if (res) return res;
            }
          }
          return null;
        };
        const rows = findRows(root);
        if (rows) total += rows.querySelectorAll('li label').length;
      });
      return total;
    },
    getSpinner: () => {
      const allSpinners = Array.from(document.querySelectorAll('spinner-sk'));
      const active = allSpinners.find(
        (s) => s.hasAttribute('active') || getComputedStyle(s).opacity !== '0'
      );
      if (active) return active;

      const app = DomUtils.getExploreApp();
      if (app && app.shadowRoot) {
        const s = app.shadowRoot.querySelector('spinner-sk');
        if (s && (s.hasAttribute('active') || getComputedStyle(s).opacity !== '0')) return s;
        const p = app.shadowRoot.querySelector(CONFIG.pickerTag);
        if (p && p.shadowRoot) {
          const sp = p.shadowRoot.querySelector('spinner-sk');
          if (sp && (sp.hasAttribute('active') || getComputedStyle(sp).opacity !== '0')) return sp;
        }
      }
      return null;
    },
    getPagination: () => {
      const app = DomUtils.getExploreApp();
      if (!app) return null;
      const root = app.shadowRoot || app;
      const page = root.querySelector('pagination-sk');
      if (page) {
        return {
          element: page,
          pageSize: parseInt(page.getAttribute('page_size') || '0', 10),
          total: parseInt(page.getAttribute('total') || '0', 10),
          offset: parseInt(page.getAttribute('offset') || '0', 10),
        };
      }
      return null;
    },
    clickNextPage: async () => {
      const p = DomUtils.getPagination();
      if (!p || !p.element) return false;
      const root = p.element.shadowRoot || p.element;
      const btn = root.querySelector('button.next');
      if (btn && !btn.disabled) {
        await Interaction.click(btn);
        return true;
      }
      return false;
    },
    getLoadAllChartsButton: () => {
      const app = DomUtils.getExploreApp();
      if (!app) return null;
      const root = app.shadowRoot || app;
      const buttons = Array.from(root.querySelectorAll('button'));
      return buttons.find((b) => b.innerText.includes('Load All Charts'));
    },
    getReduceMessage: () => {
      const picker = DomUtils.getTestPicker();
      return picker ? (picker.shadowRoot || picker).querySelector('#max-message') : null;
    },
    isChromeInternal: () => {
      const header = document.querySelector('header h1.name');
      return !!(header && header.innerText.includes('chrome-internal'));
    },
  };

  // --- INTERACTION ---
  const Interaction = {
    verifySelection: async (field, expectedItems) => {
      const combo = DomUtils.getComboBox(field);
      if (!combo) return false;
      for (let i = 0; i < 10; i++) {
        const current = combo.selectedItems || [];
        const match =
          current.length === expectedItems.length &&
          current.every((val) => expectedItems.includes(val));
        if (match) return true;
        await sleep(200);
      }
      return false;
    },
    setSelection: async (field, items) => {
      await sleep(200);

      if (!field || !field.isConnected) {
        Logger.warn(`Cannot set selection on disconnected field '${field ? field.label : 'null'}'`);
        return;
      }

      const combo = DomUtils.getComboBox(field);
      if (combo) {
        for (let attempt = 1; attempt <= 3; attempt++) {
          Logger.info(`Action: Set '${field.label}' to [${items.join(', ')}] (Attempt ${attempt})`);
          combo.selectedItems = items;
          combo.dispatchEvent(new CustomEvent('change', { bubbles: true }));
          combo.dispatchEvent(
            new CustomEvent('selected-items-changed', { detail: { value: items }, bubbles: true })
          );

          const ok = await Interaction.verifySelection(field, items);
          if (ok) return;

          Logger.warn(`Selection verify failed. Retrying...`);
          await sleep(500);
        }
        Logger.error(`Failed to set '${field.label}' to [${items.join(', ')}] after 3 attempts.`);
      } else {
        Logger.warn(`Combo box not found for '${field.label}'`);
      }
    },
    setSplit: async (field, isSplit) => {
      await sleep(500);
      const checkbox = field.shadowRoot
        ? field.shadowRoot.querySelector('#split-by')
        : field.querySelector('#split-by');
      if (checkbox) {
        Logger.info(`Action: Set Split '${field.label}' to ${isSplit}`);
        checkbox.checked = isSplit;
        checkbox.dispatchEvent(new CustomEvent('change', { bubbles: true }));
      } else {
        Logger.warn(`Split checkbox not found for ${field.label}. Dispatching event manually.`);
        field.dispatchEvent(
          new CustomEvent('split-by-changed', {
            detail: { param: field.label, split: isSplit },
            bubbles: true,
            composed: true,
          })
        );
      }
    },
    setPrimary: async (field, isPrimary) => {
      await sleep(500);
      const checkbox = (field.shadowRoot || field).querySelector('#select-primary');
      if (checkbox) {
        Logger.info(`Action: Set Primary '${field.label}' to ${isPrimary}`);
        if (checkbox.checked !== isPrimary) {
          checkbox.checked = isPrimary;
          checkbox.dispatchEvent(new CustomEvent('change', { bubbles: true }));
        }
      } else {
        Logger.warn(`Primary checkbox not found for ${field.label}.`);
      }
    },
    setSelectAll: async (field, select) => {
      await sleep(500);
      const checkbox = (field.shadowRoot || field).querySelector('#select-all');
      if (checkbox) {
        if (checkbox.checked !== select) {
          Logger.info(`Action: Click 'Select All' on '${field.label}'`);
          checkbox.checked = select;
          checkbox.dispatchEvent(new CustomEvent('change', { bubbles: true }));
        }
      } else {
        Logger.warn(`Select All checkbox not found for ${field.label}`);
      }
    },
    isSelectAllChecked: (field) => {
      const checkbox = (field.shadowRoot || field).querySelector('#select-all');
      return checkbox ? checkbox.checked : false;
    },
    click: async (element) => {
      if (!element) return;
      element.click();
      element.dispatchEvent(
        new MouseEvent('click', { bubbles: true, cancelable: true, view: window })
      );
      await sleep(500);
    },
  };

  // --- WAIT UTILS ---

  const waitForCondition = async (predicate, timeout = 10000, msg) => {
    const start = Date.now();
    while (Date.now() - start < timeout) {
      checkFailure();
      await checkPauseState(); // Check if we need to pause
      if (predicate()) return true;
      await sleep(100);
    }
    Logger.warn(`Timeout waiting for: ${msg}`);
    return false;
  };

  const waitForGraphs = async (timeout = 10000) => {
    const start = Date.now();
    while (Date.now() - start < timeout) {
      checkFailure();
      await checkPauseState();
      const spinner = DomUtils.getSpinner();
      const graphs = DomUtils.getGraphs(true).length;
      if (!spinner && graphs > 0) {
        // Ensure UI settles (No 'render-complete' event available on graph container)
        await rawSleep(1000);
        return true;
      }
      await sleep(100);
    }
    Logger.warn('Timeout waiting for graphs (and spinner inactive).');
    return false;
  };

  const waitForStableState = async (timeout = 10000) => {
    const start = Date.now();
    await sleep(200);
    return new Promise((resolve) => {
      const check = async () => {
        if (FATAL_ERROR) {
          resolve(false);
          return;
        }

        // Dynamic check for backend busy
        if (Date.now() < BACKEND_BUSY_UNTIL) {
          setTimeout(check, 500);
          return;
        }

        if (Date.now() - start > timeout) {
          Logger.warn(`Timeout waiting for stable state (${timeout}ms).`);
          resolve(false);
          return;
        }

        const traceCount = DomUtils.getTraceCount();
        const fields = DomUtils.getPickerFields();
        const anyDisabled = fields.some((f) => {
          const combo = DomUtils.getComboBox(f);
          return combo && combo.hasAttribute('readonly');
        });

        if (traceCount >= 0 && !anyDisabled) {
          resolve(true);
        } else {
          setTimeout(check, 50);
        }
      };
      check();
    });
  };

  // --- LOGGER & HUD ---
  const Logger = {
    info: (msg) =>
      console.log(`%c [${getTimestamp()}] INFO ${msg}`, 'color: #2196F3; font-weight: bold;'),
    action: (msg) => {
      const logLine = `[${getTimestamp()}] ${msg}`;
      console.log(`%c ${logLine}`, 'color: #AA00FF; font-weight: bold; font-size: 1.1em;');
      TEST_HISTORY.push(`TEST_HEADER: ${msg}`);
      TEST_HISTORY.push(logLine);
    },
    success: (msg) => {
      console.log(`%c [${getTimestamp()}] SUCCESS ${msg}`, 'color: #4CAF50; font-weight: bold;');
      TEST_HISTORY.push(`[PASS] ${msg}`);
    },
    warn: (msg) =>
      console.warn(`%c [${getTimestamp()}] WARN ${msg}`, 'color: #FF9800; font-weight: bold;'),
    error: (msg) => {
      console.error(
        `%c [${getTimestamp()}] ERROR ${msg}`,
        'color: #F44336; font-weight: bold; font-size: 1.2em;'
      );
      TEST_HISTORY.push(`[FAIL] ${msg}`);
    },
  };

  const logStats = (label) => {
    const traces = DomUtils.getTraceCount();
    const graphs = DomUtils.getGraphs().length;
    const sidebar = DomUtils.getSidebarCount();
    Logger.info(
      `[STATS] ${label} -> Traces: ${traces}, Graphs: ${graphs}, Sidebar Rows: ${sidebar}`
    );
    return { traces, graphs, sidebar };
  };

  const TestHUD = {
    el: null,
    init() {
      const existing = document.getElementById('test-hud-overlay');
      if (existing) existing.remove();
      this.el = document.createElement('div');
      this.el.id = 'test-hud-overlay';
      Object.assign(this.el.style, {
        position: 'fixed',
        bottom: '20px',
        right: '20px',
        background: 'rgba(33, 33, 33, 0.95)',
        color: '#eee',
        padding: '16px',
        borderRadius: '8px',
        zIndex: '999999',
        fontFamily: 'Consolas, monospace',
        border: '1px solid #00BCD4',
        width: '600px',
        fontSize: '12px',
      });
      document.body.appendChild(this.el);
      this.el.addEventListener('click', (e) => {
        if (e.target.id === 'monkey-pause-btn') togglePause();
      });
      this.update('Initializing...', 'Waiting...', 'PENDING');
    },
    update(action, detail, status) {
      const btnText = IS_PAUSED ? 'RESUME' : 'PAUSE';
      this.el.innerHTML =
        '<div style="border-bottom:1px solid #555; margin-bottom:5px; ' +
        'font-weight:bold; color:#00BCD4; display:flex; justify-content:space-between;">' +
        `<span>SUPER MONKEY v${VERSION}</span>` +
        '<button id="monkey-pause-btn" style="cursor:pointer; background:#555; ' +
        `color:white; border:none; padding:2px 6px;">${btnText}</button>` +
        '</div>' +
        `<div><strong style="color:#B388FF">TEST #${TEST_COUNTER}:</strong> ${action || ''}</div>` +
        `<div style="color:#aaa">${detail || ''}</div>` +
        `<div style="margin-top:5px; font-weight:bold; color:${
          status === 'FAIL' ? '#F44336' : '#4CAF50'
        }">` +
        `STATUS: ${status}</div>`;
    },
    showSummary(duration) {
      console.log('Generating Summary. History items:', TEST_HISTORY.length);
      const blocks = [];
      let currentBlock = { title: 'Initialization', logs: [], status: 'PASS' };

      TEST_HISTORY.forEach((line) => {
        if (line.startsWith('TEST_HEADER:')) {
          blocks.push(currentBlock);
          currentBlock = {
            title: line.replace('TEST_HEADER:', '').trim(),
            logs: [],
            status: 'PASS',
          };
        } else {
          currentBlock.logs.push(line);
          if (line.includes('[FAIL]')) currentBlock.status = 'FAIL';
        }
      });
      blocks.push(currentBlock);

      const tests = blocks.filter((b) => b.title.includes('Test ')).length;
      const fails = blocks.filter((b) => b.status === 'FAIL').length;
      const passes = tests - fails;
      const statusColor = fails > 0 ? '#F44336' : '#4CAF50';
      const statusText = fails > 0 ? 'FAILURES DETECTED' : 'ALL TESTS PASSED';

      const summaryHTML =
        '<div style="position:fixed; top:50%; left:50%; transform:translate(-50%, -50%); ' +
        'background:rgba(20, 20, 20, 0.98); color:#eee; padding:0; ' +
        `border:2px solid ${statusColor}; border-radius:12px; z-index:1000000; ` +
        'width:80%; max-width:900px; max-height:85vh; overflow:hidden; display:flex; ' +
        'flex-direction:column; box-shadow: 0 0 30px rgba(0,0,0,0.9);">' +
        `<div style="background:${statusColor}22; padding:20px; ` +
        `border-bottom:1px solid ${statusColor}; ` +
        'display:flex; justify-content:space-between; align-items:center;">' +
        '<div>' +
        `<h2 style="color:${statusColor}; margin:0; font-size:24px;">${statusText}</h2>` +
        `<div style="color:#aaa; margin-top:5px;">Duration: <strong>${duration}s</strong> | ` +
        `Tests: <strong>${tests}</strong></div>` +
        '</div>' +
        '<div style="text-align:right; font-size:18px;">' +
        `<span style="color:#4CAF50; margin-right:15px;">PASS: ${passes}</span>` +
        `<span style="color:#F44336;">FAIL: ${fails}</span>` +
        '</div></div>' +
        '<div style="flex:1; overflow-y:auto; padding:20px; font-family:Consolas, monospace;">' +
        '<h3 style="color:#ddd; margin-top:0;">Detailed Log</h3>' +
        '<div style="background:#111; padding:15px; border-radius:6px; border:1px solid #333; ' +
        `font-size:12px; line-height:1.5;">${blocks
          .map((block) => {
            if (block.logs.length === 0 && block.title === 'Initialization') return '';
            const hasFail = block.logs.some((l) => l.includes('[FAIL]'));
            const finalColor = hasFail ? '#F44336' : '#4CAF50';
            const statusTag = hasFail ? '[FAIL]' : '[PASS]';

            const logsHTML = block.logs
              .map((line) => {
                if (line.includes('[PASS]'))
                  return (
                    '<div style="color:#888;">PASS <span style="color:#aaa">' +
                    line.replace('[PASS]', '').trim() +
                    '</span></div>'
                  );
                if (line.includes('[FAIL]'))
                  return (
                    '<div style="color:#F44336; font-weight:bold;">FAIL ' +
                    line.replace('[FAIL]', '').trim() +
                    '</div>'
                  );
                if (line.includes('[STATS]'))
                  return `<div style="color:#2196F3; font-style:italic;">${line}</div>`;
                return `<div style="color:#666;">${line}</div>`;
              })
              .join('');
            return (
              `<div style="margin-bottom:15px; border-left:4px solid ${finalColor}; ` +
              `padding-left:15px;">` +
              `<div style="color:${finalColor}; font-weight:bold; ` +
              `font-size:16px; margin-bottom:8px; ` +
              'border-bottom:1px solid #333; padding-bottom:5px; ' +
              'display:flex; justify-content:space-between;">' +
              `<span>${block.title}</span><span style="opacity:0.8;">${statusTag}</span></div>` +
              `<div>${logsHTML}</div></div>`
            );
          })
          .join('')}</div></div>` +
        '<div style="padding:15px 20px; background:#222; ' +
        'text-align:right; border-top:1px solid #333;">' +
        '<button id="test-summary-close-btn" ' +
        'style="background:#333; color:white; border:1px solid #555; padding:10px 25px; ' +
        'font-weight:bold; cursor:pointer; border-radius:4px; font-size:14px; ' +
        'transition:background 0.2s;">Close Report</button></div></div>';
      const overlay = document.createElement('div');
      overlay.id = 'test-summary-overlay';
      overlay.innerHTML = summaryHTML;
      document.body.appendChild(overlay);
      document.getElementById('test-summary-close-btn').addEventListener('click', () => {
        const el = document.getElementById('test-summary-overlay');
        if (el) el.remove();
      });
    },
  };

  // --- TEST LOGIC ---

  const assert = (condition, msg, expected, actual) => {
    let detail = msg;
    if (expected !== undefined && actual !== undefined) {
      detail += ` (Expected: ${expected}, Actual: ${actual})`;
    }

    if (!condition) {
      Logger.error(`ASSERT FAILED: ${detail}`);
      if (msg.toLowerCase().includes('graph')) {
        console.warn('--- GRAPH DEBUG INFO ---');
        DomUtils.getGraphs(true);
      }
      TestHUD.update(null, detail, 'FAIL');
      return false;
    }
    Logger.success(detail);
    return true;
  };

  const runTest = async (name, description, testFn) => {
    TEST_COUNTER++;
    Logger.action(`Test ${TEST_COUNTER}: ${name}`);
    Logger.info(`Description: ${description}`);
    TestHUD.update(name, description, 'RUNNING');
    try {
      await testFn();
      TestHUD.update(name, 'Completed successfully', 'PASS');
    } catch (e) {
      TestHUD.update(name, e.message, 'FAIL');
      throw e;
    }
    await sleep(1000);
  };

  // --- DYNAMIC DISCOVERY ---

  const findPath = async (fieldIndex, minTraces = 1, maxTraces = 999999, deadline = 0) => {
    if (deadline === 0) deadline = Date.now() + 120000;
    if (Date.now() > deadline) {
      Logger.warn('[Discovery] Global deadline exceeded. Stopping.');
      return null;
    }

    if (fieldIndex > 0) {
      const btn = DomUtils.getPlotButton();
      if (btn && !btn.disabled) {
        Logger.info('Plot button enabled. Path found.');
        return [];
      }
    }

    let fields = DomUtils.getPickerFields();
    if (fieldIndex < fields.length) {
      const f = fields[fieldIndex];
      if (!f.options || f.options.length === 0) {
        await waitForCondition(
          () => f.options && f.options.length > 0,
          2000,
          `Options for '${f.label}'`
        );
      }
    }

    fields = DomUtils.getPickerFields();
    if (fieldIndex >= fields.length) {
      const btn = DomUtils.getPlotButton();
      if (btn && !btn.disabled) return [];
      return null;
    }

    const field = fields[fieldIndex];
    let options = field.options || [];
    Logger.info(`[Discovery] At field '${field.label}' (${options.length} options)`);

    options = shuffleArray([...options]);
    const maxAttempts = Math.min(options.length, 20);
    let failedRecursionCount = 0;
    let lowCountRejectCounter = 0;

    for (let i = 0; i < maxAttempts; i++) {
      if (Date.now() > deadline) break;
      const val = options[i];
      const initialTraceCount = DomUtils.getTraceCount();

      if (fieldIndex === 0) {
        // Root Field
        if (i > 0) {
          await Interaction.setSelection(field, []);
          await waitForStableState(CONFIG.timeouts.short);
        }
        await Interaction.setSelection(field, [val]);
        if (initialTraceCount > 0) {
          await waitForCondition(
            () => {
              const t = DomUtils.getTraceCount();
              if (t > 0 && t < 20 && minTraces > 20) return true;
              return t !== initialTraceCount && t !== -1;
            },
            CONFIG.timeouts.medium,
            'Trace Count Change'
          );
        }
      } else {
        // Child Field
        await Interaction.setSelection(field, []);
        await waitForStableState(CONFIG.timeouts.short);
        await Interaction.setSelection(field, [val]);
      }

      // Wait for count
      const startWait = Date.now();
      let traceCount = 0;
      while (Date.now() - startWait < 10000) {
        checkFailure();
        const t = DomUtils.getTraceCount();
        if (t > 0) {
          traceCount = t;
          break;
        }
        if (t === 0 && Date.now() - startWait > 2000) {
          traceCount = 0;
          break;
        }
        await sleep(50);
      }
      if (DomUtils.getTraceCount() >= 0) traceCount = DomUtils.getTraceCount();

      Logger.info(`[Discovery] '${val}' -> ${traceCount} traces.`);

      // If valid, Recurse IMMEDIATELY
      if (traceCount > 0 || DomUtils.getPickerFields().length > fieldIndex + 1) {
        if (traceCount < minTraces && minTraces >= 50) {
          Logger.warn(
            `[Discovery] '${val}' -> ${traceCount} traces. Too few (Min: ${minTraces}). Skipping.`
          );
          lowCountRejectCounter++;

          if (lowCountRejectCounter >= 5) {
            Logger.warn(
              `[Discovery] ${lowCountRejectCounter} consecutive options yielded ` +
                `too few traces. Abandoning '${field.label}' options.`
            );
            if (fieldIndex > 0) {
              await Interaction.setSelection(field, []);
              await waitForStableState(CONFIG.timeouts.short);
            }
            return null;
          }

          if (fieldIndex > 0) {
            await Interaction.setSelection(field, []);
            await waitForStableState(CONFIG.timeouts.short);
          }
          continue;
        }

        if (traceCount === 0) {
          // Double check if next field populated (implicit success)
          const fs = DomUtils.getPickerFields();
          const nextF = fs[fieldIndex + 1];
          if (!nextF || !nextF.options || nextF.options.length === 0) {
            Logger.warn(`[Discovery] '${val}' -> 0 traces and no children. Skipping.`);
            continue;
          }
        }

        Logger.info(`[Discovery] Candidate '${val}' accepted. Recursing...`);

        if (fieldIndex + 1 >= DomUtils.getPickerFields().length) {
          return [{ label: field.label, value: val, element: field }];
        }

        const res = await findPath(fieldIndex + 1, minTraces, maxTraces, deadline);
        if (res !== null) {
          return [{ label: field.label, value: val, element: field }, ...res];
        } else {
          Logger.warn(`[Discovery] Recursion failed from '${val}'. Trying next option...`);
          failedRecursionCount++;
          if (failedRecursionCount >= 2) {
            Logger.warn(
              `[Discovery] Too many recursion failures from '${field.label}'. Abandoning.`
            );
            if (fieldIndex > 0) await Interaction.setSelection(field, []);
            return null;
          }
          // Cleanup for next loop
          if (fieldIndex > 0) {
            await Interaction.setSelection(field, []);
            await waitForStableState(CONFIG.timeouts.short);
          }
        }
      }
    }

    return null;
  };

  // --- TEST SUITE ---

  try {
    TestHUD.init();

    const toast = document.querySelector('error-toast-sk') || document.querySelector('toast-sk');
    if (toast) {
      new MutationObserver(() => {
        const text = toast.innerText;

        // Handle "Pending Query" as a Warning (Wait), NOT Fatal
        if (text.includes('pending query')) {
          Logger.warn(`Backend Busy: "${text}". Pausing 5s...`);
          BACKEND_BUSY_UNTIL = Date.now() + 5000;
          return;
        }

        // Downgrade "Query must not be empty" to warning + pause
        if (text.includes('The query must not be empty')) {
          Logger.warn(`Empty Query Error: "${text}". Pausing 3s...`);
          BACKEND_BUSY_UNTIL = Date.now() + 3000;
          return;
        }

        if (text.includes('No data found')) {
          Logger.error(`Toast Error: ${text}`);
        }
      }).observe(toast, { childList: true, subtree: true, characterData: true });
    }

    // 1. Dynamic Load & Optimize
    await runTest(
      '1. Dynamic Load & Refine',
      'Drill down through options to find a dataset with ~50 traces.',
      async () => {
        const fields = DomUtils.getPickerFields();
        logStats('Before Discovery');

        // 1. Find High Volume (>20)
        Logger.info('Phase 1: Finding High Volume Path...');
        CONTEXT.highVolumePath = await findPath(0, 20, 10000000); // Reduced from 50 to 20

        if (!CONTEXT.highVolumePath || CONTEXT.highVolumePath.length === 0) {
          Logger.warn('High volume path failed. Finding ANY path.');
          CONTEXT.highVolumePath = await findPath(0, 1, 10000000);
        }

        if (!CONTEXT.highVolumePath)
          throw new Error(
            'Could not find any valid path. Check backend connectivity or increase timeouts.'
          );
        Logger.success(
          `High Volume Path: ${CONTEXT.highVolumePath
            .map((p) => `${p.label}=${p.value}`)
            .join(' > ')}`
        );

        // 2. Optimize (Reduce to 5-50)
        Logger.info('Phase 2: Optimizing for 15-80 traces...');

        const startIdx = CONTEXT.highVolumePath.length;
        const optimizedTail = await findPath(startIdx, 15, 80);

        if (optimizedTail) {
          CONTEXT.path = [...CONTEXT.highVolumePath, ...optimizedTail];
          Logger.success('Optimization successful!');
        } else {
          Logger.warn('Could not optimize further. Using High Volume path.');
          CONTEXT.path = CONTEXT.highVolumePath;
        }

        Logger.success(
          `Final Test Path: ${CONTEXT.path.map((p) => `${p.label}=${p.value}`).join(' > ')}`
        );

        const btn = DomUtils.getPlotButton();
        btn.click();

        await waitForStableState(10000);
        await waitForCondition(
          () => {
            const t = DomUtils.getTraceCount();
            return t > 0;
          },
          30000,
          'Trace Count > 0'
        );

        await waitForCondition(() => DomUtils.getGraphs().length > 0, 5000, 'Graph Render');

        logStats('After Plot');

        CONTEXT.initialTraceCount = DomUtils.getTraceCount();
        assert(
          CONTEXT.initialTraceCount > 0,
          `Trace count (${CONTEXT.initialTraceCount}) should be > 0`
        );
        CONTEXT.graphsRendered = DomUtils.getGraphs().length > 0;
      }
    );

    // 2. Graph Layout Verification
    await runTest(
      '2. Graph Layout Check',
      'Verify graphs are rendered with valid dimensions.',
      async () => {
        const graphs = DomUtils.getGraphs();
        if (!CONTEXT.graphsRendered && graphs.length === 0) {
          Logger.warn('Skipping Layout Check (No graphs rendered).');
          return;
        }
        assert(graphs.length > 0, 'Graphs should be visible');
        graphs.forEach((g, i) => {
          const rect = g.getBoundingClientRect();
          if (rect.width < 10 || rect.height < 10) {
            Logger.warn(`Graph #${i} seems collapsed: ${rect.width}x${rect.height}`);
          }
        });
        Logger.success(`Verified ${graphs.length} graphs are rendered.`);
      }
    );

    // 3. URL State Verification
    await runTest('3. URL State Sync', 'Verify URL parameters match selection.', async () => {
      const lastStep = CONTEXT.path[CONTEXT.path.length - 1];
      if (lastStep) {
        const encodedVal = encodeURIComponent(lastStep.value);
        Logger.info(`Waiting for URL to contain '${lastStep.value}' (Encoded: ${encodedVal})...`);

        await waitForCondition(
          () => {
            return window.location.href.includes(encodedVal);
          },
          5000,
          'URL Update'
        );

        const url = window.location.href;
        if (url.includes(encodedVal)) {
          Logger.success(`URL contains selected value '${lastStep.value}'`);
        } else {
          Logger.warn(`URL missing selected value '${lastStep.value}'. URL: ${url}`);
        }
      }
    });

    // Shared State for Tests 4-7
    let lifecycleField = null;
    let lifecycleItems = [];

    // 4. Split Mode Enable
    await runTest('4. Split Mode Enable', 'Find a safe field and enable split.', async () => {
      if (!CONTEXT.graphsRendered && DomUtils.getGraphs().length === 0) {
        Logger.warn('SKIPPING Split Test: No graphs rendered initially (Backend failure?).');
        return;
      }

      Logger.info("Searching for a field where 'Select All' yields 2-15 traces...");
      const allFields = DomUtils.getPickerFields();
      const lastPathFieldLabel = CONTEXT.path[CONTEXT.path.length - 1].label;
      const lastPathIdx = allFields.findIndex((f) => f.label === lastPathFieldLabel);
      const candidates = [];

      if (lastPathIdx > -1 && lastPathIdx < allFields.length - 1) {
        candidates.push(allFields[lastPathIdx + 1]);
      }
      for (let i = CONTEXT.path.length - 1; i > 0; i--) {
        const label = CONTEXT.path[i].label;
        const f = allFields.find((field) => field.label === label);
        if (f) candidates.push(f);
      }

      for (let i = 0; i < candidates.length; i++) {
        const f = candidates[i];
        if (!f.isConnected) continue;
        originalSelection = f.selectedItems;
        totalOptions = f.options || [];

        if (totalOptions.length < 2) continue;

        Logger.info(`Checking '${f.label}' (Options: ${totalOptions.length})...`);
        await Interaction.setSelection(f, []); // Select All
        await waitForStableState(10000);

        const t = DomUtils.getTraceCount();
        if (t >= 2 && t <= 15) {
          Logger.success(`Candidate '${f.label}' has ${t} traces. Testing Split...`);
          Logger.info(
            `Field '${f.label}' Options: [${(f.options || []).slice(0, 10).join(', ')}${
              f.options.length > 10 ? '...' : ''
            }]`
          );

          // Test Split
          await Interaction.setSplit(f, true);
          await waitForStableState(5000);

          const btn = DomUtils.getPlotButton();
          if (btn && !btn.disabled) {
            Logger.info('Plot button enabled. Clicking...');
            await Interaction.click(btn);
            await waitForStableState(5000);
          }
          await waitForGraphs(CONFIG.timeouts.medium);

          const g = DomUtils.getGraphs().length;

          if (g > 1) {
            Logger.success(`Split Verified! Field '${f.label}' produced ${g} graphs.`);
            lifecycleField = f;
            lifecycleItems = f.selectedItems;
            break; // Found and Verified
          } else {
            Logger.warn(
              `Split Failed for '${f.label}' (Traces: ${t}, Graphs: ${g}). Unsplitting...`
            );
            await Interaction.setSplit(f, false);
            await waitForStableState(3000);
            // Revert selection
            await Interaction.setSelection(f, originalSelection);
            await waitForStableState(3000);
          }
        } else if (t > 15) {
          Logger.info(`Field '${f.label}' yielded ${t} traces. Too many.`);

          // STRATEGY 1: Try selecting a random subset of THIS field
          let subsetSplitSuccess = false;
          if (totalOptions.length >= 3) {
            const optsCopyBase = [...totalOptions];

            // Try up to 10 different subsets (Increased from 3)
            for (let attempt = 1; attempt <= 10; attempt++) {
              if (subsetSplitSuccess) break;

              const subset = [];
              const optsCopy = [...optsCopyBase];
              // Randomly pick 2 or 3 items (Smaller subset = better chance of hitting <15 traces)
              const size = Math.floor(Math.random() * 2) + 2;

              for (let k = 0; k < size; k++) {
                if (optsCopy.length === 0) break;
                const idx = Math.floor(Math.random() * optsCopy.length);
                subset.push(optsCopy.splice(idx, 1)[0]);
              }

              Logger.info(
                `Strategy 1 (Attempt ${attempt}/10): Selecting subset on '${
                  f.label
                }': [${subset.join(', ')}]...`
              );
              await Interaction.setSelection(f, subset);
              await waitForStableState(3000); // Faster wait for check

              const tSubset = DomUtils.getTraceCount();
              if (tSubset >= 2 && tSubset <= 15) {
                Logger.success(
                  `Subset yielded ${tSubset} traces. Attempting Split on '${f.label}'...`
                );
                await Interaction.setSplit(f, true);
                await waitForStableState(5000);

                const btn = DomUtils.getPlotButton();
                if (btn && !btn.disabled) {
                  Logger.info('Plot button enabled. Clicking...');
                  await Interaction.click(btn);
                }

                await waitForGraphs(CONFIG.timeouts.medium);
                const g = DomUtils.getGraphs(true).length;
                Logger.info(`Graphs rendered: ${g}`);

                if (g > 1) {
                  Logger.success(`Split Successful on '${f.label}' (Subset)!`);
                  subsetSplitSuccess = true;
                  lifecycleField = f;
                  lifecycleItems = subset;
                  break;
                } else {
                  Logger.warn(`Split Failed (Subset) for '${f.label}'. Unsplitting...`);
                  await Interaction.setSplit(f, false);
                  await waitForStableState(2000);
                }
              } else {
                Logger.info(`Subset yielded ${tSubset} traces. Not viable for split.`);
              }
            }
          }

          if (subsetSplitSuccess) break; // Exit outer loop

          // STRATEGY 2: Drill Down (Original Logic)
          if (totalOptions.length > 0) {
            // Randomize drill-down to avoid sticking to bad data
            const randomIdx = Math.floor(Math.random() * totalOptions.length);
            const drillItem = totalOptions[randomIdx];
            Logger.info(
              `Strategy 2: Drilling down: Selecting '${drillItem}' (Random) in '${f.label}'...`
            );
            await Interaction.setSelection(f, [drillItem]);
            await waitForStableState(10000);

            const updatedFields = DomUtils.getPickerFields();
            const currentIdx = updatedFields.findIndex((field) => field.label === f.label);
            if (currentIdx > -1 && currentIdx < updatedFields.length - 1) {
              const nextField = updatedFields[currentIdx + 1];
              Logger.info(`Found next field: '${nextField.label}'. Adding to candidates.`);
              candidates.push(nextField);
            } else {
              await Interaction.setSelection(f, originalSelection);
              await waitForStableState(5000);
            }
          } else {
            await Interaction.setSelection(f, originalSelection);
            await waitForStableState(5000);
          }
        } else {
          Logger.info(`Field '${f.label}' yielded ${t} traces. Too few. Restoring...`);
          await Interaction.setSelection(f, originalSelection);
          await waitForStableState(5000);
        }
      }

      if (!lifecycleField) {
        throw new Error(
          'No safe field found for Split Test. Cannot proceed with Trace Manipulation tests.'
        );
      }
      // Split is already enabled by the loop verification.
    });

    // 5. Trace Manipulation (Remove)
    await runTest(
      '5. Remove Trace',
      'Remove a trace and verify graph count decreases.',
      async () => {
        if (!lifecycleField) {
          throw new Error('Skipping: No lifecycle context (Test 4 Failed).');
        }
        const initialGraphs = DomUtils.getGraphs().length;
        const initialTraces = DomUtils.getTraceCount();
        const opts = lifecycleField.options;
        const keep = opts.slice(1);

        Logger.info(`[Step 5 Start] Graphs: ${initialGraphs}, Traces: ${initialTraces}`);
        Logger.info(`Removing item: '${opts[0]}'...`);

        await Interaction.setSelection(lifecycleField, keep);
        await waitForStableState(10000);
        await waitForGraphs(CONFIG.timeouts.medium);

        const newGraphs = DomUtils.getGraphs().length;
        const newTraces = DomUtils.getTraceCount();
        Logger.info(`[Step 5 End] Graphs: ${newGraphs}, Traces: ${newTraces}`);

        if (initialGraphs > 1 && newGraphs === initialGraphs - 1) {
          Logger.success(`Graph count decreased (Exp: ${initialGraphs - 1}, Act: ${newGraphs})`);
        } else {
          Logger.warn(
            `Graph count check soft-failed (Exp: ${
              initialGraphs - 1
            }, Act: ${newGraphs}). Continuing...`
          );
        }
      }
    );

    // 6. Trace Manipulation (Add)
    await runTest(
      '6. Add Trace',
      'Add a trace back and verify graph count increases.',
      async () => {
        if (!lifecycleField) {
          throw new Error('Skipping: No lifecycle context (Test 4 Failed).');
        }
        const initialGraphs = DomUtils.getGraphs().length;
        const initialTraces = DomUtils.getTraceCount();
        const opts = lifecycleField.options;
        const toAdd = opts[0];
        const current = lifecycleField.selectedItems || [];

        Logger.info(`[Step 6 Start] Graphs: ${initialGraphs}, Traces: ${initialTraces}`);
        Logger.info(`Current Selection: [${current.join(', ')}]`);
        Logger.info(`Adding back item: '${toAdd}'...`);

        await Interaction.setSelection(lifecycleField, [toAdd, ...current]);
        await waitForStableState(10000);
        await waitForGraphs(CONFIG.timeouts.medium);

        const newGraphs = DomUtils.getGraphs().length;
        const newTraces = DomUtils.getTraceCount();
        Logger.info(`[Step 6 End] Graphs: ${newGraphs}, Traces: ${newTraces}`);

        if (newGraphs >= initialGraphs + 1) {
          Logger.success(`Graph count increased (Exp: >=${initialGraphs + 1}, Act: ${newGraphs})`);
        } else {
          Logger.warn(
            `Graph count check soft-failed (Exp: >=${
              initialGraphs + 1
            }, Act: ${newGraphs}). Continuing...`
          );
        }
      }
    );

    // 7. Sidebar Consistency & Unsplit
    await runTest(
      '7. Sidebar & Unsplit',
      'Verify sidebar matches traces, then unsplit.',
      async () => {
        if (!lifecycleField) {
          throw new Error('Skipping: No lifecycle context (Test 4 Failed).');
        }
        const traces = DomUtils.getTraceCount();
        const sidebar = DomUtils.getSidebarCount();
        Logger.info(`Traces: ${traces}, Sidebar: ${sidebar}`);
        if (traces !== sidebar) {
          Logger.warn(
            `Sidebar count mismatch (Traces: ${traces}, Sidebar: ${sidebar}). ` +
              `Known issue (b/485902011).`
          );
        }

        Logger.info('Unsplitting...');
        await Interaction.setSplit(lifecycleField, false);
        await waitForStableState(10000);
        assert(DomUtils.getGraphs().length === 1, 'Should revert to 1 graph');
        if (lifecycleField) {
          await Interaction.setSelection(lifecycleField, [lifecycleField.options[0]]);
          await waitForStableState(5000);
        }
      }
    );

    // 8. Overflow (Dynamic)
    await runTest(
      '8. Overflow Check',
      "Clear a field to trigger 'All' selection and verify 'Reduce Traces' warning.",
      async () => {
        logStats('Start');
        const target = CONTEXT.path[CONTEXT.path.length - 1];
        const fields = DomUtils.getPickerFields();
        const field = fields.find((f) => f.label === target.label);
        Logger.info(`Clearing field '${target.label}'...`);
        await Interaction.setSelection(field, []);
        await waitForStableState(10000);
        const reduceMsg = DomUtils.getReduceMessage();
        const isVisible = reduceMsg && reduceMsg.offsetParent !== null;
        const graphs = DomUtils.getGraphs().length;
        if (isVisible) {
          Logger.success('Overflow triggered. Reduce Message visible.');
          assert(graphs === 0, 'Graphs should be hidden');
        } else {
          Logger.info(`No overflow. Graphs: ${graphs}`);
        }
      }
    );

    // 9. Recovery
    await runTest(
      '9. Subset Recovery',
      'Restore the cleared field and verify graph/traces recover.',
      async () => {
        const target = CONTEXT.path[CONTEXT.path.length - 1];
        const fields = DomUtils.getPickerFields();
        const field = fields.find((f) => f.label === target.label);
        Logger.info(`Restoring '${target.value}'...`);
        await Interaction.setSelection(field, [target.value]);
        await waitForStableState(10000);
        const btn = DomUtils.getPlotButton();
        if (btn && !btn.disabled) btn.click();
        await waitForStableState(10000);
        assert(
          DomUtils.getTraceCount() === CONTEXT.initialTraceCount,
          'Trace count should recover'
        );
      }
    );

    // 10. Pagination Test (Backtracking & Iteration)
    await runTest(
      '10. Pagination Test',
      'Backtrack to high-cardinality field, split, and iterate all pages.',
      async () => {
        Logger.info('Starting Pagination Test...');

        // --- RESET UI HELPER ---
        const resetUI = async () => {
          Logger.info('Resetting UI to blank state...');
          const fields = DomUtils.getPickerFields();

          // Clear from bottom up, BUT SKIP ROOT (index 0)
          for (let i = fields.length - 1; i > 0; i--) {
            const f = fields[i];
            if (f.selectedItems && f.selectedItems.length > 0) {
              await Interaction.setSelection(f, []);
              await waitForStableState(1000); // Slower clear to avoid race
            }
          }

          // Wait for graphs to disappear
          await waitForCondition(() => DomUtils.getGraphs().length === 0, 5000, 'Graphs Removed');

          // Verify
          const graphs = DomUtils.getGraphs().length;
          if (graphs === 0) Logger.success('UI Reset Successful. No graphs.');
          else Logger.warn(`UI Reset Partial. ${graphs} graphs remaining.`);
        };

        // Exec Reset
        await resetUI();

        const path = CONTEXT.highVolumePath;
        let splitTarget = null;
        let foundHighVolume = false;

        // Restore path until we find a high cardinality field
        for (let i = path.length - 1; i >= 0; i--) {
          const step = path[i];
          let f = null;
          for (let a = 0; a < 5; a++) {
            const fs = DomUtils.getPickerFields();
            f = fs.find((field) => field.label === step.label);
            if (f) break;
            await sleep(1000);
          }
          if (!f) continue;

          const opts = f.options || [];
          Logger.info(`Field '${step.label}' has ${opts.length} options.`);

          if (opts.length > 5) {
            // High Volume Threshold
            foundHighVolume = true;
            Logger.success(`Found Split Target: '${step.label}' (${opts.length} options).`);
            splitTarget = f;

            Logger.info(`Selecting ALL on '${step.label}'...`);
            await Interaction.setSelection(f, []);
            await waitForStableState(20000);

            Logger.info(`Splitting '${step.label}'...`);
            await Interaction.setSplit(f, true);
            await waitForStableState(20000);
            break;
          } else {
            Logger.info(`'${step.label}' too small. Selecting '${step.value}'...`);
            await Interaction.setSelection(f, [step.value]);
            await waitForStableState(5000);
          }
        }

        if (!foundHighVolume) {
          Logger.warn('Could not find a High Volume field. Using last selected field.');
          // Just stick with current state
        }

        // Verify Pagination Active or Force It
        await waitForCondition(
          () => {
            const p = DomUtils.getPagination();
            return p && p.total > 0;
          },
          10000,
          'Pagination Init'
        );

        let page = DomUtils.getPagination();

        if (page && page.total > 0 && page.total <= page.pageSize) {
          Logger.warn(
            `Total items (${page.total}) <= PageSize (${page.pageSize}). ` +
              `Forcing PageSize=5 for test.`
          );
          if (page.element) {
            page.element.setAttribute('page_size', '5');
            await sleep(2000);
            // Refresh object
            page = DomUtils.getPagination();
            Logger.info(`Forced PageSize: ${page.pageSize}`);
          }
        }

        if (page && page.total > 0) {
          Logger.success(`Pagination Active: Total ${page.total}, PageSize ${page.pageSize}`);

          if (page.total > page.pageSize) {
            const totalPages = Math.ceil(page.total / page.pageSize);
            Logger.info(`Iterating through ~${totalPages} pages...`);

            let currentPage = 1;
            while (true) {
              await waitForGraphs(10000);
              const graphs = DomUtils.getGraphs().length;
              const p = DomUtils.getPagination();

              Logger.info(`Page ${currentPage}: Offset ${p.offset}, Graphs Rendered: ${graphs}`);
              assert(graphs > 0, `Page ${currentPage} should have graphs`);
              assert(
                graphs <= p.pageSize,
                `Page ${currentPage} graph count (${graphs}) <= PageSize (${p.pageSize})`
              );

              // Try to click Next
              const clicked = await DomUtils.clickNextPage();
              if (!clicked) {
                Logger.info('Next button disabled or not found. End of pagination.');
                break;
              }

              // Wait for Offset to change
              const previousOffset = p.offset;
              await waitForCondition(
                () => {
                  const newP = DomUtils.getPagination();
                  return newP && newP.offset !== previousOffset;
                },
                10000,
                'Page Offset Update'
              );

              await waitForStableState(5000);
              currentPage++;

              if (currentPage > 20) {
                Logger.warn('Safety break: Exceeded 20 pages. Stopping iteration.');
                break;
              }
            }
            Logger.success(`Successfully iterated through ${currentPage} pages.`);
          } else {
            Logger.warn(
              `Pagination failed to activate even after force ` +
                `(Total ${page.total}, PageSize ${page.pageSize}).`
            );
          }
        } else {
          const graphs = DomUtils.getGraphs().length;
          if (graphs > 1) Logger.success(`Generated ${graphs} graphs (No pagination control).`);
          else Logger.warn('Split did not generate multiple graphs.');
        }
      }
    );

    // 11. Chrome Internal Speedometer
    if (DomUtils.isChromeInternal()) {
      await runTest(
        '11. Chrome Internal Speedometer',
        'Specific scenario: speedometer3.crossbench on pixel9, split by test.',
        async () => {
          Logger.info('Chrome Internal detected. Starting specific scenario...');

          // Reset UI
          Logger.info('Resetting UI...');
          const fields = DomUtils.getPickerFields();
          for (let i = fields.length - 1; i >= 0; i--) {
            await Interaction.setSelection(fields[i], []);
            await waitForStableState(1000);
          }
          await waitForCondition(() => DomUtils.getTraceCount() === 0, 5000, 'Trace Clear');

          // 1. Benchmark
          const fBenchmark = fields.find((f) => f.label === 'benchmark');
          if (!fBenchmark) throw new Error("Field 'benchmark' not found");
          await Interaction.setSelection(fBenchmark, ['speedometer3.crossbench']);
          await waitForStableState(5000);

          // 2. Bot
          Logger.info("Waiting for 'bot' field...");
          await waitForCondition(
            () => {
              const fs = DomUtils.getPickerFields();
              return fs.some((f) => f.label === 'bot');
            },
            10000,
            "Field 'bot' appearance"
          );

          const fields2 = DomUtils.getPickerFields();
          const fBot = fields2.find((f) => f.label === 'bot');
          if (!fBot) {
            Logger.error(
              `Field 'bot' not found after wait. Available: ${fields2
                .map((f) => f.label)
                .join(', ')}. Aborting Test 11.`
            );
            return;
          }
          await Interaction.setSelection(fBot, ['android-pixel9-perf']);
          await waitForStableState(5000);

          // 3. Test (Selection + Split)
          const fields3 = DomUtils.getPickerFields();
          const fTest = fields3.find((f) => f.label === 'test');
          if (!fTest) {
            Logger.error("Field 'test' not found. Aborting Test 11.");
            return;
          }

          const tests = [
            'Charts-chartjs',
            'Charts-observable-plot',
            'Editor-CodeMirror',
            'Editor-TipTap',
            'NewsSite-Next',
            'TodoMVC-jQuery',
          ];
          Logger.info(`Selecting ${tests.length} tests...`);
          await Interaction.setSelection(fTest, tests);
          await waitForStableState(5000);

          Logger.info("Enabling Split on 'test'...");
          await Interaction.setSplit(fTest, true);
          await waitForStableState(5000);

          // 4. Plot
          const btn = DomUtils.getPlotButton();
          if (btn && !btn.disabled) {
            Logger.info('Clicking Plot...');
            await Interaction.click(btn);
          }

          // 5. Verify
          Logger.info('Waiting for loader to disappear (this can take a while)...');
          await waitForGraphs(30000);

          const graphs = DomUtils.getGraphs().length;
          Logger.info(`Graphs rendered: ${graphs}`);

          if (graphs === tests.length) {
            Logger.success(`All ${graphs} graphs rendered correctly.`);
          } else {
            Logger.warn(
              `Graph count mismatch! Expected: ${tests.length}, Actual: ${graphs}.` +
                ' (Potential Bug b/485457450: "Show all... doesn\'t show all graphs")'
            );
          }

          assert(graphs > 0, 'Should render at least one graph');
        }
      );
    } else {
      Logger.info('Skipping Test 11 (Not chrome-internal environment).');
    }

    // 12. Load All Charts Test
    if (DomUtils.isChromeInternal()) {
      await runTest(
        '12. Load All Charts',
        'Select All Tests, Split, Verify Pagination limit, then Load All.',
        async () => {
          Logger.info('Starting Load All Charts Test...');

          // Reset UI
          Logger.info('Resetting UI (Clear & Unsplit)...');
          const fields = DomUtils.getPickerFields();
          for (let i = fields.length - 1; i >= 0; i--) {
            const f = fields[i];
            const root = f.shadowRoot || f;
            const splitBox = root.querySelector('#split-by');
            if (splitBox && splitBox.checked) {
              Logger.info(`Unsplitting '${f.label}'...`);
              await Interaction.setSplit(f, false);
              await waitForStableState(3000);
            }
            if (f.selectedItems && f.selectedItems.length > 0) {
              await Interaction.setSelection(f, []);
              await waitForStableState(2000);
            }
          }
          await waitForCondition(() => DomUtils.getTraceCount() === 0, 10000, 'Trace Clear');

          // Setup: Benchmark -> Bot -> Test(All)
          const fBenchmark = fields.find((f) => f.label === 'benchmark');
          if (fBenchmark) {
            await Interaction.setSelection(fBenchmark, ['speedometer3.crossbench']);
            await waitForStableState(5000);
          }

          const fs2 = DomUtils.getPickerFields();
          const fBot = fs2.find((f) => f.label === 'bot');
          if (fBot) {
            await Interaction.setSelection(fBot, ['android-pixel9-perf']);
            await waitForStableState(5000);
          }

          const fs3 = DomUtils.getPickerFields();
          const fTest = fs3.find((f) => f.label === 'test');
          if (!fTest) {
            Logger.warn('Test 12: Field "test" not found. Aborting.');
            return;
          }

          Logger.info('Selecting ALL tests (Empty selection)...');
          await Interaction.setSelection(fTest, []); // Select All
          await waitForStableState(5000);

          Logger.info('Enabling Split on "test"...');
          await Interaction.setSplit(fTest, true);
          await waitForStableState(5000);

          const btn = DomUtils.getPlotButton();
          if (btn && !btn.disabled) {
            await Interaction.click(btn);
            await waitForGraphs(30000);
          }

          // Verify
          const app = DomUtils.getExploreApp();
          const root = app.shadowRoot || app;
          const pageSizeInput = root.querySelector('input[type="number"]');
          const pageSize = pageSizeInput ? parseInt(pageSizeInput.value, 10) : 50; // Default

          const graphs = DomUtils.getGraphs().length;
          const totalTraces = DomUtils.getTraceCount();

          Logger.info(
            `Initial State: ${graphs} Graphs, ${totalTraces} Traces. PageSize: ${pageSize}`
          );

          if (totalTraces > pageSize) {
            // Graph count should be limited
            // Note: Sometimes pageSize is 30, but graphs might be slightly off if some failed?
            // But strictly it should be <= pageSize.
            assert(graphs <= pageSize, `Graphs (${graphs}) should be <= PageSize (${pageSize})`);

            const loadAllBtn = DomUtils.getLoadAllChartsButton();
            if (loadAllBtn) {
              Logger.info('Found "Load All Charts". Clicking...');
              await Interaction.click(loadAllBtn);

              Logger.info('Waiting for all graphs (60s)...');
              // Wait for graph count to increase
              await waitForCondition(
                () => DomUtils.getGraphs().length > pageSize,
                60000,
                'Graph Count Increase'
              );
              await waitForStableState(10000);

              const newGraphs = DomUtils.getGraphs().length;
              Logger.info(`Post-Load State: ${newGraphs} Graphs.`);

              if (newGraphs === totalTraces) {
                Logger.success(`Successfully loaded all ${newGraphs} charts.`);
              } else {
                Logger.warn(
                  `Loaded ${newGraphs} charts, expected ${totalTraces}. (Potential Bug b/485457450)`
                );
              }
              assert(newGraphs >= totalTraces, 'Should load all traces');

              // Verify Bug b/485457450: "Refreshing the page does only list the subset again"
              Logger.info('Simulating page refresh via popstate (Bug b/485457450)...');
              window.dispatchEvent(new PopStateEvent('popstate', { state: window.history.state }));
              await waitForGraphs(30000);
              await waitForStableState(10000);

              const refreshedGraphs = DomUtils.getGraphs().length;
              Logger.info(`Graphs after refresh: ${refreshedGraphs}`);
              if (refreshedGraphs === newGraphs) {
                Logger.success('Graph count maintained after refresh.');
              } else {
                Logger.warn(
                  `Graph count changed after refresh! (Bug b/485457450). ` +
                    `Expected ${newGraphs}, got ${refreshedGraphs}`
                );
              }
            } else {
              Logger.warn('"Load All Charts" button NOT found (Total > PageSize).');
            }
          } else {
            Logger.info('Total traces fit on one page. Skipping Load All click.');
          }
        }
      );
    }

    // 13. Primary Checkbox Split
    await runTest('13. Primary & Split', 'Select Primary, then Split.', async () => {
      const fields = DomUtils.getPickerFields();
      let primaryField = null;

      // Find field with Primary checkbox
      for (const f of fields) {
        const root = f.shadowRoot || f;
        const box = root.querySelector('#select-primary');
        const style = getComputedStyle(box || f); // Fallback
        // Check if hidden (attribute or style)
        if (box && !box.hidden && style.display !== 'none') {
          primaryField = f;
          break;
        }
      }

      if (!primaryField) {
        Logger.warn('No visible "Primary" checkbox found on any field. Skipping.');
        return;
      }

      Logger.info(`Found field with Primary: '${primaryField.label}'`);

      // Reset
      await Interaction.setSelection(primaryField, []);
      await waitForStableState(2000);

      // Select Primary
      await Interaction.setPrimary(primaryField, true);
      await waitForStableState(5000);

      const t = DomUtils.getTraceCount();
      Logger.info(`Traces after Primary: ${t}`);

      if (t > 1) {
        // Split
        await Interaction.setSplit(primaryField, true);
        await waitForStableState(5000);

        const btn = DomUtils.getPlotButton();
        if (btn && !btn.disabled) {
          await Interaction.click(btn);
          await waitForGraphs(CONFIG.timeouts.medium);
        }

        const g = DomUtils.getGraphs().length;
        Logger.info(`Graphs after Split Primary: ${g}`);

        if (g > 1) {
          Logger.success('Primary Split Verified.');
        } else {
          Logger.warn(`Primary Split produced ${g} graphs (Expected > 1).`);
        }
      } else {
        Logger.warn(`Primary selection yielded ${t} traces. Cannot test split.`);
      }
    });

    // 14. Select All & Recovery
    await runTest(
      '14. Select All & Recovery',
      'Select All -> Split -> Remove 1 -> Add Back -> Verify All Checked',
      async () => {
        if (!lifecycleField) {
          Logger.warn('Skipping Test 14: No lifecycle field available (Test 4 Failed).');
          return;
        }

        Logger.info(`Using field '${lifecycleField.label}'`);

        // 1. Select All via Checkbox
        // First clear
        await Interaction.setSelection(lifecycleField, []);
        await waitForStableState(2000);

        await Interaction.setSelectAll(lifecycleField, true);
        await waitForStableState(5000);

        // 2. Split
        await Interaction.setSplit(lifecycleField, true);
        await waitForStableState(5000);

        // 3. Verify Graphs
        await waitForGraphs(CONFIG.timeouts.medium);
        const graphs1 = DomUtils.getGraphs().length;
        Logger.info(`Graphs (All): ${graphs1}`);

        // 4. Remove 1 item
        const opts = lifecycleField.options;
        const toRemove = opts[0];
        const keep = opts.slice(1);
        Logger.info(`Removing '${toRemove}'...`);
        await Interaction.setSelection(lifecycleField, keep);
        await waitForStableState(5000);
        await waitForGraphs(CONFIG.timeouts.medium);

        const graphs2 = DomUtils.getGraphs().length;
        Logger.info(`Graphs (All - 1): ${graphs2}`);
        if (graphs1 > 1) {
          assert(
            graphs2 === graphs1 - 1,
            `Graph count should decrease by 1 (Exp: ${graphs1 - 1}, Act: ${graphs2})`
          );
        }

        // Verify Select All is UNCHECKED
        assert(
          !Interaction.isSelectAllChecked(lifecycleField),
          'Select All should be unchecked after removing item'
        );

        // 5. Add Back
        Logger.info(`Adding back '${toRemove}'...`);
        // We select ALL options manually
        await Interaction.setSelection(lifecycleField, opts);
        await waitForStableState(5000);
        await waitForGraphs(CONFIG.timeouts.medium);

        const graphs3 = DomUtils.getGraphs().length;
        assert(graphs3 === graphs1, 'Graph count should restore');

        // 6. Verify Select All is CHECKED (Automatic)
        const isAll = Interaction.isSelectAllChecked(lifecycleField);
        Logger.info(`Select All Checkbox State: ${isAll ? 'CHECKED' : 'UNCHECKED'}`);
        assert(
          isAll,
          'Select All checkbox should be automatically checked when all items are selected'
        );
      }
    );

    const duration = ((Date.now() - START_TIME) / 1000).toFixed(1);
    Logger.success(`[DONE] Suite Completed in ${duration}s!`);
    TestHUD.update('TEST SUITE COMPLETED', `All tests finished in ${duration}s.`, 'PASS');
    TestHUD.showSummary(duration);
  } catch (e) {
    console.error(e);
    TestHUD.update('CRASH', e.message, 'FAIL');
  }
  console.log('--- END OF SCRIPT ---');
})();
