package services

import (
	"time"

	"github.com/eyoba-bisru/overtime-backend/internal/models"
	"github.com/eyoba-bisru/overtime-backend/internal/repository"
	"github.com/google/uuid"
)

func CreateOvertimeService(input models.Overtime) (uuid.UUID, error) {

	input.Status = models.OvertimePending
	input.ID = uuid.New()
	date, err := time.Parse("2006-01-02", input.Date)
	input.Date = date.Format("2006-01-02")
	if err != nil {
		return uuid.Nil, err
	}

	data, err := repository.CreateOvertimeRepo(&input)
	if err != nil {
		return uuid.Nil, err
	}
	return data, nil
}
