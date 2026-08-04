package repo

import (
	"database/sql"
	"errors"
	"weblog/model"
)

type BoardRepo struct {
	db *sql.DB
}

func NewBoardRepo(n *sql.DB) *BoardRepo {
	return &BoardRepo{db: n}
}

func (r *BoardRepo) Create(n *model.Board) (int, error) {
	var id int
	query := `INSERT INTO boards (user_id, title, content, image_path, is_private)
	          VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := r.db.QueryRow(query, n.UserID, n.Title, n.Content, n.Imagepath, n.Is_Private).Scan(&id)
	return id, err
}

func (r *BoardRepo) Delete(id int, userID int) error {
	res, err := r.db.Exec("DELETE FROM boards WHERE id = $1 AND user_id = $2", id, userID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("board not found or you are not the owner")
	}
	return nil
}

func (r *BoardRepo) AddShares(boardID int, userIDs []int) error {
	for _, uid := range userIDs {
		_, err := r.db.Exec(
			`INSERT INTO board_shares (board_id, user_id) VALUES ($1, $2)
			 ON CONFLICT (board_id, user_id) DO NOTHING`,
			boardID, uid,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

const visibleWhere = `(
	b.is_private = false		
	OR b.user_id = $1
	OR EXISTS (SELECT 1 FROM board_shares bs WHERE bs.board_id = b.id AND bs.user_id = $1)
)`

func (r *BoardRepo) GetVisibleBoards(userID int) ([]model.Board, error) {
	query := `
			SELECT b.id, b.user_id, u.username, b.title, b.content, b.image_path, b.is_private, b.created_at
			FROM boards b
			JOIN users u ON u.id = b.user_id
			WHERE ` + visibleWhere + `
			ORDER BY b.created_at DESC `

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var boards []model.Board
	for rows.Next() {
		var b model.Board
		if err := rows.Scan(&b.ID, &b.UserID, &b.AuthorUsername, &b.Title, &b.Content, &b.Imagepath, &b.Is_Private, &b.CreatedAt); err != nil {
			return nil, err
		}
		boards = append(boards, b)
	}
	return boards, rows.Err()
}

func (r *BoardRepo) GetByIDVisible(id, userID int) (*model.Board, error) {
	query := `
		SELECT b.id, b.user_id, u.username, b.title, b.content, b.image_path, b.is_private, b.created_at
		FROM boards b
		JOIN users u ON u.id = b.user_id
		WHERE b.id = $2 AND ` + visibleWhere

	var b model.Board
	err := r.db.QueryRow(query, userID, id).Scan(
		&b.ID, &b.UserID, &b.AuthorUsername, &b.Title, &b.Content, &b.Imagepath, &b.Is_Private, &b.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &b, nil
}
