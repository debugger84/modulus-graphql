package graphql

import (
	application "github.com/debugger84/modulus-application"
	"go.uber.org/dig"
)

type ModuleConfig struct {
	apiUrl          string
	playgroundUrl   string
	complexityLimit int
	container       *dig.Container
}

func (s *ModuleConfig) SetContainer(container *dig.Container) {
	s.container = container
}

func (s *ModuleConfig) ModuleRoutes() []application.RouteInfo {
	var dp *GraphQlServer
	err := s.container.Invoke(func(dep *GraphQlServer) error {
		dp = dep
		return nil
	})
	if err != nil {
		panic("Cannot init GraphQL server:" + err.Error() + ". Try to copy graph.dist from this module to your projects and generate resolvers")
	}
	routes := application.NewRoutes()
	routes.Post(
		s.apiUrl,
		dp.ApiHandler(),
	)
	routes.Get(
		s.playgroundUrl,
		dp.PlaygroundHandler("Playground", s.apiUrl),
	)
	return routes.GetRoutesInfo()
}

func NewModuleConfig() *ModuleConfig {
	return &ModuleConfig{}
}

func (s *ModuleConfig) ProvidedServices() []interface{} {
	return []interface{}{
		NewGraphQLServer,
		func() *ModuleConfig {
			return s
		},
	}
}

func (s *ModuleConfig) InitConfig(config application.Config) error {
	if s.playgroundUrl == "" {
		s.playgroundUrl = config.GetEnv("GQL_PLAYGROUND_URL")
	}

	if s.complexityLimit == 0 {
		s.complexityLimit = config.GetEnvAsInt("GQL_COMPLEXITY_LIMIT")
	}
	if s.apiUrl == "" {
		s.apiUrl = config.GetEnv("GQL_API_URL")
	}

	return nil
}

func (s *ModuleConfig) OnStart() error {
	return nil
}
