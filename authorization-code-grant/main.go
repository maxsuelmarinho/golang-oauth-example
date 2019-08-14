package main

import (
	"net/http"
	"os"
	"fmt"
	"encoding/json"
)

func main() {
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)

	httpClient := http.Client{}

	http.HandleFunc("/oauth/redirect", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			fmt.Fprintf(os.Stdout, "could not parse query: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		clientID := os.Getenv("CLIENT_ID")
		clientSecret := os.Getenv("CLIENT_SECRET")
		code := r.FormValue("code")
		fmt.Println("Code: " + code)
		fmt.Println("Client ID: " + clientID)
		fmt.Println("Client Secret: " + clientSecret)
		URL := fmt.Sprintf("https://github.com/login/oauth/access_token?client_id=%s&client_secret%s&code=%s", clientID, clientSecret, code)
		req, err := http.NewRequest(http.MethodPost, URL, nil)
		if err != nil {
			fmt.Fprintf(os.Stdout, "could not create HTTP request: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		req.Header.Set("accept", "application/json")
		res, err := httpClient.Do(req)
		if err != nil {
			fmt.Fprintf(os.Stdout, "could not send HTTP request: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer res.Body.Close()

		var o OAuthAccessResponse
		if err := json.NewDecoder(res.Body).Decode(&o); err != nil {
			fmt.Fprintf(os.Stdout, "could not parse JSON response: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		fmt.Println("access token: " + o.AccessToken)
		w.Header().Set("Location", "/welcome.html?access_token="+o.AccessToken)
		w.WriteHeader(http.StatusFound)
	})

	http.ListenAndServe(":8080", nil)
}

type OAuthAccessResponse struct {
	AccessToken string `json:"access_token"`
}