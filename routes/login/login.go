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
		log.Fatal("Random Read Error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	state := base64.StdEncoding.EncodeToString(b)
	log.Printf("Generated random state: %s", state)

	session, err := app.Store.Get(r, "auth-session")
	if err != nil {
		log.Fatal("Session Get Error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Values["state"] = state
	err = session.Save(r, w)
	if err != nil {
		log.Fatal("Session Save Error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Print("Got Session")

	authenticator, err := auth.NewAuthenticator()
	if err != nil {
		log.Fatal("Generate Authenticator Error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Print("Got Authenticator")

	log.Printf("Redirect %s", authenticator.Config.AuthCodeURL(state))
	http.Redirect(w, r, authenticator.Config.AuthCodeURL(state), http.StatusTemporaryRedirect)
}
