package services

import (
	"errors"

	"github.com/rycln/loyalsys/internal/models"
)

const (
	testUserID      = models.UserID(1)
	testOtherUserID = models.UserID(2)
)

var (
	errTest = errors.New("test error")
)
