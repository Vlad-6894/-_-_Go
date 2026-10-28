package service

import (
	"context"
	"security/internal/core/domains"
)

type DangerService struct {
	repository DangerRepository
}

//go:generate mockgen -source=service.go -destination=mocks/mock_repo.go -package=mocks
type DangerRepository interface {
	CreateDanger(
		ctx context.Context,
		danger domains.Danger,
	) (domains.Danger, error)

	GetDanger(
		ctx context.Context,
		md5 string,
	) (domains.Danger, error)

	UpdateDanger(
		ctx context.Context,
		danger domains.Danger,
		md5 string,
	) (domains.Danger, error)
}

func NewDangerService(
	repository DangerRepository,
) *DangerService {
	return &DangerService{
		repository: repository,
	}
}
