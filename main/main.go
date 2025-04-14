package main

import (
	"encoding/json"
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

func register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { //PostMethod pour envoyé des info et créer des nouvelles reessources
		er := http.StatusMethodNotAllowed
		http.Error(w, "Invalid method", er)
		return
	}

	// Récupération des données du formulaire
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	if len(username) < 4 || len(password) < 6 { //longueur nécessaire
		er := http.StatusNotAcceptable
		http.Error(w, "Username/MotDePasse Invalid", er)
		return
	}

	// Vérifier si l'utilisateur existe déjà
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = ? OR email = ?)", username, email).Scan(&exists)
	if err != nil {
		http.Error(w, "Erreur lors de la vérification", http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "Nom d'utilisateur ou email déjà utilisé", http.StatusConflict)
		return
	}

	// Hacher le mot de passe
	hashedPassword, err := hashPassword(password)
	if err != nil {
		http.Error(w, "Erreur lors du hash du mot de passe", http.StatusInternalServerError)
		return
	}

	// Insérer l'utilisateur dans la base
	_, err = db.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)", username, email, hashedPassword)
	if err != nil {
		http.Error(w, "Erreur lors de l'insertion", http.StatusInternalServerError)
		return
	}

	fmt.Println("Inscription réussie pour", username)
	http.Redirect(w, r, "/index?success=1", http.StatusSeeOther)
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

	var id int
	var hashedPassword string

	// Récupérer l'utilisateur dans la base
	err := db.QueryRow("SELECT id, password FROM users WHERE username = ?", username).Scan(&id, &hashedPassword)
	if err != nil {
		http.Error(w, "Nom d'utilisateur ou mot de passe incorrect", http.StatusUnauthorized)
		return
	}

	if !checkPasswordHash(password, hashedPassword) {
		http.Error(w, "Nom d'utilisateur ou mot de passe incorrect", http.StatusUnauthorized)
		return
	}

	sessionToken := generateToken(32)
	csrfToken := generateToken(32)

	// Mettre à jour les tokens en base
	_, err = db.Exec("UPDATE users SET session_token = ?, csrf_token = ? WHERE id = ?", sessionToken, csrfToken, id)
	if err != nil {
		http.Error(w, "Erreur lors de la mise à jour des tokens", http.StatusInternalServerError)
		return
	}

	fmt.Println("→ Tokens stockés en BDD pour", username)
	fmt.Println("session_token =", sessionToken)
	fmt.Println("csrf_token =", csrfToken)

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

	users[username] = Login{
		HashedPassword: hashedPassword,
		SessionToken:   sessionToken,
		CSRFToken:      csrfToken,
	}

	http.Redirect(w, r, "/index", http.StatusSeeOther)
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

	username := r.FormValue("username")

	// On vérifie que l'utilisateur est bien autorisé à se déconnecter
	if err := Authorize(r); err != nil {
		http.Error(w, "Non autorisé", http.StatusUnauthorized)
		return
	}

	// On vide les tokens en base
	_, err := db.Exec("UPDATE users SET session_token = '', csrf_token = '' WHERE username = ?", username)
	if err != nil {
		http.Error(w, "Erreur lors de la déconnexion", http.StatusInternalServerError)
		return
	}

	// On supprime les cookies côté clienty
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

	fmt.Fprintln(w, "Déconnexion réussie !")
}

/*-------------Fonction gestion des posts----------------*/
func createPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// Autorisation CSRF + session
	if err := Authorize(r); err != nil {
		http.Error(w, "Non autorisé", http.StatusUnauthorized)
		return
	}

	username := r.FormValue("username")
	title := r.FormValue("title")
	content := r.FormValue("content")

	if title == "" || content == "" {
		http.Error(w, "Titre ou contenu vide", http.StatusBadRequest)
		return
	}

	// Récupérer l’ID de l’utilisateur depuis son pseudo
	var userID int
	err := db.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&userID)
	if err != nil {
		http.Error(w, "Utilisateur non trouvé", http.StatusInternalServerError)
		return
	}

	// Insérer le post dans la base
	_, err = db.Exec("INSERT INTO posts (user_id, title, content) VALUES (?, ?, ?)", userID, title, content)
	if err != nil {
		http.Error(w, "Erreur lors de la création du post", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "Post créé avec succès !")
}

type Post struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Date     string `json:"created_at"`
}

// Fonction pour récupérer tous les posts sur  la page index
func getPosts(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
		SELECT posts.id, users.username, posts.title, posts.content, posts.created_at
		FROM posts
		JOIN users ON posts.user_id = users.id
		ORDER BY posts.created_at DESC
	`)
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.ID, &p.Username, &p.Title, &p.Content, &p.Date); err != nil {
			http.Error(w, "Erreur lors du scan", http.StatusInternalServerError)
			return
		}
		posts = append(posts, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

// Fonction pour récupérer un post spécifique
func getPost(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID du post manquant", http.StatusBadRequest)
		return
	}

	var p Post
	err := db.QueryRow(`
		SELECT posts.id, users.username, posts.title, posts.content, posts.created_at
		FROM posts
		JOIN users ON posts.user_id = users.id
		WHERE posts.id = ?
	`, id).Scan(&p.ID, &p.Username, &p.Title, &p.Content, &p.Date)

	if err != nil {
		http.Error(w, "Post introuvable", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func getCurrentUser(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie("session_token")
	if err != nil || sessionCookie.Value == "" {
		http.Error(w, "Non connecté", http.StatusUnauthorized)
		return
	}

	var username string
	err = db.QueryRow("SELECT username FROM users WHERE session_token = ?", sessionCookie.Value).Scan(&username)
	if err != nil {
		http.Error(w, "Utilisateur non trouvé", http.StatusUnauthorized)
		return
	}

	w.Write([]byte(username))
}

func main() {
	InitDB("../forum.db") // Connexion SQLite

	http.HandleFunc("/register", register)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/protected", protected)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../static"))))
	http.HandleFunc("/createPost", createPost)
	http.HandleFunc("/me", getCurrentUser)
	http.HandleFunc("/posts", getPosts)
	http.HandleFunc("/getPost", getPost)

	http.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../auth.html")
	})

	http.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../index.html")
	})

	http.HandleFunc("/profil", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../profil.html")
	})

	http.HandleFunc("/post", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../post.html")
	})

	http.ListenAndServe(":8080", nil)
}
