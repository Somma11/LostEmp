package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"

	"lostemp/internal/config"
	"lostemp/internal/db"
	"lostemp/internal/handlers"
	"lostemp/internal/middleware"

	_ "lostemp/docs"
)

// @title LostEmp API
// @version 1.0
// @description API de gestão de empréstimos de equipamentos escolares (mockup/MVP). Papéis: admin, professor, technician.
// @contact.name LostEmp
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfg := config.Load()

	database, err := db.Open(cfg.DBPath)
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})

	app.Get("/docs/*", swagger.HandlerDefault)
	handlers.RegisterRoutes(app, cfg, database)

	log.Printf("LostEmp API ouvindo em http://localhost:%s (docs em /docs)", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
