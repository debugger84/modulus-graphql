package graphql

import (
	"context"
	"errors"
	application "github.com/debugger84/modulus-application"
	"html/template"
	"net/http"
	"regexp"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/gorilla/websocket"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

const GraphQlPanic application.ActionErrorKey = `GraphQlPanic`

type GraphQlServer struct {
	config     *ModuleConfig
	logger     application.Logger
	corsRegexp *regexp.Regexp
	srv        *handler.Server
}

func NewGraphQLServer(
	config *ModuleConfig,
	logger application.Logger,
	es graphql.ExecutableSchema,
) *GraphQlServer {
	gQlServer := &GraphQlServer{
		config: config,
		logger: logger,
	}
	gQlServer.initServer(es)

	return gQlServer
}

func (f *GraphQlServer) initServer(es graphql.ExecutableSchema) {
	//loadersMutex := new(sync.Mutex)
	var mb int64 = 1 << 20

	srv := handler.New(es)

	srv.AddTransport(
		transport.Websocket{
			KeepAlivePingInterval: 10 * time.Second,
			Upgrader: websocket.Upgrader{
				CheckOrigin: func(r *http.Request) bool {
					return true
				},
			},
			InitFunc: func(ctx context.Context, initPayload transport.InitPayload) (context.Context, error) {
				//token := initPayload.Authorization()
				//reg := regexp.MustCompile(AuthorizationRegexp)
				//ctx = f.addCurrentUserToContext(
				//	f.authenticator,
				//	ctx,
				//	reg,
				//	token,
				//)
				//loaders := graph.NewDataLoaderContainer(ctx, container, loadersMutex)
				//ctx = graph.SetDataLoaders(ctx, loaders)

				return ctx, nil
			},
		},
	)
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(
		transport.MultipartForm{
			MaxUploadSize: mb * 101,
			MaxMemory:     mb * 151,
		},
	)

	srv.SetQueryCache(lru.New(1000))
	srv.Use(extension.AutomaticPersistedQuery{Cache: lru.New(1000)})

	srv.Use(extension.Introspection{})
	srv.Use(extension.FixedComplexityLimit(f.config.complexityLimit))

	srv.SetRecoverFunc(
		func(ctx context.Context, err interface{}) error {
			if tErr, ok := err.(error); ok {
				f.logger.Error(ctx, tErr.Error())
			}

			return errors.New("GraphQl panic")
		},
	)

	srv.SetErrorPresenter(
		func(ctx context.Context, e error) *gqlerror.Error {
			if gqlErr, ok := e.(*gqlerror.Error); ok {
				if iErr, ok := gqlErr.Unwrap().(ExtendedError); ok {
					return ConvertToGraphQlError(ctx, iErr)
				}
			}

			return graphql.DefaultErrorPresenter(ctx, e)
		},
	)

	f.srv = srv
}

func (f GraphQlServer) GetServer() *handler.Server {
	return f.srv
}

func (f *GraphQlServer) ApiHandler() http.HandlerFunc {
	handlerFunc := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// and call the next with our new context
			r = r.WithContext(ctx)

			f.srv.ServeHTTP(w, r)
		},
	)

	return handlerFunc
}

func (f *GraphQlServer) PlaygroundHandler(title string, endpoint string) http.HandlerFunc {
	var page = template.Must(template.New("graphql-playground").Parse(`<!DOCTYPE html>
<html>
<head>
	<meta charset=utf-8/>
	<meta name="viewport" content="user-scalable=no, initial-scale=1.0, minimum-scale=1.0, maximum-scale=1.0, minimal-ui">
	<link rel="shortcut icon" href="https://graphcool-playground.netlify.com/favicon.png">
	<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/graphql-playground-react@{{ .version }}/build/static/css/index.css"
		integrity="{{ .cssSRI }}" crossorigin="anonymous"/>
	<link rel="shortcut icon" href="https://cdn.jsdelivr.net/npm/graphql-playground-react@{{ .version }}/build/favicon.png"
		integrity="{{ .faviconSRI }}" crossorigin="anonymous"/>
	<script src="https://cdn.jsdelivr.net/npm/graphql-playground-react@{{ .version }}/build/static/js/middleware.js"
		integrity="{{ .jsSRI }}" crossorigin="anonymous"></script>
	<title>{{.title}}</title>
</head>
<body>
<style type="text/css">
	html { font-family: "Open Sans", sans-serif; overflow: hidden; }
	body { margin: 0; background: #172a3a; }
</style>
<div id="root"/>
<script type="text/javascript">
	window.addEventListener('load', function (event) {
		const root = document.getElementById('root');
		root.classList.add('playgroundIn');
		const wsProto = location.protocol == 'https:' ? 'wss:' : 'ws:'
		GraphQLPlayground.init(root, {
			endpoint: location.protocol + '//' + location.host + '{{.endpoint}}',
			subscriptionsEndpoint: wsProto + '//' + location.host + '{{.endpoint }}',
           shareEnabled: true,
			settings: {
				'request.credentials': 'same-origin'
			}
		})
	})
</script>
</body>
</html>
`))

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html")
		err := page.Execute(w, map[string]string{
			"title":      title,
			"endpoint":   endpoint,
			"version":    "1.7.28",
			"cssSRI":     "sha256-dKnNLEFwKSVFpkpjRWe+o/jQDM6n/JsvQ0J3l5Dk3fc=",
			"faviconSRI": "sha256-GhTyE+McTU79R4+pRO6ih+4TfsTOrpPwD8ReKFzb3PM=",
			"jsSRI":      "sha256-VVwEZwxs4qS5W7E+/9nXINYgr/BJRWKOi/rTMUdmmWg=",
		})
		if err != nil {
			panic(err)
		}
	}
}
