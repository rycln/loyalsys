package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	testServerAddr  = ":8081"
	testDatabaseURI = "test_dsn"
	testAccrualAddr = "test_addr"
	testTimeout     = time.Duration(3) * time.Minute
	testKey         = "secret_key"
	testLoggerLevel = "info"
)

func TestConfigBuilder_WithEnvParsing(t *testing.T) {
	testCfg := &Cfg{
		RunAddr:     testAccrualAddr,
		DatabaseURI: testDatabaseURI,
		AccrualAddr: testAccrualAddr,
		Timeout:     testTimeout,
		Key:         testKey,
		LogLevel:    testLoggerLevel,
	}

	t.Setenv("RUN_ADDRESS", testCfg.RunAddr)
	t.Setenv("DATABASE_URI", testCfg.DatabaseURI)
	t.Setenv("ACCRUAL_SYSTEM_ADDRESS", testCfg.AccrualAddr)
	t.Setenv("TIMEOUT_DUR", testCfg.Timeout.String())
	t.Setenv("JWT_KEY", testCfg.Key)
	t.Setenv("LOG_LEVEL", testCfg.LogLevel)

	t.Run("valid test", func(t *testing.T) {
		cfg, err := NewConfigBuilder().
			WithEnvParsing().
			Build()
		assert.NoError(t, err)
		assert.Equal(t, testCfg, cfg)
	})
}

func TestConfigBuilder_WithDefaultJWTKey(t *testing.T) {
	t.Run("valid test", func(t *testing.T) {
		cfg, err := NewConfigBuilder().
			WithDefaultJWTKey().
			Build()
		assert.NoError(t, err)
		assert.NotEmpty(t, cfg.Key)
	})
}

func TestConfigBuilder_WithFlagParsing(t *testing.T) {
	oldArgs := os.Args
	defer func() {
		os.Args = oldArgs
	}()

	testCfg := &Cfg{
		RunAddr:     testServerAddr,
		DatabaseURI: testDatabaseURI,
		AccrualAddr: testAccrualAddr,
		Timeout:     testTimeout,
		Key:         testKey,
		LogLevel:    testLoggerLevel,
	}

	t.Run("valid test", func(t *testing.T) {
		os.Args = []string{
			"./gophermart",
			"-a=" + testCfg.RunAddr,
			"-d=" + testCfg.DatabaseURI,
			"-r=" + testCfg.AccrualAddr,
			"-t=" + testCfg.Timeout.String(),
			"-k=" + testCfg.Key,
			"-l=" + testCfg.LogLevel,
		}

		cfg, err := NewConfigBuilder().
			WithFlagParsing().
			Build()
		assert.NoError(t, err)
		assert.Equal(t, testCfg, cfg)
	})
}
