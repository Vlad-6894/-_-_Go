package postgres_repository

import (
	core_postgres_pool "security/internal/core/repository/postgres/pool"
	"time"

	"github.com/google/uuid"
)

type DangersPostgresRepository struct {
	pool core_postgres_pool.Pool
}

type DangerModel struct {
	ID          uuid.UUID
	Version     int
	Name        string
	MD5         string
	DangerType  string
	DangerLevel string
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}

func NewDangersPostgresRepository(
	pool core_postgres_pool.Pool,
) *DangersPostgresRepository {
	return &DangersPostgresRepository{
		pool: pool,
	}
}
