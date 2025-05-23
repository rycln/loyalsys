package handlers

import (
	"errors"
	"time"

	"github.com/rycln/loyalsys/internal/models"
)

const (
	testUserID       = models.UserID(1)
	testUserLogin    = "login"
	testUserPassword = "password"
	testKey          = "secret_key"
	testTimeout      = time.Duration(5) * time.Second
	testJWTString    = "abc.def.ghi"
	validLuhnString  = "4512812345678909"
)

var (
	errTest = errors.New("test error")
)
