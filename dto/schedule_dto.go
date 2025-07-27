package dto

import (
	"time"

	"github.com/google/uuid"
)

type (
	ScheduleFullResponse struct {
		ScheduleID uuid.UUID `json:"schedule_id"`
		DayOfWeek  int       `json:"day_of_week"`
		DayName    string    `json:"day_name"`
		OpenTime   string    `json:"open_time"`
		CloseTime  string    `json:"close_time"`
		Field      FieldCompactResponse
	}

	ScheduleResponse struct {
		ScheduleID uuid.UUID `json:"schedule_id"`
		DayOfWeek  int       `json:"day_of_week"`
		DayName    string    `json:"day_name"`
		OpenTime   string    `json:"open_time"`
		CloseTime  string    `json:"close_time"`
		FieldID    uuid.UUID `json:"field_id"`
	}

	FieldCompactResponse struct {
		FieldID      uuid.UUID               `json:"field_id"`
		FieldName    string                  `json:"field_name"`
		FieldAddress string                  `json:"field_address"`
		FieldPrice   int                     `json:"field_price"`
		FieldImage   string                  `json:"field_image"`
		Category     CategoryCompactResponse `json:"category"`
	}

	CreateScheduleRequest struct {
		FieldID   string    `json:"field_id" binding:"required"`
		DayOfWeek int       `json:"day_of_week" binding:"required,min=1,max=7"`
		OpenTime  time.Time `json:"open_time" binding:"required"`
		CloseTime time.Time `json:"close_time" binding:"required"`
	}

	UpdateScheduleRequest struct {
		ScheduleID string     `json:"-"`        
		FieldID    *string    `json:"field_id"`
		DayOfWeek  *int       `json:"day_of_week"`
		OpenTime   *time.Time `json:"open_time"`
		CloseTime  *time.Time `json:"close_time"`
	}

	DeleteScheduleRequest struct {
		ScheduleID string `json:"-"`
	}
)
