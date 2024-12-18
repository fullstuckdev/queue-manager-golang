package models

import (
    "time"
)

type Tenant struct {
    ID        int64      `json:"id"`
    ClientID  string     `json:"client_id"`
    Name      string     `json:"name"`
    CreatedAt time.Time  `json:"created_at"`
    UpdatedAt time.Time  `json:"updated_at"`
    DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type CreateTenantRequest struct {
    ClientID string `json:"client_id" validate:"required"`
    Name     string `json:"name" validate:"required"`
}

type ProcessPayloadRequest struct {
    ClientID string      `json:"client_id" validate:"required"`
    Payload  interface{} `json:"payload" validate:"required"`
} 