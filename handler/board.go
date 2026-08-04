package handler


import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"weblog/repo"
	"weblog/service"
)

const (
	uploadsDir    = "uploads"
	maxImageBytes = 5 << 20 // 5 MB
)

var allowedImageExts = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".webp": true,
}

type BoardHandler struct {
	boardService   *service.BoardService
	commentService *service.CommentService
	userRepo       *repo.UserRepo
}

func NewBoardHandler(b *service.BoardService, c *service.CommentService, u *repo.UserRepo) *BoardHandler {
	return &BoardHandler{boardService: b, commentService: c, userRepo: u}
}

func (h *BoardHandler) currentUsername(userID int) string {
	user, err := h.userRepo.FindUserByID(userID)
	if err != nil || user == nil {
		return ""
	}
	return user.Username
}

func (h *BoardHandler) GetHome(c echo.Context) error {
	userID := c.Get("user_id").(int)

	boards, err := h.boardService.GetVisibleBoards(userID)
	if err != nil {
		// return c.String(http.StatusInternalServerError, "failed to fetch boards")÷
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.Render(http.StatusOK, "index.html", map[string]interface{}{
		"Username": h.currentUsername(userID),
		"Boards":   boards,
	})
}

func (h *BoardHandler) GetBoard(c echo.Context) error {
	userID := c.Get("user_id").(int)

	boardID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "invalid board id")
	}

	board, err := h.boardService.GetBoard(boardID, userID)
	if err != nil {

		return c.String(http.StatusNotFound, "board not found")
	}

	comments, err := h.commentService.GetComments(boardID)
	if err != nil {
		return c.String(http.StatusInternalServerError, "failed to fetch comments")
	}

	return c.Render(http.StatusOK, "board.html", map[string]interface{}{
		"Username": h.currentUsername(userID),
		"Board":    board,
		"Comments": comments,
		"IsOwner":  board.UserID == userID,
	})
}

func (h *BoardHandler) CreatBoard(c echo.Context) error {
	userID := c.Get("user_id").(int)

	title := c.FormValue("title")
	content := c.FormValue("content")
	isPrivate := c.FormValue("is_private") == "on" || c.FormValue("is_private") == "true"
	sharedUsernames := c.FormValue("shared_usernames")

	imagePath, err := h.saveUploadedImage(c)

	if err != nil {
		return c.Render(http.StatusOK, "index.html", map[string]interface{}{
			"Username": h.currentUsername(userID),
			"Error":    err.Error(),
		})
	}

	_, err = h.boardService.CreatBoard(userID, title, content, imagePath, isPrivate, sharedUsernames)

	if err != nil {
		return c.Render(http.StatusOK, "index.html", map[string]interface{}{
			"Username": h.currentUsername(userID),
			"Error":    err.Error(),
		})
	}

	return c.Redirect(http.StatusSeeOther, "/")
}

func (h *BoardHandler) saveUploadedImage(c echo.Context) (string, error) {

	fileHeader, err := c.FormFile("image")

	if err != nil {
		return "", nil // no image submitted 
	}

	if fileHeader.Size > maxImageBytes {
		return "", fmt.Errorf("image is too large (max %d MB)", maxImageBytes/(1<<20))
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if !allowedImageExts[ext] {
		return "", fmt.Errorf("unsupported image type: %s", ext)
	}

	src, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to read uploaded image")
	}
	defer src.Close()

	filename := uuid.NewString() + ext
	dstPath := filepath.Join(uploadsDir, filename)

	dst, err := os.Create(dstPath)
	if err != nil {
		return "", fmt.Errorf("failed to save uploaded image")
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("failed to save uploaded image")
	}

	return "/" + uploadsDir + "/" + filename, nil
}


func (h *BoardHandler) DeleteBoard(c echo.Context) error {
	userID := c.Get("user_id").(int)

	boardID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "invalid board id")
	}

	if err := h.boardService.DeleteBoard(boardID, userID); err != nil {
		return c.String(http.StatusForbidden, err.Error())
	}

	return c.Redirect(http.StatusSeeOther, "/")
}
