package main

import (
	"database/sql"
)

func GetExistingLikeDislike(userID string, postID int) (*LikesDislikes, error) {
	query := `
        SELECT id, user_id, post_id, value
        FROM Likes_Dislikes
        WHERE user_id = ? AND post_id = ?
    `

	row := db.QueryRow(query, userID, postID)

	var like LikesDislikes
	err := row.Scan(&like.ID, &like.UserID, &like.PostID, &like.Value)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &like, nil
}
