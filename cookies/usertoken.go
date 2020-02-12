package cookies

import (
	"net/http"
)

// SetUserAuthToken sets a cookie containing users auth token ("access token")
func SetUserAuthToken(w http.ResponseWriter, userAuthToken, domain string) {
	set(w, florenceCookieKey, userAuthToken, domain, maxAgeBrowserSession)
}

// GetUserAuthToken reads access_token  cookie and returns it's value
func GetUserAuthToken(req *http.Request) (string, error) {
	return get(req, florenceCookieKey)
}
