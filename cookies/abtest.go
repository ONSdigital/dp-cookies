package cookies

import (
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"github.com/ONSdigital/log.go/v2/log"
)

type ABTestCookieAspect struct {
	New CookieTime `json:"new,omitempty"`
	Old CookieTime `json:"old,omitempty"`
}

const (
	errGettingABTestCookieAspect = "error getting a/b test cookie aspect"
)

type abTestCookie map[string]ABTestCookieAspect

// ErrABTestCookieNotFound is used when a/b test cookie isn't found
var ErrABTestCookieNotFound = errors.New("a/b test cookie not found")

func GetABTestCookieAspect(req *http.Request, aspectID string) ABTestCookieAspect {
	aBTestCookie, err := getABTestCookie(req)
	switch {
	case errors.Is(err, ErrABTestCookieNotFound):
		log.Info(req.Context(), "a/b test cookie not found", log.Data{"aspectID": aspectID})
		return ABTestCookieAspect{}
	case err != nil:
		log.Error(req.Context(), errGettingABTestCookieAspect, err, log.Data{"aspectID": aspectID})
		return ABTestCookieAspect{}
	}

	return aBTestCookie[aspectID]
}

func SetABTestCookieAspect(w http.ResponseWriter, req *http.Request, aspectID, domain string, aspect ABTestCookieAspect) {
	cookie, err := getABTestCookie(req)
	switch {
	case errors.Is(err, ErrABTestCookieNotFound):
		cookie = make(abTestCookie)
	case err != nil:
		log.Error(req.Context(), errGettingABTestCookieAspect, err, log.Data{"aspectID": aspectID})
		return
	}
	cookie[aspectID] = aspect

	if err = setABTestCookie(w, cookie, domain); err != nil {
		log.Error(req.Context(), "error updating a/b test cookie aspect", err, log.Data{"aspectID": aspectID, "aspect": aspect})
	}
}

func RemoveABTestCookieAspect(w http.ResponseWriter, req *http.Request, aspectID, domain string) {
	cookie, err := getABTestCookie(req)
	switch {
	case errors.Is(err, ErrABTestCookieNotFound):
		return
	case err != nil:
		log.Error(req.Context(), errGettingABTestCookieAspect, err, log.Data{"aspectID": aspectID})
		return
	}

	delete(cookie, aspectID)

	if err = setABTestCookie(w, cookie, domain); err != nil {
		log.Error(req.Context(), "error removing a/b test cookie aspect", err, log.Data{"aspectID": aspectID})
	}
}

func ServABTest(w http.ResponseWriter, req *http.Request, n, o http.Handler, aspect ABTestCookieAspect) {
	now := time.Now()
	if aspect.New.After(now) {
		n.ServeHTTP(w, req)
		return
	}

	if aspect.Old.After(now) {
		o.ServeHTTP(w, req)
		return
	}
}

type Randomiser = func() ABTestCookieAspect

func HandleCookieAndServ(w http.ResponseWriter, req *http.Request, n, o http.Handler, aspectID, domain string, randomiser Randomiser) {
	aspect := randomiser()
	SetABTestCookieAspect(w, req, aspectID, domain, aspect)

	ServABTest(w, req, n, o, aspect)
}

func HandleABTestExit(w http.ResponseWriter, req *http.Request, o http.Handler, aspectID, domain string) {
	now := Now()
	aspect := ABTestCookieAspect{New: now, Old: now.Add(time.Hour * 24)}

	SetABTestCookieAspect(w, req, aspectID, domain, aspect)

	o.ServeHTTP(w, req)
}

func setABTestCookie(w http.ResponseWriter, cookie abTestCookie, domain string) error {
	b, err := json.Marshal(cookie)
	if err != nil {
		return err
	}
	path := "/"
	httpOnly := false

	set(w, aBTestKey, string(b), domain, path, maxAgeOneYear, http.SameSiteLaxMode, httpOnly)

	return nil
}

func getABTestCookie(req *http.Request) (abTestCookie, error) {
	rawABTestCookie, err := req.Cookie(aBTestKey)
	switch {
	case errors.Is(err, http.ErrNoCookie):
		return abTestCookie{}, ErrABTestCookieNotFound
	case err != nil:
		return abTestCookie{}, err
	}

	unescapedCookie, err := url.QueryUnescape(rawABTestCookie.Value)
	if err != nil {
		return abTestCookie{}, err
	}

	var cookie abTestCookie
	err = json.Unmarshal([]byte(unescapedCookie), &cookie)
	if err != nil {
		return abTestCookie{}, err
	}

	return cookie, nil
}

var DefaultABTestRandomiser = func(percentage int) Randomiser {
	return func() ABTestCookieAspect {
		now := Now()

		//nolint:gosec //does not need to be cryptographically secure
		if rand.Intn(100) < percentage {
			return ABTestCookieAspect{New: now.Add(time.Hour * 24), Old: now}
		}

		return ABTestCookieAspect{New: now, Old: now.Add(time.Hour * 24)}
	}
}

// Handler returns the relevant handler on the basis of the supplied parameters.
// It delegates to both abTestHandler and abTestPurgeHandler on the basis the abTest parameter, but it is really
// an encapsulation of the decision-making process as to what handler is used.
// Important - if AbTest is switched off it returns the new by default - this is to match router functionality.
func Handler(abTest bool, newHandler, oldHandler http.Handler, percentage int, aspectID, domain, exitNew string) http.HandlerFunc {
	if abTest {
		return abTestHandler(newHandler, oldHandler, percentage, aspectID, domain, exitNew)
	}
	return abTestPurgeHandler(newHandler, aspectID, domain)
}

// abTestHandler routes requests to either the old or new handler, for a given aspectID, according to the given percentage
// i.e. for the given percentage of calls X, X% will be routed to the new handler, and the remainder to the old handler.
// Most of the functionality is provided by the dp-cookies library, which uses a single ab_test cookie to embed all aspects
// If the aspect does not exist or has expired, it is created/renewed according to a particular randomiser - in general
// the DefaultABTestRandomiser in the library is sufficient
// A well known string - the exitNew string -  can be used as a query parameter to the call, in order to definitively chose
// the old handler
func abTestHandler(newHandler, oldHandler http.Handler, percentage int, aspectID, domain, exitNew string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		now := time.Now().UTC()

		if _, ok := req.URL.Query()[exitNew]; ok {
			HandleABTestExit(w, req, oldHandler, aspectID, domain)
			return
		}

		aspect := GetABTestCookieAspect(req, aspectID)

		if (aspect.New.IsZero() && aspect.Old.IsZero()) || (aspect.New.Before(now) && aspect.Old.Before(now)) {
			HandleCookieAndServ(w, req, newHandler, oldHandler, aspectID, domain, DefaultABTestRandomiser(percentage))
			return
		}

		ServABTest(w, req, newHandler, oldHandler, aspect)
	})
}

// abTestPurgeHandler is used to remove a given AspectID from the single ab_test cookie handled by the dp-cookies library
// It is useful when AB Testing for a particular aspect has finished, but the aspect is still embedded in client's ab_test
// cookie - this handler will remove the aspect and can be left in use for several weeks after testing has finished to 'clean'
// the underlying ab_test cookie
func abTestPurgeHandler(newHandler http.Handler, aspectID, domain string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		RemoveABTestCookieAspect(w, req, aspectID, domain)
		newHandler.ServeHTTP(w, req)
	})
}
