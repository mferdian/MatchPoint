package migrations

import (
	"fieldreserve/model"

	"gorm.io/gorm"
)

func Seed(db *gorm.DB) error {
	if err := SeedFromJSON[model.User](db, "./migrations/json/users.json", model.User{}, "Email"); err != nil {
		return err
	}

	if err := SeedFromJSON[model.Category](db, "./migrations/json/category.json", model.Category{}, "Name"); err != nil {
		return err
	}

	return nil
}
