package model

import "time"

type Board struct {
	ID             int       `db:"id"`
	UserID         int       `db:"user_id"`
	AuthorUsername string    `db:"author_username"`
	Title          string    `db:"title"`
	Imagepath      string    `db:"image_path"`
	Is_Private     bool      `db:"is_private"`
	Content        string    `db:"content"`
	CreatedAt      time.Time `db:"created_at"`
}
