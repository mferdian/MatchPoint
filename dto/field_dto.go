package dto

import (
	"fieldreserve/model"
	"mime/multipart"

	"github.com/google/uuid"
)

type (
	CreateFieldRequest struct {
		FieldName    string                `form:"field_name" binding:"required"`
		CategoryID   string                `form:"category_id" binding:"required"`
		FieldAddress string                `form:"field_address" binding:"required"`
		FieldPrice   int                   `form:"field_price" binding:"required"`
		FieldImage   *multipart.FileHeader `form:"field_image" binding:"required"`
	}

	CategoryCompactResponse struct {
		CategoryID  uuid.UUID `json:"category_id"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
	}

	FieldResponse struct {
		FieldID      uuid.UUID `json:"field_id"`
		FieldName    string    `json:"field_name"`
		FieldAddress string    `json:"field_address"`
		FieldPrice   int       `json:"field_price"`
		FieldImage   string    `json:"field_image"`
		CategoryID   uuid.UUID `json:"category_id"`
	}

	FieldFullResponse struct {
		FieldID      uuid.UUID               `json:"field_id"`
		FieldName    string                  `json:"field_name"`
		FieldAddress string                  `json:"field_address"`
		FieldPrice   int                     `json:"field_price"`
		FieldImage   string                  `json:"field_image"`
		Category     CategoryCompactResponse `json:"category"`
	}

	UpdateFieldRequest struct {
		FieldID      string                `form:"-"`
		FieldName    string                `form:"field_name"`
		FieldAddress string                `form:"field_address"`
		FieldPrice   int                   `form:"field_price"`
		FieldImage   *multipart.FileHeader `form:"field_image"`
	}
	DeleteFieldRequest struct {
		FieldID string `json:"-"`
	}


	// Paginate
	
	FieldPaginationRequest struct {
		PaginationRequest
		FieldID string `form:"field_id"`
	}

	FieldPaginationResponse struct {
		Data []FieldResponse `json:"data"`
		PaginationResponse 
	}

	FieldPaginationRepositoryResponse struct {
		PaginationResponse
		Fields []model.Field
	}
)
