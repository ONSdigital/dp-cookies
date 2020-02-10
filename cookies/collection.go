package cookies

import (
	"fmt"
	"net/http"
)

// SetCollection sets a cookie containing collection ID
func SetCollection(w http.ResponseWriter, value, domain string) {
	set(w, collectionIDCookieKey, value, domain, -1)
}

// GetCollection reads collection_id cookie and returns it's value
func GetCollection(req *http.Request) (string, error) {
	collectionID, err := req.Cookie(collectionIDCookieKey)
	if err != nil {
		return "", fmt.Errorf("could not find %v cookie", collectionIDCookieKey)
	}
	return collectionID.Value, nil
}
