package core_http_server

import (
	"fmt"
	"net/http"
)

type ApiRouter struct {
	*http.ServeMux
}

func NewApiVersionRouter() *ApiRouter {
	return &ApiRouter{
		ServeMux: http.NewServeMux(),
	}
}

func (r *ApiRouter) RegisterRoutes(routes ...Route) {
	for _, route := range routes {
		prefix := fmt.Sprintf("%s /%s", route.Method, route.Path)

		r.ServeMux.HandleFunc(prefix, route.Handler)
	}
}
