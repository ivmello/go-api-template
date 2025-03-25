package errors

import (
	"errors"
	"fmt"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ApplicationError represents an application error
type ApplicationError struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

// Error returns the error message
func (e *ApplicationError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap returns the wrapped error
func (e *ApplicationError) Unwrap() error {
	return e.Err
}

// NewBadRequestError creates a new bad request error
func NewBadRequestError(message string, err error) *ApplicationError {
	return &ApplicationError{
		Code:    http.StatusBadRequest,
		Message: message,
		Err:     err,
	}
}

// NewUnauthorizedError creates a new unauthorized error
func NewUnauthorizedError(message string, err error) *ApplicationError {
	return &ApplicationError{
		Code:    http.StatusUnauthorized,
		Message: message,
		Err:     err,
	}
}

// NewForbiddenError creates a new forbidden error
func NewForbiddenError(message string, err error) *ApplicationError {
	return &ApplicationError{
		Code:    http.StatusForbidden,
		Message: message,
		Err:     err,
	}
}

// NewNotFoundError creates a new not found error
func NewNotFoundError(message string, err error) *ApplicationError {
	return &ApplicationError{
		Code:    http.StatusNotFound,
		Message: message,
		Err:     err,
	}
}

// NewConflictError creates a new conflict error
func NewConflictError(message string, err error) *ApplicationError {
	return &ApplicationError{
		Code:    http.StatusConflict,
		Message: message,
		Err:     err,
	}
}

// NewInternalServerError creates a new internal server error
func NewInternalServerError(message string, err error) *ApplicationError {
	return &ApplicationError{
		Code:    http.StatusInternalServerError,
		Message: message,
		Err:     err,
	}
}

// ToGRPCError converts an ApplicationError to a gRPC error
func ToGRPCError(err error) error {
	var appErr *ApplicationError
	if errors.As(err, &appErr) {
		switch appErr.Code {
		case http.StatusBadRequest:
			return status.Error(codes.InvalidArgument, appErr.Message)
		case http.StatusUnauthorized:
			return status.Error(codes.Unauthenticated, appErr.Message)
		case http.StatusForbidden:
			return status.Error(codes.PermissionDenied, appErr.Message)
		case http.StatusNotFound:
			return status.Error(codes.NotFound, appErr.Message)
		case http.StatusConflict:
			return status.Error(codes.AlreadyExists, appErr.Message)
		default:
			return status.Error(codes.Internal, appErr.Message)
		}
	}
	return status.Error(codes.Unknown, err.Error())
}