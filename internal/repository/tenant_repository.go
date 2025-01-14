package repository

import (
	"context"
	"time"

	"tenant/internal/model"
	"tenant/pkg/derrors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TenantRepository interface {
	CreateTenant(ctx context.Context, tenant *model.Tenant) error
	SoftDeleteTenant(ctx context.Context, clientID string) error
	GetTenantByClientID(ctx context.Context, clientID string) (*model.Tenant, error)
}

type tenantRepository struct {
	db *pgxpool.Pool
}

func NewTenantRepository(db *pgxpool.Pool) TenantRepository {
	return &tenantRepository{db: db}
}

func (r *tenantRepository) CreateTenant(ctx context.Context, tenant *model.Tenant) error {
	query := `
        INSERT INTO tenants (client_id, name, created_at, updated_at)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `
	err := r.db.QueryRow(ctx, query, tenant.ClientID, tenant.Name, time.Now(), time.Now()).Scan(&tenant.ID)
	return err
}

func (r *tenantRepository) SoftDeleteTenant(ctx context.Context, clientID string) error {
	query := `
        UPDATE tenants
        SET deleted_at = $1, updated_at = $2
        WHERE client_id = $3 AND deleted_at IS NULL
    `
	cmdTag, err := r.db.Exec(ctx, query, time.Now(), time.Now(), clientID)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return derrors.New(derrors.NotFound, "tenant not found or already deleted")
	}
	return nil
}

func (r *tenantRepository) GetTenantByClientID(ctx context.Context, clientID string) (*model.Tenant, error) {
	query := `
        SELECT id, client_id, name, created_at, updated_at
        FROM tenants
        WHERE client_id = $1 AND deleted_at IS NULL
    `
	var tenant model.Tenant
	err := r.db.QueryRow(ctx, query, clientID).Scan(
		&tenant.ID,
		&tenant.ClientID,
		&tenant.Name,
		&tenant.CreatedAt,
		&tenant.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}
