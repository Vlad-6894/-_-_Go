package security_http_transport

import (
	"net/http"
	core_logger "security/internal/core/logger"
	core_http_response "security/internal/core/transport/http/response"
)

func (h *HTTPSecurityTransport) HandleCreateDanger(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)
}
