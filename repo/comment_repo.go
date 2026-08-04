package repo

import(
	"weblog/model"
	"database/sql"

)

type CommentRepo struct{
	db * sql.DB
}

func NewCommentRepo (D *sql.DB) (*CommentRepo){
	return &CommentRepo{db:D}
}

func (r *CommentRepo) Add(c *model.Comment) error{
	query := "INSERT INTO Comments (board_id, user_id, content) VALUES ($1, $2, $3)"
	_, err := r.db.Exec(query, c.BoardID, c.UserID, c.Content)
	return err
}

func (r *CommentRepo) GetByBoardID(boardID int) ([]model.Comment, error) {
	query := `
		SELECT c.id, c.board_id, c.user_id, u.username, c.content, c.created_at
		FROM Comments c
		JOIN users u ON u.id = c.user_id
		WHERE c.board_id = $1
		ORDER BY c.created_at ASC`

	rows, err := r.db.Query(query, boardID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var Comments []model.Comment
	for rows.Next() {
		var c model.Comment
		if err := rows.Scan(&c.ID, &c.BoardID, &c.UserID, &c.Username, &c.Content, &c.CreatedAt); err != nil {
			return nil, err
		}
		Comments = append(Comments, c)
	}
	return Comments, rows.Err()
}