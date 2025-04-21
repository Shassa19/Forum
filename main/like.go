package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

type LikesDislikes struct {
	ID        int
	UserID    int
	PostID    int
	Value     int
	CreatedAt time.Time
}

func handlerLikeDislike(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		log.Println("❌ Erreur ParseForm:", err)
		http.Error(w, "Erreur de parsing", http.StatusBadRequest)
		return
	}

	action := r.FormValue("action")
	postIDStr := r.FormValue("post_id")
	log.Println(" Requête reçue à /api/like")
	log.Println(" action =", action)
	log.Println(" post_id =", postIDStr)

	if postIDStr == "" || action == "" {
		http.Error(w, "Paramètres manquants", http.StatusBadRequest)
		return
	}

	userID, err := getUserID(r)
	if err != nil {
		http.Error(w, "Non connecté", http.StatusUnauthorized)
		return
	}

	postID, err := strconv.Atoi(postIDStr)
	if err != nil || postID == 0 {
		http.Error(w, "Post ID invalide", http.StatusBadRequest)
		return
	}

	if action == "like" {
		// Supprime un éventuel like
		db.Exec("DELETE FROM likes WHERE user_id = ? AND post_id = ?", userID, postID)
		// Ajoute un nouveau
		_, err = db.Exec("INSERT OR IGNORE INTO likes(user_id, post_id) VALUES (?, ?)", userID, postID)
	} else if action == "dislike" {
		// Supprime un éventuel dislike
		db.Exec("DELETE FROM dislikes WHERE user_id = ? AND post_id = ?", userID, postID)
		// Ajoute un nouveau
		_, err = db.Exec("INSERT OR IGNORE INTO dislikes(user_id, post_id) VALUES (?, ?)", userID, postID)
	} else {
		http.Error(w, "Action invalide", http.StatusBadRequest)
		return
	}

	if err != nil {
		log.Println("Erreur BDD:", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
	})
}

func SaveOrToggleLike(like *LikesDislikes) error {
	existing, err := likeUser(like.UserID, like.PostID)
	if err != nil {
		return err
	}

	if existing == nil {
		_, err = db.Exec(`
			INSERT INTO likes(user_id, post_id, value)
			VALUES (?, ?, ?)`,
			like.UserID, like.PostID, like.Value)
		return err
	}

	if existing.Value == like.Value {
		_, err = db.Exec("DELETE FROM likes WHERE id = ?", existing.ID)
	} else {
		_, err = db.Exec("UPDATE likes SET value = ? WHERE id = ?", like.Value, existing.ID)
	}
	return err
}

func likeUser(userID int, postID int) (*LikesDislikes, error) {
	query := `SELECT id, user_id, post_id, value, created_at FROM likes WHERE user_id = ? AND post_id = ?`
	row := db.QueryRow(query, userID, postID)
	var l LikesDislikes
	err := row.Scan(&l.ID, &l.UserID, &l.PostID, &l.Value, &l.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &l, err
}

func getUserID(r *http.Request) (int, error) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return 0, err
	}
	var id int
	err = db.QueryRow("SELECT id FROM users WHERE session_token = ?", cookie.Value).Scan(&id)
	return id, err
}
