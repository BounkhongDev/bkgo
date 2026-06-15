package errs

import "net/http"

// AppError is a structured error with an HTTP status code and error code string.
// Data carries optional structured payload (e.g. validation field errors).
type AppError struct {
	Status  int    `json:"-"`
	Code    string `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// Error implements the error interface on a value receiver so both
// AppError and *AppError satisfy error without risk of nil-pointer dereference.
func (e AppError) Error() string { return e.Message }

// New creates a new *AppError with a given HTTP status, error code, and message.
func New(status int, code, message string) *AppError {
	return &AppError{Status: status, Code: code, Message: message}
}

// IsAppError checks if err is an *AppError or AppError and returns a pointer to it.
func IsAppError(err error) (*AppError, bool) {
	if ae, ok := err.(*AppError); ok {
		return ae, true
	}
	if ae, ok := err.(AppError); ok {
		return &ae, true
	}
	return nil, false
}

// Predefined sentinels are value types (not pointers) so external callers
// cannot mutate the global state — only the constructors below return *AppError.
var (
	ErrNotFound      = AppError{Status: http.StatusNotFound, Code: "NOT_FOUND", Message: "resource not found"}
	ErrUnauthorized  = AppError{Status: http.StatusUnauthorized, Code: "UNAUTHORIZED", Message: "unauthorized"}
	ErrForbidden     = AppError{Status: http.StatusForbidden, Code: "FORBIDDEN", Message: "forbidden"}
	ErrBadRequest    = AppError{Status: http.StatusBadRequest, Code: "BAD_REQUEST", Message: "bad request"}
	ErrConflict      = AppError{Status: http.StatusConflict, Code: "CONFLICT", Message: "resource already exists"}
	ErrInternal      = AppError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "internal server error"}
	ErrUnprocessable = AppError{Status: http.StatusUnprocessableEntity, Code: "UNPROCESSABLE", Message: "unprocessable entity"}
)

// Constructor helpers return *AppError for custom messages.

func NotFound(msg string) *AppError {
	return New(http.StatusNotFound, "NOT_FOUND", msg)
}

func Unauthorized(msg string) *AppError {
	return New(http.StatusUnauthorized, "UNAUTHORIZED", msg)
}

func Forbidden(msg string) *AppError {
	return New(http.StatusForbidden, "FORBIDDEN", msg)
}

func BadRequest(msg string) *AppError {
	return New(http.StatusBadRequest, "BAD_REQUEST", msg)
}

func Conflict(msg string) *AppError {
	return New(http.StatusConflict, "CONFLICT", msg)
}

func Internal(msg string) *AppError {
	return New(http.StatusInternalServerError, "INTERNAL_ERROR", msg)
}

func Unprocessable(msg string) *AppError {
	return New(http.StatusUnprocessableEntity, "UNPROCESSABLE", msg)
}

// UnprocessableFields returns a 422 error carrying structured field-level errors.
func UnprocessableFields(msg string, fields any) *AppError {
	return &AppError{Status: http.StatusUnprocessableEntity, Code: "UNPROCESSABLE", Message: msg, Data: fields}
}
