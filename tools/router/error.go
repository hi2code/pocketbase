package router

import (
	"context"
	"database/sql"
	"errors"
	"io/fs"
	"net/http"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/tools/i18n"
	"github.com/pocketbase/pocketbase/tools/inflector"
)

// SafeErrorItem defines a common error interface for a printable public safe error.
type SafeErrorItem interface {
	// Code represents a fixed unique identifier of the error (usually used as translation key).
	Code() string

	// Error is the default English human readable error message that will be returned.
	Error() string
}

// SafeErrorParamsResolver defines an optional interface for specifying dynamic error parameters.
type SafeErrorParamsResolver interface {
	// Params defines a map with dynamic parameters to return as part of the public safe error view.
	Params() map[string]any
}

// SafeErrorResolver defines an error interface for resolving the public safe error fields.
type SafeErrorResolver interface {
	// Resolve allows modifying and returning a new public safe error data map.
	Resolve(errData map[string]any) any
}

// ApiError defines the struct for a basic api error response.
type ApiError struct {
	rawData any

	Data    map[string]any `json:"data"`
	Message string         `json:"message"`
	Status  int            `json:"status"`
}

// Error makes it compatible with the `error` interface.
func (e *ApiError) Error() string {
	return e.Message
}

// RawData returns the unformatted error data (could be an internal error, text, etc.)
func (e *ApiError) RawData() any {
	return e.rawData
}

// Is reports whether the current ApiError wraps the target.
func (e *ApiError) Is(target error) bool {
	err, ok := e.rawData.(error)
	if ok {
		return errors.Is(err, target)
	}

	apiErr, ok := target.(*ApiError)

	return ok && e == apiErr
}

// NewNotFoundError creates and returns 404 ApiError.
func NewNotFoundError(message string, rawErrData any) *ApiError {
	if message == "" {
		message = "The requested resource wasn't found."
	}

	return NewApiError(http.StatusNotFound, message, rawErrData)
}

// NewBadRequestError creates and returns 400 ApiError.
func NewBadRequestError(message string, rawErrData any) *ApiError {
	if message == "" {
		message = "Something went wrong while processing your request."
	}

	return NewApiError(http.StatusBadRequest, message, rawErrData)
}

// NewForbiddenError creates and returns 403 ApiError.
func NewForbiddenError(message string, rawErrData any) *ApiError {
	if message == "" {
		message = "You are not allowed to perform this request."
	}

	return NewApiError(http.StatusForbidden, message, rawErrData)
}

// NewUnauthorizedError creates and returns 401 ApiError.
func NewUnauthorizedError(message string, rawErrData any) *ApiError {
	if message == "" {
		message = "Missing or invalid authentication."
	}

	return NewApiError(http.StatusUnauthorized, message, rawErrData)
}

// NewInternalServerError creates and returns 500 ApiError.
func NewInternalServerError(message string, rawErrData any) *ApiError {
	if message == "" {
		message = "Something went wrong while processing your request."
	}

	return NewApiError(http.StatusInternalServerError, message, rawErrData)
}

func NewTooManyRequestsError(message string, rawErrData any) *ApiError {
	if message == "" {
		message = "Too Many Requests."
	}

	return NewApiError(http.StatusTooManyRequests, message, rawErrData)
}

// NewNotFoundErrorWithContext creates and returns 404 ApiError with i18n support.
func NewNotFoundErrorWithContext(ctx context.Context, message string, rawErrData any) *ApiError {
	if message == "" {
		message = i18n.T(ctx, "error_not_found")
	}

	return NewApiErrorWithContext(ctx, http.StatusNotFound, message, rawErrData)
}

// NewBadRequestErrorWithContext creates and returns 400 ApiError with i18n support.
func NewBadRequestErrorWithContext(ctx context.Context, message string, rawErrData any) *ApiError {
	if message == "" {
		message = i18n.T(ctx, "error_bad_request")
	}

	return NewApiErrorWithContext(ctx, http.StatusBadRequest, message, rawErrData)
}

// NewForbiddenErrorWithContext creates and returns 403 ApiError with i18n support.
func NewForbiddenErrorWithContext(ctx context.Context, message string, rawErrData any) *ApiError {
	if message == "" {
		message = i18n.T(ctx, "error_forbidden")
	}

	return NewApiErrorWithContext(ctx, http.StatusForbidden, message, rawErrData)
}

