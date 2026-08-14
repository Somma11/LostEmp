package apperror

import "fmt"

type AppError struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *AppError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func New(status int, code, message string) *AppError {
	return &AppError{Status: status, Code: code, Message: message}
}

func BadRequest(message string) *AppError {
	return New(400, "bad_request", message)
}

func Unauthorized(message string) *AppError {
	return New(401, "unauthorized", message)
}

func Forbidden(message string) *AppError {
	return New(403, "forbidden", message)
}

func NotFound(message string) *AppError {
	return New(404, "not_found", message)
}

func Conflict(message string) *AppError {
	return New(409, "conflict", message)
}

func Internal(message string) *AppError {
	return New(500, "internal_error", message)
}
