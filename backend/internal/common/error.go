package common

import (
	"errors"
	"net/http"
)

const (
	CodeOK                  = 0
	CodeInvalidRequest      = 1001
	CodeUnauthorized        = 1002
	CodeForbidden           = 1003
	CodeNotFound            = 1004
	CodeConflict            = 1005
	CodeInternal            = 1006
	CodeLoginFailed         = 2001
	CodeProjectKeygenFailed = 3001
	CodeLicenseIssueFailed  = 4001
)

var (
	ErrInvalidParam       = errors.New("invalid param")
	ErrInvalidRequestBody = errors.New("invalid request body")
	ErrNotFound           = errors.New("not found")
)

type AppError struct {
	HTTPStatus int
	Code       int
	Message    string
}

func (e *AppError) Error() string {
	return e.Message
}

func NewInvalidRequest(message string) *AppError {
	return &AppError{HTTPStatus: http.StatusBadRequest, Code: CodeInvalidRequest, Message: message}
}

func NewUnauthorized(message string) *AppError {
	return &AppError{HTTPStatus: http.StatusUnauthorized, Code: CodeUnauthorized, Message: message}
}

func NewForbidden(message string) *AppError {
	return &AppError{HTTPStatus: http.StatusForbidden, Code: CodeForbidden, Message: message}
}

func NewNotFound(message string) *AppError {
	return &AppError{HTTPStatus: http.StatusNotFound, Code: CodeNotFound, Message: message}
}

func NewConflict(message string) *AppError {
	return &AppError{HTTPStatus: http.StatusConflict, Code: CodeConflict, Message: message}
}

func NewInternal(code int, message string) *AppError {
	return &AppError{HTTPStatus: http.StatusInternalServerError, Code: code, Message: message}
}
