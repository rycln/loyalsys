package api

import (
	"fmt"
	"time"
)

type errorTooManyRequests struct {
	err      error
	duration time.Duration
}

func (err *errorTooManyRequests) Error() string {
	return fmt.Sprintf("%v Retry after: %v", err.err, err.duration.String())
}

func (err *errorTooManyRequests) IsTooManyRequests() bool {
	return true
}

func (err *errorTooManyRequests) GetRetryAfterDuration() time.Duration {
	return err.duration
}

func newErrorTooManyRequests(dur time.Duration, err error) error {
	return &errorTooManyRequests{
		err:      err,
		duration: dur,
	}
}

type errorNoContent struct {
	err error
}

func (err *errorNoContent) Error() string {
	return err.err.Error()
}

func (err *errorNoContent) IsNoContent() bool {
	return true
}

func newErrorNoContent(err error) error {
	return &errorNoContent{
		err: err,
	}
}
