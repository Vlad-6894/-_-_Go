package security_http_transport

import core_http_server "security/internal/core/transport/http/server"

type HTTPSecurityTransport struct {
	service SecurityService
}

type SecurityService interface{}

func NewHTTPSecurityTransport(
	service SecurityService,
) *HTTPSecurityTransport {
	return &HTTPSecurityTransport{
		service: service,
	}
}

func (h *HTTPSecurityTransport) Routes() []core_http_server.Route {
	routes := []core_http_server.Route{
		{
			Path:    "info",
			Method:  "POST",
			Handler: h.HandleCreateDanger,
		},
	}

	return routes
}
