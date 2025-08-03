package service

import (
	"context"
	"fieldreserve/constants"
	"fieldreserve/dto"
	"fieldreserve/helpers"
	"fieldreserve/model"
	"fieldreserve/repository"
	"fieldreserve/utils"

	"github.com/google/uuid"
)

type (
	IFieldService interface {
		CreateField(ctx context.Context, req dto.CreateFieldRequest) (dto.FieldResponse, error)
		GetAllFieldWithPagination(ctx context.Context, req dto.FieldPaginationRequest) (dto.FieldPaginationResponse, error)
		GetFieldByID(ctx context.Context, fieldID string) (dto.FieldFullResponse, error)
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
	utils.Log.Info("Creating new field")

	imageName, err := helpers.SaveImage(req.FieldImage, "./assets/fields", "field")
	if err != nil {
		utils.Log.Errorf("Failed to save image: %v", err)
		return dto.FieldResponse{}, constants.ErrSaveImages
	}

	categoryUUID, err := uuid.Parse(req.CategoryID)
	if err != nil {
		utils.Log.Errorf("Invalid category UUID: %v", err)
		return dto.FieldResponse{}, constants.ErrInvalidUUID
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
		utils.Log.Errorf("Failed to create field in repository: %v", err)
		return dto.FieldResponse{}, err
	}

	utils.Log.Infof("Field created successfully: %v", field.FieldID)

	return dto.FieldResponse{
		FieldID:      field.FieldID,
		FieldName:    field.FieldName,
		FieldAddress: field.FieldAddress,
		FieldPrice:   field.FieldPrice,
		FieldImage:   field.FieldImage,
		CategoryID:   field.CategoryID,
	}, nil
}

func (fs *FieldService) GetAllFieldWithPagination(ctx context.Context, req dto.FieldPaginationRequest) (dto.FieldPaginationResponse, error) {
	utils.Log.Infof("Fetching all fields with pagination: page=%d, perPage=%d", req.Page, req.PerPage)

	dataWithPaginate, err := fs.fieldRepo.GetAllFieldWithPagination(ctx, nil, req)
	if err != nil {
		utils.Log.Errorf("Failed to fetch paginated fields: %v", err)
		return dto.FieldPaginationResponse{}, constants.ErrGetAllField
	}

	utils.Log.Infof("Fetched %d fields", len(dataWithPaginate.Fields))

	var datas []dto.FieldResponse
	for _, field := range dataWithPaginate.Fields {
		data := dto.FieldResponse{
			FieldID:      field.FieldID,
			FieldName:    field.FieldName,
			FieldAddress: field.FieldAddress,
			FieldPrice:   field.FieldPrice,
			FieldImage:   field.FieldImage,
			CategoryID:   field.CategoryID,
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

func (fs *FieldService) GetFieldByID(ctx context.Context, fieldID string) (dto.FieldFullResponse, error) {
	utils.Log.Infof("Fetching field by ID: %s", fieldID)

	if _, err := uuid.Parse(fieldID); err != nil {
		utils.Log.Errorf("Invalid field UUID: %v", err)
		return dto.FieldFullResponse{}, constants.ErrInvalidUUID
	}

	field, _, err := fs.fieldRepo.GetFieldByID(ctx, nil, fieldID)
	if err != nil {
		utils.Log.Errorf("Failed to fetch field: %v", err)
		return dto.FieldFullResponse{}, constants.ErrGetFieldByID
	}

	utils.Log.Infof("Field fetched successfully: %s", fieldID)

	res := dto.FieldFullResponse{
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
	utils.Log.Infof("Updating field: %s", req.FieldID)

	if _, err := uuid.Parse(req.FieldID); err != nil {
		utils.Log.Errorf("Invalid field UUID: %v", err)
		return dto.FieldResponse{}, constants.ErrInvalidUUID
	}

	field, _, err := fs.fieldRepo.GetFieldByID(ctx, nil, req.FieldID)
	if err != nil {
		utils.Log.Errorf("Field not found: %v", err)
		return dto.FieldResponse{}, constants.ErrGetFieldByID
	}

	if req.FieldName != "" {
		field.FieldName = req.FieldName
	}
	if req.FieldAddress != "" {
		field.FieldAddress = req.FieldAddress
	}
	if req.FieldPrice < 0 {
		utils.Log.Warn("Invalid field price: < 0")
		return dto.FieldResponse{}, constants.ErrInvalidFieldPrice
	}
	if req.FieldPrice != 0 {
		field.FieldPrice = req.FieldPrice
	}

	if req.FieldImage != nil {
		imageName, err := helpers.SaveImage(req.FieldImage, "./assets/fields", "field")
		if err != nil {
			utils.Log.Errorf("Failed to save new image: %v", err)
			return dto.FieldResponse{}, constants.ErrSaveImages
		}
		field.FieldImage = imageName
	}

	if err := fs.fieldRepo.UpdateField(ctx, nil, field); err != nil {
		utils.Log.Errorf("Failed to update field: %v", err)
		return dto.FieldResponse{}, constants.ErrUpdateField
	}

	utils.Log.Infof("Field updated successfully: %s", req.FieldID)

	res := dto.FieldResponse{
		FieldID:      field.FieldID,
		FieldName:    field.FieldName,
		FieldAddress: field.FieldAddress,
		FieldPrice:   field.FieldPrice,
		FieldImage:   field.FieldImage,
		CategoryID:   field.CategoryID,
	}

	return res, nil
}

func (fs *FieldService) DeleteField(ctx context.Context, req dto.DeleteFieldRequest) (dto.FieldResponse, error) {
	utils.Log.Infof("Deleting field: %s", req.FieldID)

	if _, err := uuid.Parse(req.FieldID); err != nil {
		utils.Log.Errorf("Invalid UUID for deletion: %v", err)
		return dto.FieldResponse{}, constants.ErrInvalidUUID
	}

	deletedField, _, err := fs.fieldRepo.GetFieldByID(ctx, nil, req.FieldID)
	if err != nil {
		utils.Log.Errorf("Field not found for deletion: %v", err)
		return dto.FieldResponse{}, constants.ErrGetFieldByID
	}

	err = fs.fieldRepo.DeleteField(ctx, nil, req.FieldID)
	if err != nil {
		utils.Log.Errorf("Failed to delete field: %v", err)
		return dto.FieldResponse{}, constants.ErrDeleteFieldByID
	}

	utils.Log.Infof("Field deleted successfully: %s", req.FieldID)

	res := dto.FieldResponse{
		FieldID:      deletedField.FieldID,
		FieldName:    deletedField.FieldName,
		FieldAddress: deletedField.FieldAddress,
		FieldPrice:   deletedField.FieldPrice,
		FieldImage:   deletedField.FieldImage,
		CategoryID:   deletedField.CategoryID,
	}

	return res, nil
}
