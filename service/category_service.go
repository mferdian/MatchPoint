package service

import (
	"context"
	"fieldreserve/constants"
	"fieldreserve/dto"
	"fieldreserve/model"
	"fieldreserve/repository"

	"github.com/google/uuid"
)

type (
	ICategoryService interface {
		CreateCategory(ctx context.Context, req dto.CreateCategoryRequest) (dto.CategoryResponse, error)
		GetAllCategoryWithPagination(ctx context.Context, req dto.CategoryPaginationRequest) (dto.CategoryPaginationResponse, error)
		GetCategoryByID(ctx context.Context, categoryID string) (dto.CategoryResponse, error)
		UpdateCategory(ctx context.Context, req dto.UpdateCategoryRequest) (dto.CategoryResponse, error)
		DeleteCategory(ctx context.Context, req dto.DeleteCategoryRequest) (dto.CategoryResponse, error)
	}

	CategoryService struct {
		categoryRepo repository.ICategoryRepository
	}
)

func NewCategoryService(categoryRepo repository.ICategoryRepository) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
	}
}

func (cs *CategoryService) CreateCategory(ctx context.Context, req dto.CreateCategoryRequest) (dto.CategoryResponse, error) {
	category := model.Category{
		CategoryID:  uuid.New(),
		Name:        req.Name,
		Description: req.Description,
	}

	err := cs.categoryRepo.CreateCategory(ctx, nil, category)
	if err != nil {
		return dto.CategoryResponse{}, constants.ErrCreateCategory
	}

	res := dto.CategoryResponse{
		ID:          category.CategoryID,
		Name:        category.Name,
		Description: category.Description,
	}

	return res, nil
}

func (cs *CategoryService) GetAllCategoryWithPagination(ctx context.Context, req dto.CategoryPaginationRequest) (dto.CategoryPaginationResponse, error) {
	dataWithPaginate, err := cs.categoryRepo.GetAllCategoryWithPagination(ctx, nil, req)
	if err != nil {
		return dto.CategoryPaginationResponse{}, constants.ErrGetAllCategory
	}

	var datas []dto.CategoryResponse
	for _, category := range dataWithPaginate.Categorys {
		data := dto.CategoryResponse{
			ID:          category.CategoryID,
			Name:        category.Name,
			Description: category.Description,
		}

		datas = append(datas, data)
	}

	return dto.CategoryPaginationResponse{
		Data: datas,
		PaginationResponse: dto.PaginationResponse{
			Page:    dataWithPaginate.Page,
			PerPage: dataWithPaginate.PerPage,
			MaxPage: dataWithPaginate.MaxPage,
			Count:   dataWithPaginate.Count,
		},
	}, nil
}

func (cs *CategoryService) GetCategoryByID(ctx context.Context, categoryID string) (dto.CategoryResponse, error) {
	if _, err := uuid.Parse(categoryID); err != nil {
		return dto.CategoryResponse{}, constants.ErrInvalidUUID
	}

	category, _, err := cs.categoryRepo.GetCategoryByID(ctx, nil, categoryID)

	if err != nil {
		return dto.CategoryResponse{}, constants.ErrGetCategoryByID
	}

	res := dto.CategoryResponse{
		ID:          category.CategoryID,
		Name:        category.Name,
		Description: category.Description,
	}

	return res, nil
}

func (cs *CategoryService) UpdateCategory(ctx context.Context, req dto.UpdateCategoryRequest) (dto.CategoryResponse, error) {
	if _, err := uuid.Parse(req.ID); err != nil {
		return dto.CategoryResponse{}, constants.ErrInvalidUUID
	}

	category, _, err := cs.categoryRepo.GetCategoryByID(ctx, nil, req.ID)
	if err != nil {
		return dto.CategoryResponse{}, constants.ErrGetCategoryByID
	}

	if req.Name != "" {
		category.Name = req.Name
	}

	if req.Description != "" {
		category.Description = req.Description
	}

	err = cs.categoryRepo.UpdateCategory(ctx, nil, category)
	if err != nil {
		return dto.CategoryResponse{}, constants.ErrUpdateCategory
	}

	res := dto.CategoryResponse{
		ID:          category.CategoryID,
		Name:        category.Name,
		Description: category.Description,
	}

	return res, nil
}

func (cs *CategoryService) DeleteCategory(ctx context.Context, req dto.DeleteCategoryRequest) (dto.CategoryResponse, error) {

	if _, err := uuid.Parse(req.ID); err != nil {
		return dto.CategoryResponse{}, constants.ErrInvalidUUID
	}

	deletedCategory, _, err := cs.categoryRepo.GetCategoryByID(ctx, nil, req.ID)
	if err != nil {
		return dto.CategoryResponse{}, constants.ErrGetCategoryByID
	}

	err = cs.categoryRepo.Deletecategory(ctx, nil, req.ID)
	if err != nil {
		return dto.CategoryResponse{}, constants.ErrDeleteCategoryByID
	}

	res := dto.CategoryResponse{
		ID:          deletedCategory.CategoryID,
		Name:        deletedCategory.Name,
		Description: deletedCategory.Description,
	}

	return res, nil
}
