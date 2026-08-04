package service

import (
	"errors"
	"strings"

	"weblog/model"
	"weblog/repo"
)

type CommentService struct {
	commentRepo *repo.CommentRepo
	boardRepo   *repo.BoardRepo
}

func NewCommentService(c *repo.CommentRepo, b *repo.BoardRepo) *CommentService {
	return &CommentService{commentRepo: c, boardRepo: b}
}

func (s *CommentService) AddComment(boardID, userID int, content string) error {
	content = strings.TrimSpace(content)
	if content == "" {
		return errors.New("comment cannot be empty")
	}

	board, err := s.boardRepo.GetByIDVisible(boardID, userID)
	if err != nil {
		return err
	}
	if board == nil {
		return errors.New("board not found")
	}

	comment := &model.Comment{
		BoardID: boardID,
		UserID:  userID,
		Content: content,
	}
	return s.commentRepo.Add(comment)
}

func (s *CommentService) GetComments(boardID int) ([]model.Comment, error) {
	return s.commentRepo.GetByBoardID(boardID)
}
