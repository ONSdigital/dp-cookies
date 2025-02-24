package cookies

import (
	"net/http"
)

// SetRefreshToken sets a cookie containing users refresh token ("refresh_token")
func SetRefreshToken(w http.ResponseWriter, refreshToken, domain string) {
	path := "/api/v1/tokens/self"
	httpOnly := true
	set(w, refreshCookieKey, refreshToken, domain, path, maxAgeBrowserSession, http.SameSiteStrictMode, httpOnly)
}

// GetRefreshToken reads refresh_token cookie and returns it's value
func GetRefreshToken(req *http.Request) (string, error) {
	return get(req, refreshCookieKey)
}
