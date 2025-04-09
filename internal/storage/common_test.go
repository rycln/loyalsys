package storage

import (
	"errors"

	"github.com/rycln/loyalsys/internal/models"
)

const (
	testUserID = models.UserID(1)
)

var (
	errTest = errors.New("test error")
)
