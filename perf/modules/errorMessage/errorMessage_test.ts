import { errorMessage } from './index';
import { assert } from 'chai';
import sinon from 'sinon';
import { telemetry, CountMetric } from '../telemetry/telemetry';

describe('errorMessage', () => {
  let increaseCounterStub: sinon.SinonStub;
  let reportErrorToServerStub: sinon.SinonStub;

  beforeEach(() => {
    increaseCounterStub = sinon.stub(telemetry, 'increaseCounter');
    reportErrorToServerStub = sinon.stub(telemetry, 'reportErrorToServer');
  });

  afterEach(() => {
    sinon.restore();
  });

  it('dispatches error-sk event with default duration 0', (done) => {
    const message = 'test message';
    const onErrorMessage = (e: Event) => {
      const detail = (e as CustomEvent).detail;
      assert.equal(detail.message, message);
      assert.equal(detail.duration, 0);
      document.removeEventListener('error-sk', onErrorMessage);
      done();
    };
    document.addEventListener('error-sk', onErrorMessage);
    errorMessage(message);
  });

  it('dispatches error-sk event with provided duration', (done) => {
    const message = 'another message';
    const duration = 5000;
    const onErrorMessage = (e: Event) => {
      const detail = (e as CustomEvent).detail;
      assert.equal(detail.message, message);
      assert.equal(detail.duration, duration);
      document.removeEventListener('error-sk', onErrorMessage);
      done();
    };
    document.addEventListener('error-sk', onErrorMessage);
    errorMessage(message, duration);
  });

  describe('with Telemetry', () => {
    it('dispatches error-sk event with default duration 0', (done) => {
      const message = 'telemetry message';
      const onErrorMessage = (e: Event) => {
        const detail = (e as CustomEvent).detail;
        assert.equal(detail.message, message);
        assert.equal(detail.duration, 0);
        document.removeEventListener('error-sk', onErrorMessage);
        done();
      };
      document.addEventListener('error-sk', onErrorMessage);
      errorMessage(message);
    });

    it('extracts errorCode from Response object', () => {
      const resp = new Response('error', { status: 404, statusText: 'Not Found' });
      errorMessage({ resp: resp }, 0, {});

      assert.isTrue(
        increaseCounterStub.calledWith(CountMetric.FrontendErrorReported, {
          source: 'default',
          errorCode: '404',
        })
      );
    });

    it('uses provided errorCode even if Response is present', () => {
      const resp = new Response('error', { status: 404, statusText: 'Not Found' });
      errorMessage({ resp: resp }, 0, {
        errorCode: 'CUSTOM_ERROR',
      });

      assert.isTrue(
        increaseCounterStub.calledWith(CountMetric.FrontendErrorReported, {
          source: 'default',
          errorCode: 'CUSTOM_ERROR',
        })
      );
    });

    it('uses countMetricSource if provided', () => {
      errorMessage('test message', 0, {
        countMetricSource: CountMetric.DataFetchFailure,
        source: 'custom-source',
        errorCode: '404',
      });

      assert.isTrue(
        increaseCounterStub.calledWith(CountMetric.DataFetchFailure, {
          source: 'custom-source',
          errorCode: '404',
        })
      );
    });
  });

  describe('logging error to server', () => {
    it('calls reportErrorToServer with message string', () => {
      const errorBody = 'test error';
      const errorSource = 'test source';
      errorMessage(errorBody, 0, { source: errorSource });
      assert.isTrue(reportErrorToServerStub.calledWith(errorBody, { source: errorSource }));
    });

    it('calls reportErrorToServer with object containing message', () => {
      const errorObj = { message: 'object error' };
      const errorSource = 'test source';

      errorMessage(errorObj, 0, { source: errorSource });
      assert.isTrue(reportErrorToServerStub.calledWith(errorObj.message, { source: errorSource }));
    });

    it('calls reportErrorToServer with response object (statusText)', () => {
      const errorBody = {
        resp: new Response(null, { statusText: 'Not Found' }),
      };
      const errorSource = 'test source';
      errorMessage(errorBody, 0, { source: errorSource });
      assert.isTrue(reportErrorToServerStub.calledWith('Not Found', { source: errorSource }));
    });
  });
});
