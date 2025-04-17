package config

import (
	"crypto/rand"
	"flag"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/rycln/loyalsys/internal/logger"
)

const (
	defaultServerAddr  = ":8080"
	defaultTimeout     = time.Duration(2) * time.Minute
	defaultKeyLength   = 32
	defaultLoggerLevel = "debug"
)

type Cfg struct {
	RunAddr     string        `env:"RUN_ADDRESS"`
	DatabaseURI string        `env:"DATABASE_URI"`
	AccrualAddr string        `env:"ACCRUAL_SYSTEM_ADDRESS"`
	Timeout     time.Duration `env:"TIMEOUT_DUR"`
	Key         string        `env:"JWT_KEY"`
	LogLevel    string        `env:"LOG_LEVEL"`
}

type ConfigBuilder struct {
	cfg *Cfg
	err error
}

func NewConfigBuilder() *ConfigBuilder {
	return &ConfigBuilder{
		cfg: &Cfg{
			RunAddr:  defaultServerAddr,
			Timeout:  defaultTimeout,
			LogLevel: defaultLoggerLevel,
		},
		err: nil,
	}
}

func (b *ConfigBuilder) WithFlagParsing() *ConfigBuilder {
	if b.err != nil {
		return b
	}

	flag.StringVar(&b.cfg.RunAddr, "a", b.cfg.RunAddr, "Address and port to start the server")
	flag.StringVar(&b.cfg.DatabaseURI, "d", b.cfg.DatabaseURI, "Database connection address")
	flag.StringVar(&b.cfg.AccrualAddr, "r", b.cfg.AccrualAddr, "Accrual connection address")
	flag.DurationVar(&b.cfg.Timeout, "t", b.cfg.Timeout, "Timeout duration in seconds")
	flag.StringVar(&b.cfg.Key, "k", b.cfg.Key, "Key for jwt autorization")
	flag.StringVar(&b.cfg.LogLevel, "l", b.cfg.LogLevel, "Logger level")
	flag.Parse()

	return b
}

func (b *ConfigBuilder) WithEnvParsing() *ConfigBuilder {
	if b.err != nil {
		return b
	}

	err := env.Parse(b.cfg)
	if err != nil {
		b.cfg = nil
		b.err = err
		return b
	}

	return b
}

func (b *ConfigBuilder) WithDefaultJWTKey() *ConfigBuilder {
	if b.err != nil {
		return b
	}

	if b.cfg.Key == "" {
		key, err := generateKey(defaultKeyLength)
		if err != nil {
			b.cfg = nil
			b.err = err
			return b
		}
		b.cfg.Key = key
		logger.Log.Warn("Default JWT key used!")
	}

	return b
}

func generateKey(n int) (string, error) {
	key := make([]byte, n)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	return string(key), nil
}

func (b *ConfigBuilder) Build() (*Cfg, error) {
	return b.cfg, b.err
}
