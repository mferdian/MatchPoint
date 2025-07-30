package repository

import (
	"context"
	"fieldreserve/dto"
	"fieldreserve/model"
	"math"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	IFieldRepository interface {
		CreateField(ctx context.Context, tx *gorm.DB, field model.Field) error
		GetAllFieldWithPagination(ctx context.Context, tx *gorm.DB, req dto.FieldPaginationRequest) (dto.FieldPaginationRepositoryResponse, error)
		GetFieldByID(ctx context.Context, tx *gorm.DB, fieldID string) (model.Field, bool, error)
		UpdateField(ctx context.Context, tx *gorm.DB, field model.Field) error
		GetCategoryByID(ctx context.Context, tx *gorm.DB, categoryID uuid.UUID) (model.Category, error)
		DeleteField(ctx context.Context, tx *gorm.DB, fieldID string) error
		GetAllWithSchedules(ctx context.Context, tx *gorm.DB) ([]model.Field, error)
	}

	FieldRepository struct {
		db *gorm.DB
	}
)

func NewFieldRepository(db *gorm.DB) *FieldRepository {
	return &FieldRepository{
		db: db,
	}
}

func (fr *FieldRepository) CreateField(ctx context.Context, tx *gorm.DB, field model.Field) error {
	if tx == nil {
		tx = fr.db
	}

	return tx.WithContext(ctx).Create(&field).Error
}

func (fr *FieldRepository) GetAllFieldWithPagination(ctx context.Context, tx *gorm.DB, req dto.FieldPaginationRequest) (dto.FieldPaginationRepositoryResponse, error) {
	if tx == nil {
		tx = fr.db
	}

	var fields []model.Field
	var count int64

	// Default pagination jika kosong
	if req.PaginationRequest.PerPage == 0 {
		req.PaginationRequest.PerPage = 10
	}
	if req.PaginationRequest.Page == 0 {
		req.PaginationRequest.Page = 1
	}

	// Base query
	query := tx.WithContext(ctx).Model(&model.Field{})

	// Search
	if req.PaginationRequest.Search != "" {
		searchValue := "%" + strings.ToLower(req.PaginationRequest.Search) + "%"
		query = query.Where("LOWER(field_name) LIKE ?", searchValue) // <== diperbaiki dari "feild_name"
	}

	// Filter FieldID jika diberikan
	if req.FieldID != "" {
		query = query.Where("field_id = ?", req.FieldID)
	}

	// Hitung total data
	if err := query.Count(&count).Error; err != nil {
		return dto.FieldPaginationRepositoryResponse{}, err
	}

	// Ambil data sesuai pagination
	if err := query.Order("created_at DESC").Scopes(Paginate(req.PaginationRequest.Page, req.PaginationRequest.PerPage)).Find(&fields).Error; err != nil {
		return dto.FieldPaginationRepositoryResponse{}, err
	}

	totalPage := int64(math.Ceil(float64(count) / float64(req.PaginationRequest.PerPage)))

	return dto.FieldPaginationRepositoryResponse{
		Fields: fields,
		PaginationResponse: dto.PaginationResponse{
			Page:    req.PaginationRequest.Page,
			PerPage: req.PaginationRequest.PerPage,
			MaxPage: totalPage,
			Count:   count,
		},
	}, nil
}

func (fr *FieldRepository) GetFieldByID(ctx context.Context, tx *gorm.DB, fieldID string) (model.Field, bool, error) {
	if tx == nil {
		tx = fr.db
	}

	var field model.Field
	if err := tx.WithContext(ctx).Preload("Category").Where("field_id = ?", fieldID).Take(&field).Error; err != nil {
		return model.Field{}, false, err
	}

	return field, true, nil
}
func (fr *FieldRepository) UpdateField(ctx context.Context, tx *gorm.DB, field model.Field) error {
	if tx == nil {
		tx = fr.db
	}

	return tx.WithContext(ctx).Where("field_id = ?", field.FieldID).Updates(&field).Error
}
func (fr *FieldRepository) DeleteField(ctx context.Context, tx *gorm.DB, fieldID string) error {
	if tx == nil {
		tx = fr.db
	}

	return tx.WithContext(ctx).Where("field_id = ?", fieldID).Delete(&model.Field{}).Error
}

func (fr *FieldRepository) GetCategoryByID(ctx context.Context, tx *gorm.DB, categoryID uuid.UUID) (model.Category, error) {
	if tx == nil {
		tx = fr.db
	}

	var category model.Category
	if err := tx.WithContext(ctx).Where("category_id = ?", categoryID).First(&category).Error; err != nil {
		return model.Category{}, err
	}

	return category, nil
}

func (fr *FieldRepository) GetAllWithSchedules(ctx context.Context, tx *gorm.DB) ([]model.Field, error) {
	if tx == nil {
		tx = fr.db
	}

	var fields []model.Field
	err := tx.WithContext(ctx).
		Preload("Category").
		Preload("Schedules").
		Find(&fields).Error
	return fields, err
}
