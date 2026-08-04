package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"weblog/service"
)

type CommentHandler struct {
	commentService *service.CommentService
}

func NewCommentHandler(c *service.CommentService) *CommentHandler {
	return &CommentHandler{commentService: c}
}

func (h *CommentHandler) CreateComment(c echo.Context) error {
	userID := c.Get("user_id").(int)

	boardID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "invalid board id")
	}

	content := c.FormValue("content")

	if err := h.commentService.AddComment(boardID, userID, content); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.Redirect(http.StatusSeeOther, "/weblog/" + c.Param("id"))
}
