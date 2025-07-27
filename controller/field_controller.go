package controller

import (
	"fieldreserve/constants"
	"fieldreserve/dto"
	"fieldreserve/service"
	"fieldreserve/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type (
	IFieldController interface {
		CreateField(ctx *gin.Context)
		GetAllField(ctx *gin.Context)
		GetFieldByID(ctx *gin.Context)
		UpdateField(ctx *gin.Context)
		DeleteField(ctx *gin.Context)
	}

	FieldController struct {
		fieldService service.IFieldService
	}
)

func NewFieldController(fieldService service.IFieldService) *FieldController {
	return &FieldController{
		fieldService: fieldService,
	}
}

func (fc *FieldController) CreateField(ctx *gin.Context) {
	var payload dto.CreateFieldRequest
	if err := ctx.ShouldBind(&payload); err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := fc.fieldService.CreateField(ctx, payload)
	if err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_CREATE_FIELD, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(constants.MESSAGE_SUCCESS_CREATE_FIELD, result)
	ctx.JSON(http.StatusCreated, res)
}
func (fc *FieldController) GetAllField(ctx *gin.Context) {
	var payload dto.FieldPaginationRequest
	if err := ctx.ShouldBind(&payload); err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := fc.fieldService.GetAllFieldWithPagination(ctx.Request.Context(), payload)
	if err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_GET_ALL_FIELD, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	
	res := utils.BuildResponseSuccess(constants.MESSAGE_SUCCESS_GET_ALL_FIELD, result)
	ctx.JSON(http.StatusOK, res)
}
func (fc *FieldController) GetFieldByID(ctx *gin.Context) {
	fieldID := ctx.Param("id")

	if _, err := uuid.Parse(fieldID); err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_UUID_FORMAT, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := fc.fieldService.GetFieldByID(ctx.Request.Context(), fieldID)
	if err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_GET_DETAIL_FIELD, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusNotFound, res)
		return
	}

	res := utils.BuildResponseSuccess(constants.MESSAGE_SUCCESS_GET_DETAIL_FIELD, result)
	ctx.JSON(http.StatusOK, res)
}
func (fc *FieldController) UpdateField(ctx *gin.Context) {
	fieldID := ctx.Param("id")
	role := ctx.GetString("role")

	if role != constants.ENUM_ROLE_ADMIN {
		res := utils.BuildResponseFailed(constants.ErrDeniedAccess.Error(), "only admin can update field", nil)
		ctx.AbortWithStatusJSON(http.StatusForbidden, res)
		return
	}

	var payload dto.UpdateFieldRequest
	payload.FieldID = fieldID

	if err := ctx.ShouldBind(&payload); err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := fc.fieldService.UpdateField(ctx.Request.Context(), payload)
	if err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_UPDATE_FIELD, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(constants.MESSAGE_SUCCESS_UPDATE_FIELD, result)
	ctx.JSON(http.StatusOK, res)
}
func (fc *FieldController) DeleteField(ctx *gin.Context) {
	fieldID := ctx.Param("id")

	role := ctx.GetString("role")

	if role != constants.ENUM_ROLE_ADMIN {
		res := utils.BuildResponseFailed(constants.ErrDeniedAccess.Error(), "only admin can delete field", nil)
		ctx.AbortWithStatusJSON(http.StatusForbidden, res)
		return
	}

	var payload dto.DeleteFieldRequest
	payload.FieldID = fieldID
	if err := ctx.ShouldBind(&payload); err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := fc.fieldService.DeleteField(ctx.Request.Context(), payload)
	if err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_DELETE_FIELD, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(constants.MESSAGE_SUCCESS_DELETE_FIELD, result)
	ctx.JSON(http.StatusOK, res)
}
