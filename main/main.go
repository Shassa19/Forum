package main

import (
	"fmt"
	"net/http"
	"time"
)

type Login struct {
	HashedPassword string
	SessionToken   string
	CSRFToken      string
}

var users = map[string]Login{}

func main() {
	http.HandleFunc("/register", register)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/protected", protected)
	http.ListenAndServe(":8080", nil)
}

func register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { //PostMethod pour envoyé des info et créer des nouvelles reessources
		er := http.StatusMethodNotAllowed
		http.Error(w, "Invalid method", er)
		return
	}

	// Récupération des données du formulaire
	username := r.FormValue("username")
	password := r.FormValue("password")
	if len(username) < 4 || len(password) < 6 { //longueur nécessaire
		er := http.StatusNotAcceptable
		http.Error(w, "Username/MotDePasse Invalid", er)
		return
	}

	if _, ok := users[username]; ok {
		er := http.StatusConflict
		http.Error(w, "Username déjà utilisé", er)
		return
	}

	hashedPassword, _ := hashPassword(password)
	users[username] = Login{
		HashedPassword: hashedPassword,
	}

	fmt.Fprintln(w, "Inscription d'utilisateur réussie !")
}

// Fonction de connexion avec la méthode POST en prenant l'username et le mot de passe
func login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		er := http.StatusMethodNotAllowed
		http.Error(w, "Invalid request method", er) // username inexistant -> erreur
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	user, ok := users[username]
	if !ok || !checkPasswordHash(password, user.HashedPassword) {
		er := http.StatusUnauthorized
		http.Error(w, "Utilisateur ou mot de passe incorecte", er)
		return
	}

	sessionToken := generateToken(32)
	csrfToken := generateToken(32)

	// On installe un cookie pour la session à chaque fois qu'une requete est faite
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  time.Now().Add(24 * time.Hour), //Durée de la session
		HttpOnly: true,                           //sécurise pour empecher le vol de cookie par des scripts
	})

	//On met le token CSRF dans un cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    csrfToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: false, //false car besoin d'accès coté client pour le rendre accessible au reste
	})

	//Maintenant on stock le token dans la database
	user.SessionToken = sessionToken
	user.CSRFToken = csrfToken
	users[username] = user

	fmt.Fprintln(w, "Connexion réussie !")
}

// Fonction pour protéger les données de l'utilisateur
// On vérifie si l'utilisateur est connecté et si le token CSRF est valide
func protected(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		er := http.StatusMethodNotAllowed
		http.Error(w, "Invalid request method", er)
		return
	}

	if err := Authorize(r); err != nil {
		er := http.StatusUnauthorized
		http.Error(w, "Unauthorized", er)
		return
	}

	username := r.FormValue("username")
	fmt.Fprintf(w, "CSRF validé ! Bienvenue %s !", username)
}

func logout(w http.ResponseWriter, r *http.Request) {
	if err := Authorize(r); err != nil {
		er := http.StatusUnauthorized
		http.Error(w, "Unauthorized", er)
		return
	}

	// On supprime les cookies
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: false,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: false,
	})

	// Et on supprime les tokens de la database
	username := r.FormValue("username")
	user, _ := users[username]
	user.SessionToken = ""
	user.CSRFToken = ""
	users[username] = user

	fmt.Fprintln(w, "Déconnexion réussie !")
}
