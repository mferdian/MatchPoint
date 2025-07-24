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
	ICategoryController interface {
		CreateCategory(ctx *gin.Context)
		GetAllCatgory(ctx *gin.Context)
		GetCategoryByID(ctx *gin.Context)
		UpdateCategory(ctx *gin.Context)
		DeleteCategory(ctx *gin.Context)
	}

	CategoryController struct {
		categoryService service.ICategoryService
	}
)

func NewCategoryController(categoryService service.ICategoryService) *CategoryController {
	return &CategoryController{
		categoryService: categoryService,
	}
}

func (cc *CategoryController) CreateCategory(ctx *gin.Context) {
	var payload dto.CreateCategoryRequest
	if err := ctx.ShouldBind(&payload); err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := cc.categoryService.CreateCategory(ctx, payload)
	if err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_CREATE_CATEGORY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(constants.MESSAGE_SUCCESS_CREATE_CATEGORY, result)
	ctx.AbortWithStatusJSON(http.StatusCreated, res)
}

func (cc *CategoryController) GetAllCatgory(ctx *gin.Context) {
	var payload dto.CategoryPaginationRequest
	if err := ctx.ShouldBind(&payload); err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := cc.categoryService.GetAllCategoryWithPagination(ctx.Request.Context(), payload)
	if err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_GET_ALL_CATEGORY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.Response{
		Status:   true,
		Messsage: constants.MESSAGE_SUCCESS_GET_ALL_CATEGORY,
		Data:     result.Data,
		Meta:     result.PaginationResponse,
	}

	ctx.JSON(http.StatusOK, res)
}
func (cc *CategoryController) GetCategoryByID(ctx *gin.Context) {
	categoryID := ctx.Param("id")

	if _, err := uuid.Parse(categoryID); err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_UUID_FORMAT, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := cc.categoryService.GetCategoryByID(ctx.Request.Context(), categoryID)
	if err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_GET_DETAIL_CATEGORY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusNotFound, res)
		return
	}

	res := utils.BuildResponseSuccess(constants.MESSAGE_SUCCESS_GET_DETAIL_CATEGORY, result)
	ctx.JSON(http.StatusOK, res)
}
func (cc *CategoryController) UpdateCategory(ctx *gin.Context) {
	categoryID := ctx.Param("id")

	role := ctx.GetString("role")

	if role != constants.ENUM_ROLE_ADMIN {
		res := utils.BuildResponseFailed(constants.ErrDeniedAccess.Error(), "only admin can update category", nil)
		ctx.AbortWithStatusJSON(http.StatusForbidden, res)
		return
	}

	var payload dto.UpdateCategoryRequest
	payload.ID = categoryID

	if err := ctx.ShouldBind(&payload); err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := cc.categoryService.UpdateCategory(ctx.Request.Context(), payload)
	if err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_UPDATE_CATEGORY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(constants.MESSAGE_SUCCESS_UPDATE_CATEGORY, result)
	ctx.JSON(http.StatusOK, res)
}
func (cc *CategoryController) DeleteCategory(ctx *gin.Context) {
	categoryID := ctx.Param("id")
	

	role := ctx.GetString("role")

	if role != constants.ENUM_ROLE_ADMIN {
		res := utils.BuildResponseFailed(constants.ErrDeniedAccess.Error(), "only admin can delete category", nil)
		ctx.AbortWithStatusJSON(http.StatusForbidden, res)
		return
	}

	var payload dto.DeleteCategoryRequest
	payload.ID = categoryID
	if err := ctx.ShouldBind(&payload); err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := cc.categoryService.DeleteCategory(ctx.Request.Context(), payload)
	if err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_DELETE_CATEGORY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(constants.MESSAGE_SUCCESS_DELETE_CATEGORY, result)
	ctx.JSON(http.StatusOK, res)
}
