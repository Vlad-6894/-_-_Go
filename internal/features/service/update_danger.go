package service

import (
	"context"
	"fmt"
	"security/internal/core/domains"
)

func (s *DangerService) PutDanger(
	ctx context.Context,
	md5 string,
	put domains.DangerPut,
) (domains.Danger, error) {
	danger, err := s.repository.GetDanger(ctx, md5)
	if err != nil {
		return domains.Danger{}, fmt.Errorf("fail get danger: %w", err)
	}

	if err := danger.ApplyPut(put); err != nil {
		return domains.Danger{}, fmt.Errorf("apply danger put error: %w", err)
	}

	putDanger, err := s.repository.UpdateDanger(ctx, danger, md5)
	if err != nil {
		return domains.Danger{}, fmt.Errorf("fail update danger: %w", err)
	}

	return putDanger, nil
}
