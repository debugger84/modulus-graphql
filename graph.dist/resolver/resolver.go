package resolver

//go:generate go run github.com/99designs/gqlgen generate

import (
	application "github.com/debugger84/modulus-application"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	logger application.Logger
}

func NewResolver(
	logger application.Logger,
) *Resolver {
	return &Resolver{
		logger: logger,
	}
}
