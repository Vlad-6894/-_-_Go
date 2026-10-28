package postgres_repository

import (
	"context"
	"errors"
	"fmt"
	"security/internal/core/domains"
	core_errors "security/internal/core/errors"

	"github.com/jackc/pgx/v5/pgconn"
)

func (r *DangersPostgresRepository) CreateDanger(
	ctx context.Context,
	danger domains.Danger,
) (domains.Danger, error) {
	ctxWithTime, cancel := context.WithTimeout(ctx, r.pool.GetTimeot())
	defer cancel()

	sqlReq := `
	INSERT INTO security.dangers (id, name, md, danger_type, danger_level, created_at, updated_at)
	VALUES($1,$2,$3,$4,$5,$6,$7)
	RETURNING id, version, name, md, danger_type, danger_level, created_at, updated_at;
	`

	row := r.pool.QueryRow(
		ctxWithTime,
		sqlReq,
		danger.ID,
		danger.Name,
		danger.MD5,
		danger.DangerType,
		danger.DangerLevel,
		danger.CreatedAt,
		danger.UpdatedAt,
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
		var pgErr *pgconn.PgError
		if !errors.As(err, &pgErr) {
			return domains.Danger{}, fmt.Errorf("fail database: %w", err)
		}

		switch pgErr.Code {
		case "23505":
			return domains.Danger{}, fmt.Errorf("danger already exists: %v: %w", err, core_errors.ErrConflict)

		case "23514":
			return domains.Danger{}, fmt.Errorf("dates is invalid: %v: %w", err, core_errors.ErrInvalidArgument)
		default:
			return domains.Danger{}, fmt.Errorf("databes fail: %w", err)
		}
	}

	newDanger := dangerDomainFromModel(model)

	return newDanger, nil
}

func dangerDomainFromModel(model DangerModel) domains.Danger {
	return domains.NewDanger(
		model.ID,
		model.Version,
		model.Name,
		model.MD5,
		model.DangerType,
		model.DangerLevel,
		model.CreatedAt,
		model.UpdatedAt,
	)
}
