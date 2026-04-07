/**
 * Function that returns the root domain of a sub-domain.
 *
 * I.e. it will return "skia.org" if the host is "perf.skia.org".
 *
 * For internal "corp.goog" domains, it returns the full host to preserve
 * the specific instance name.
 */
export function rootDomain(host: string = window.location.host): string {
  const ret = host.split('.').slice(-2).join('.');
  if (ret === 'corp.goog') {
    return host;
  }
  return ret;
}
