package repository

import (
	"context"
	"fieldreserve/dto"
	"fieldreserve/model"
	"math"
	"strings"

	"gorm.io/gorm"
)

type (
	ICategoryRepository interface {
		CreateCategory(ctx context.Context, tx *gorm.DB, category model.Category) error
		GetAllCategoryWithPagination(ctx context.Context, tx *gorm.DB, req dto.CategoryPaginationRequest) (dto.CategoryPaginationRepositoryResponse, error)
		GetCategoryByID(ctx context.Context, tx *gorm.DB, categoryID string) (model.Category, bool, error)
		UpdateCategory(ctx context.Context, tx *gorm.DB, category model.Category) error
		Deletecategory(ctx context.Context, tx *gorm.DB, categoryID string) error
	}

	CategoryRepository struct {
		db *gorm.DB
	}
)

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{
		db: db,
	}
}

func (cr *CategoryRepository) CreateCategory(ctx context.Context, tx *gorm.DB, category model.Category) error {
	if tx == nil {
		tx = cr.db
	}

	return tx.WithContext(ctx).Create(&category).Error
}

func (cr *CategoryRepository) GetAllCategoryWithPagination(ctx context.Context, tx *gorm.DB, req dto.CategoryPaginationRequest) (dto.CategoryPaginationRepositoryResponse, error) {
	if tx == nil {
		tx = cr.db
	}

	var categorys []model.Category
	var err error
	var count int64

	if req.PaginationRequest.PerPage == 0 {
		req.PaginationRequest.PerPage = 10
	}

	if req.PaginationRequest.Page == 0 {
		req.PaginationRequest.Page = 1
	}

	query := tx.WithContext(ctx).Model(&model.Category{})

	if req.PaginationRequest.Search != "" {
		searchValue := "%" + strings.ToLower(req.PaginationRequest.Search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ?",
			searchValue, searchValue)
	}

	if req.CategoryID != "" {
		query = query.Where("category_id = ?", req.CategoryID)
	}

	if err := query.Count(&count).Error; err != nil {
		return dto.CategoryPaginationRepositoryResponse{}, err
	}

	if err := query.Order("created_at DESC").Scopes(Paginate(req.PaginationRequest.Page, req.PaginationRequest.PerPage)).Find(&categorys).Error; err != nil {
		return dto.CategoryPaginationRepositoryResponse{}, err
	}

	totalPage := int64(math.Ceil(float64(count) / float64(req.PaginationRequest.PerPage)))

	return dto.CategoryPaginationRepositoryResponse{
		Categorys: categorys,
		PaginationResponse: dto.PaginationResponse{
			Page:    req.PaginationRequest.Page,
			PerPage: req.PaginationRequest.PerPage,
			MaxPage: totalPage,
			Count:   count,
		},
	}, err
}
func (cr *CategoryRepository) GetCategoryByID(ctx context.Context, tx *gorm.DB, categoryID string) (model.Category, bool, error) {
	if tx == nil {
		tx = cr.db
	}

	var category model.Category
	if err := tx.WithContext(ctx).Where("category_id = ?", categoryID).Take(&category).Error; err != nil {
		return model.Category{}, false, err
	}

	return category, true, nil
}

func (cr *CategoryRepository) UpdateCategory(ctx context.Context, tx *gorm.DB, category model.Category) error {
	if tx == nil {
		tx = cr.db
	}

	return tx.WithContext(ctx).Where("category_id = ?", category.CategoryID).Updates(&category).Error
}
func (cr *CategoryRepository) Deletecategory(ctx context.Context, tx *gorm.DB, categoryID string) error {
	if tx == nil {
		tx = cr.db
	}

	return tx.WithContext(ctx).Where("category_id = ?", categoryID).Delete(&model.Category{}).Error
}
