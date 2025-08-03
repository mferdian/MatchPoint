package service

import (
	"context"
	"fieldreserve/constants"
	"fieldreserve/dto"
	"fieldreserve/model"
	"fieldreserve/repository"
	"fieldreserve/utils"
	
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
		utils.Log.Errorf("Failed to create category: %v", err)
		return dto.CategoryResponse{}, constants.ErrCreateCategory
	}

	utils.Log.Infof("Created category: %s (%s)", category.Name, category.CategoryID.String())

	return dto.CategoryResponse{
		ID:          category.CategoryID,
		Name:        category.Name,
		Description: category.Description,
	}, nil
}

func (cs *CategoryService) GetAllCategoryWithPagination(ctx context.Context, req dto.CategoryPaginationRequest) (dto.CategoryPaginationResponse, error) {
	dataWithPaginate, err := cs.categoryRepo.GetAllCategoryWithPagination(ctx, nil, req)
	if err != nil {
		utils.Log.Errorf("Failed to get categories: %v", err)
		return dto.CategoryPaginationResponse{}, constants.ErrGetAllCategory
	}

	var datas []dto.CategoryResponse
	for _, category := range dataWithPaginate.Categorys {
		datas = append(datas, dto.CategoryResponse{
			ID:          category.CategoryID,
			Name:        category.Name,
			Description: category.Description,
		})
	}

	utils.Log.Infof("Fetched %d categories (Page %d)", len(datas), dataWithPaginate.Page)

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
		utils.Log.Warnf("Invalid UUID when get category: %s", categoryID)
		return dto.CategoryResponse{}, constants.ErrInvalidUUID
	}

	category, _, err := cs.categoryRepo.GetCategoryByID(ctx, nil, categoryID)
	if err != nil {
		utils.Log.Errorf("Failed to get category by ID %s: %v", categoryID, err)
		return dto.CategoryResponse{}, constants.ErrGetCategoryByID
	}

	utils.Log.Infof("Fetched category by ID: %s", categoryID)

	return dto.CategoryResponse{
		ID:          category.CategoryID,
		Name:        category.Name,
		Description: category.Description,
	}, nil
}

func (cs *CategoryService) UpdateCategory(ctx context.Context, req dto.UpdateCategoryRequest) (dto.CategoryResponse, error) {
	if _, err := uuid.Parse(req.ID); err != nil {
		utils.Log.Warnf("Invalid UUID when update category: %s", req.ID)
		return dto.CategoryResponse{}, constants.ErrInvalidUUID
	}

	category, _, err := cs.categoryRepo.GetCategoryByID(ctx, nil, req.ID)
	if err != nil {
		utils.Log.Errorf("Failed to find category to update (ID: %s): %v", req.ID, err)
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
		utils.Log.Errorf("Failed to update category %s: %v", category.CategoryID, err)
		return dto.CategoryResponse{}, constants.ErrUpdateCategory
	}

	utils.Log.Infof("Updated category: %s (%s)", category.Name, category.CategoryID.String())

	return dto.CategoryResponse{
		ID:          category.CategoryID,
		Name:        category.Name,
		Description: category.Description,
	}, nil
}

func (cs *CategoryService) DeleteCategory(ctx context.Context, req dto.DeleteCategoryRequest) (dto.CategoryResponse, error) {
	if _, err := uuid.Parse(req.ID); err != nil {
		utils.Log.Warnf("Invalid UUID when delete category: %s", req.ID)
		return dto.CategoryResponse{}, constants.ErrInvalidUUID
	}

	deletedCategory, _, err := cs.categoryRepo.GetCategoryByID(ctx, nil, req.ID)
	if err != nil {
		utils.Log.Errorf("Failed to find category to delete (ID: %s): %v", req.ID, err)
		return dto.CategoryResponse{}, constants.ErrGetCategoryByID
	}

	err = cs.categoryRepo.Deletecategory(ctx, nil, req.ID)
	if err != nil {
		utils.Log.Errorf("Failed to delete category %s: %v", req.ID, err)
		return dto.CategoryResponse{}, constants.ErrDeleteCategoryByID
	}

	utils.Log.Infof("Deleted category: %s (%s)", deletedCategory.Name, deletedCategory.CategoryID.String())

	return dto.CategoryResponse{
		ID:          deletedCategory.CategoryID,
		Name:        deletedCategory.Name,
		Description: deletedCategory.Description,
	}, nil
}
