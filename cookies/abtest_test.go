package cookies

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetABTestCookieAspect(t *testing.T) {
	Convey("Given a http abTestCookie does not exist", t, func() {
		req := httptest.NewRequest("GET", "/", nil)

		Convey("When GetABTestCookieAspect is called with the requested aspect ID", func() {
			aspect := GetABTestCookieAspect(req, "test-aspect")

			Convey("A zero value ABTestCookieAspect is returned", func() {
				So(aspect, ShouldBeZeroValue)
			})
		})
	})

	Convey("Given a http abTestCookie exists, but without the requested aspect ID", t, func() {
		c := &http.Cookie{
			Name:  aBTestKey,
			Value: url.QueryEscape(`{"second-aspect":{"new":"2021-12-31T09:30:00","old":"2022-01-01T09:30:00"}}`),
		}
		req := httptest.NewRequest("GET", "/", nil)
		req.AddCookie(c)

		Convey("When GetABTestCookieAspect is called with the requested aspect ID", func() {
			aspect := GetABTestCookieAspect(req, "test-aspect")

			Convey("A zero value ABTestCookieAspect is returned", func() {
				So(aspect, ShouldBeZeroValue)
			})
		})
	})

	Convey("Given a http abTestCookie exists with the requested aspect ID", t, func() {
		c := &http.Cookie{
			Name:  aBTestKey,
			Value: url.QueryEscape(`{"test-aspect":{"new":"2020-06-16T17:28:45","old":"2020-06-15T17:28:45"},"second-aspect":{"new":"2021-12-31T09:30:00","old":"2022-01-01T09:30:00"}}`),
		}
		req := httptest.NewRequest("GET", "/", nil)
		req.AddCookie(c)

		Convey("When GetABTestCookieAspect is called with the requested aspect ID", func() {
			aspect := GetABTestCookieAspect(req, "test-aspect")

			Convey("The requested ABTestCookieAspect is returned", func() {
				So(aspect, ShouldResemble, ABTestCookieAspect{New: MustParseCookieTime("2020-06-15T17:28:45").Add(time.Hour * 24), Old: MustParseCookieTime("2020-06-15T17:28:45")})
			})
		})
	})
}

func TestSetABTestCookieAspect(t *testing.T) {
	Convey("Given a http abTestCookie does not exist", t, func() {
		req := httptest.NewRequest("GET", "/", nil)
		rec := httptest.NewRecorder()

		Convey("When SetABTestCookieAspect is called with a valid Aspect", func() {
			aspect := ABTestCookieAspect{New: MustParseCookieTime("2020-06-15T17:28:45").Add(time.Hour * 24), Old: MustParseCookieTime("2020-06-15T17:28:45")}
			SetABTestCookieAspect(rec, req, "test-aspect", "test-domain", aspect)

			Convey("The http abTestCookie is created and set, and contains the new aspect", func() {
				So(rec.Result().Cookies(), ShouldHaveLength, 1)
				cookieSetInResponse := rec.Result().Cookies()[0]

				So(cookieSetInResponse.Value, ShouldEqual, url.QueryEscape(`{"test-aspect":{"new":"2020-06-16T17:28:45","old":"2020-06-15T17:28:45"}}`))
				So(cookieSetInResponse.Path, ShouldEqual, "/")
				So(cookieSetInResponse.Domain, ShouldEqual, "test-domain")
				So(cookieSetInResponse.MaxAge, ShouldEqual, maxAgeOneYear)
				So(cookieSetInResponse.SameSite, ShouldEqual, http.SameSiteLaxMode)
			})
		})
	})

	Convey("Given a http abTestCookie does exist", t, func() {
		c := &http.Cookie{
			Name:  aBTestKey,
			Value: url.QueryEscape(`{"second-aspect":{"new":"2021-12-31T09:30:00","old":"2022-01-01T09:30:00"}}`),
		}
		req := httptest.NewRequest("GET", "/", nil)
		req.AddCookie(c)
		rec := httptest.NewRecorder()

		Convey("When SetABTestCookieAspect is called with a valid Aspect that does NOT currently exist in the http abTestCookie", func() {
			aspect := ABTestCookieAspect{New: MustParseCookieTime("2020-06-15T17:28:45").Add(time.Hour * 24), Old: MustParseCookieTime("2020-06-15T17:28:45")}
			SetABTestCookieAspect(rec, req, "test-aspect", "test-domain", aspect)

			Convey("The new aspect is added to the http abTestCookie", func() {
				So(rec.Result().Cookies(), ShouldHaveLength, 1)
				cookieSetInResponse := rec.Result().Cookies()[0]

				So(cookieSetInResponse.Value, ShouldContainSubstring, url.QueryEscape(`"test-aspect":{"new":"2020-06-16T17:28:45","old":"2020-06-15T17:28:45"}`))
				So(cookieSetInResponse.Value, ShouldContainSubstring, url.QueryEscape(`"second-aspect":{"new":"2021-12-31T09:30:00","old":"2022-01-01T09:30:00"}`))
				So(cookieSetInResponse.Path, ShouldEqual, "/")
				So(cookieSetInResponse.Domain, ShouldEqual, "test-domain")
				So(cookieSetInResponse.MaxAge, ShouldEqual, maxAgeOneYear)
				So(cookieSetInResponse.SameSite, ShouldEqual, http.SameSiteLaxMode)
			})
		})

		Convey("When SetABTestCookieAspect is called with a valid Aspect that does currently exist in the http abTestCookie", func() {
			aspect := ABTestCookieAspect{New: MustParseCookieTime("2022-01-12T17:28:45").Add(time.Hour * 24), Old: MustParseCookieTime("2022-01-12T17:28:45")}
			SetABTestCookieAspect(rec, req, "second-aspect", "test-domain", aspect)

			Convey("The existing aspect is overwritten in the http abTestCookie", func() {
				So(rec.Result().Cookies(), ShouldHaveLength, 1)
				cookieSetInResponse := rec.Result().Cookies()[0]

				So(cookieSetInResponse.Value, ShouldContainSubstring, url.QueryEscape(`"second-aspect":{"new":"2022-01-13T17:28:45","old":"2022-01-12T17:28:45"}`))
				So(cookieSetInResponse.Path, ShouldEqual, "/")
				So(cookieSetInResponse.Domain, ShouldEqual, "test-domain")
				So(cookieSetInResponse.MaxAge, ShouldEqual, maxAgeOneYear)
				So(cookieSetInResponse.SameSite, ShouldEqual, http.SameSiteLaxMode)
			})
		})

	})
}

