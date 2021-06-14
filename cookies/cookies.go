package cookies

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
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

	// idCookieKey is the name of cookie set by Florence to store users id token used for refreshing an access_token
	idCookieKey = "id_token"

	// idCookieKey is the name of cookie set by Florence to store users refresh token used for refreshing an access_token
	refreshCookieKey = "refresh_token"

	// collectionIDCookieKey is the name of cookie set by Florence to store currenct active collection
	collectionIDCookieKey = "collection"

	// maxAgeOneYear is length of time to expire a cookie in a year
	maxAgeOneYear = 31622400

	// maxAgeBrowserSession is length of time to expire cookie on browser close
	maxAgeBrowserSession = 0
)

var isRunningLocalDev bool

func init() {
	// Set a LIBRARY_TEST environment variable to TRUE when running locally or testing.
	// Concourse test will run 'make debug' which includes setting this variable automatically.
	// Note, this is required as we don't have the means to test secure cookies.
	IsRunningLocal := os.Getenv("LIBRARY_TEST")
	var err error
	isRunningLocalDev, err = strconv.ParseBool(IsRunningLocal)
	if err != nil {
		isRunningLocalDev = false
	}
}

func set(w http.ResponseWriter, name, value, domain string, maxAge int) {
	encodedValue := url.QueryEscape(value)
	cookie := &http.Cookie{
		Name:     name,
		Value:    encodedValue,
		Path:     "/",
		Domain:   domain,
		HttpOnly: false,
		Secure:   isRunningLocalDev,
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
