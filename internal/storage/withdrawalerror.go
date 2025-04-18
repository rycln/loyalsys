package storage

import "errors"

var (
	ErrNoWithdrawal = errors.New("no withdrawals")
)

type errNoWithdrawal struct {
	err error
}

func (err *errNoWithdrawal) Error() string {
	return err.err.Error()
}

func (err *errNoWithdrawal) Unwrap() error {
	return err.err
}

func (err *errNoWithdrawal) IsErrNoWithdrawal() bool {
	return true
}

func newErrNoWithdrawal(err error) error {
	return &errNoWithdrawal{
		err: err,
	}
}
