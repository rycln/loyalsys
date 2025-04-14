package api

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrTooManyRequests = errors.New("too many requests")
	ErrNoContent       = errors.New("no content")
)

type errRetryAfter struct {
	err      error
	duration time.Duration
}

func (err *errRetryAfter) Error() string {
	return fmt.Sprintf("%v Retry after: %v", err.err, err.duration.String())
}

func (err *errRetryAfter) GetRetryAfterDuration() time.Duration {
	return err.duration
}

func newErrRetryAfter(dur time.Duration, err error) error {
	return &errRetryAfter{
		err:      err,
		duration: dur,
	}
}
