package cookies

import (
	"net/http"
)

// SetLang sets a cookie containing locale code
func SetLang(w http.ResponseWriter, lang, domain string) {
	path := "/"
	httpOnly := false
	set(w, localeCookieKey, lang, domain, path, maxAgeOneYear, http.SameSiteLaxMode, httpOnly)
}

// GetLang reads lang cookie and returns it's value
func GetLang(req *http.Request) (string, error) {
	return get(req, localeCookieKey)
}
