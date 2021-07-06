package cookies

import (
	"net/http"
)

// SetIDToken sets a cookie containing users id token ("id_token")
func SetIDToken(w http.ResponseWriter, idToken, domain string) {
	path := "/"
	httpOnly := false
	set(w, idCookieKey, idToken, domain, path, maxAgeBrowserSession, http.SameSiteLaxMode, httpOnly)
}

// GetIDToken reads id_token cookie and returns it's value
func GetIDToken(req *http.Request) (string, error) {
	return get(req, idCookieKey)
}
