package dto

import (
	"fieldreserve/model"

	"github.com/google/uuid"
)

type (
	CategoryResponse struct {
		ID          uuid.UUID `json:"category_id"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
	}

	CreateCategoryRequest struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	UpdateCategoryRequest struct {
		ID          string `json:"-"`
		Name        string `json:"name,omitempty"`
		Description string `json:"description,omitempty"`
	}

	DeleteCategoryRequest struct {
		ID string `json:"-"`
	}

	// Pagination
	CategoryPaginationRequest struct {
		PaginationRequest
		CategoryID string `form:"category_id"`
	}

	CategoryPaginationResponse struct {
		PaginationResponse
		Data []CategoryResponse `json:"data"`
	}

	CategoryPaginationRepositoryResponse struct {
		PaginationResponse
		Categorys []model.Category
	}
)
