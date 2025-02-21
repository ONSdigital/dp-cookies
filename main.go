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
		cookies.SetPolicy(w, policy, domain)  //nolint:staticcheck // To be removed in future iteration
		cookies.SetPreferenceIsSet(w, domain) //nolint:staticcheck // To be removed in future iteration
		cookies.SetLang(w, "en", domain)
		cookies.SetCollection(w, "test-collection-id-123456789", domain)
		cookies.SetUserAuthToken(w, "test-user-auth-token", domain)
	})

	http.HandleFunc("/cookies", func(w http.ResponseWriter, r *http.Request) {
		cookiesResponse := cookies.GetCookiePreferences(r)
		_, err := fmt.Fprintf(w, "%+v \n", cookiesResponse)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/lang", func(w http.ResponseWriter, r *http.Request) {
		cookiesResponse, err := cookies.GetLang(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, wErr := fmt.Fprintf(w, "%+v \n", cookiesResponse)
		if wErr != nil {
			http.Error(w, wErr.Error(), http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/collection", func(w http.ResponseWriter, r *http.Request) {
		cookiesResponse, err := cookies.GetCollection(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, wErr := fmt.Fprintf(w, "%+v \n", cookiesResponse)
		if wErr != nil {
			http.Error(w, wErr.Error(), http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/user-auth-token", func(w http.ResponseWriter, r *http.Request) {
		cookiesResponse, err := cookies.GetUserAuthToken(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, wErr := fmt.Fprintf(w, "%+v \n", cookiesResponse)
		if wErr != nil {
			http.Error(w, wErr.Error(), http.StatusInternalServerError)
			return
		}
	})

	fmt.Println("Running on port 22888")
	http.ListenAndServe(":22888", nil) //nolint:all // local dev server
}
