package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

const (
	authServerURL = "http://localhost:8081"
)

var (
	config = oauth2.Config{
		ClientID:     "12345",
		ClientSecret: "0123456",
		Scopes:       []string{"all"},
		RedirectURL:  "http://localhost:8080/oauth2",
		Endpoint: oauth2.Endpoint{
			AuthURL:  authServerURL + "/authorize",
			TokenURL: authServerURL + "/token",
		},
	}
	globalToken *oauth2.Token // Non-concurrent security
)

func main() {
	http.HandleFunc("/", homeHandler)

	http.HandleFunc("/oauth2", authorizeHandler)
	http.HandleFunc("/refresh", refreshTokenHandler)
	http.HandleFunc("/try", tryAccessTokenHandler)
	http.HandleFunc("/pwd", passwordCredentialsGrantHandler)
	http.HandleFunc("/client", clientCredentialsGrantHandler)

	log.Println("Client is running at 8080 port.")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("home handler")
	u := config.AuthCodeURL("xyz")
	http.Redirect(w, r, u, http.StatusFound)
}

func authorizeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("authorize handler")
	r.ParseForm()
	state := r.Form.Get("state")
	if state != "xyz" {
		http.Error(w, "State invalid", http.StatusBadRequest)
		return
	}

	code := r.Form.Get("code")
	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
		return
	}

	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	globalToken = token

	e := json.NewEncoder(w)
	e.SetIndent("", "	")
	e.Encode(*token)
}

func refreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("refresh handler")
	if globalToken == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	globalToken.Expiry = time.Now()
	token, err := config.TokenSource(context.Background(), globalToken).Token()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	globalToken = token
	e := json.NewEncoder(w)
	e.SetIndent("", "	")
	e.Encode(token)
}

func tryAccessTokenHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("try handler")
	if globalToken == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	resp, err := http.Get(fmt.Sprintf("%s/test?access_token=%s", authServerURL, globalToken.AccessToken))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()
	io.Copy(w, resp.Body)
}

func passwordCredentialsGrantHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("pwd handler")
	token, err := config.PasswordCredentialsToken(context.Background(), "test", "test")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	globalToken = token
	e := json.NewEncoder(w)
	e.SetIndent("", "	")
	e.Encode(token)
}

func clientCredentialsGrantHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("client handler")
	cfg := clientcredentials.Config{
		ClientID: config.ClientID,
		ClientSecret: config.ClientSecret,
		TokenURL: config.Endpoint.TokenURL,
	}

	token, err := cfg.Token(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	e := json.NewEncoder(w)
	e.SetIndent("", "	")
	e.Encode(token)
}