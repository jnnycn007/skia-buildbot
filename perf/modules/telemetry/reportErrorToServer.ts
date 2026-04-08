/**
 * @fileoverview This file defines a function to report frontend errors to the backend.
 */

import { TelemetryErrorOptions } from './types';

/**
 * Logs an error message to the backend.
 * Sends the data immediately to the `/_/fe_error_log` endpoint.
 */
export async function reportErrorToServer(errorBody: string, options: TelemetryErrorOptions = {}) {
  const errorLog = {
    message: errorBody,
    source: options.source || 'default',
    errorCode: options.errorCode || '500',
    endpoint: options.endpoint || '',
    method: options.method || '',
    url: options.url || '',
    stack: options.stack || '',
  };

  try {
    const response = await fetch('/_/fe_error_log', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(errorLog),
    });

    if (!response.ok) {
      console.error('Failed to send frontend error log. Status:', response.status);
    }
  } catch (e) {
    console.error(e, 'Failed to send frontend error log:', errorLog);
  }
}
