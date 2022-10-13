package cookies

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/ONSdigital/log.go/v2/log"
)

type CookieTime struct {
	time.Time
}

const cookieTimeFormat = "2006-01-02T15:04:05"

type InvalidCookieTimeString struct {
	value, err string
}

func (ict InvalidCookieTimeString) Error() string {
	return fmt.Sprintf("invalid date string (%q): %s", ict.value, ict.err)
}

// ParseCookieTime constructs a CookieTime from the given datetime parameter. If the parameter value does not
// conform to the 'cookieTimeFormat', an error is returned
func ParseCookieTime(datetime string) (CookieTime, error) {
	if datetime == "" {
		return CookieTime{}, nil
	}
	d, err := time.Parse(cookieTimeFormat, datetime)
	if err != nil {
		return CookieTime{}, InvalidCookieTimeString{datetime, err.Error()}
	}

	return CookieTime{Time: d}, nil
}

// MustParseCookieTime is a convenience function which is only for use in tests
func MustParseCookieTime(date string) CookieTime {
	d, err := ParseCookieTime(date)
	if err != nil {
		log.Fatal(context.Background(), "MustParseCookieTime", InvalidCookieTimeString{value: date})
	}

	return d
}

func Now() CookieTime {
	return CookieTime{Time: time.Now()}
}

func (ct CookieTime) Add(d time.Duration) CookieTime {
	return CookieTime{Time: ct.Time.Add(d)}
}

func (ct CookieTime) String() string {
	return ct.Time.UTC().Format(cookieTimeFormat)
}

func (ct CookieTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", ct.String())), nil
}

func (ct *CookieTime) UnmarshalJSON(js []byte) error {
	uqt, err := strconv.Unquote(string(js))
	if err != nil {
		return err
	}
	t, e := ParseCookieTime(uqt)
	if e != nil {
		return e
	}

	*ct = t

	return nil
}
