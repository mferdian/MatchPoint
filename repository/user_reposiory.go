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
	IUserRepository interface {
		GetUserByEmail(ctx context.Context, tx *gorm.DB, email string) (model.User, bool, error)
		GetUserByID(ctx context.Context, tx *gorm.DB, userID string) (model.User, bool, error)
		GetAllUserWithPagination(ctx context.Context, tx *gorm.DB, req dto.UserPaginationRequest) (dto.UserPaginationRepositoryResponse, error)
		CreateUser(ctx context.Context, tx *gorm.DB, user model.User) error
		UpdateUser(ctx context.Context, tx *gorm.DB, user model.User) error
		DeleteUserByID(ctx context.Context, tx *gorm.DB, userID string) error
	}

	UserRepository struct {
		db *gorm.DB
	}
)

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (ur *UserRepository) GetUserByEmail(ctx context.Context, tx *gorm.DB, email string) (model.User, bool, error) {
	if tx == nil {
		tx = ur.db
	}

	var user model.User
	if err := tx.WithContext(ctx).Where("email = ?", email).Take(&user).Error; err != nil {
		return model.User{}, false, err
	}

	return user, true, nil
}

func (ur *UserRepository) GetUserByID(ctx context.Context, tx *gorm.DB, userID string) (model.User, bool, error) {
	if tx == nil {
		tx = ur.db
	}

	var user model.User
	if err := tx.WithContext(ctx).Where("user_id = ?", userID).Take(&user).Error; err != nil {
		return model.User{}, false, err
	}

	return user, true, nil
}

func (ur *UserRepository) GetAllUserWithPagination(ctx context.Context, tx *gorm.DB, req dto.UserPaginationRequest) (dto.UserPaginationRepositoryResponse, error) {
	if tx == nil {
		tx = ur.db
	}

	var users []model.User
	var err error
	var count int64

	if req.PaginationRequest.PerPage == 0 {
		req.PaginationRequest.PerPage = 10
	}

	if req.PaginationRequest.Page == 0 {
		req.PaginationRequest.Page = 1
	}

	query := tx.WithContext(ctx).Model(&model.User{})

	if req.PaginationRequest.Search != "" {
		searchValue := "%" + strings.ToLower(req.PaginationRequest.Search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(email) LIKE ?",
			searchValue, searchValue, searchValue)
	}

	if req.UserID != "" {
		query = query.Where("user_id = ?", req.UserID)
	}

	if err := query.Count(&count).Error; err != nil {
		return dto.UserPaginationRepositoryResponse{}, err
	}

	if err := query.Order("created_at DESC").Scopes(Paginate(req.PaginationRequest.Page, req.PaginationRequest.PerPage)).Find(&users).Error; err != nil {
		return dto.UserPaginationRepositoryResponse{}, err
	}

	totalPage := int64(math.Ceil(float64(count) / float64(req.PaginationRequest.PerPage)))

	return dto.UserPaginationRepositoryResponse{
		Users: users,
		PaginationResponse: dto.PaginationResponse{
			Page:    req.PaginationRequest.Page,
			PerPage: req.PaginationRequest.PerPage,
			MaxPage: totalPage,
			Count:   count,
		},
	}, err
}

func (ur *UserRepository) CreateUser(ctx context.Context, tx *gorm.DB, user model.User) error {
	if tx == nil {
		tx = ur.db
	}

	return tx.WithContext(ctx).Create(&user).Error
}

func (ur *UserRepository) UpdateUser(ctx context.Context, tx *gorm.DB, user model.User) error {
	if tx == nil {
		tx = ur.db
	}

	return tx.WithContext(ctx).Where("user_id = ?", user.UserID).Updates(&user).Error
}

func (ur *UserRepository) DeleteUserByID(ctx context.Context, tx *gorm.DB, userID string) error {
	if tx == nil {
		tx = ur.db
	}

	return tx.WithContext(ctx).Where("user_id = ?", userID).Delete(&model.User{}).Error
}
