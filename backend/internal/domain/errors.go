package domain

import "errors"

var (
	ErrDeviceNotFound      = errors.New("device not found")
	ErrInvalidDeviceType   = errors.New("invalid device type")
	ErrInvalidPayload      = errors.New("invalid payload")
	ErrValidationFailed    = errors.New("validation failed")
	ErrInternal            = errors.New("internal error")
)
