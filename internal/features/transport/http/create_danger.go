package security_http_transport

import (
	"encoding/json"
	"net/http"
	"security/internal/core/domains"
	core_logger "security/internal/core/logger"
	core_http_response "security/internal/core/transport/http/response"
)

func (h *HTTPSecurityTransport) HandleCreateDanger(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	var request CreateDangerRequestDTO

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		responseHandler.ErrorResponse(err, "Fail to decode create danger DTO!")
		return
	}

	danger := dangerDomainFromCreateDTO(request)

	newDanger, err := h.service.CreateDanger(ctx, danger)
	if err != nil {
		responseHandler.ErrorResponse(err, "fail to create danger!")
		return
	}

	response := dangerRequestDtoFromDomain(newDanger)

	responseHandler.ToJSONResponse(response, http.StatusCreated)
}

func dangerDomainFromCreateDTO(dto CreateDangerRequestDTO) domains.Danger {
	return domains.NewUnitializedDanger(dto.Name, dto.MD5, dto.DangerType, dto.DangerLevel)
}

func dangerRequestDtoFromDomain(danger domains.Danger) DangerResponseDTO {
	return DangerResponseDTO{
		ID:          danger.ID,
		Version:     danger.Version,
		Name:        danger.Name,
		MD5:         danger.MD5,
		DangerType:  danger.DangerType,
		DangerLevel: danger.DangerLevel,
		CreatedAt:   danger.CreatedAt,
		UpdatedAt:   danger.UpdatedAt,
	}
}
