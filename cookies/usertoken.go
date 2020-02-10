package cookies

import (
	"fmt"
	"net/http"
)

// SetUserAuthToken sets a cookie containing users auth token ("access token")
func SetUserAuthToken(w http.ResponseWriter, userAuthToken, domain string) {
	set(w, florenceCookieKey, userAuthToken, domain, -1)
}

// GetUserAuthToken reads access_token  cookie and returns it's value
func GetUserAuthToken(req *http.Request) (string, error) {
	userAccessToken, err := req.Cookie(florenceCookieKey)
	if err != nil {
		return "", fmt.Errorf("could not find %v cookie", florenceCookieKey)
	}
	return userAccessToken.Value, nil
}
