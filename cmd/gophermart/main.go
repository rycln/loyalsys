package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/rycln/loyalsys/internal/app"
	"github.com/rycln/loyalsys/internal/config"
)

const shutdownTimeout = 5 * time.Second

func main() {
	cfg, err := config.NewConfigBuilder().
		WithFlagParsing().
		WithEnvParsing().
		WithDefaultJWTKey().
		Build()
	if err != nil {
		log.Fatalf("Can't initialize the configuration: %v", err)
	}

	app, err := app.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	shutdownDone := make(chan struct{})

	go func() {
		<-stop

		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		err := app.Shutdown(ctx)
		if err != nil {
			log.Fatalf("Shutdown error: %v", err)
		}

		shutdownDone <- struct{}{}
	}()

	err = app.Listen(cfg.RunAddr)
	if err != nil {
		log.Fatalf("Can't start the server: %v", err)
	}

	<-shutdownDone

	err = app.Cleanup()
	if err != nil {
		log.Fatalf("Cleanup error: %v", err)
	}

	log.Println(strings.TrimPrefix(os.Args[0], "./") + " shutted down gracefully")
}
