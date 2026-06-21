package v1

import (
	"net/http"
)

func RouterGroup(
	mux *http.ServeMux,
	prefix string,
	middlewares []func(http.Handler) http.Handler,
	init func(register func(method, pattern string, handler http.HandlerFunc)),
) {
	register := func(method, pattern string, handler http.HandlerFunc) {
		var finalHandler http.Handler = handler
		if middlewares != nil {
			for i := len(middlewares) - 1; i >= 0; i-- {
				finalHandler = middlewares[i](finalHandler)
			}
		}
		finalPattern := method + " " + prefix + pattern
		mux.Handle(finalPattern, finalHandler)
	}
	init(register)
}
