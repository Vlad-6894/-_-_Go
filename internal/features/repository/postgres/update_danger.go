package postgres_repository

import (
	"context"
	"errors"
	"fmt"
	"security/internal/core/domains"
	core_errors "security/internal/core/errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (r *DangersPostgresRepository) UpdateDanger(
	ctx context.Context,
	danger domains.Danger,
	md5 string,
) (domains.Danger, error) {
	ctxWithTime, cancel := context.WithTimeout(ctx, r.pool.GetTimeot())
	defer cancel()

	sqlReq := `
	UPDATE security.dangers
	SET name = $1, md = $2, danger_type = $3, danger_level = $4, version = version + 1, updated_at = $5
	WHERE md = $6
	RETURNING id, version, name, md, danger_type, danger_level, created_at, updated_at;
	`

	row := r.pool.QueryRow(
		ctxWithTime,
		sqlReq,
		danger.Name,
		danger.MD5,
		danger.DangerType,
		danger.DangerLevel,
		time.Now(),
		md5,
	)

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

	newDanger := dangerDomainFromModel(model)

	return newDanger, nil
}
