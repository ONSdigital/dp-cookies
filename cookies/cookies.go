package cookies

import "net/http"

import "net/url"

const (
	// name of cookie used to determine the user selected cookie preferences
	cookiesPolicyCookieKey = "cookies_policy"

	// name of cookie set once a user has made a choice preference decision
	cookiesPreferencesSetCookieKey = "cookies_preferences_set"

	// name of cookie with user choosen language of website
	localeCookieKey = "lang"

	// name of cookie set by Florence to store users access token
	florenceCookieKey = "access_token"

	// name of cookie set by Florence to store currenct active collection
	collectionIDCookieKey = "collection"
)

func set(w http.ResponseWriter, name, value string) {
	encodedValue := url.QueryEscape(value)
	cookie := &http.Cookie{
		Name:     name,
		Value:    encodedValue,
		HttpOnly: false,
	}
	http.SetCookie(w, cookie)
}
