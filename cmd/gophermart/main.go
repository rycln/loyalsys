package main

import (
	"log"

	"github.com/gofiber/contrib/fiberzap/v2"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/timeout"
	"github.com/rycln/loyalsys/internal/config"
	"github.com/rycln/loyalsys/internal/handlers"
	"github.com/rycln/loyalsys/internal/logger"
	"github.com/rycln/loyalsys/internal/middleware"
	"github.com/rycln/loyalsys/internal/services"
	"github.com/rycln/loyalsys/internal/storage"
	"go.uber.org/zap/zapcore"
)

func main() {
	cfg := config.NewCfg()

	err := logger.LogInit(cfg.LogLevel)
	if err != nil {
		log.Fatalf("Can't initialize the logger: %v", err)
	}
	defer logger.Log.Sync()

	db, err := storage.NewDB(cfg.DatabaseURI)
	if err != nil {
		log.Fatalf("Can't open database: %v", err)
	}

	userstrg := storage.NewUserStorage(db)
	orderstrg := storage.NewOrderStorage(db)

	userservice := services.NewUserService(userstrg)
	orderservice := services.NewOrderService(orderstrg)

	registerHandler := handlers.NewRegisterHandler(userservice, cfg)
	loginHandler := handlers.NewLoginHandler(userservice, cfg)
	postOrderHandler := handlers.NewPostOrderHandler(orderservice, cfg)

	app := fiber.New()
	app.Use(fiberzap.New(fiberzap.Config{
		Logger: logger.Log,
		Fields: []string{"url", "method", "latency", "status", "bytesSent"},
		Levels: []zapcore.Level{zapcore.InfoLevel},
	}))
	app.Post("/api/user/register", middleware.ContentTypeChecker("application/json"), timeout.NewWithContext(registerHandler, cfg.Timeout))
	app.Post("/api/user/login", middleware.ContentTypeChecker("application/json"), timeout.NewWithContext(loginHandler, cfg.Timeout))
	app.Use(middleware.NoTokenChecker(), jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(cfg.Key)},
	}))
	app.Post("/api/user/orders", middleware.ContentTypeChecker("text/plain"), timeout.NewWithContext(postOrderHandler, cfg.Timeout))

	err = app.Listen(cfg.RunAddr)
	if err != nil {
		log.Fatalf("Can't start the server: %v", err)
	}
}
