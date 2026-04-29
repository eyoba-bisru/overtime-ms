package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrUnauthorized      = errors.New("unauthorized action")
	ErrNotFound          = errors.New("record not found")
	ErrInvalidTransition = errors.New("invalid status transition")
	ErrImmutable         = errors.New("record is immutable and cannot be modified")
)

// APIResponse is the standard response envelope for all API endpoints.
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Meta    *PaginationMeta `json:"meta,omitempty"`
}

// PaginationMeta holds pagination metadata for list endpoints.
type PaginationMeta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int64 `json:"total_pages"`
}

type Base struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
