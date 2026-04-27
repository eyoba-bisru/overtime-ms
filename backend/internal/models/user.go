package models

type Role string

const (
	Admin     Role = "admin"
	Checker   Role = "checker"
	Approver  Role = "approver"
	Applicant Role = "applicant"
)

type User struct {
	Base      `json:",inline"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Password  string `json:"password"`
	Role      Role   `json:"role"`
	IsBlocked bool   `json:"is_blocked"`
}
