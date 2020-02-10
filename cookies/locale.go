package cookies

import (
	"fmt"
	"net/http"
)

// SetLang sets a cookie containing locale code
func SetLang(w http.ResponseWriter, lang, domain string) {
	set(w, localeCookieKey, lang, domain, 31622400)
}

// GetLang reads lang cookie and returns it's value
func GetLang(req *http.Request) (string, error) {
	userAccessToken, err := req.Cookie(localeCookieKey)
	if err != nil {
		return "", fmt.Errorf("could not find %v cookie", localeCookieKey)
	}
	return userAccessToken.Value, nil
}
