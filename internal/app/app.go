package app

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/gofiber/contrib/fiberzap/v2"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/timeout"
	"github.com/rycln/loyalsys/internal/api"
	"github.com/rycln/loyalsys/internal/config"
	"github.com/rycln/loyalsys/internal/db"
	"github.com/rycln/loyalsys/internal/handlers"
	"github.com/rycln/loyalsys/internal/logger"
	"github.com/rycln/loyalsys/internal/middleware"
	"github.com/rycln/loyalsys/internal/services"
	"github.com/rycln/loyalsys/internal/storage"
	"github.com/rycln/loyalsys/internal/strategies/password"
	"github.com/rycln/loyalsys/internal/worker"
	"go.uber.org/zap/zapcore"
)

type App struct {
	*fiber.App
	db           *sql.DB
	workerStopCh chan struct{}
	workerDoneCh chan struct{}
}

func New(cfg *config.Cfg) (*App, error) {
	err := logger.LogInit(cfg.LogLevel)
	if err != nil {
		return nil, fmt.Errorf("can't initialize the logger: %v", err)
	}

	database, err := storage.NewDB(cfg.DatabaseURI)
	if err != nil {
		return nil, fmt.Errorf("can't open database: %v", err)
	}
	err = db.RunMigrations(database)
	if err != nil {
		return nil, fmt.Errorf("can't apply migrations: %v", err)
	}

	userStrg := storage.NewUserStorage(database)
	orderStrg := storage.NewOrderStorage(database)

	restyClient := resty.New()
	client := api.NewOrderUpdateClient(restyClient, cfg.AccrualAddr, cfg.Timeout)
	workerCfg := worker.NewSyncWorkerConfigBuilder().
		WithTimeout(cfg.Timeout).
		Build()
	orderUpdater := worker.NewOrderSyncWorker(client, orderStrg, workerCfg)

	stopCh := make(chan struct{})
	doneCh := orderUpdater.Run(stopCh)

	passwordStrategy := password.NewBCryptHasher()
	userService := services.NewUserService(userStrg, passwordStrategy)
	orderService := services.NewOrderService(orderStrg)
	jwtService := services.NewJWTService(cfg.Key)

	registerHandler := handlers.NewRegisterHandler(userService, jwtService)
	loginHandler := handlers.NewLoginHandler(userService, jwtService)
	postOrderHandler := handlers.NewPostOrderHandler(orderService, jwtService)
	getOrdersHandler := handlers.NewGetOrdersHandler(orderService, jwtService)

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
	app.Get("/api/user/orders", timeout.NewWithContext(getOrdersHandler, cfg.Timeout))

	return &App{
		App:          app,
		db:           database,
		workerStopCh: stopCh,
		workerDoneCh: doneCh,
	}, nil
}

func (app *App) Shutdown(ctx context.Context) error {
	close(app.workerStopCh)

	select {
	case <-app.workerDoneCh:
	case <-ctx.Done():
		return fmt.Errorf("worker shutdown timeout: %w", ctx.Err())
	}

	if err := app.App.Shutdown(); err != nil {
		return err
	}
	return nil
}

func (app *App) Cleanup() error {
	defer logger.Log.Sync()

	if err := app.db.Close(); err != nil {
		return err
	}

	return nil
}
