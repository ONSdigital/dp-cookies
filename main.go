package main

import (
	"fmt"
	"net/http"

	"github.com/ONSdigital/dp-cookies/cookies"
)

var policy = cookies.Policy{
	Essential: true,
	Usage:     true,
}

func main() {
	var domain = "localhost"
	http.HandleFunc("/set-cookies", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Setting all cookies")
		cookies.SetPolicy(w, policy, domain)
		cookies.SetPreferenceIsSet(w, domain)
		cookies.SetLang(w, "en", domain)
		cookies.SetCollection(w, "test-collection-id-123456789", domain)
		cookies.SetUserAuthToken(w, "test-user-auth-token", domain)
	})

	http.HandleFunc("/cookies", func(w http.ResponseWriter, r *http.Request) {
		cookiesResponse := cookies.GetCookiePreferences(r)
		w.Write([]byte(fmt.Sprintf("%+v \n", cookiesResponse)))
	})

	http.HandleFunc("/lang", func(w http.ResponseWriter, r *http.Request) {
		cookiesResponse, err := cookies.GetLang(r)
		if err != nil {
			w.Write([]byte(err.Error()))
		}
		w.Write([]byte(fmt.Sprintf("%+v \n", cookiesResponse)))
	})

	http.HandleFunc("/collection", func(w http.ResponseWriter, r *http.Request) {
		cookiesResponse, err := cookies.GetCollection(r)
		if err != nil {
			w.Write([]byte(err.Error()))
		}
		w.Write([]byte(fmt.Sprintf("%+v \n", cookiesResponse)))
	})

	http.HandleFunc("/user-auth-token", func(w http.ResponseWriter, r *http.Request) {
		cookiesResponse, err := cookies.GetUserAuthToken(r)
		if err != nil {
			w.Write([]byte(err.Error()))
		}
		w.Write([]byte(fmt.Sprintf("%+v \n", cookiesResponse)))
	})

	fmt.Println("Running on port 22888")
	http.ListenAndServe(":22888", nil)
}
