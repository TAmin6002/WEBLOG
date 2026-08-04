package service

import (
	"errors"
	"fmt"
	"strings"

	"weblog/model"
	"weblog/repo"
)

type BoardService struct {
	boardRepo *repo.BoardRepo
	userRepo  *repo.UserRepo
}

func NewBoardService(b *repo.BoardRepo, u *repo.UserRepo) *BoardService {
	return &BoardService{boardRepo: b, userRepo: u}
}

func (r *BoardService) CreatBoard(userID int, title, content, image_path string, is_private bool, sharedusernamesraw string) (int, error) {

	title = strings.TrimSpace(title)
	content = strings.TrimSpace(content)

	if title == "" || content == "" {
		return 0, errors.New("title and content cannot be empty")
	}

	board := &model.Board{
		UserID:     userID,
		Title:      title,
		Content:    content,
		Imagepath:  image_path,
		Is_Private: is_private,
	}

	board_id, err := r.boardRepo.Create(board)
	
	if err != nil {
		return 0, err
	}

	username := ParsingUsernameList(sharedusernamesraw)
	if len(username) == 0 {
		return board_id, nil
	}

	found, missing, err := r.userRepo.FindUserIDsByUsernames(username)
	if err != nil {
		return board_id, err
	}

	if len(missing) > 0 {
		return board_id, fmt.Errorf("board was created, but these usernames don't exist and weren't shared with: %s", strings.Join(missing, ", "))
	}

	userIDs := make([]int, 0, len(found))
	for _, id := range found {
		userIDs = append(userIDs, id)
	}

	if err := r.boardRepo.AddShares(board_id, userIDs); err != nil {
		return board_id, err
	}

	return board_id, nil

}

func ParsingUsernameList(raw string) []string {

	parts := strings.Split(raw, ",")
	seen := make(map[string]bool)
	var out []string

	for _, p := range parts {

		p = strings.TrimSpace(p)

		if p == "" || seen[p] {
			continue
		}
		seen[p] = true
		out = append(out, p)
	}

	return out

}

func (r *BoardService) GetVisibleBoards(userID int) ([]model.Board, error) {
	return r.boardRepo.GetVisibleBoards(userID)
}

func (r *BoardService) GetBoard(id, userID int) (*model.Board, error) {
	board, err := r.boardRepo.GetByIDVisible(id, userID)
	if err != nil {
		return nil, err
	}
	if board == nil {
		return nil, errors.New("board not found")
	}
	return board, nil
}

func (r *BoardService) DeleteBoard(id, userID int) error {
	return r.boardRepo.Delete(id, userID)
}
