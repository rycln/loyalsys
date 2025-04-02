package config

import (
	"crypto/rand"
	"flag"
	"time"

	"github.com/caarlos0/env/v11"
)

const (
	defaultServerAddr  = ":8080"
	defultTimeout      = time.Duration(2) * time.Minute
	defaultKeyLength   = 10
	defaultLoggerLevel = "info"
)

type Cfg struct {
	RunAddr     string        `env:"RUN_ADDRESS"`
	DatabaseDsn string        `env:"DATABASE_DSN"`
	AccrualAddr string        `env:"ACCRUAL_SYSTEM_ADDRESS"`
	Timeout     time.Duration `env:"TIMEOUT_DUR"`
	Key         string        `env:"KEY"`
	LogLevel    string        `env:"-"`
}

func NewCfg() *Cfg {
	cfg := &Cfg{}

	flag.StringVar(&cfg.RunAddr, "a", defaultServerAddr, "Address and port to start the server (environment variable RUN_ADDRESS has higher priority)")
	flag.StringVar(&cfg.DatabaseDsn, "d", "", "Database connection address (environment variable DATABASE_DSN has higher priority)")
	flag.StringVar(&cfg.AccrualAddr, "r", "", "Accrual connection address (environment variable ACCRUAL_SYSTEM_ADDRESS has higher priority)")
	flag.DurationVar(&cfg.Timeout, "t", defultTimeout, "Timeout duration in seconds (environment variable TIMEOUT_DUR has higher priority)")
	flag.StringVar(&cfg.Key, "k", "", "Key for jwt autorization (environment variable KEY has higher priority)")
	flag.StringVar(&cfg.LogLevel, "l", defaultLoggerLevel, "Logger level")
	flag.Parse()

	err := env.Parse(cfg)
	if err != nil {
		panic(err)
	}

	if cfg.Key == "" {
		cfg.Key = generateKey(defaultKeyLength)
	}

	return cfg
}

func generateKey(n int) string {
	key := make([]byte, n)
	_, err := rand.Read(key)
	if err != nil {
		panic(err)
	}
	return string(key)
}
