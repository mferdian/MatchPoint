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
	BookingDate   time.Time `json:"booking_date"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	TotalPayment  float64   `json:"total_payment"`
	ProofPayment  string    `json:"proof_payment"`
	Status        string    `json:"status"`

	PaymentUploadedAt *time.Time `json:"payment_uploaded_at"`
	PaymentVerifiedAt *time.Time `json:"payment_verified_at"`
	CancelledAt       *time.Time `json:"cancelled_at"`

	User  User  `gorm:"foreignKey:UserID;references:UserID"`
	Field Field `gorm:"foreignKey:FieldID;references:FieldID"`

	TimeStamp
}
