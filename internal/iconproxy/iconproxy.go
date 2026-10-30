// Package iconproxy names where magpie downloads a remote picture from. A
// plugin's, an import link's or the library's icon sits on a host its author
// picked, and fetching it straight from there would hand that author the
// user's IP each time a page showed it. Every such fetch goes through
// usemagpie.ai's /api/icon instead (site/worker.js), which asks the host
// itself and passes on only an image; the host sees Cloudflare's address.
package iconproxy

import (
	"net/url"
	"os"
	"strings"
)

// Base is the proxy's endpoint. MAGPIE_ICON_PROXY replaces it, for tests.
const Base = "https://usemagpie.ai/api/icon"

// URL is the address to fetch the picture at u through. One on usemagpie.ai
// already is asked for as it is: the worker needn't fetch its own site.
func URL(u string) string {
	if p, err := url.Parse(u); err == nil && p.Scheme == "https" && ownSite(p.Hostname()) {
		return u
	}
	base := Base
	if v := os.Getenv("MAGPIE_ICON_PROXY"); v != "" {
		base = v
	}
	return base + "?url=" + url.QueryEscape(u)
}

// ownSite reports whether hostname is usemagpie.ai or a name under it, which
// the worker refuses to fetch (site/worker.js): it compares the parsed
// hostname, in lower case, so a subdomain, a spelled-out :443 and another
// case of the name are the site's own too — asked for here, each went to the
// worker and came back 400.
func ownSite(hostname string) bool {
	h := strings.ToLower(hostname)
	return h == "usemagpie.ai" || strings.HasSuffix(h, ".usemagpie.ai")
}
