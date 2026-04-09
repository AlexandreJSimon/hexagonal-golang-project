package http

import (
	"net/http"

	"github.com/AlexandreJSimon/hexagonal-golang-project/internal/infra/http/middleware"
)

func Group(
	parent *http.ServeMux,
	prefix string,
	register func(mux *http.ServeMux),
	middlewares ...middleware.MiddlewareFunc,
) {
	subMux := http.NewServeMux()

	register(subMux)

	var handler http.Handler = subMux

	handler = middleware.Chain(middlewares...)(handler)

	parent.Handle(prefix+"/", http.StripPrefix(prefix, handler))
}
