package models

import (
	"github.com/google/uuid"
)

type OvertimeStatus string

const (
	OvertimePending  OvertimeStatus = "pending"
	OvertimeChecked  OvertimeStatus = "checked"
	OvertimeApproved OvertimeStatus = "approved"
	OvertimeRejected OvertimeStatus = "rejected"
)

type OvertimeProgram string

const (
	Night   OvertimeProgram = "night"
	Weekend OvertimeProgram = "weekend"
	Holiday OvertimeProgram = "holiday"
)

// Overtime represents overtime request entity
type Overtime struct {
	Base
	UserID         uuid.UUID       `json:"user_id" db:"user_id"`
	UserName       string          `json:"user_name" db:"user_name"`
	DepartmentID   uuid.UUID       `json:"department_id" db:"department_id"`
	DepartmentName string          `json:"department_name" db:"department_name"`
	Date           string          `json:"date" db:"date"`
	StartTime      string          `json:"start_time" db:"start_time"`
	EndTime        string          `json:"end_time" db:"end_time"`
	JobDone        string          `json:"job_done" db:"job_done"`
	Status         OvertimeStatus  `json:"status" db:"status"`
	Program        OvertimeProgram `json:"program" db:"program"`
	Duration       float64         `json:"duration" db:"duration"`
}
