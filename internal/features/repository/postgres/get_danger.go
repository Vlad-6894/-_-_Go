package postgres_repository

import (
	"context"
	"errors"
	"fmt"
	"security/internal/core/domains"
	core_errors "security/internal/core/errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (r *DangersPostgresRepository) GetDanger(
	ctx context.Context,
	md5 string,
) (domains.Danger, error) {
	ctxWithTime, cancel := context.WithTimeout(ctx, r.pool.GetTimeot())
	defer cancel()

	sqlReq := `
	SELECT id, version, name, md, danger_type, danger_level, created_at, updated_at FROM security.dangers
	WHERE md = $1;
	`

	row := r.pool.QueryRow(ctxWithTime, sqlReq, md5)

	var model DangerModel

	if err := row.Scan(
		&model.ID,
		&model.Version,
		&model.Name,
		&model.MD5,
		&model.DangerType,
		&model.DangerLevel,
		&model.CreatedAt,
		&model.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domains.Danger{}, fmt.Errorf("not found value: %v: %w", err, core_errors.ErrNotFound)
		}

		var pgErr *pgconn.PgError
		if !errors.As(err, &pgErr) {
			return domains.Danger{}, fmt.Errorf("fail database: %w", err)
		}

		switch pgErr.Code {
		case "23514":
			return domains.Danger{}, fmt.Errorf("dates is invalid: %v: %w", err, core_errors.ErrInvalidArgument)
		default:
			return domains.Danger{}, fmt.Errorf("databes fail: %w", err)
		}
	}

	danger := dangerDomainFromModel(model)

	return danger, nil
}
