package router

import (
	"log"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"gorm.io/gorm"

	"github.com/khanhnp-2797/echo-realworld-api/internal/config"
	"github.com/khanhnp-2797/echo-realworld-api/internal/handler"
	"github.com/khanhnp-2797/echo-realworld-api/internal/middleware"
	"github.com/khanhnp-2797/echo-realworld-api/internal/repository"
	"github.com/khanhnp-2797/echo-realworld-api/internal/service"
)

// RegisterRoutes wires all repositories, services, and handlers onto an
// existing Echo instance. Intended for simple main() usage:
//
//	db := database.InitDB()
//	router.RegisterRoutes(e, db)
func RegisterRoutes(e *echo.Echo, db *gorm.DB) {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Repositories (data access layer)
	userRepo := repository.NewUserRepository(db)
	articleRepo := repository.NewArticleRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	tagRepo := repository.NewTagRepository(db)

	// Services (business logic layer)
	userSvc := service.NewUserService(userRepo, cfg.JWT)
	articleSvc := service.NewArticleService(articleRepo)
	commentSvc := service.NewCommentService(commentRepo, articleRepo)
	tagSvc := service.NewTagService(tagRepo)

	// Handlers (HTTP layer)
	userHandler := handler.NewUserHandler(userSvc)
	articleHandler := handler.NewArticleHandler(articleSvc, commentSvc)
	tagHandler := handler.NewTagHandler(tagSvc)

	registerAPIRoutes(e, cfg.JWT.Secret, userHandler, articleHandler, tagHandler)
}

// New configures and returns the Echo router (for use without RegisterRoutes).
func New(
	jwtSecret string,
	userHandler *handler.UserHandler,
	articleHandler *handler.ArticleHandler,
	tagHandler *handler.TagHandler,
) *echo.Echo {
	e := echo.New()
	e.Use(echomiddleware.CORS())
	registerAPIRoutes(e, jwtSecret, userHandler, articleHandler, tagHandler)
	return e
}

// registerAPIRoutes attaches all /api/* routes to an Echo instance.
func registerAPIRoutes(
	e *echo.Echo,
	jwtSecret string,
	userHandler *handler.UserHandler,
	articleHandler *handler.ArticleHandler,
	tagHandler *handler.TagHandler,
) {
	e.Use(echomiddleware.CORS())

	// Swagger UI — GET /swagger/*
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	auth := middleware.JWTAuth(jwtSecret, false)   // required auth
	optAuth := middleware.JWTAuth(jwtSecret, true) // optional auth (parse token if present)

	api := e.Group("/api")

	// ── Task 4: Auth ─────────────────────────────────────────────────────────
	api.POST("/users", userHandler.Register)           // POST /api/users
	api.POST("/users/login", userHandler.Login)        // POST /api/users/login
	api.GET("/user", userHandler.GetCurrentUser, auth) // GET  /api/user

	// ── Task 3: Profiles ──────────────────────────────────────────────────────
	api.GET("/profiles/:username", userHandler.GetProfile) // GET /api/profiles/:username

	// ── Task 2: Articles (read-only CRUD) ────────────────────────────────────
	api.GET("/articles", articleHandler.ListArticles, optAuth)     // GET /api/articles
	api.GET("/articles/:slug", articleHandler.GetArticle, optAuth) // GET /api/articles/:slug

	// ── Task 3: Comments ──────────────────────────────────────────────────────
	api.POST("/articles/:slug/comments", articleHandler.AddComment, auth) // POST /api/articles/:slug/comments
	api.GET("/articles/:slug/comments", articleHandler.GetComments)       // GET  /api/articles/:slug/comments

	// ── Task 1: Tags ──────────────────────────────────────────────────────────
	api.GET("/tags", tagHandler.ListTags) // GET /api/tags
}
