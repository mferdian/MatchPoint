package model

import (
	"fieldreserve/helpers"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	UserID   uuid.UUID `json:"user_id" gorm:"type:uuid;primaryKey;column:user_id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
	NoTelp   string    `json:"no_telp" gorm:"column:no_telp"`
	Address  string    `json:"address"`
	Role     string    `json:"role"`

	TimeStamp
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()

	var err error
	u.Password, err = helpers.HashPassword(u.Password)
	if err != nil {
		return err
	}

	return nil
}
