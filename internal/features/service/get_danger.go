package service

import (
	"context"
	"fmt"
	"security/internal/core/domains"
)

func (s *DangerService) GetDanger(
	ctx context.Context,
	md5 string,
) (domains.Danger, error) {
	danger, err := s.repository.GetDanger(ctx, md5)
	if err != nil {
		return domains.Danger{}, fmt.Errorf("fail get danger: %w", err)
	}

	return danger, nil
}
