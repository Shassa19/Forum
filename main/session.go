package main

import (
	"errors"
	"net/http"
)

var AuthError = errors.New("Unauthorized")

func Authorize(r *http.Request) error { //On prend les infos de l'utilisateur et on les vérifie
	username := r.FormValue("username")
	if username == "" {
		return AuthError
	}

	//On prend le token via le cookie
	sessionCookie, err := r.Cookie("session_token")
	if err != nil || sessionCookie.Value == "" {
		return AuthError
	}

	// On prend le CSRF token via le header
	csrfHeader := r.Header.Get("X-CSRF-Token")
	if csrfHeader == "" {
		return AuthError
	}

	// Vérifie que les tokens correspondent à ceux en base
	var storedSession, storedCSRF string
	err = db.QueryRow("SELECT session_token, csrf_token FROM users WHERE username = ?", username).
		Scan(&storedSession, &storedCSRF)
	if err != nil || storedSession != sessionCookie.Value || storedCSRF != csrfHeader {
		return AuthError
	}

	return nil
}
