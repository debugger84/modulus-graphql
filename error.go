package graphql

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

type ExtendedError interface {
	Msg() string
	Status() string
	ErrCode() int
}

type ExtendedValidationError interface {
	ExtendedError
	Errors() map[string]string
	Extra() map[string]interface{}
}

type ClientError struct {
	msg     string
	status  string
	errCode int
}

func (m *ClientError) Msg() string {
	return m.msg
}

func (m *ClientError) Status() string {
	return m.status
}

func (m *ClientError) ErrCode() int {
	return m.errCode
}

func (m *ClientError) Error() string {
	return m.status
}

func CreateGraphQlError(
	ctx context.Context,
	msg string,
	status string,
	errCode int,
) *gqlerror.Error {
	return &gqlerror.Error{
		Path:    graphql.GetPath(ctx),
		Message: msg,
		Extensions: map[string]interface{}{
			"status":    status,
			"errorCode": errCode,
			"message":   msg,
		},
	}
}

func CreateGraphQlErrorWithPlaceholders(
	ctx context.Context,
	msg string,
	status string,
	errCode int,
	placeHolders map[string]interface{},
) *gqlerror.Error {
	return &gqlerror.Error{
		Path:    graphql.GetPath(ctx),
		Message: msg,
		Extensions: map[string]interface{}{
			"status":    status,
			"errorCode": errCode,
			"message":   msg,
			"details": map[string]interface{}{
				"placeholders": placeHolders,
			},
		},
	}
}

func ConvertToGraphQlError(
	ctx context.Context,
	err ExtendedError,
) *gqlerror.Error {
	return CreateGraphQlError(ctx, err.Msg(), err.Status(), err.ErrCode())
}

func ConvertToGraphQlBackendValidationError(
	ctx context.Context,
	err ExtendedValidationError,
) *gqlerror.Error {
	return &gqlerror.Error{
		Path:    graphql.GetPath(ctx),
		Message: err.Msg(),
		Extensions: map[string]interface{}{
			"status":    "validation.error",
			"errorCode": 400,
			"message":   err.Msg(),
			"errors":    err.Errors(),
			"extra":     err.Extra(),
		},
	}
}

func GetError(errs []error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}

	return nil
}
