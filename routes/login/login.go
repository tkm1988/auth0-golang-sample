package login

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"

	"app"
	"auth"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("Open /login")

	// Generate random state
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal("Random Read Error")
		return
	}
	state := base64.StdEncoding.EncodeToString(b)

	session, err := app.Store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal("Session Get Error")
		return
	}
	session.Values["state"] = state
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal("Session Save Error")
		return
	}

	authenticator, err := auth.NewAuthenticator()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal("Generate Authenticator Error")
		return
	}

	http.Redirect(w, r, authenticator.Config.AuthCodeURL(state), http.StatusTemporaryRedirect)
}
