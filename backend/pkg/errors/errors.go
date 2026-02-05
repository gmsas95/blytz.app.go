package errors

import (
	"errors"
	"fmt"
)

type ErrorCode string

const (
	ErrInternal            ErrorCode = "INTERNAL_ERROR"
	ErrValidation          ErrorCode = "VALIDATION_ERROR"
	ErrNotFound            ErrorCode = "NOT_FOUND"
	ErrUnauthorized        ErrorCode = "UNAUTHORIZED"
	ErrForbidden           ErrorCode = "FORBIDDEN"
	ErrConflict            ErrorCode = "CONFLICT"
	ErrTimeout             ErrorCode = "TIMEOUT"
	ErrRateLimit           ErrorCode = "RATE_LIMIT_EXCEEDED"
	ErrInvalidBid          ErrorCode = "INVALID_BID"
	ErrAuctionNotLive      ErrorCode = "AUCTION_NOT_LIVE"
	ErrAuctionEnded        ErrorCode = "AUCTION_ENDED"
	ErrBidTooLow           ErrorCode = "BID_TOO_LOW"
	ErrInsufficientFunds   ErrorCode = "INSUFFICIENT_FUNDS"
	ErrPaymentFailed       ErrorCode = "PAYMENT_FAILED"
	ErrOrderNotCancellable ErrorCode = "ORDER_NOT_CANCELLABLE"
)

type AppError struct {
	Code       ErrorCode
	Message    string
	Details    map[string]interface{}
	Cause      error
	HTTPStatus int
}

func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (cause: %v)", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Cause
}

func New(code ErrorCode, message string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		Details:    make(map[string]interface{}),
		HTTPStatus: HTTPStatusFromCode(code),
	}
}

func Wrap(err error, code ErrorCode, message string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		Cause:      err,
		Details:    make(map[string]interface{}),
		HTTPStatus: HTTPStatusFromCode(code),
	}
}

func (e *AppError) WithDetails(key string, value interface{}) *AppError {
	e.Details[key] = value
	return e
}

func (e *AppError) WithHTTPStatus(status int) *AppError {
	e.HTTPStatus = status
	return e
}

func HTTPStatusFromCode(code ErrorCode) int {
	switch code {
	case ErrNotFound:
		return 404
	case ErrValidation:
		return 400
	case ErrUnauthorized:
		return 401
	case ErrForbidden:
		return 403
	case ErrConflict:
		return 409
	case ErrTimeout:
		return 504
	case ErrRateLimit:
		return 429
	default:
		return 500
	}
}

func IsNotFound(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code == ErrNotFound
	}
	return false
}

func IsValidation(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code == ErrValidation
	}
	return false
}

func IsUnauthorized(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code == ErrUnauthorized
	}
	return false
}

var (
	ErrUserNotFound       = New(ErrNotFound, "user not found")
	ErrAuctionNotFound    = New(ErrNotFound, "auction not found")
	ErrProductNotFound    = New(ErrNotFound, "product not found")
	ErrOrderNotFound      = New(ErrNotFound, "order not found")
	ErrInvalidInput       = New(ErrValidation, "invalid input")
	ErrUnauthorizedAccess = New(ErrUnauthorized, "unauthorized access")
)