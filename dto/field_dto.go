package dto

import "fieldreserve/model"

type (
	CreateFieldRequest struct {
	}

	FieldResponse struct {
	}

	UpdateFiledRequest struct {
	}

	DeleteFieldRequest struct {
	}

	FieldPaginationRequest struct {
	}

	FieldPaginationResponse struct {
		PaginationResponse
		Data []FieldResponse
	}

	FieldPaginationRepositoryResponse struct {
		PaginationResponse
		Fielss []model.Field
	}
)