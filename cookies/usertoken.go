package cookies

import (
	"net/http"
)

// SetUserAuthToken sets a cookie containing users auth token ("access token")
func SetUserAuthToken(w http.ResponseWriter, userAuthToken, domain string) {
	path := "/"
	httpOnly := true
	set(w, florenceCookieKey, userAuthToken, domain, path, maxAgeBrowserSession, http.SameSiteStrictMode, httpOnly)
}

// GetUserAuthToken reads access_token  cookie and returns it's value
func GetUserAuthToken(req *http.Request) (string, error) {
	return get(req, florenceCookieKey)
}
