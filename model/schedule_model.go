package model

import (
	"time"

	"github.com/google/uuid"
)

type Schedule struct {
	ScheduleID uuid.UUID `gorm:"type:uuid;primaryKey;column:schedule_id"`
	FieldID    uuid.UUID `gorm:"type:uuid;not null"`
	DayOfWeek  int       `json:"day_of_week"`
	OpenTime   time.Time `json:"open_time"`
	CloseTime  time.Time `json:"close_time"`

	Field Field `gorm:"foreignKey:FieldID;references:FieldID"`

	TimeStamp
}
