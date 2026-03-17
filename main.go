// Package main is the entry point for the RealWorld API.
//
// @title           RealWorld Echo API
// @version         1.0
// @description     Conduit backend implementation using Go + Echo + GORM (PostgreSQL).
// @contact.name    API Support
//
// @host            localhost:8080
// @BasePath        /api
//
// @securityDefinitions.apikey BearerAuth
// @in                         header
// @name                       Authorization
// @description                Use format: Token <jwt>
package main

import (
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "github.com/khanhnp-2797/echo-realworld-api/docs" // swagger generated docs
	"github.com/khanhnp-2797/echo-realworld-api/internal/database"
	"github.com/khanhnp-2797/echo-realworld-api/internal/router"
	appvalidator "github.com/khanhnp-2797/echo-realworld-api/pkg/validator"
)

const port = "8080"

func main() {
	// 1. Khởi tạo Echo
	e := echo.New()
	e.HideBanner = true

	// 2. Thêm Middleware cơ bản (Logger, Recover)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// 3. Validator
	e.Validator = appvalidator.New()

	// 4. Kết nối Database (load .env + connect PostgreSQL + auto-migrate)
	db := database.InitDB()

	// 5. Đăng ký Routes (wire repos → services → handlers → routes)
	router.RegisterRoutes(e, db)

	// 6. Khởi động Server
	fmt.Printf("\n🚀 Server running at:  http://localhost:%s\n", port)
	fmt.Printf("📖 Swagger UI at:      http://localhost:%s/swagger/index.html\n\n", port)
	if err := e.Start(":" + port); err != nil {
		log.Fatal(err)
	}
}
