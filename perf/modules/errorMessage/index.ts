import { errorMessage as elementsErrorMessage } from '../../../elements-sk/modules/errorMessage';
import { CountMetric, telemetry } from '../telemetry/telemetry';

export interface TelemetryErrorOptions {
  countMetricSource?: CountMetric;
  source?: string;
  errorCode?: string;
}

/**
 * This is the same function as element-sk errorMessage, but also
 * track error occurrences via a telemetry system.
 * duration default to 0, which means the toast doesn't close automatically.
 * countMetricSource is to identify the metric name.
 * source indicating the origin of the error, defaulting to 'default'.
 * errorCode representing the error code, defaulting to '500'.
 */
export const errorMessageWithTelemetry = (
  message: string | { message: string } | { resp: Response } | object,
  duration: number = 0,
  options: TelemetryErrorOptions = {}
): void => {
  if (options.countMetricSource) {
    telemetry.increaseCounter(options.countMetricSource, {
      source: options.source || 'default',
      errorCode: options.errorCode || '500',
    });
  }
  elementsErrorMessage(message, duration);
};

/**
 * This is the same function as element-sk errorMessage, but defaults to a 0s
 * delay, which means the toast doesn't close automatically.
 */
export const errorMessage = (
  message: string | { message: string } | { resp: Response } | object,
  duration: number = 0
): void => {
  elementsErrorMessage(message, duration);
};
