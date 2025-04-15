package storage

import (
	"errors"
)

var (
	ErrLoginConflict = errors.New("login already registered")
	ErrNoUser        = errors.New("user does not exist")
)

type loginConflict struct {
	err error
}

func (err *loginConflict) Error() string {
	return err.err.Error()
}

func (err *loginConflict) IsLoginConflict() bool {
	return true
}

func newErrLoginConflict(err error) error {
	return &loginConflict{
		err: err,
	}
}
