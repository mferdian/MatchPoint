package service

import (
	"context"
	"fieldreserve/constants"
	"fieldreserve/dto"
	"fieldreserve/helpers"
	"fieldreserve/model"
	"fieldreserve/repository"

	"github.com/google/uuid"
)

type (
	IFieldService interface {
		CreateField(ctx context.Context, req dto.CreateFieldRequest) (dto.FieldResponse, error)
		GetAllFieldWithPagination(ctx context.Context, req dto.FieldPaginationRequest) (dto.FieldPaginationResponse, error)
		GetFieldByID(ctx context.Context, fieldID string) (dto.FieldResponse, error)
		UpdateField(ctx context.Context, req dto.UpdateFieldRequest) (dto.FieldResponse, error)
		DeleteField(ctx context.Context, req dto.DeleteFieldRequest) (dto.FieldResponse, error)
	}

	FieldService struct {
		fieldRepo repository.IFieldRepository
	}
)

func NewFieldService(fieldRepo repository.IFieldRepository) *FieldService {
	return &FieldService{
		fieldRepo: fieldRepo,
	}
}

func (fs *FieldService) CreateField(ctx context.Context, req dto.CreateFieldRequest) (dto.FieldResponse, error) {
	imageName, err := helpers.SaveImage(req.FieldImage, "./assets/fields", "field")
	if err != nil {
		return dto.FieldResponse{}, constants.ErrSaveImages
	}

	categoryUUID, err := uuid.Parse(req.CategoryID)
	if err != nil {
		return dto.FieldResponse{}, constants.ErrInvalidUUID
	}

	category, err := fs.fieldRepo.GetCategoryByID(ctx, nil, categoryUUID)
	if err != nil {
		return dto.FieldResponse{}, constants.ErrGetCategoryByID
	}

	field := model.Field{
		FieldID:      uuid.New(),
		CategoryID:   categoryUUID,
		FieldName:    req.FieldName,
		FieldAddress: req.FieldAddress,
		FieldPrice:   req.FieldPrice,
		FieldImage:   imageName,
	}

	if err := fs.fieldRepo.CreateField(ctx, nil, field); err != nil {
		return dto.FieldResponse{}, err
	}

	return dto.FieldResponse{
		FieldID:      field.FieldID,
		FieldName:    field.FieldName,
		FieldAddress: field.FieldAddress,
		FieldPrice:   field.FieldPrice,
		FieldImage:   field.FieldImage,
		Category: dto.CategoryCompactResponse{
			CategoryID:  category.CategoryID,
			Name:        category.Name,
			Description: category.Description,
		},
	}, nil
}



func (fs *FieldService) GetAllFieldWithPagination(ctx context.Context, req dto.FieldPaginationRequest) (dto.FieldPaginationResponse, error) {
	dataWithPaginate, err := fs.fieldRepo.GetAllFieldWithPagination(ctx, nil, req)
	if err != nil {
		return dto.FieldPaginationResponse{}, constants.ErrGetAllField
	}

	var datas []dto.FieldResponse
	for _, field := range dataWithPaginate.Fields {
		data := dto.FieldResponse{
			FieldID:      field.FieldID,
			FieldName:    field.FieldName,
			FieldAddress: field.FieldAddress,
			FieldPrice:   field.FieldPrice,
			FieldImage:   field.FieldImage,
			Category: dto.CategoryCompactResponse{
				CategoryID:  field.Category.CategoryID,
				Name:        field.Category.Name,
				Description: field.Category.Description,
			},
		}

		datas = append(datas, data)
	}

	return dto.FieldPaginationResponse{
		Data: datas,
		PaginationResponse: dto.PaginationResponse{
			Page:    dataWithPaginate.Page,
			PerPage: dataWithPaginate.PerPage,
			MaxPage: dataWithPaginate.MaxPage,
			Count:   dataWithPaginate.Count,
		},
	}, nil
}

func (fs *FieldService) GetFieldByID(ctx context.Context, fieldID string) (dto.FieldResponse, error) {
	if _, err := uuid.Parse(fieldID); err != nil {
		return dto.FieldResponse{}, constants.ErrInvalidUUID
	}

	field, _, err := fs.fieldRepo.GetFieldByID(ctx, nil, fieldID)

	if err != nil {
		return dto.FieldResponse{}, constants.ErrGetFieldByID
	}

	res := dto.FieldResponse{
		FieldID:      field.FieldID,
		FieldName:    field.FieldName,
		FieldAddress: field.FieldAddress,
		FieldPrice:   field.FieldPrice,
		FieldImage:   field.FieldImage,
		Category: dto.CategoryCompactResponse{
			CategoryID:  field.Category.CategoryID,
			Name:        field.Category.Name,
			Description: field.Category.Description,
		},
	}

	return res, nil
}

func (fs *FieldService) UpdateField(ctx context.Context, req dto.UpdateFieldRequest) (dto.FieldResponse, error) {
	if _, err := uuid.Parse(req.FieldID); err != nil {
		return dto.FieldResponse{}, constants.ErrInvalidUUID
	}

	field, _, err := fs.fieldRepo.GetFieldByID(ctx, nil, req.FieldID)
	if err != nil {
		return dto.FieldResponse{}, constants.ErrGetFieldByID
	}

	if req.FieldName != "" {
		field.FieldName = req.FieldName
	}
	if req.FieldAddress != "" {
		field.FieldAddress = req.FieldAddress
	}
	if req.FieldPrice < 0 {
		return dto.FieldResponse{}, constants.ErrInvalidFieldPrice
	}
	if req.FieldPrice != 0 {
		field.FieldPrice = req.FieldPrice
	}

	if req.FieldImage != nil {
		imageName, err := helpers.SaveImage(req.FieldImage, "./assets/fields", "field")
		if err != nil {
			return dto.FieldResponse{}, constants.ErrSaveImages
		}
		field.FieldImage = imageName
	}

	if err := fs.fieldRepo.UpdateField(ctx, nil, field); err != nil {
		return dto.FieldResponse{}, constants.ErrUpdateField
	}
	res := dto.FieldResponse{
		FieldID:      field.FieldID,
		FieldName:    field.FieldName,
		FieldAddress: field.FieldAddress,
		FieldPrice:   field.FieldPrice,
		FieldImage:   field.FieldImage,
		Category: dto.CategoryCompactResponse{
			CategoryID:  field.Category.CategoryID,
			Name:        field.Category.Name,
			Description: field.Category.Description,
		},
	}

	return res, nil
}


func (fs *FieldService) DeleteField(ctx context.Context, req dto.DeleteFieldRequest) (dto.FieldResponse, error) {

	if _, err := uuid.Parse(req.FieldID); err != nil {
		return dto.FieldResponse{}, constants.ErrInvalidUUID
	}

	deletedField, _, err := fs.fieldRepo.GetFieldByID(ctx, nil, req.FieldID)
	if err != nil {
		return dto.FieldResponse{}, constants.ErrGetFieldByID
	}

	err = fs.fieldRepo.DeleteField(ctx, nil, req.FieldID)
	if err != nil {
		return dto.FieldResponse{}, constants.ErrDeleteFieldByID
	}

	res := dto.FieldResponse{
		FieldID:      deletedField.FieldID,
		FieldName:    deletedField.FieldName,
		FieldAddress: deletedField.FieldAddress,
		FieldPrice:   deletedField.FieldPrice,
		FieldImage:   deletedField.FieldImage,
		Category: dto.CategoryCompactResponse{
			CategoryID:  deletedField.Category.CategoryID,
			Name:        deletedField.Category.Name,
			Description: deletedField.Category.Description,
		},
	}

	return res, nil
}
