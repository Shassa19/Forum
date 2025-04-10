package main

import (
	"errors"
	"fmt"
	"net/http"
)

var AuthError = errors.New("Unauthorized")

func Authorize(r *http.Request) error {
	username := r.FormValue("username")
	if username == "" {
		fmt.Println("⚠️ username manquant")
		return AuthError
	}

	// On prend le token via le cookie
	sessionCookie, err := r.Cookie("session_token")
	if err != nil || sessionCookie.Value == "" {
		fmt.Println("⚠️ cookie session_token absent ou vide")
		return AuthError
	}

	// On prend le CSRF token via le header
	csrfHeader := r.Header.Get("X-CSRF-Token")
	if csrfHeader == "" {
		fmt.Println("⚠️ header CSRF manquant")
		return AuthError
	}

	// Vérifie que les tokens correspondent à ceux en base
	var storedSession, storedCSRF string
	err = db.QueryRow("SELECT session_token, csrf_token FROM users WHERE username = ?", username).
		Scan(&storedSession, &storedCSRF)
	if err != nil {
		fmt.Println("⚠️ utilisateur non trouvé dans la BDD")
		return AuthError
	}

	if sessionCookie.Value != storedSession {
		fmt.Println("⚠️ session_token ne correspond pas")
	}
	if csrfHeader != storedCSRF {
		fmt.Println("⚠️ csrf_token ne correspond pas")
	}

	fmt.Println("→ Tokens récupérés depuis BDD pour", username)
	fmt.Println("session_cookie =", sessionCookie.Value)
	fmt.Println("csrf_header =", csrfHeader)
	fmt.Println("stored_session =", storedSession)
	fmt.Println("stored_csrf =", storedCSRF)

	if sessionCookie.Value != storedSession || csrfHeader != storedCSRF {
		return AuthError
	}

	return nil
}
