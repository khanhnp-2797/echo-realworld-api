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
	articleHandler := handler.NewArticleHandler(articleSvc, commentSvc, userSvc)
	tagHandler := handler.NewTagHandler(tagSvc)

	registerAPIRoutes(e, cfg.JWT.Secret, commentRepo, userHandler, articleHandler, tagHandler)
}

// New configures and returns the Echo router (for use without RegisterRoutes).
func New(
	jwtSecret string,
	commentRepo repository.CommentRepository,
	userHandler *handler.UserHandler,
	articleHandler *handler.ArticleHandler,
	tagHandler *handler.TagHandler,
) *echo.Echo {
	e := echo.New()
	e.Use(echomiddleware.CORS())
	registerAPIRoutes(e, jwtSecret, commentRepo, userHandler, articleHandler, tagHandler)
	return e
}

// registerAPIRoutes attaches all /api/* routes to an Echo instance.
func registerAPIRoutes(
	e *echo.Echo,
	jwtSecret string,
	commentRepo repository.CommentRepository,
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
	api.GET("/profiles/:username", userHandler.GetProfile, optAuth)          // GET    /api/profiles/:username
	api.POST("/profiles/:username/follow", userHandler.FollowUser, auth)     // POST   /api/profiles/:username/follow
	api.DELETE("/profiles/:username/follow", userHandler.UnfollowUser, auth) // DELETE /api/profiles/:username/follow

	// ── Task 2+5: Articles ────────────────────────────────────────────────────
	api.GET("/articles/feed", articleHandler.Feed, auth)                           // GET    /api/articles/feed   (must be before :slug)
	api.GET("/articles", articleHandler.ListArticles, optAuth)                     // GET    /api/articles
	api.GET("/articles/:slug", articleHandler.GetArticle, optAuth)                 // GET    /api/articles/:slug
	api.POST("/articles/:slug/favorite", articleHandler.FavoriteArticle, auth)     // POST   /api/articles/:slug/favorite
	api.DELETE("/articles/:slug/favorite", articleHandler.UnfavoriteArticle, auth) // DELETE /api/articles/:slug/favorite

	// ── Task 3+6: Comments ────────────────────────────────────────────────────
	api.POST("/articles/:slug/comments", articleHandler.AddComment, auth)                                                // POST   /api/articles/:slug/comments
	api.GET("/articles/:slug/comments", articleHandler.GetComments, optAuth)                                             // GET    /api/articles/:slug/comments
	api.DELETE("/articles/:slug/comments/:id", articleHandler.DeleteComment, auth, middleware.CommentOwner(commentRepo)) // DELETE /api/articles/:slug/comments/:id

	// ── Task 1: Tags ──────────────────────────────────────────────────────────
	api.GET("/tags", tagHandler.ListTags) // GET /api/tags
}
