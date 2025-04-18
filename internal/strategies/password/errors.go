package password

import "errors"

var (
	ErrWrongPassword = errors.New("wrong password")
)

type errWrongPassword struct {
	err error
}

func (err *errWrongPassword) Error() string {
	return err.err.Error()
}

func (err *errWrongPassword) Unwrap() error {
	return err.err
}

func (err *errWrongPassword) IsErrWrongPassword() bool {
	return true
}

func newErrWrongPassword(err error) error {
	return &errWrongPassword{
		err: err,
	}
}
