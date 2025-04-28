package services

import "errors"

var (
	ErrWrongNum      = errors.New("luhn algorithm validation failed")
	ErrOrderExists   = errors.New("order already registered by user")
	ErrOrderConflict = errors.New("order already registered by other user")
)

type errWrongNum struct {
	err error
}

func (err *errWrongNum) Error() string {
	return err.err.Error()
}

func (err *errWrongNum) Unwrap() error {
	return err.err
}

func (err *errWrongNum) IsErrWrongNum() bool {
	return true
}

func newErrWrongNum(err error) error {
	return &errWrongNum{
		err: err,
	}
}

type errOrderExists struct {
	err error
}

func (err *errOrderExists) Error() string {
	return err.err.Error()
}

func (err *errOrderExists) Unwrap() error {
	return err.err
}

func (err *errOrderExists) IsErrOrderExists() bool {
	return true
}

func newErrOrderExists(err error) error {
	return &errOrderExists{
		err: err,
	}
}

type errOrderConflict struct {
	err error
}

func (err *errOrderConflict) Error() string {
	return err.err.Error()
}

func (err *errOrderConflict) Unwrap() error {
	return err.err
}

func (err *errOrderConflict) IsErrOrderConflict() bool {
	return true
}

func newErrOrderConflict(err error) error {
	return &errOrderConflict{
		err: err,
	}
}
