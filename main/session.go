package main

import (
	"errors"
	"net/http"
)

var AuthError = errors.New("Unauthorized")

func Authorize(r *http.Request) error { //On prend les infos de l'utilisateur et on les vérifie
	username := r.FormValue("username")
	user, ok := users[username]
	if !ok {
		return AuthError
	}

	//On prend le token via le cookie
	st, err := r.Cookie("session_token")
	if err != nil || st.Value == "" || st.Value != user.SessionToken {
		return AuthError
	}

	// On prend le CSRF token via le header
	csrf := r.Header.Get("X-CSRF-Token") //en gros permet qu'uniquement notre site puisse envoyer des requêtes et lire les cookies
	if csrf != user.CSRFToken || csrf == "" {
		return AuthError
	}

	return nil
}
