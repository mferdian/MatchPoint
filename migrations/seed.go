package migrations

import (
	"fieldreserve/model"

	"gorm.io/gorm"
)

func Seed(db *gorm.DB) error {
	err := SeedFromJSON[model.User](db, "./migrations/json/users.json", model.User{}, "Email")
	if err != nil {
		return err
	}

	return nil
}