// NewUnauthorizedErrorWithContext creates and returns 401 ApiError with i18n support.
func NewUnauthorizedErrorWithContext(ctx context.Context, message string, rawErrData any) *ApiError {
	if message == "" {
		message = i18n.T(ctx, "error_unauthorized")
	}

	return NewApiErrorWithContext(ctx, http.StatusUnauthorized, message, rawErrData)
}

// NewInternalServerErrorWithContext creates and returns 500 ApiError with i18n support.
func NewInternalServerErrorWithContext(ctx context.Context, message string, rawErrData any) *ApiError {
	if message == "" {
		message = i18n.T(ctx, "error_internal_server")
	}

	return NewApiErrorWithContext(ctx, http.StatusInternalServerError, message, rawErrData)
}

// NewTooManyRequestsErrorWithContext creates and returns 429 ApiError with i18n support.
func NewTooManyRequestsErrorWithContext(ctx context.Context, message string, rawErrData any) *ApiError {
	if message == "" {
		message = i18n.T(ctx, "error_too_many_requests")
	}

	return NewApiErrorWithContext(ctx, http.StatusTooManyRequests, message, rawErrData)
}

// NewApiError creates and returns new normalized ApiError instance.
func NewApiError(status int, message string, rawErrData any) *ApiError {
	if message == "" {
		message = http.StatusText(status)
	}

	return &ApiError{
		rawData: rawErrData,
		Data:    safeErrorsData(rawErrData),
		Status:  status,
		Message: strings.TrimSpace(inflector.Sentenize(message)),
	}
}

// NewApiErrorWithContext creates and returns new normalized ApiError instance with i18n support.
func NewApiErrorWithContext(ctx context.Context, status int, message string, rawErrData any) *ApiError {
	if message == "" {
		// Use i18n for default messages based on status
		switch status {
		case http.StatusNotFound:
			message = i18n.T(ctx, "error_not_found")
		case http.StatusBadRequest:
			message = i18n.T(ctx, "error_bad_request")
		case http.StatusForbidden:
			message = i18n.T(ctx, "error_forbidden")
		case http.StatusUnauthorized:
			message = i18n.T(ctx, "error_unauthorized")
		case http.StatusInternalServerError:
			message = i18n.T(ctx, "error_internal_server")
		case http.StatusTooManyRequests:
			message = i18n.T(ctx, "error_too_many_requests")
		default:
			message = http.StatusText(status)
		}
	}

	return &ApiError{
		rawData: rawErrData,
		Data:    safeErrorsDataWithContext(ctx, rawErrData),
		Status:  status,
		Message: strings.TrimSpace(inflector.Sentenize(message)),
	}
}

// ToApiError wraps err into ApiError instance (if not already).
func ToApiError(err error) *ApiError {
	var apiErr *ApiError

	if !errors.As(err, &apiErr) {
		// no ApiError found -> assign a generic one
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, fs.ErrNotExist) {
			apiErr = NewNotFoundError("", err)
		} else {
			apiErr = NewBadRequestError("", err)
		}
	}

	return apiErr
}

// ToApiErrorWithContext wraps err into ApiError instance (if not already) with i18n support.
func ToApiErrorWithContext(ctx context.Context, err error) *ApiError {
	var apiErr *ApiError

	if !errors.As(err, &apiErr) {
		// no ApiError found -> assign a generic one
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, fs.ErrNotExist) {
			apiErr = NewNotFoundErrorWithContext(ctx, "", err)
		} else {
			apiErr = NewBadRequestErrorWithContext(ctx, "", err)
		}
	}

	return apiErr
}

// -------------------------------------------------------------------

func safeErrorsData(data any) map[string]any {
	switch v := data.(type) {
	case validation.Errors:
		return resolveSafeErrorsData(v)
	case error:
		validationErrors := validation.Errors{}
		if errors.As(v, &validationErrors) {
			return resolveSafeErrorsData(validationErrors)
		}
		return map[string]any{} // not nil to ensure that is json serialized as object
	case map[string]validation.Error:
		return resolveSafeErrorsData(v)
	case map[string]SafeErrorItem:
		return resolveSafeErrorsData(v)
	case map[string]error:
		return resolveSafeErrorsData(v)
	case map[string]string:
		return resolveSafeErrorsData(v)
	case map[string]any:
		return resolveSafeErrorsData(v)
	default:
		return map[string]any{} // not nil to ensure that is json serialized as object
	}
}

