package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/gofiber/contrib/fiberzap/v2"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/timeout"
	"github.com/rycln/loyalsys/internal/client"
	"github.com/rycln/loyalsys/internal/config"
	"github.com/rycln/loyalsys/internal/handlers"
	"github.com/rycln/loyalsys/internal/logger"
	"github.com/rycln/loyalsys/internal/middleware"
	"github.com/rycln/loyalsys/internal/services"
	"github.com/rycln/loyalsys/internal/storage"
	"github.com/rycln/loyalsys/internal/strategies/password"
	"github.com/rycln/loyalsys/internal/worker"
	"go.uber.org/zap/zapcore"
)

const shutdownTimeout = 5 * time.Second

type App struct {
	*fiber.App
	cfg    *config.Cfg
	db     *sql.DB
	worker *worker.OrderSyncWorker
}

func New() (*App, error) {
	cfg, err := config.NewConfigBuilder().
		WithFlagParsing().
		WithEnvParsing().
		WithDefaultJWTKey().
		Build()
	if err != nil {
		return nil, fmt.Errorf("can't initialize the configuration: %v", err)
	}

	err = logger.LogInit(cfg.LogLevel)
	if err != nil {
		return nil, fmt.Errorf("can't initialize the logger: %v", err)
	}

	database, err := storage.NewDB(cfg.DatabaseURI)
	if err != nil {
		return nil, fmt.Errorf("can't open database: %v", err)
	}

	userStrg := storage.NewUserStorage(database)
	orderStrg := storage.NewOrderStorage(database)
	withdrawalStrg := storage.NewWithdrawalStorage(database)
	balanceStrg := storage.NewBalanceStorage(database)

	restyClient := resty.New()
	client := client.NewOrderUpdateClient(restyClient, cfg.AccrualAddr, cfg.Timeout)
	workerCfg := worker.NewSyncWorkerConfigBuilder().
		WithTimeout(cfg.Timeout).
		Build()
	orderUpdater := worker.NewOrderSyncWorker(client, orderStrg, workerCfg)

	passwordStrategy := password.NewBCryptHasher()
	userService := services.NewUserService(userStrg, passwordStrategy)
	orderService := services.NewOrderService(orderStrg)
	balanceService := services.NewBalanceService(balanceStrg)
	withdrawalService := services.NewWithdrawalService(withdrawalStrg, balanceService)
	jwtService := services.NewJWTService(cfg.Key)

	registerHandler := handlers.NewRegisterHandler(userService, jwtService)
	loginHandler := handlers.NewLoginHandler(userService, jwtService)
	postOrderHandler := handlers.NewPostOrderHandler(orderService, jwtService)
	getOrdersHandler := handlers.NewGetOrdersHandler(orderService, jwtService)
	getBalanceHandler := handlers.NewGetBalanceHandler(balanceService, jwtService)
	postWithdrawalHandler := handlers.NewPostWithdrawalHandler(withdrawalService, jwtService)
	getWithdrawalsHandler := handlers.NewGetWithdrawalsHandler(withdrawalService, jwtService)

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
	app.Get("/api/user/balance", timeout.NewWithContext(getBalanceHandler, cfg.Timeout))
	app.Post("/api/user/balance/withdraw", timeout.NewWithContext(postWithdrawalHandler, cfg.Timeout))
	app.Get("/api/user/withdrawals", timeout.NewWithContext(getWithdrawalsHandler, cfg.Timeout))

	return &App{
		App:    app,
		cfg:    cfg,
		db:     database,
		worker: orderUpdater,
	}, nil
}

func (app *App) Run() error {
	workerCtx, workerCancel := context.WithCancel(context.Background())
	defer workerCancel()

	doneCh := app.worker.Run(workerCtx)

	go func() {
		err := app.Listen(app.cfg.RunAddr)
		if err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	<-shutdown

	workerCancel()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err := app.shutdown(shutdownCtx, doneCh)
	if err != nil {
		return fmt.Errorf("shutdown error: %v", err)
	}

	err = app.cleanup()
	if err != nil {
		return fmt.Errorf("cleanup error: %v", err)
	}

	log.Println(strings.TrimPrefix(os.Args[0], "./") + " shutted down gracefully")

	return nil
}

func (app *App) shutdown(ctx context.Context, doneCh <-chan struct{}) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("worker shutdown timeout: %w", ctx.Err())
	case <-doneCh:
	}

	if err := app.App.Shutdown(); err != nil {
		return err
	}
	return nil
}

func (app *App) cleanup() error {
	defer logger.Log.Sync()

	if err := app.db.Close(); err != nil {
		return err
	}

	return nil
}
