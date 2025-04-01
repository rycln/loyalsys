package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/timeout"
	"github.com/rycln/loyalsys/internal/config"
	"github.com/rycln/loyalsys/internal/handlers"
	"github.com/rycln/loyalsys/internal/logger"
	"github.com/rycln/loyalsys/internal/services"
	"github.com/rycln/loyalsys/internal/storage"
)

func main() {
	err := logger.LogInit()
	if err != nil {
		log.Fatalf("Can't initialize the logger: %v", err)
	}
	defer logger.Log.Sync()

	cfg := config.NewCfg()

	db, err := storage.NewDB(cfg.DatabaseDsn)
	if err != nil {
		log.Fatalf("Can't open database: %v", err)
	}
	userstrg := storage.NewUserStorage(db)

	userservice := services.NewUserService(userstrg)

	registerHandler := handlers.NewRegisterHandler(userservice, cfg)

	app := fiber.New()

	app.Post("/api/user/register", timeout.NewWithContext(registerHandler, cfg.Timeout))

	err = app.Listen(cfg.RunAddr)
	if err != nil {
		log.Fatalf("Can't start the server: %v", err)
	}
}