func TestRemoveABTestCookieAspect(t *testing.T) {
	Convey("Given a http abTestCookie does not exist", t, func() {
		req := httptest.NewRequest("GET", "/", nil)
		rec := httptest.NewRecorder()

		Convey("Calling RemoveABTestCookieAspect is a noop", func() {
			RemoveABTestCookieAspect(rec, req, "test-aspect", "test-domain")
		})
	})

	Convey("Given a http abTestCookie exists, but without the requested aspect", t, func() {
		c := &http.Cookie{
			Name:  aBTestKey,
			Value: url.QueryEscape(`{"second-aspect":{"new":"2021-12-31T09:30:00","old":"2022-01-01T09:30:00"}}`),
		}
		req := httptest.NewRequest("GET", "/", nil)
		req.AddCookie(c)
		rec := httptest.NewRecorder()

		Convey("Calling RemoveABTestCookieAspect is a noop", func() {
			RemoveABTestCookieAspect(rec, req, "test-aspect", "test-domain")

			So(rec.Result().Cookies(), ShouldHaveLength, 1)
			So(rec.Result().Cookies()[0].Name, ShouldEqual, c.Name)
			So(rec.Result().Cookies()[0].Value, ShouldEqual, c.Value)
		})
	})

	Convey("Given a http abTestCookie exists with the requested aspect", t, func() {
		c := &http.Cookie{
			Name:  aBTestKey,
			Value: url.QueryEscape(`{"test-aspect":{"new":"2020-06-16T17:28:45","old":"2020-06-15T17:28:45"},"second-aspect":{"new":"2021-12-31T09:30:00","old":"2022-01-01T09:30:00"}}`),
		}
		req := httptest.NewRequest("GET", "/", nil)
		req.AddCookie(c)
		rec := httptest.NewRecorder()

		Convey("Calling RemoveABTestCookieAspect removes the requested aspect from the http abTestCookie", func() {
			RemoveABTestCookieAspect(rec, req, "test-aspect", "test-domain")

			So(rec.Result().Cookies(), ShouldHaveLength, 1)
			So(rec.Result().Cookies()[0].Name, ShouldEqual, c.Name)
			So(rec.Result().Cookies()[0].Value, ShouldEqual, url.QueryEscape(`{"second-aspect":{"new":"2021-12-31T09:30:00","old":"2022-01-01T09:30:00"}}`))
		})
	})
}

