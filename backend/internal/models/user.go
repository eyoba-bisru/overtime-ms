package models

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	Admin     Role = "admin"
	Checker   Role = "checker"
	Approver  Role = "approver"
	Applicant Role = "applicant"
	Finance   Role = "finance"
)

type User struct {
	Base
	Email           string     `json:"email" db:"email"`
	Name            string     `json:"name" db:"name"`
	PasswordHash    string     `json:"-" db:"password_hash"`
	Role                Role       `json:"role" db:"role"`
	DepartmentID        uuid.UUID  `json:"department_id" db:"department_id"`
	Department          *Department `json:"department,omitempty"`
	IsBlocked           bool       `json:"is_blocked" db:"is_blocked"`
	EmailVerified       bool       `json:"email_verified" db:"email_verified"`
	EmailVerifiedAt     *time.Time `json:"email_verified_at,omitempty" db:"email_verified_at"`
	ForcePasswordChange bool       `json:"force_password_change" db:"force_password_change"`
	LastLoginAt         *time.Time `json:"last_login_at,omitempty" db:"last_login_at"`
}
