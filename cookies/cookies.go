package cookies

import (
	"fmt"
	"net/http"
	"net/url"
)

const (
	// cookiesPolicyCookieKey is the  name of cookie used to determine the user selected cookie preferences
	cookiesPolicyCookieKey = "cookies_policy"

	// cookiesPreferencesSetCookieKey is the name of cookie set once a user has made a choice preference decision
	cookiesPreferencesSetCookieKey = "cookies_preferences_set"

	// localeCookieKey is the name of cookie with user choosen language of website
	localeCookieKey = "lang"

	// florenceCookieKey is the name of cookie set by Florence to store users access token
	florenceCookieKey = "access_token"

	// collectionIDCookieKey is the name of cookie set by Florence to store currenct active collection
	collectionIDCookieKey = "collection"

	// maxAgeOneYear is length of time to expire a cookie in a year
	maxAgeOneYear = 31622400

	// maxAgeBrowserSession is length of time to expire cookie on browser close
	maxAgeBrowserSession = 0
)

func set(w http.ResponseWriter, name, value, domain string, maxAge int) {
	encodedValue := url.QueryEscape(value)
	cookie := &http.Cookie{
		Name:     name,
		Value:    encodedValue,
		Path:     "/",
		Domain:   domain,
		HttpOnly: false,
		Secure:   true,
		MaxAge:   maxAge,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)
}

func get(req *http.Request, name string) (string, error) {
	cookie, err := req.Cookie(name)
	if err != nil {
		return "", fmt.Errorf("could not find cookie named '%v'", name)
	}
	return cookie.Value, nil
}
