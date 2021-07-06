package cookies

import (
	"net/http"
)

// SetCollection sets a cookie containing collection ID
func SetCollection(w http.ResponseWriter, value, domain string) {
	path := "/"
	httpOnly := false
	set(w, collectionIDCookieKey, value, domain, path, maxAgeBrowserSession, http.SameSiteLaxMode, httpOnly)
}

// GetCollection reads collection_id cookie and returns it's value
func GetCollection(req *http.Request) (string, error) {
	return get(req, collectionIDCookieKey)
}
