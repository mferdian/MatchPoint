package model

import "github.com/google/uuid"

type Field struct {
	FieldID      uuid.UUID `gorm:"type:uuid;primaryKey;column:field_id"`
	CategoryID   uuid.UUID `gorm:"type:uuid;not null"`
	FieldName    string    `json:"field_name"`
	FieldAddress string    `json:"field_address"`
	FieldPrice   int       `json:"field_price"`
	FieldImage   string    `json:"field_image"`

	TimeStamp

	Category Category `gorm:"foreignKey:CategoryID;references:CategoryID"`
}