func TestServABTest(t *testing.T) {
	Convey("Given an old request handler and a new request handler", t, func() {
		oldHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) { _, _ = w.Write([]byte("Old Served")) })
		newHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) { _, _ = w.Write([]byte("New Served")) })
		req := httptest.NewRequest("GET", "/", nil)

		Convey("When handling a test aspect where the 'new' aspect time is in the future", func() {
			testAspect := ABTestCookieAspect{New: Now().Add(time.Hour * 24)}

			Convey("The new request handler is used to handle the request regardless of the value of the 'old' aspect time", func() {
				rec := httptest.NewRecorder()
				ServABTest(rec, req, newHandler, oldHandler, testAspect)
				b, _ := ioutil.ReadAll(rec.Result().Body)
				So(b, ShouldResemble, []byte("New Served"))

				testAspect.Old = Now()
				rec = httptest.NewRecorder()
				ServABTest(rec, req, newHandler, oldHandler, testAspect)
				b, _ = ioutil.ReadAll(rec.Result().Body)
				So(b, ShouldResemble, []byte("New Served"))

				testAspect.Old = Now().Add(time.Hour * 24)
				rec = httptest.NewRecorder()
				ServABTest(rec, req, newHandler, oldHandler, testAspect)
				b, _ = ioutil.ReadAll(rec.Result().Body)
				So(b, ShouldResemble, []byte("New Served"))
			})
		})

		Convey("When handling a test aspect where the 'old' aspect time is in the future", func() {
			testAspect := ABTestCookieAspect{Old: Now().Add(time.Hour * 24)}

			Convey("The old request handler is used to handle the request only when the new aspect time is in the past", func() {
				rec := httptest.NewRecorder()
				ServABTest(rec, req, newHandler, oldHandler, testAspect)
				b, _ := ioutil.ReadAll(rec.Result().Body)
				So(b, ShouldResemble, []byte("Old Served"))

				testAspect.New = Now()
				rec = httptest.NewRecorder()
				ServABTest(rec, req, newHandler, oldHandler, testAspect)
				b, _ = ioutil.ReadAll(rec.Result().Body)
				So(b, ShouldResemble, []byte("Old Served"))

				testAspect.New = Now().Add(time.Hour * 24)
				rec = httptest.NewRecorder()
				ServABTest(rec, req, newHandler, oldHandler, testAspect)
				b, _ = ioutil.ReadAll(rec.Result().Body)
				So(b, ShouldResemble, []byte("New Served"))
			})
		})
	})
}

func TestHandleCookieAndServ(t *testing.T) {
	Convey("Given an old request handler, new request handler and a test randomiser", t, func() {
		testOld := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) { _, _ = w.Write([]byte("Old Served")) })
		testNew := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) { _, _ = w.Write([]byte("New Served")) })
		testAspect := ABTestCookieAspect{New: Now().Add(time.Hour * 24), Old: Now()}
		testRandomiser := func() ABTestCookieAspect { return testAspect }

		Convey("When HandleCookieAndServ is called with an aspectID", func() {
			req := httptest.NewRequest("GET", "/", nil)
			rec := httptest.NewRecorder()
			HandleCookieAndServ(rec, req, testNew, testOld, "test-aspect", "test-domain", testRandomiser)

			Convey("The the new aspect is set in the http abTestCookie", func() {
				So(rec.Result().Cookies(), ShouldHaveLength, 1)
				cookieSetInResponse := rec.Result().Cookies()[0]

				testAspectMarshalled, err := json.Marshal(abTestCookie{"test-aspect": testAspect})
				if err != nil {
					t.Fatalf("error marshalling abTestCookie: %v", err)
				}
				So(cookieSetInResponse.Value, ShouldEqual, url.QueryEscape(string(testAspectMarshalled)))
			})

			Convey("And the appropriate handler serves the request", func() {
				b, _ := ioutil.ReadAll(rec.Result().Body)
				So(b, ShouldResemble, []byte("New Served"))
			})
		})
	})
}

