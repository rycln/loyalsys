package services

import "errors"

var (
	ErrWrongOrderNum     = errors.New("luhn algorithm validation failed")
	ErrNotEnoughCurrency = errors.New("not enough currency")
)

type errWrongOrderNum struct {
	err error
}

func (err *errWrongOrderNum) Error() string {
	return err.err.Error()
}

func (err *errWrongOrderNum) Unwrap() error {
	return err.err
}

func (err *errWrongOrderNum) IsErrWrongOrderNum() bool {
	return true
}

func newErrWrongOrderNum(err error) error {
	return &errWrongOrderNum{
		err: err,
	}
}

type errNotEnoughCurrency struct {
	err error
}

func (err *errNotEnoughCurrency) Error() string {
	return err.err.Error()
}

func (err *errNotEnoughCurrency) Unwrap() error {
	return err.err
}

func (err *errNotEnoughCurrency) IsErrNotEnoughCurrency() bool {
	return true
}

func newErrNotEnoughCurrency(err error) error {
	return &errNotEnoughCurrency{
		err: err,
	}
}
