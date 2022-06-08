package graphql

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	application "github.com/debugger84/modulus-application"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func FromValidationErr(
	ctx context.Context,
	errors []application.ValidationError,
) *gqlerror.Error {
	invalidInputs := make([]map[string]string, len(errors))
	msg := "Unknown error"
	for i, validationError := range errors {
		if i == 0 {
			msg = validationError.Error()
		}
		invalidInputs[i] = map[string]string{
			"field":   validationError.Field,
			"message": validationError.Error(),
			"id":      string(validationError.Identifier),
		}
	}

	return &gqlerror.Error{
		Message: errors[0].Err,
		Path:    graphql.GetPath(ctx),
		Extensions: map[string]interface{}{
			"code":          "validation.error",
			"message":       msg,
			"invalidInputs": invalidInputs,
		},
	}
}

func FromCommonErr(
	ctx context.Context,
	error *application.CommonError,
) *gqlerror.Error {
	if error == nil {
		return nil
	}

	return &gqlerror.Error{
		Message: error.Err,
		Path:    graphql.GetPath(ctx),
		Extensions: map[string]interface{}{
			"code":    error.Identifier,
			"message": error.Err,
		},
	}
}
