package apperrors

import "errors"

var (
	ErrInvalidSubscriptionData = errors.New("invalid subscription data")
	ErrSubscriptionNotFound    = errors.New("subscription not found")
)
