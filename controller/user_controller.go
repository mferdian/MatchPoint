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
	IUserController interface {
		CreateUser(ctx *gin.Context)
		GetUserByEmail(ctx *gin.Context)
		GetUserByID(ctx *gin.Context)
		GetAllUser(ctx *gin.Context)
		UpdateUser(ctx *gin.Context)
		DeleteUser(ctx *gin.Context)
	}

	UserController struct {
		userService service.IUserService
	}
)

func NewUserController(userService service.IUserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (uh *UserController) CreateUser(ctx *gin.Context) {
	var payload dto.CreateUserRequest
	if err := ctx.ShouldBind(&payload); err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := uh.userService.CreateUser(ctx, payload)
	if err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_CREATE_USER, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(constants.MESSAGE_SUCCESS_CREATE_USER, result)
	ctx.AbortWithStatusJSON(http.StatusCreated, res)
}

func (uh *UserController) GetUserByEmail(ctx *gin.Context) {
	var payload dto.LoginUserRequest
	if err := ctx.ShouldBind(&payload); err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := uh.userService.GetUserByEmail(ctx, payload)
	if err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_LOGIN_USER, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(constants.MESSAGE_SUCCESS_CREATE_USER, result)
	ctx.AbortWithStatusJSON(http.StatusOK, res)
}

func (uc *UserController) GetUserByID(ctx *gin.Context) {
	idStr := ctx.Param("id")

	userID := ctx.GetString("user_id")
	role := ctx.GetString("role")

	if role == constants.ENUM_ROLE_USER && userID != idStr {
		res := utils.BuildResponseFailed("unauthorized", "you can only get your own account", nil)
		ctx.AbortWithStatusJSON(http.StatusForbidden, res)
		return
	}

	if _, err := uuid.Parse(idStr); err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_UUID_FORMAT, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := uc.userService.GetuserByID(ctx.Request.Context(), idStr)
	if err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_GET_DETAIL_USER, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusNotFound, res)
		return
	}

	res := utils.BuildResponseSuccess(constants.MESSAGE_SUCCESS_GET_DETAIL_USER, result)
	ctx.JSON(http.StatusOK, res)
}

func (uh *UserController) GetAllUser(ctx *gin.Context) {
	var payload dto.UserPaginationRequest
	if err := ctx.ShouldBind(&payload); err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := uh.userService.GetAllUserWithPagination(ctx.Request.Context(), payload)
	if err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_GET_LIST_USER, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.Response{
		Status:   true,
		Messsage: constants.MESSAGE_SUCCESS_GET_LIST_USER,
		Data:     result.Data,
		Meta:     result.PaginationResponse,
	}

	ctx.JSON(http.StatusOK, res)
}

func (uh *UserController) UpdateUser(ctx *gin.Context) {
	idStr := ctx.Param("id")

	userID := ctx.GetString("user_id")
	role := ctx.GetString("role")

	if role == constants.ENUM_ROLE_USER && userID != idStr {
		res := utils.BuildResponseFailed("unauthorized", "you can only update your own account", nil)
		ctx.AbortWithStatusJSON(http.StatusForbidden, res)
		return
	}

	var payload dto.UpdateUserRequest
	payload.ID = idStr

	if err := ctx.ShouldBind(&payload); err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := uh.userService.UpdateUser(ctx.Request.Context(), payload)
	if err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_UPDATE_USER, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(constants.MESSAGE_SUCCESS_UPDATE_USER, result)
	ctx.JSON(http.StatusOK, res)
}

func (uh *UserController) DeleteUser(ctx *gin.Context) {
	idStr := ctx.Param("id")

	userID := ctx.GetString("user_id")
	role := ctx.GetString("role")

	if role == constants.ENUM_ROLE_USER && userID != idStr {
		res := utils.BuildResponseFailed("unauthorized", "you can only delete your own account", nil)
		ctx.AbortWithStatusJSON(http.StatusForbidden, res)
		return
	}

	var payload dto.DeleteUserRequest
	payload.UserID = idStr
	if err := ctx.ShouldBind(&payload); err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := uh.userService.DeleteUser(ctx.Request.Context(), payload)
	if err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_DELETE_USER, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(constants.MESSAGE_SUCCESS_DELETE_USER, result)
	ctx.JSON(http.StatusOK, res)
}
