package vk

import (
	"fmt"
	"strconv"
)

type ServerError int

type Error struct {
	Code    ServerError `json:"error_code"`
	Message string      `json:"error_msg"`
}

func (e Error) Error() string {
	return e.Code.Error()
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

func (s ServerError) String() string {
	v, ok := map[ServerError]string{
		ErrUnknown:                 "Unknown error occured, try again later",
		ErrApplicationDisabled:     "Application disabled",
		ErrInsufficientPermissions: "Insufficient permissions, use account.getAppPermissions",
	}[s]
	if !ok {
		return strconv.Itoa(int(s))
	}
	return fmt.Sprintf("%s (%d)", v, s)
}

func (s ServerError) Error() string {
	return s.String()
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
)
