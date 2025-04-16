package main

import (
	"errors"
	"log"
	"net/http"
)

var AuthError = errors.New("Unauthorized")

func Authorize(r *http.Request) error {
	// ⚠️ Important pour récupérer les données multipart/form-data
	if err := r.ParseMultipartForm(1024); err != nil {
		log.Println("Erreur ParseMultipartForm:", err)
		return AuthError
	}

	username := r.FormValue("username")
	if username == "" {
		log.Println("Username vide")
		return AuthError
	}
	log.Println("USERNAME AUTH:", username)

	// Session
	sessionCookie, err := r.Cookie("session_token")
	if err != nil || sessionCookie.Value == "" {
		log.Println("Cookie de session manquant ou vide")
		return AuthError
	}
	log.Println("SESSION COOKIE:", sessionCookie.Value)

	// CSRF
	csrfHeader := r.Header.Get("X-CSRF-Token")
	if csrfHeader == "" {
		log.Println("CSRF token manquant dans le header")
		return AuthError
	}
	log.Println("CSRF HEADER:", csrfHeader)

	// Vérifie dans la base
	var storedSession, storedCSRF string
	err = db.QueryRow("SELECT session_token, csrf_token FROM users WHERE username = ?", username).
		Scan(&storedSession, &storedCSRF)
	if err != nil {
		log.Println("Erreur SQL lors de la récupération des tokens:", err)
		return AuthError
	}
	log.Println("STORED TOKENS:", storedSession, storedCSRF)

	if storedSession != sessionCookie.Value || storedCSRF != csrfHeader {
		log.Println("Tokens non correspondants")
		return AuthError
	}

	return nil
}
