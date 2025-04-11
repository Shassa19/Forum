package main

import (
	"errors"
	"net/http"
)

var AuthError = errors.New("Unauthorized")

func Authorize(r *http.Request) error {
	username := r.FormValue("username")
	if username == "" {
		return AuthError
	}

	// Récupérer le cookie de session
	sessionCookie, err := r.Cookie("session_token")
	if err != nil || sessionCookie.Value == "" {
		return AuthError
	}

	// Récupérer le token CSRF envoyé dans le header
	csrfHeader := r.Header.Get("X-CSRF-Token")
	if csrfHeader == "" {
		return AuthError
	}

	// Comparer avec les valeurs enregistrées en base
	var storedSession, storedCSRF string
	err = db.QueryRow("SELECT session_token, csrf_token FROM users WHERE username = ?", username).
		Scan(&storedSession, &storedCSRF)
	if err != nil || storedSession != sessionCookie.Value || storedCSRF != csrfHeader {
		return AuthError
	}

	return nil
}
