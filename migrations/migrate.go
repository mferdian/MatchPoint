package migrations

import (
	"fieldreserve/model"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	db = db.Debug()

	if err := db.AutoMigrate(&model.User{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&model.Category{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&model.Field{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&model.Schedule{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&model.Booking{}); err != nil {
		return err
	}

	return nil
}