func TestHandleABTestExit(t *testing.T) {
	Convey("Given an old request handler", t, func() {
		testOld := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) { _, _ = w.Write([]byte("Old Served")) })

		Convey("When HandleABTestExit is called with an aspectID", func() {
			req := httptest.NewRequest("GET", "/", nil)
			rec := httptest.NewRecorder()
			HandleABTestExit(rec, req, testOld, "test-aspect", "test-domain")

			Convey("A new aspect for the requested aspect ID is set in the http abTestCookie", func() {
				So(rec.Result().Cookies(), ShouldHaveLength, 1)
				cookieSetInResponse := rec.Result().Cookies()[0]
				So(cookieSetInResponse.Value, ShouldContainSubstring, url.QueryEscape(`{"test-aspect":{"new":"`))
			})

			Convey("And the old request handler serves the request", func() {
				b, _ := ioutil.ReadAll(rec.Result().Body)
				So(b, ShouldResemble, []byte("Old Served"))
			})
		})
	})
}

func Test_getABTestCookie(t *testing.T) {
	Convey("given a valid http abTestCookie exists", t, func() {
		c := &http.Cookie{
			Name:  aBTestKey,
			Value: url.QueryEscape(`{"test-aspect":{"new":"2020-06-16T17:28:45","old":"2020-06-15T17:28:45"},"second-aspect":{"new":"2021-12-31T09:30:00","old":"2022-01-01T09:30:00"}}`),
		}
		req := httptest.NewRequest("GET", "/", nil)
		req.AddCookie(c)

		Convey("when getABTestCookie is called, the expected abTestCookie is returned", func() {
			receivedCookie, err := getABTestCookie(req)
			So(err, ShouldBeNil)
			So(receivedCookie, ShouldResemble, abTestCookie{
				"test-aspect":   ABTestCookieAspect{New: MustParseCookieTime("2020-06-15T17:28:45").Add(time.Hour * 24), Old: MustParseCookieTime("2020-06-15T17:28:45")},
				"second-aspect": ABTestCookieAspect{New: MustParseCookieTime("2021-12-31T09:30:00"), Old: MustParseCookieTime("2021-12-31T09:30:00").Add(time.Hour * 24)}})
		})
	})

	Convey("given a valid abTestCookie does not exist", t, func() {
		c := &http.Cookie{Name: "some-cookie", Value: "some value"}
		req := httptest.NewRequest("GET", "/", nil)
		req.AddCookie(c)

		Convey("when getABTestCookie is called, an error is returned", func() {
			receivedCookie, err := getABTestCookie(req)
			So(err, ShouldEqual, ErrABTestCookieNotFound)
			So(receivedCookie, ShouldResemble, abTestCookie{})
		})
	})
}

func Test_setABTestCookie(t *testing.T) {
	Convey("given a valid abTestCookie", t, func() {
		c := abTestCookie{
			"test-aspect":   ABTestCookieAspect{New: MustParseCookieTime("2020-06-15T17:28:45").Add(time.Hour * 24), Old: MustParseCookieTime("2020-06-15T17:28:45")},
			"second-aspect": ABTestCookieAspect{New: MustParseCookieTime("2021-12-31T09:30:00"), Old: MustParseCookieTime("2021-12-31T09:30:00").Add(time.Hour * 24)},
		}
		rec := httptest.NewRecorder()

		Convey("when setABTestCookie is called, the corresponding http cookie is correctly set in a http Response", func() {
			err := setABTestCookie(rec, c, "test-domain")
			So(err, ShouldBeNil)
			So(rec.Result().Cookies(), ShouldHaveLength, 1)

			cookieSetInResponse := rec.Result().Cookies()[0]
			So(cookieSetInResponse.Value, ShouldContainSubstring, url.QueryEscape(`"test-aspect":{"new":"2020-06-16T17:28:45","old":"2020-06-15T17:28:45"}`))
			So(cookieSetInResponse.Value, ShouldContainSubstring, url.QueryEscape(`"second-aspect":{"new":"2021-12-31T09:30:00","old":"2022-01-01T09:30:00"}`))
			So(cookieSetInResponse.Path, ShouldEqual, "/")
			So(cookieSetInResponse.Domain, ShouldEqual, "test-domain")
			So(cookieSetInResponse.MaxAge, ShouldEqual, maxAgeOneYear)
			So(cookieSetInResponse.SameSite, ShouldEqual, http.SameSiteLaxMode)
		})
	})
}

func TestDefaultABTestRandomiserAt100(t *testing.T) {
	var iterations, n = 100, 0

	randomiser := DefaultABTestRandomiser(100)
	for i := 0; i < iterations; i++ {
		a := randomiser()
		if a.New.After(time.Now()) {
			n++
		}
	}

	if n != iterations {
		t.Errorf("a percentage of 100%% requires ALL generated aspects to favour the New value. Expected: %d Got %d", iterations, n)
	}
}
