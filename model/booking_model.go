package model

import (
	"time"

	"github.com/google/uuid"
)

type Booking struct {
	BookingID     uuid.UUID `gorm:"type:uuid;primaryKey;column:booking_id"`
	UserID        uuid.UUID `gorm:"type:uuid;not null"`
	FieldID       uuid.UUID `gorm:"type:uuid;not null"`
	PaymentMethod string    `json:"payment_method"`
	BookingDate   time.Time	`json:"booking_date"`
	StartTime     time.Time	`json:"start_time"`
	EndTime       time.Time	`json:"end_time"`
	Status        string    `json:"status"`

	User    User  `gorm:"foreignKey:UserID;references:UserID"`
	Field   Field `gorm:"foreignKey:FieldID;references:FieldID"`

	TimeStamp
}

