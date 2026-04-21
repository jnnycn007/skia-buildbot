import './trace-chart-sk';
import { TraceChartSk } from './trace-chart-sk';
import { expect } from 'chai';

describe('trace-chart-sk', () => {
  let element: TraceChartSk;

  beforeEach(async () => {
    element = document.createElement('trace-chart-sk') as TraceChartSk;
    document.body.appendChild(element);
    await element.updateComplete;
  });

  afterEach(() => {
    document.body.removeChild(element);
  });

  it('renders chart-tooltip-sk when hovered', async () => {
    // Simulate hover by setting the internal state
    (element as any)['_hoveredPoint'] = {
      series: { id: 'test', color: '#fff', rows: [] },
      row: { commit_number: 100, val: 10.0, createdat: 1000 },
      x: 100,
      y: 100,
    };
    await element.updateComplete;

    const tooltip = element.shadowRoot!.querySelector('.hover-tooltip');
    expect(tooltip).to.not.be.null;
  });

  it('computes subrepo rolls correctly', async () => {
    element.series = [
      {
        id: 'test',
        color: '#fff',
        rows: [
          { commit_number: 1, val: 10.0, createdat: 1000, metadata: { V8: 'v1' } },
          { commit_number: 2, val: 11.0, createdat: 2000, metadata: { V8: 'v2' } },
        ],
      },
    ];
    element.selectedSubrepo = 'V8';
    await element.updateComplete;

    const rolls = (element as any)['_subrepoRolls'];
    expect(rolls.length).to.equal(1);
    expect(rolls[0].oldVer).to.equal('v1');
    expect(rolls[0].newVer).to.equal('v2');
  });

  it('reads CSS variables for chart colors', async () => {
    const oldGetComputedStyle = window.getComputedStyle;
    let called = false;
    window.getComputedStyle = (el: Element) => {
      if (el === element) {
        called = true;
      }
      return oldGetComputedStyle(el);
    };

    try {
      (element as any)['_processedSeries'] = [
        { id: 'test', color: '#fff', rows: [{ commit_number: 1, val: 1, createdat: 1 }] },
      ];
      (element as any)['_drawBackground']();

      expect(called).to.be.true;
    } finally {
      window.getComputedStyle = oldGetComputedStyle;
    }
  });

  it('uses commit numbers by default on X axis', async () => {
    const canvas = element.shadowRoot!.querySelector('#chart-canvas') as HTMLCanvasElement;
    const ctx = canvas.getContext('2d')!;
    const oldFillText = ctx.fillText;
    const texts: string[] = [];
    ctx.fillText = function (text: string, x: number, y: number) {
      texts.push(text);
      oldFillText.call(this, text, x, y);
    };

    try {
      (element as any)['_processedSeries'] = [
        {
          id: 'test',
          color: '#fff',
          rows: [{ commit_number: 100, val: 1, createdat: 1234567890 }],
        },
      ];
      (element as any)['_drawBackground']();

      const hasDate = texts.some((t) => t.includes('-'));
      expect(hasDate).to.be.false;

      const hasCommit = texts.some((t) => t.includes('100'));
      expect(hasCommit).to.be.true;
    } finally {
      ctx.fillText = oldFillText;
    }
  });

  it('uses dates on X axis when dateMode is enabled', async () => {
    element.dateMode = true;
    await element.updateComplete;

    const canvas = element.shadowRoot!.querySelector('#chart-canvas') as HTMLCanvasElement;
    const ctx = canvas.getContext('2d')!;
    const oldFillText = ctx.fillText;
    const texts: string[] = [];
    ctx.fillText = function (text: string, x: number, y: number) {
      texts.push(text);
      oldFillText.call(this, text, x, y);
    };

    try {
      (element as any)['_processedSeries'] = [
        {
          id: 'test',
          color: '#fff',
          rows: [{ commit_number: 100, val: 1, createdat: 1234567890 }],
        },
      ];
      (element as any)['_drawBackground']();

      const hasDate = texts.some((t) => t.includes('-'));
      expect(hasDate).to.be.true;
    } finally {
      ctx.fillText = oldFillText;
    }
  });

  it('emits range-selected event on Ctrl+Drag', async () => {
    let eventDetail: any = null;
    element.addEventListener('range-selected', (e: any) => {
      eventDetail = e.detail;
    });

    // Setup mapping and dimensions so calculations don't fail
    (element as any)['_processedSeries'] = [
      { id: 'test', color: '#fff', rows: [{ commit_number: 100, val: 1, createdat: 1 }] },
    ];

    // Mock getChartBoundsAndMapping to return expected values
    const oldGetMapping = (element as any)['_getChartBoundsAndMapping'];
    (element as any)['_getChartBoundsAndMapping'] = () => ({
      minX: 0,
      maxX: 1000,
      padding: { left: 0, top: 0, right: 0, bottom: 0 },
      graphWidth: 1000,
      graphHeight: 400,
    });

    try {
      // Simulate pointer down with Ctrl
      (element as any)['_dragCtx'] = {
        isDragging: true,
        dragStartX: 100,
        dragStartY: 100,
        isCtrl: true,
        currentX: 200,
      };

      const upEvent = new PointerEvent('pointerup', { clientX: 200, clientY: 100, ctrlKey: true });
      // Mock canvas.getBoundingClientRect
      const canvas = element.shadowRoot!.querySelector('#chart-canvas') as HTMLCanvasElement;
      canvas.getBoundingClientRect = () =>
        ({ left: 0, top: 0, width: 1000, height: 400 }) as DOMRect;

      (element as any)['_handlePointerUp'](upEvent);

      expect(eventDetail).to.not.be.null;
      expect(eventDetail.startX).to.equal(100);
      expect(eventDetail.endX).to.equal(200);
    } finally {
      (element as any)['_getChartBoundsAndMapping'] = oldGetMapping;
    }
  });

  it('synthesizes stdMin and stdMax when showStd is enabled', async () => {
    element.showStd = true;
    element.series = [
      {
        id: 'test',
        color: '#fff',
        rows: [{ commit_number: 100, val: 10.0, createdat: 1000 }],
        allStats: {
          error: [{ commit_number: 100, val: 2.0, createdat: 1000 }],
        },
      },
    ];
    await element.updateComplete;

    const processed = (element as any)['_processedSeries'];
    expect(processed).to.not.be.empty;
    const s = processed[0];
    expect(s.allStats).to.not.be.undefined;
    expect(s.allStats['stdMin']).to.not.be.undefined;
    expect(s.allStats['stdMax']).to.not.be.undefined;
    expect(s.allStats['stdMin'][0].val).to.equal(8.0); // 10.0 - 2.0
    expect(s.allStats['stdMax'][0].val).to.equal(12.0); // 10.0 + 2.0
  });

  it('calculates countMaxY when showCount is enabled', async () => {
    element.showCount = true;
    (element as any)['_processedSeries'] = [
      {
        id: 'test',
        color: '#fff',
        rows: [{ commit_number: 100, val: 10.0, createdat: 1000 }],
        allStats: {
          count: [{ commit_number: 100, val: 50.0, createdat: 1000 }],
        },
      },
    ];

    const canvas = element.shadowRoot!.querySelector('#chart-canvas') as HTMLCanvasElement;
    const rect = canvas.getBoundingClientRect();
    const mapping = (element as any)['_getChartBoundsAndMapping'](rect);

    expect(mapping.countMaxY).to.equal(50.0);
  });

  it('triggers background redraw when showStd changes', async () => {
    let redrawCalled = false;
    const oldDrawBackground = (element as any)['_drawBackground'];
    (element as any)['_drawBackground'] = function () {
      redrawCalled = true;
      oldDrawBackground.apply(this);
    };

    try {
      element.showStd = true;
      await element.updateComplete;
      expect(redrawCalled).to.be.true;
    } finally {
      (element as any)['_drawBackground'] = oldDrawBackground;
    }
  });
  it('draws No Data message when minX is Infinity', async () => {
    const canvas = element.shadowRoot!.querySelector('#chart-canvas') as HTMLCanvasElement;
    const ctx = canvas.getContext('2d')!;
    const oldFillText = ctx.fillText;
    const texts: string[] = [];
    ctx.fillText = function (text: string, x: number, y: number) {
      texts.push(text);
      oldFillText.call(this, text, x, y);
    };

    try {
      (element as any)['_processedSeries'] = [{ id: 'test', color: '#fff', rows: [] }];
      (element as any)['_drawBackground']();

      const hasMessage = texts.some((t) => t.includes('No data available'));
      expect(hasMessage).to.be.true;
    } finally {
      ctx.fillText = oldFillText;
    }
  });
});
