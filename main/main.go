package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

// InitDB initialise la base de données et crée les tables
func InitDB() {
	db, err := sql.Open("sqlite3", "forum.db")
	if err != nil {
		log.Fatal("Erreur lors de l'ouverture de la base de données :", err)
	}
	defer db.Close()

	// Création des tables
	createTables := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL
	);
	CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		user_id INTEGER,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(user_id) REFERENCES users(id)
	);
	CREATE TABLE IF NOT EXISTS comments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		post_id INTEGER,
		user_id INTEGER,
		content TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(post_id) REFERENCES posts(id),
		FOREIGN KEY(user_id) REFERENCES users(id)
	);
	CREATE TABLE IF NOT EXISTS likes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		post_id INTEGER,
		user_id INTEGER,
		type TEXT CHECK(type IN ('like', 'dislike')),
		FOREIGN KEY(post_id) REFERENCES posts(id),
		FOREIGN KEY(user_id) REFERENCES users(id)
	);`

	_, err = db.Exec(createTables)
	if err != nil {
		log.Fatal("Erreur lors de la création des tables :", err)
	}

	fmt.Println("Base de données initialisée avec succès !")
}

// RegisterHandler gère l'inscription des utilisateurs
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	if username == "" || email == "" || password == "" {
		http.Error(w, "Tous les champs sont requis", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite3", "forum.db")
	if err != nil {
		http.Error(w, "Erreur de connexion à la base de données", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)", username, email, password)
	if err != nil {
		http.Error(w, "Erreur lors de l'inscription", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Inscription réussie !")
}

// LoginHandler gère la connexion des utilisateurs
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		http.Error(w, "Email et mot de passe requis", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite3", "forum.db")
	if err != nil {
		http.Error(w, "Erreur de connexion à la base de données", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var id int
	err = db.QueryRow("SELECT id FROM users WHERE email = ? AND password = ?", email, password).Scan(&id)
	if err != nil {
		http.Error(w, "Identifiants incorrects", http.StatusUnauthorized)
		return
	}

	fmt.Fprintf(w, "Connexion réussie !")
}

func main() {
	// Initialisation de la base de données
	InitDB()

	fmt.Println("Serveur Go lancé sur http://localhost:8080")

	// Servir les fichiers statiques (CSS, JS, images)
	fs := http.FileServer(http.Dir("./static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Routes pour l'inscription et la connexion
	http.HandleFunc("/register", RegisterHandler)
	http.HandleFunc("/login", LoginHandler)

	// Lancer le serveur
	log.Fatal(http.ListenAndServe(":8080", nil))
}
