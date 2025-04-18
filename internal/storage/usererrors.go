package storage

import (
	"errors"
)

var (
	ErrLoginConflict = errors.New("login already registered")
	ErrNoUser        = errors.New("user does not exist")
)

type errLoginConflict struct {
	err error
}

func (err *errLoginConflict) Error() string {
	return err.err.Error()
}

func (err *errLoginConflict) Unwrap() error {
	return err.err
}

func (err *errLoginConflict) IsErrLoginConflict() bool {
	return true
}

func newErrLoginConflict(err error) error {
	return &errLoginConflict{
		err: err,
	}
}

type errNoUser struct {
	err error
}

func (err *errNoUser) Error() string {
	return err.err.Error()
}

func (err *errNoUser) Unwrap() error {
	return err.err
}

func (err *errNoUser) IsErrNoUser() bool {
	return true
}

func newErrNoUser(err error) error {
	return &errNoUser{
		err: err,
	}
}
