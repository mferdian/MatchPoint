package migrations

import (
	"fieldreserve/model"

	"gorm.io/gorm"
)

func Rollback(db *gorm.DB) error {
	tables := []interface{}{
		&model.User{},
		&model.Category{},
		&model.Field{},
		&model.Schedule{},
		&model.Booking{},
	}

	for _, table := range tables {
		if err := db.Migrator().DropTable(table); err != nil {
			return err
		}
	}

	return nil
}
