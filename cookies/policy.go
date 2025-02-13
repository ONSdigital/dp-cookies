package cookies

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// PreferencesResponse is a combination of cookie policy and whether they have be set by user
type PreferencesResponse struct {
	IsPreferenceSet bool
	Policy          Policy
}

// Policy is cookie policy setting choosen by a user
type Policy struct {
	Essential bool `json:"essential"`
	Usage     bool `json:"usage"`
}

var defaultPolicy = Policy{
	Essential: true,
	Usage:     false,
}

// GetCookiePreferences returns a struct with all cookie preferences
func GetCookiePreferences(req *http.Request) PreferencesResponse {
	isPreferenceSet := getPreferencesIsSet(req)
	cookiePolicy := getPolicy(req)
	return PreferencesResponse{
		IsPreferenceSet: isPreferenceSet,
		Policy:          cookiePolicy,
	}
}

// SetPreferenceIsSet sets a cookie to record a user has set cookie preferences
func SetPreferenceIsSet(w http.ResponseWriter, domain string) {
	path := "/"
	httpOnly := false
	set(w, cookiesPreferencesSetCookieKey, "true", domain, path, maxAgeOneYear, http.SameSiteLaxMode, httpOnly)
}

func getPreferencesIsSet(req *http.Request) bool {
	cookiesPreferencesSetCookie, err := req.Cookie(cookiesPreferencesSetCookieKey)
	if err != nil {
		return false
	}

	cookieIsPreferenceSet, err := strconv.ParseBool(cookiesPreferencesSetCookie.Value)
	if err != nil {
		return false
	}

	return cookieIsPreferenceSet
}

// SetPolicy sets a cookie with the users preferences, or sets default preferences on error
func SetPolicy(w http.ResponseWriter, policy Policy, domain string) {
	b, err := json.Marshal(policy)
	if err != nil {
		b, err = json.Marshal(defaultPolicy)
	}
	path := "/"
	httpOnly := false
	set(w, cookiesPolicyCookieKey, string(b), domain, path, maxAgeOneYear, http.SameSiteLaxMode, httpOnly)
}

// SetONSPolicy sets the ONS cookie with the users preferences, or sets default preferences on error
func SetONSPolicy(w http.ResponseWriter, policy ONSPolicy, domain string) {
	b, err := json.Marshal(policy)
	if err != nil {
		b, _ = json.Marshal(defaultONSPolicy)
	}
	path := "/"
	httpOnly := false
	setCookieWithUnencodedValue(w, onsCookiePolicyCookieKey, string(b), domain, path, maxAgeOneYear, http.SameSiteLaxMode, httpOnly)
}

func getPolicy(req *http.Request) Policy {
	cookiePolicyCookie, err := req.Cookie(cookiesPolicyCookieKey)
	if err != nil {
		return defaultPolicy
	}

	unescapedPolicy, err := url.QueryUnescape(cookiePolicyCookie.Value)
	if err != nil {
		return defaultPolicy
	}

	cookiePolicy := Policy{}
	json.Unmarshal([]byte(unescapedPolicy), &cookiePolicy)
	return cookiePolicy
}

func getONSPolicy(req *http.Request) ONSPolicy {
	cookiePolicyCookie, err := req.Cookie(onsCookiePolicyCookieKey)
	if err != nil {
		return defaultONSPolicy
	}

	unescapedPolicy, err := url.QueryUnescape(cookiePolicyCookie.Value)
	if err != nil {
		return defaultONSPolicy
	}

	// Replace single quotes with double quotes to make it valid JSON
	validJSONPolicy := strings.ReplaceAll(unescapedPolicy, "'", "\"")

	cookiePolicy := ONSPolicy{}
	json.Unmarshal([]byte(validJSONPolicy), &cookiePolicy)
	return cookiePolicy
}