func resolveSafeErrorsData[T any](data map[string]T) map[string]any {
	result := map[string]any{}

	for name, err := range data {
		if isNestedError(err) {
			result[name] = safeErrorsData(err)
		} else {
			result[name] = resolveSafeErrorItem(err)
		}
	}

	return result
}

func isNestedError(err any) bool {
	switch err.(type) {
	case validation.Errors,
		map[string]validation.Error,
		map[string]SafeErrorItem,
		map[string]error,
		map[string]string,
		map[string]any:
		return true
	}

	return false
}

// resolveSafeErrorItem extracts from each validation error its
// public safe error code and message.
func resolveSafeErrorItem(err any) any {
	data := map[string]any{}

	if obj, ok := err.(SafeErrorItem); ok {
		// extract the specific error code and message
		data["code"] = obj.Code()
		data["message"] = inflector.Sentenize(obj.Error())
	} else {
		// fallback to the default public safe values
		data["code"] = "validation_invalid_value"
		data["message"] = "Invalid value."
	}

	if s, ok := err.(SafeErrorParamsResolver); ok {
		params := s.Params()
		if len(params) > 0 {
			data["params"] = params
		}
	}

	if s, ok := err.(SafeErrorResolver); ok {
		return s.Resolve(data)
	}

	return data
}

// -------------------------------------------------------------------
// i18n support functions

func safeErrorsDataWithContext(ctx context.Context, data any) map[string]any {
	switch v := data.(type) {
	case validation.Errors:
		return resolveSafeErrorsDataWithContext(ctx, v)
	case error:
		validationErrors := validation.Errors{}
		if errors.As(v, &validationErrors) {
			return resolveSafeErrorsDataWithContext(ctx, validationErrors)
		}
		return map[string]any{} // not nil to ensure that is json serialized as object
	case map[string]validation.Error:
		return resolveSafeErrorsDataWithContext(ctx, v)
	case map[string]SafeErrorItem:
		return resolveSafeErrorsDataWithContext(ctx, v)
	case map[string]error:
		return resolveSafeErrorsDataWithContext(ctx, v)
	case map[string]string:
		return resolveSafeErrorsDataWithContext(ctx, v)
	case map[string]any:
		return resolveSafeErrorsDataWithContext(ctx, v)
	default:
		return map[string]any{} // not nil to ensure that is json serialized as object
	}
}

func resolveSafeErrorsDataWithContext[T any](ctx context.Context, data map[string]T) map[string]any {
	result := map[string]any{}

	for name, err := range data {
		if isNestedError(err) {
			result[name] = safeErrorsDataWithContext(ctx, err)
		} else {
			result[name] = resolveSafeErrorItemWithContext(ctx, err)
		}
	}

	return result
}

// resolveSafeErrorItemWithContext extracts from each validation error its
// public safe error code and message with i18n support.
func resolveSafeErrorItemWithContext(ctx context.Context, err any) any {
	data := map[string]any{}

	if obj, ok := err.(SafeErrorItem); ok {
		// extract the specific error code and message
		code := obj.Code()
		data["code"] = code

		// Use i18n translation for the message
		defaultMsg := obj.Error()
		var params map[string]any

		if s, ok := err.(SafeErrorParamsResolver); ok {
			params = s.Params()
		}

		translatedMsg := i18n.TranslateError(ctx, code, defaultMsg, params)
		data["message"] = inflector.Sentenize(translatedMsg)
	} else {
		// fallback to the default public safe values
		data["code"] = "validation_invalid_value"
		data["message"] = i18n.T(ctx, "validation_invalid_value")
	}

	if s, ok := err.(SafeErrorParamsResolver); ok {
		params := s.Params()
		if len(params) > 0 {
			data["params"] = params
		}
	}

	if s, ok := err.(SafeErrorResolver); ok {
		return s.Resolve(data)
	}

	return data
}
