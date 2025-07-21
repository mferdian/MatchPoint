package model

import "github.com/google/uuid"

type Category struct {
	CategoryID  uuid.UUID `gorm:"type:uuid;primaryKey;column:category_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	
	TimeStamp
}