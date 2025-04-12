package api

import (
	"fmt"
	"time"
)

type errorTooManyRequests struct {
	message  string
	duration time.Duration
}

func (err *errorTooManyRequests) Error() string {
	return fmt.Sprintf("%v %v", err.message, err.duration.String())
}

func (err *errorTooManyRequests) IsTooManyRequests() bool {
	return true
}

func (err *errorTooManyRequests) GetRetryAfterDuration() time.Duration {
	return err.duration
}

func newErrorTooManyRequests(dur time.Duration) error {
	return &errorTooManyRequests{
		message:  "too many requests. Retry after:",
		duration: dur,
	}
}

type errorNoContent struct {
	message string
}

func (err *errorNoContent) Error() string {
	return err.message
}

func (err *errorNoContent) IsNoContent() bool {
	return true
}

func newErrorNoContent() error {
	return &errorNoContent{
		message: "order not registered",
	}
}
