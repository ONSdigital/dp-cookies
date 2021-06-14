package cookies

import (
"net/http"
)

// SetRefreshToken sets a cookie containing users refresh token ("refresh_token")
func SetRefreshToken(w http.ResponseWriter, refreshToken, domain string) {
	set(w, refreshCookieKey, refreshToken, domain, maxAgeBrowserSession)
}

// GetRefreshToken reads refresh_token cookie and returns it's value
func GetRefreshToken(req *http.Request) (string, error) {
	return get(req, refreshCookieKey)
}
