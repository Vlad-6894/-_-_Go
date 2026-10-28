package security_http_transport

import (
	"net/http"
	core_logger "security/internal/core/logger"
	core_http_response "security/internal/core/transport/http/response"
	core_http_utils "security/internal/core/transport/http/utils"
)

const (
	md5Key = "md5"
)

func (h *HTTPSecurityTransport) HandleGetDanger(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	md5, err := core_http_utils.GetStringPathValue(r, md5Key)
	if err != nil {
		responseHandler.ErrorResponse(err, "fail to get md5 from path!")
		return
	}

	danger, err := h.service.GetDanger(ctx, md5)
	if err != nil {
		responseHandler.ErrorResponse(err, "fail to get danger from service!!")
		return
	}

	response := dangerRequestDtoFromDomain(danger)

	responseHandler.ToJSONResponse(response, http.StatusOK)
}
