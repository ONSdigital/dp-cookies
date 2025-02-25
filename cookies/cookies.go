package cookies

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const (
	// cookiesPolicyCookieKey is the  name of cookie used to determine the user selected cookie preferences
	//
	// Deprecated: cookiesPolicyCookieKey should only be used for maintaining legacy systems. Use onsCookiesPolicyCookieKey instead.
	cookiesPolicyCookieKey = "cookies_policy"

	// onsCookiePolicyCookieKey is the name of cookie used to determine the user selected ONS cookie preferences
	onsCookiePolicyCookieKey = "ons_cookie_policy"

	// cookiesPreferencesSetCookieKey is the name of cookie set once a user has made a choice preference decision
	//
	// Deprecated: cookiesPreferencesSetCookieKey should only be used for maintaining legacy systems. Use onsCookiesPreferencesSetCookieKey instead.
	cookiesPreferencesSetCookieKey = "cookies_preferences_set"

	// onsCookiePreferencesSetCookieKey is the name of cookie set once a user has made a choice preference decision
	onsCookiePreferencesSetCookieKey = "ons_cookie_message_displayed"

	// localeCookieKey is the name of cookie with user choosen language of website
	localeCookieKey = "lang"

	// florenceCookieKey is the name of cookie set by Florence to store users access token
	florenceCookieKey = "access_token"

	// idCookieKey is the name of cookie set by Florence to store users id token used for refreshing an access_token
	idCookieKey = "id_token"

	// idCookieKey is the name of cookie set by Florence to store users refresh token used for refreshing an access_token
	refreshCookieKey = "refresh_token"

	// aBTestKey is the name of the cookie set to control a/b tests
	aBTestKey = "ab_test"

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

func set(w http.ResponseWriter, name, value, domain, path string, maxAge int, sameSite http.SameSite, httpOnly bool) {
	encodedValue := url.QueryEscape(value)

	cookie := &http.Cookie{
		Name:     name,
		Value:    encodedValue,
		Path:     path,
		Domain:   domain,
		HttpOnly: httpOnly,
		Secure:   isRunningLocalDev,
		MaxAge:   maxAge,
		SameSite: sameSite,
	}
	http.SetCookie(w, cookie)
}

// setUnencoded sets a cookie with the value not encoded
func setUnencoded(w http.ResponseWriter, name, value, domain, path string, maxAge int, sameSite http.SameSite, httpOnly bool) {
	convertedValue := strings.ReplaceAll(value, "\"", "'")

	cookie := &http.Cookie{
		Name:     name,
		Value:    convertedValue,
		Path:     path,
		Domain:   domain,
		HttpOnly: httpOnly,
		Secure:   isRunningLocalDev,
		MaxAge:   maxAge,
		SameSite: sameSite,
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

func update(req *http.Request, w http.ResponseWriter, name, value, domain, path string, maxAge int, sameSite http.SameSite, httpOnly bool) {
	exisitingCookieValue, err := get(req, name)
	if err != nil {
		exisitingCookieValue = ""
	}
	newValue := exisitingCookieValue + value
	set(w, name, newValue, domain, path, maxAge, sameSite, httpOnly)
}
