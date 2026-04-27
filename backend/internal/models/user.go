package models

import "time"

type Role string

const (
	Admin     Role = "admin"
	Checker   Role = "checker"
	Approver  Role = "approver"
	Applicant Role = "applicant"
)

type User struct {
	Base            `json:",inline"`
	Email           string    `json:"email"`
	Name            string    `json:"name"`
	PasswordHash    string    `json:"password_hash"`
	Role            Role      `json:"role"`
	IsBlocked       bool      `json:"is_blocked"`
	LastLoginAt     time.Time `json:"last_login_at"`
	EmailVerified   bool      `json:"email_verified"`
	EmailVerifiedAt time.Time `json:"email_verified_at"`
	DeletedAt       time.Time `json:"deleted_at"`
}

type UserResponse struct {
	Base
	Email           string    `json:"email"`
	Name            string    `json:"name"`
	Role            Role      `json:"role"`
	IsBlocked       bool      `json:"is_blocked"`
	LastLoginAt     time.Time `json:"last_login_at"`
	EmailVerified   bool      `json:"email_verified"`
	EmailVerifiedAt time.Time `json:"email_verified_at"`
	DeletedAt       time.Time `json:"deleted_at"`
}
