package migrations

import (
	"encoding/json"
	"errors"
	"fieldreserve/helpers"
	"fmt"
	"io"
	"os"

	"gorm.io/gorm"
)

func SeedFromJSON[T any](db *gorm.DB, filePath string, model T, uniqueFields ...string) error {
	jsonFile, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer jsonFile.Close()

	jsonData, err := io.ReadAll(jsonFile)
	if err != nil {
		return fmt.Errorf("failed to read JSON data: %w", err)
	}

	var listData []T
	if err := json.Unmarshal(jsonData, &listData); err != nil {
		return fmt.Errorf("failed to unmarshal JSON data: %w", err)
	}

	if err := db.AutoMigrate(&model); err != nil {
		return fmt.Errorf("failed to migrate model: %w", err)
	}

	for _, data := range listData {
		query := db.Model(&model)

		for _, field := range uniqueFields {
			val, err := helpers.GetFieldValue(data, field)
			if err != nil {
				return err
			}
			query = query.Where(fmt.Sprintf("%s = ?", helpers.SnakeCase(field)), val)
		}

		var existing T
		if err := query.First(&existing).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
			continue
		}

		if err := db.Create(&data).Error; err != nil {
			return fmt.Errorf("failed to insert data: %w", err)
		}
	}

	return nil
}
