package vk

import (
	"fmt"
	"strconv"
)

type ServerError int

type RequestParam struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ExecuteError struct {
	Method  string      `json:"method"`
	Code    ServerError `json:"error_code"`
	Message string      `json:"error_msg"`
}

type Error struct {
	Code    ServerError    `json:"error_code,omitempty"`
	Message string         `json:"error_msg,omitempty"`
	Params  []RequestParam `json:"request_params,omitempty"`
	Request Request        `json:"-"`
}

func (e *Error) setRequest(r Request) {
	e.Request = r
}

func (e Error) Error() string {
	return fmt.Sprintf("%s (%d)", e.Message, e.Code)
}

func (e ExecuteError) Error() string {
	return fmt.Sprintf("%s: %s (%d)", e.Method, e.Message, e.Code)
}

// Is returns true, if err equals (or is Error with code equals) e
func (e ServerError) Is(err error) bool {
	if error(e) == err {
		return true
	}
	if another, ok := err.(ServerError); ok {
		return another == e
	}
	if another, ok := err.(Error); ok {
		return another.Code == e
	}
	return false
}

func IsServerError(err error) bool {
	if _, ok := err.(Error); ok {
		return true
	}
	return false
}

func GetServerError(err error) Error {
	if s, ok := err.(Error); ok {
		return s
	}
	panic("not a server error")
}

type ErrorResponse struct {
	Error Error `json:"error"`
}

func (s ServerError) Error() string {
	return strconv.Itoa(int(s))
}

const (
	// Possible error codes
	// https://vk.com/dev/errors
	ErrZero ServerError = iota
	ErrUnknown
	ErrApplicationDisabled
	ErrUnknownMethod
	ErrInvalidSignature
	ErrAuthFailed
	ErrTooManyRequests
	ErrInsufficientPermissions
	ErrInvalidRequest
	ErrTooManyOneTypeRequests
	ErrInternalServerError
	ErrAppInTestMode
	ErrCaptchaNeeded
	ErrNotAllowed
	ErrHttpsOnly
	ErrNeedValidation
	ErrStandaloneOnly
	ErrStandaloneOpenAPIOnly
	ErrMethodDisabled
	ErrNeedConfirmation
	ErrOneOfParametersInvalid    ServerError = 100
	ErrInvalidAPIID              ServerError = 101
	ErrInvalidAUserID            ServerError = 113
	ErrInvalidTimestamp          ServerError = 150
	ErrAlbumAccessProhibited     ServerError = 200
	ErrGroupAccessProhibited     ServerError = 203
	ErrAlbumOverflow             ServerError = 300
	ErrMoneyTransferNotAllowed   ServerError = 500
	ErrInsufficientPermissionsAd ServerError = 600
	ErrInternalServerErrorAd     ServerError = 603

	ErrBadResponseCode ServerError = -1
)
