package service

import (
	"context"
	"fmt"
	"security/internal/core/domains"
)

func (s *DangerService) CreateDanger(
	ctx context.Context,
	danger domains.Danger,
) (domains.Danger, error) {
	if err := danger.Validate(); err != nil {
		return domains.Danger{}, fmt.Errorf("fail validate danger: %w", err)
	}

	newDanger, err := s.repository.CreateDanger(ctx, danger)
	if err != nil {
		return domains.Danger{}, fmt.Errorf("fail create danger: %w", err)
	}

	return newDanger, nil
}
