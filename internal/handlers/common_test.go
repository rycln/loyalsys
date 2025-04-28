package handlers

import (
	"errors"
	"time"

	"github.com/rycln/loyalsys/internal/config"
	"github.com/rycln/loyalsys/internal/models"
)

const (
	testUserID       = models.UserID(1)
	testUserLogin    = "login"
	testUserPassword = "password"
	testKey          = "secret_key"
	testTimeout      = time.Duration(5) * time.Second
)

var (
	testCfg = &config.Cfg{
		Timeout: testTimeout,
		Key:     testKey,
	}
	errTest = errors.New("test error")
)
