package model

import "time"

type Comment struct {
	ID        int       `db:"id"`
	UserID    int       `db:"user_id"`
	BoardID   int       `db:"board_id"`
	Username  string    `db:"username"`
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"created_at"`
}
