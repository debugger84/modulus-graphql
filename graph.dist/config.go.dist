package graph

import (
	"boilerplate/internal/graph/generated"
	"boilerplate/internal/graph/resolver"
	"github.com/99designs/gqlgen/graphql"
	"go.uber.org/dig"
)

type ModuleConfig struct {
	container *dig.Container
}

func NewModuleConfig() *ModuleConfig {
	return &ModuleConfig{}
}

func (s *ModuleConfig) ProvidedServices() []interface{} {
	return []interface{}{
		resolver.NewResolver,
		func(resolverObj *resolver.Resolver) graphql.ExecutableSchema {
			c := generated.Config{Resolvers: resolverObj}

			return generated.NewExecutableSchema(c)
		},
	}
}
