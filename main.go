package main

import (
	"database/sql"
	"html/template"
	"io"
	"log"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "github.com/lib/pq"

	appmw "weblog/middleware"

	"weblog/handler"
	"weblog/repo"
	"weblog/service"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplateRenderer() *Template {
	return &Template{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}
}

func getDatabaseURL() string {

	if url := os.Getenv("DATABASE_URL"); url != "" {
		return url
	}
	if url := os.Getenv("DB_CONN"); url != "" {
		return url
	}
	return "postgres://user:password@localhost:5432/weblog_db?sslmode=disable"
}

func main() {

	db, err := sql.Open("postgres", getDatabaseURL())
	if err != nil {
		log.Fatalf("database connection error: %v", err)
	}
	
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("database is not responding: %v", err)
	}

	if err := os.MkdirAll("uploads", 0o755); err != nil {
		log.Fatalf("failed to create uploads directory: %v", err)
	}

	userRepo := repo.NewUserRepo(db)
	boardRepo := repo.NewBoardRepo(db)
	commentRepo := repo.NewCommentRepo(db)

	authService := service.NewAuthService(userRepo)
	boardService := service.NewBoardService(boardRepo, userRepo)
	commentService := service.NewCommentService(commentRepo, boardRepo)

	authHandler := handler.NewAuthHandler(authService)
	boardHandler := handler.NewBoardHandler(boardService, commentService, userRepo)
	commentHandler := handler.NewCommentHandler(commentService)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Renderer = newTemplateRenderer()

	e.Static("/uploads", "uploads")

	e.GET("/login", authHandler.HandleLogin)
	e.POST("/login", authHandler.HandleLogin)
	e.GET("/signup", authHandler.HandleSignup)
	e.POST("/signup", authHandler.HandleSignup)
	e.POST("/logout", authHandler.HandleLogout)

	e.GET("/", boardHandler.GetHome, appmw.AuthMiddleware)
	e.GET("/weblog/:id", boardHandler.GetBoard, appmw.AuthMiddleware)
	e.POST("/weblog", boardHandler.CreatBoard, appmw.AuthMiddleware)
	e.POST("/weblog/:id/delete", boardHandler.DeleteBoard, appmw.AuthMiddleware)
	e.POST("/weblog/:id/comments", commentHandler.CreateComment, appmw.AuthMiddleware)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	e.Logger.Fatal(e.Start(":" + port))
}
