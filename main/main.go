package main

import (
	"encoding/json"
	"fmt"
	"log"
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

// Structure Post pour le JSON
type Post struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Date     string `json:"date"`
	Likes    int    `json:"likes"`
	Dislikes int    `json:"dislikes"`
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
		var rawDate time.Time

		err := rows.Scan(&p.ID, &p.Username, &p.Title, &p.Content, &rawDate)
		if err != nil {
			http.Error(w, "Erreur lors du scan", http.StatusInternalServerError)
			return
		}

		p.Date = rawDate.Format(time.RFC3339) // Renvoie : 2024-04-15T13:45:00Z
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
	var rawDate time.Time

	// Récupération du post avec l'auteur
	err := db.QueryRow(`
		SELECT posts.id, users.username, posts.title, posts.content, posts.created_at
		FROM posts
		JOIN users ON posts.user_id = users.id
		WHERE posts.id = ?
	`, id).Scan(&p.ID, &p.Username, &p.Title, &p.Content, &rawDate)
	if err != nil {
		http.Error(w, "Post introuvable", http.StatusNotFound)
		return
	}
	p.Date = rawDate.Format(time.RFC3339)

	// Récupération des likes
	err = db.QueryRow("SELECT COUNT(*) FROM likes WHERE post_id = ?", id).Scan(&p.Likes)
	if err != nil {
		p.Likes = 0
	}

	// Récupération des dislikes
	err = db.QueryRow("SELECT COUNT(*) FROM dislikes WHERE post_id = ?", id).Scan(&p.Dislikes)
	if err != nil {
		p.Dislikes = 0
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

func getUserPosts(w http.ResponseWriter, r *http.Request) {
	// Récupération du token de session via le cookie
	sessionCookie, err := r.Cookie("session_token")
	if err != nil || sessionCookie.Value == "" {
		http.Error(w, "Non connecté", http.StatusUnauthorized)
		return
	}

	// Récupération du username associé au token
	var username string
	err = db.QueryRow("SELECT username FROM users WHERE session_token = ?", sessionCookie.Value).Scan(&username)
	if err != nil {
		http.Error(w, "Utilisateur non trouvé", http.StatusUnauthorized)
		return
	}

	// Requête des posts de l'utilisateur
	rows, err := db.Query(`
		SELECT posts.id, posts.title
		FROM posts
		JOIN users ON posts.user_id = users.id
		WHERE users.username = ?
		ORDER BY posts.created_at DESC
	`, username)
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type PostData struct {
		ID    int    `json:"id"`
		Title string `json:"title"`
	}

	var posts []PostData
	for rows.Next() {
		var p PostData
		if err := rows.Scan(&p.ID, &p.Title); err != nil {
			http.Error(w, "Erreur lors du scan", http.StatusInternalServerError)
			return
		}
		posts = append(posts, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func updateAvatar(w http.ResponseWriter, r *http.Request) {
	if err := Authorize(r); err != nil {
		http.Error(w, "Non autorisé", http.StatusUnauthorized)
		return
	}

	sessionCookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "Session manquante", http.StatusUnauthorized)
		return
	}

	var username string
	err = db.QueryRow("SELECT username FROM users WHERE session_token = ?", sessionCookie.Value).Scan(&username)
	if err != nil {
		http.Error(w, "Utilisateur introuvable", http.StatusUnauthorized)
		return
	}

	avatar := r.FormValue("avatar")
	if avatar == "" {
		http.Error(w, "Aucun avatar fourni", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("UPDATE users SET avatar = ? WHERE username = ?", avatar, username)
	if err != nil {
		http.Error(w, "Erreur BDD", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Avatar mis à jour"))
}

func getUserAvatar(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie("session_token")
	if err != nil || sessionCookie.Value == "" {
		http.Error(w, "Non connecté", http.StatusUnauthorized)
		return
	}

	var username, avatar string
	err = db.QueryRow("SELECT username, avatar FROM users WHERE session_token = ?", sessionCookie.Value).Scan(&username, &avatar)
	if err != nil {
		http.Error(w, "Utilisateur introuvable", http.StatusNotFound)
		return
	}

	user := map[string]string{
		"username": username,
		"avatar":   avatar,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// Fonction pour gérer les commentaires sur les posts

type Comment struct {
	ID      int    `json:"id,omitempty"`
	PostID  int    `json:"post_id,omitempty"`
	User    string `json:"username"`
	Content string `json:"content"`
	Date    string `json:"date,omitempty"`
}

// Handler pour ajouter un commentaire
func addComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	if err := Authorize(r); err != nil {
		log.Println("Erreur CSRF:", err)
		http.Error(w, "Non autorisé", http.StatusUnauthorized)
		return
	}

	postID := r.FormValue("post_id")
	content := r.FormValue("content")

	// Récupérer l'utilisateur connecté via le token
	session, err := r.Cookie("session_token")
	if err != nil {

		http.Error(w, "Session manquante", http.StatusUnauthorized)
		return
	}

	var userID int
	err = db.QueryRow("SELECT id FROM users WHERE session_token = ?", session.Value).Scan(&userID)
	if err != nil {

		http.Error(w, "Utilisateur introuvable", http.StatusUnauthorized)
		return
	}

	_, err = db.Exec("INSERT INTO comments (post_id, user_id, content) VALUES (?, ?, ?)", postID, userID, content)
	if err != nil {

		http.Error(w, "Erreur lors de l’ajout du commentaire", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Handler pour récupérer les commentaires d'un post
func getComments(w http.ResponseWriter, r *http.Request) {
	postID := r.URL.Query().Get("id")
	if postID == "" {
		http.Error(w, "ID du post manquant", http.StatusBadRequest)
		return
	}

	rows, err := db.Query(`
		SELECT users.username, comments.content
		FROM comments
		JOIN users ON comments.user_id = users.id
		WHERE comments.post_id = ?
		ORDER BY comments.created_at ASC
	`, postID)
	if err != nil {
		http.Error(w, "Erreur lors de la récupération", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var comments []map[string]string
	for rows.Next() {
		var username, content string
		if err := rows.Scan(&username, &content); err != nil {
			http.Error(w, "Erreur lors du scan", http.StatusInternalServerError)
			return
		}
		comments = append(comments, map[string]string{
			"username": username,
			"content":  content,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}

// Handler pour récupérer les 3 derniers commentaires sur SES posts
func getLastComments(w http.ResponseWriter, r *http.Request) {
	// Vérifie session + récupère username
	sessionCookie, err := r.Cookie("session_token")
	if err != nil || sessionCookie.Value == "" {
		http.Error(w, "Non connecté", http.StatusUnauthorized)
		return
	}

	var username string
	err = db.QueryRow("SELECT username FROM users WHERE session_token = ?", sessionCookie.Value).Scan(&username)
	if err != nil {
		http.Error(w, "Utilisateur introuvable", http.StatusUnauthorized)
		return
	}

	// Récupère les 3 derniers commentaires sur SES posts
	rows, err := db.Query(`
		SELECT comments.content, comments.created_at, u.username, comments.post_id
		FROM comments
		JOIN posts ON comments.post_id = posts.id
		JOIN users u ON comments.user_id = u.id
		WHERE posts.user_id = (SELECT id FROM users WHERE username = ?)
		ORDER BY comments.created_at DESC
		LIMIT 3
	`, username)
	if err != nil {
		http.Error(w, "Erreur SQL", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var replies []map[string]string

	for rows.Next() {
		var content, date, author string
		var postID int

		if err := rows.Scan(&content, &date, &author, &postID); err != nil {
			http.Error(w, "Erreur lecture", http.StatusInternalServerError)
			return
		}

		replies = append(replies, map[string]string{
			"author":  author,
			"content": content,
			"date":    date,
			"post_id": fmt.Sprint(postID), // 🔁 Converti int → string
		})
	}

	// 🔹 Important : header doit être défini **avant** d'écrire quoi que ce soit
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(replies)
}

func updateProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	if err := Authorize(r); err != nil {
		http.Error(w, "Non autorisé", http.StatusUnauthorized)
		return
	}

	if err := r.ParseMultipartForm(1024); err != nil {
		http.Error(w, "Erreur de parsing", http.StatusBadRequest)
		return
	}

	oldUsername := r.FormValue("username")
	newUsername := r.FormValue("new_username")
	newPassword := r.FormValue("new_password")

	// Vérifie si l'utilisateur existe bien
	var userID int
	err := db.QueryRow("SELECT id FROM users WHERE username = ?", oldUsername).Scan(&userID)
	if err != nil {
		http.Error(w, "Utilisateur introuvable", http.StatusUnauthorized)
		return
	}

	if newUsername == "" && newPassword == "" {
		http.Error(w, "Aucune donnée à modifier", http.StatusBadRequest)
		return
	}

	// Mise à jour selon les champs remplis
	if newUsername != "" && newPassword != "" {
		hashed, _ := hashPassword(newPassword)
		_, err = db.Exec("UPDATE users SET username = ?, password = ? WHERE id = ?", newUsername, hashed, userID)
	} else if newUsername != "" {
		_, err = db.Exec("UPDATE users SET username = ? WHERE id = ?", newUsername, userID)
	} else if newPassword != "" {
		hashed, _ := hashPassword(newPassword)
		_, err = db.Exec("UPDATE users SET password = ? WHERE id = ?", hashed, userID)
	}

	if err != nil {
		http.Error(w, "Erreur lors de la mise à jour", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Profil mis à jour avec succès !"))
}

func main() {
	InitDB("../forum.db") // Connexion SQLite

	http.HandleFunc("/register", register)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/protected", protected)

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("../assets"))))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../static"))))
	http.HandleFunc("/createPost", createPost)
	http.HandleFunc("/me", getCurrentUser)
	http.HandleFunc("/posts", getPosts)
	http.HandleFunc("/getPost", getPost)
	http.HandleFunc("/userPosts", getUserPosts)
	http.HandleFunc("/update-avatar", updateAvatar)
	http.HandleFunc("/user-info", getUserAvatar)
	http.HandleFunc("/add-comment", addComment)
	http.HandleFunc("/comments", getComments)
	http.HandleFunc("/last-comments", getLastComments)
	http.HandleFunc("/update-profile", updateProfile)
	http.HandleFunc("/api/like", handlerLikeDislike)

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
