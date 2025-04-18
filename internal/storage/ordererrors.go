package storage

import "errors"

var (
	ErrNoOrder = errors.New("order does not exist")
)

type errNoOrder struct {
	err error
}

func (err *errNoOrder) Error() string {
	return err.err.Error()
}

func (err *errNoOrder) Unwrap() error {
	return err.err
}

func (err *errNoOrder) IsErrNoOrder() bool {
	return true
}

func newErrNoOrder(err error) error {
	return &errNoOrder{
		err: err,
	}
}
