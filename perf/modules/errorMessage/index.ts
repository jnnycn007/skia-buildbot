import { errorMessage as elementsErrorMessage } from '../../../elements-sk/modules/errorMessage';
import { CountMetric, telemetry } from '../telemetry/telemetry';
import { TelemetryErrorOptions } from '../telemetry/types';

/**
 * Helper method to convert different error body types into a string.
 */
export const convertToErrorString = (
  errorBody: string | { message: string } | { resp: Response } | object
): string => {
  if (typeof errorBody === 'string') {
    return errorBody;
  }
  if (typeof errorBody === 'object' && errorBody !== null) {
    if ('message' in errorBody) {
      return (errorBody as { message: string }).message;
    }
    if ('resp' in errorBody && (errorBody as { resp: Response }).resp instanceof window.Response) {
      return (errorBody as { resp: Response }).resp.statusText;
    }
  }
  try {
    return JSON.stringify(errorBody);
  } catch (e) {
    return `Failed to report log message from frontend: ${(e as Error).message}`;
  }
};

/**
 * errorMessage dispatches an event with the error message in it.
 * It also optionally tracks error occurrences via a telemetry system
 * and logs the error to the server if a source is provided.
 *
 * duration default to 0, which means the toast doesn't close automatically.
 */
export const errorMessage = (
  message: string | { message: string } | { resp: Response } | object,
  duration: number = 0,
  options?: TelemetryErrorOptions
): void => {
  if (options) {
    const source = options.source || 'default';
    const errorCode =
      options.errorCode ||
      (isMessageWithResponse(message) ? message.resp.status.toString() : '500');

    // 1. Log the full high-cardinality metadata to the backend
    telemetry.reportErrorToServer(convertToErrorString(message), options);

    // 2. Increment the time-series metric counter with safe tags
    const metricName = options.countMetricSource || CountMetric.FrontendErrorReported;
    telemetry.increaseCounter(metricName, {
      source: source,
      errorCode: errorCode,
    });
  }
  // 3. Display the UI toast
  elementsErrorMessage(message, duration);
};

/**
 * Type guard to check if an unknown object contains a valid Fetch API Response.
 */
function isMessageWithResponse(msg: unknown): msg is { resp: Response } {
  return (
    typeof msg === 'object' &&
    msg !== null &&
    'resp' in msg &&
    (msg as Record<string, unknown>).resp instanceof Response
  );
}
