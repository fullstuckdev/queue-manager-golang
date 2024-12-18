package repository

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
	"queue-manager/internal/models"
)

func setupTestDB(t *testing.T) *pgxpool.Pool {
	dbpool, err := pgxpool.Connect(context.Background(), "postgres://postgres:postgres@localhost:5432/tenants_test")
	if err != nil {
		t.Fatalf("Unable to connect to database: %v", err)
	}
	return dbpool
}

func TestTenantRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewTenantRepository(db)

	tests := []struct {
		name    string
		tenant  *models.CreateTenantRequest
		wantErr bool
	}{
		{
			name: "Valid tenant",
			tenant: &models.CreateTenantRequest{
				ClientID: "test-client",
				Name:     "Test Tenant",
			},
			wantErr: false,
		},
		{
			name: "Duplicate client ID",
			tenant: &models.CreateTenantRequest{
				ClientID: "test-client",
				Name:     "Duplicate Tenant",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.Create(context.Background(), tt.tenant)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.tenant.ClientID, result.ClientID)
			assert.Equal(t, tt.tenant.Name, result.Name)
		})
	}
} 