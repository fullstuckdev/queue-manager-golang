package repository

import (
	"context"
	"errors"

	"queue-manager/internal/models"

	"gorm.io/gorm"
)

type TenantRepository struct {
	db *gorm.DB
}

func NewTenantRepository(db *gorm.DB) *TenantRepository {
	return &TenantRepository{db: db}
}

func (r *TenantRepository) Create(ctx context.Context, tenant *models.CreateTenantRequest) (*models.Tenant, error) {
	newTenant := &models.Tenant{
		ClientID: tenant.ClientID,
		Name:     tenant.Name,
	}
	
	result := r.db.WithContext(ctx).Create(newTenant)
	if result.Error != nil {
		return nil, result.Error
	}

	return newTenant, nil
}

func (r *TenantRepository) Delete(ctx context.Context, clientID string) error {
	result := r.db.WithContext(ctx).
		Where("client_id = ? AND deleted_at IS NULL", clientID).
		Delete(&models.Tenant{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("tenant not found")
	}

	return nil
}

func (r *TenantRepository) GetByClientID(ctx context.Context, clientID string) (*models.Tenant, error) {
	var tenant models.Tenant
	result := r.db.WithContext(ctx).
		Where("client_id = ? AND deleted_at IS NULL", clientID).
		First(&tenant)

	if result.Error != nil {
		return nil, result.Error
	}

	return &tenant, nil
} 