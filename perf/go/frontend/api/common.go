package api

import (
	"net/http"
	"regexp"

	"go.skia.org/infra/perf/go/config"
)

// getOverrideNonProdHost removes the specified suffixes from the host string if they are followed by .*.goog or .*.app.
// This is to ensure that requests from different non-prod environments (autopush, lts, qa, staging) are routed to the main environment.
func getOverrideNonProdHost(host string) string {
	re := regexp.MustCompile(`(-autopush|-lts|-qa|-staging)(\.corp\.goog|\.luci\.app)$`)
	return re.ReplaceAllString(host, "$2")
}

func preferLegacy(r *http.Request) bool {
	cookie, err := r.Cookie("fetch_anomalies_from_sql")
	if err == nil {
		return cookie.Value != "true"
	}
	return !config.Config.FetchAnomaliesFromSql
}
