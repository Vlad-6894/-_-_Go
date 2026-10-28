package security_http_transport

import (
	"context"
	"security/internal/core/domains"
	core_http_server "security/internal/core/transport/http/server"
	"time"

	"github.com/google/uuid"
)

type HTTPSecurityTransport struct {
	service SecurityService
}

type SecurityService interface {
	CreateDanger(
		ctx context.Context,
		danger domains.Danger,
	) (domains.Danger, error)

	GetDanger(
		ctx context.Context,
		md5 string,
	) (domains.Danger, error)

	PutDanger(
		ctx context.Context,
		md5 string,
		put domains.DangerPut,
	) (domains.Danger, error)
}

type CreateDangerRequestDTO struct {
	Name        string `json:"name"`
	MD5         string `json:"md5"`
	DangerType  string `json:"danger_type"`
	DangerLevel string `json:"danger_level"`
}

type DangerResponseDTO struct {
	ID          uuid.UUID  `json:"id"`
	Version     int        `json:"version"`
	Name        string     `json:"name"`
	MD5         string     `json:"md5"`
	DangerType  string     `json:"danger_type"`
	DangerLevel string     `json:"danger_level"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

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

		{
			Path:    "info/{md5}",
			Method:  "GET",
			Handler: h.HandleGetDanger,
		},

		{
			Path:    "info/{md5}",
			Method:  "PUT",
			Handler: h.HandleUpdateDanger,
		},
	}

	return routes
}
