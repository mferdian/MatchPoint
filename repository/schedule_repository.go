package repository

import (
	"context"
	"fieldreserve/model"

	"gorm.io/gorm"
)

type (
	IScheduleRepository interface {
		CreateSchedule(ctx context.Context, tx *gorm.DB, schedule model.Schedule) error
		GetAllSchedule(ctx context.Context, tx *gorm.DB) ([]model.Schedule, error)
		GetScheduleByID(ctx context.Context, tx *gorm.DB, id string) (model.Schedule, error)
		UpdateSchedule(ctx context.Context, tx *gorm.DB, schedule model.Schedule) error
		DeleteScheduleByID(ctx context.Context, tx *gorm.DB, id string) error

		GetSchedulesByFieldID(ctx context.Context, tx *gorm.DB, fieldID string) ([]model.Schedule, error)
		GetScheduleByFieldIDAndDay(ctx context.Context, tx *gorm.DB, fieldID string, day int) (model.Schedule, error)
	}
	ScheduleRepository struct {
		db *gorm.DB
	}
)

func NewScheduleRepository(db *gorm.DB) *ScheduleRepository {
	return &ScheduleRepository{db: db}
}


func (sr *ScheduleRepository) CreateSchedule(ctx context.Context, tx *gorm.DB, schedule model.Schedule) error {
	if tx == nil {
		tx = sr.db
	}
	return tx.WithContext(ctx).Create(&schedule).Error
}

func (sr *ScheduleRepository) GetAllSchedule(ctx context.Context, tx *gorm.DB) ([]model.Schedule, error) {
	if tx == nil {
		tx = sr.db
	}

	var schedules []model.Schedule
	err := tx.WithContext(ctx).Preload("Field").Find(&schedules).Error
	return schedules, err
}

func (r *ScheduleRepository) GetScheduleByID(ctx context.Context, tx *gorm.DB, id string) (model.Schedule, error) {
	var schedule model.Schedule
	if tx == nil {
		tx = r.db
	}

	err := tx.WithContext(ctx).
		Preload("Field.Category").
		First(&schedule, "schedule_id = ?", id).Error

	return schedule, err
}


func (sr *ScheduleRepository) UpdateSchedule(ctx context.Context, tx *gorm.DB, schedule model.Schedule) error {
	if tx == nil {
		tx = sr.db
	}
	return tx.WithContext(ctx).Where("schedule_id = ?", schedule.ScheduleID).Updates(&schedule).Error
}

func (sr *ScheduleRepository) DeleteScheduleByID(ctx context.Context, tx *gorm.DB, id string) error {
	if tx == nil {
		tx = sr.db
	}
	return tx.WithContext(ctx).Where("schedule_id = ?", id).Delete(&model.Schedule{}).Error
}

func (sr *ScheduleRepository) GetSchedulesByFieldID(ctx context.Context, tx *gorm.DB, fieldID string) ([]model.Schedule, error) {
	if tx == nil {
		tx = sr.db
	}
	var schedules []model.Schedule
	err := tx.WithContext(ctx).
		Where("field_id = ?", fieldID).
		Order("day_of_week asc").
		Find(&schedules).Error
	return schedules, err
}

func (sr *ScheduleRepository) GetScheduleByFieldIDAndDay(ctx context.Context, tx *gorm.DB, fieldID string, day int) (model.Schedule, error) {
	if tx == nil {
		tx = sr.db
	}
	var schedule model.Schedule
	err := tx.WithContext(ctx).
		Where("field_id = ? AND day_of_week = ?", fieldID, day).
		First(&schedule).Error
	return schedule, err
}
