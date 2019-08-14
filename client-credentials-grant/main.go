package main

import (
	"encoding/json"
	"log"
	"fmt"
	"net/http"
	"time"
	"github.com/google/uuid"

	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/store"
	"gopkg.in/oauth2.v3/server"
	"gopkg.in/oauth2.v3/errors"
	"gopkg.in/oauth2.v3/models"
)

func main() {
	manager := manage.NewDefaultManager()
	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)
	manager.MustTokenStorage(store.NewMemoryTokenStore())
	clientStore := store.NewClientStore()
	manager.MapClientStorage(clientStore)

	srv := server.NewDefaultServer(manager)
	srv.SetAllowGetAccessRequest(true)
	srv.SetClientInfoHandler(server.ClientFormHandler)
	manager.SetRefreshTokenCfg(manage.DefaultRefreshTokenCfg)

	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("Internal Error:", err.Error())
		return
	})

	srv.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Response Error:", re.Error.Error())
	})

	http.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		srv.HandleTokenRequest(w, r)
	})

	http.HandleFunc("/credentials", func(w http.ResponseWriter, r *http.Request) {
		clientID := uuid.New().String()[:8]
		clientSecret := uuid.New().String()[:8]
		err := clientStore.Set(clientID, &models.Client{
			ID: clientID,
			Secret: clientSecret,
			Domain: "http://localhost:8081"
		})
		if err != nil {
			fmt.Println(err.Error())
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[String]string{"client_id": clientId, "client_secret": clientSecret})
	})

	http.HandleFunc("/protected", ValidateToken(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hi, I'm a protected data"))
	}, srv))

	log.Fatal(http.ListenAndServe(":8081", nil))
}
