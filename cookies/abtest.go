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
		log.Error(req.Context(), "error getting a/b test cookie aspect", err, log.Data{"aspectID": aspectID})
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
		log.Error(req.Context(), "error getting a/b test cookie aspect", err, log.Data{"aspectID": aspectID})
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
		log.Error(req.Context(), "error getting a/b test cookie aspect", err, log.Data{"aspectID": aspectID})
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

	return
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
		rand.Seed(time.Now().UnixNano())

		if rand.Intn(100) < percentage {
			return ABTestCookieAspect{New: now.Add(time.Hour * 24), Old: now}
		} else {
			return ABTestCookieAspect{New: now, Old: now.Add(time.Hour * 24)}
		}
	}
}
