package security_http_transport

import (
	"encoding/json"
	"net/http"
	"security/internal/core/domains"
	core_logger "security/internal/core/logger"
	core_http_response "security/internal/core/transport/http/response"
	core_http_types "security/internal/core/transport/http/types"
	core_http_utils "security/internal/core/transport/http/utils"
)

type PatchDangerRequest struct {
	Name        core_http_types.Put[string]
	MD5         core_http_types.Put[string]
	DangerType  core_http_types.Put[string]
	DangerLevel core_http_types.Put[string]
}

func (h *HTTPSecurityTransport) HandleUpdateDanger(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	md5, err := core_http_utils.GetStringPathValue(r, md5Key)
	if err != nil {
		responseHandler.ErrorResponse(err, "fail to get md5 from path!")
		return
	}

	var request PatchDangerRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		responseHandler.ErrorResponse(err, "Fail to decode update danger DTO!")
		return
	}

	put := dangerPutFromRequest(request)

	dangerDomain, err := h.service.PutDanger(ctx, md5, put)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to put danger")
		return
	}

	response := dangerRequestDtoFromDomain(dangerDomain)

	responseHandler.ToJSONResponse(response, http.StatusOK)
}

func dangerPutFromRequest(request PatchDangerRequest) domains.DangerPut {
	return domains.NewDangerPut(
		request.Name.ToDomain(),
		request.MD5.ToDomain(),
		request.DangerType.ToDomain(),
		request.DangerLevel.ToDomain(),
	)
}
