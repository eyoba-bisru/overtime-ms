package models

import (
	"time"
)

type Role string

const (
	Admin     Role = "admin"
	Checker   Role = "checker"
	Approver  Role = "approver"
	Applicant Role = "applicant"
)

type User struct {
	Base
	Email           string     `json:"email" db:"email"`
	Name            string     `json:"name" db:"name"`
	PasswordHash    string     `json:"-" db:"password_hash"`
	Role            Role       `json:"role" db:"role"`
	IsBlocked       bool       `json:"is_blocked" db:"is_blocked"`
	EmailVerified   bool       `json:"email_verified" db:"email_verified"`
	EmailVerifiedAt *time.Time `json:"email_verified_at,omitempty" db:"email_verified_at"`
	LastLoginAt     *time.Time `json:"last_login_at,omitempty" db:"last_login_at"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}
