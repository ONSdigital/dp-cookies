package cookies

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"
)

// ABServices contains all services in A/B test and their expiry date
type ABServices struct {
	NewSearch *time.Time `json:"new_search,omitempty"`
	OldSearch *time.Time `json:"old_search,omitempty"`
}

// ErrABTestCookieNotFound is used when a/b test cookie isn't found
var ErrABTestCookieNotFound = errors.New("a/b test cookie not found")

// SetABTest sets a cookie containing collection ID
func SetABTest(w http.ResponseWriter, servs ABServices, domain string) error {
	b, err := json.Marshal(servs)
	if err != nil {
		return err
	}
	path := "/"
	httpOnly := false
	set(w, aBTestKey, string(b), domain, path, maxAgeOneYear, http.SameSiteLaxMode, httpOnly)
	return nil
}

// GetABTest reads ab_test cookie and returns it's value
func GetABTest(req *http.Request) (ABServices, error) {
	aBTestCookie, err := req.Cookie(aBTestKey)
	if err != nil {
		return ABServices{}, ErrABTestCookieNotFound
	}

	unescapedABTest, err := url.QueryUnescape(aBTestCookie.Value)
	if err != nil {
		return ABServices{}, err
	}

	ABTestServices := ABServices{}
	err = json.Unmarshal([]byte(unescapedABTest), &ABTestServices)
	if err != nil {
		return ABServices{}, err
	}
	return ABTestServices, nil
}

// UpdateNewSearch updates new search value in A/B test cookie
func UpdateNewSearch(req *http.Request, w http.ResponseWriter, newValue time.Time, domain string) error {
	cookie, err := GetABTest(req)
	if err != nil {
		return err
	}

	cookie.NewSearch = &newValue
	SetABTest(w, cookie, domain)
	return nil
}

// UpdateOldSearch updates old search value in A/B test cookie
func UpdateOldSearch(req *http.Request, w http.ResponseWriter, newValue time.Time, domain string) error {
	cookie, err := GetABTest(req)
	if err != nil {
		return err
	}

	cookie.OldSearch = &newValue
	SetABTest(w, cookie, domain)
	return nil
}

// UpdateSearch updates search values in A/B test cookie
func UpdateSearch(req *http.Request, w http.ResponseWriter, newSearchValue, oldSearchValue time.Time, domain string) error {
	cookie, err := GetABTest(req)
	if err != nil {
		return err
	}

	cookie.NewSearch = &newSearchValue
	cookie.OldSearch = &oldSearchValue

	SetABTest(w, cookie, domain)
	return nil
}
