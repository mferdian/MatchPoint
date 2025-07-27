package controller

import (
	"fieldreserve/constants"
	"fieldreserve/dto"
	"fieldreserve/service"
	"fieldreserve/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type (
	IScheduleController interface {
		CreateSchedule(ctx *gin.Context)
		GetAllSchedule(ctx *gin.Context)
		UpdateSchedule(ctx *gin.Context)
		DeleteSchedule(ctx *gin.Context)
		GetScheduleByID(ctx *gin.Context)
		GetSchedulesByFieldID(ctx *gin.Context)
		GetScheduleByFieldIDAndDay(ctx *gin.Context)
	}

	ScheduleController struct {
		scheduleService service.IScheduleService
	}
)

func NewScheduleController(scheduleService service.IScheduleService) *ScheduleController {
	return &ScheduleController{
		scheduleService: scheduleService,
	}
}

func (sc *ScheduleController) CreateSchedule(ctx *gin.Context) {
	var payload dto.CreateScheduleRequest
	if err := ctx.ShouldBind(&payload); err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := sc.scheduleService.CreateSchedule(ctx.Request.Context(), payload)
	if err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_CREATE_SCHEDULE, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(constants.MESSAGE_SUCCESS_CREATE_SCHEDULE, result)
	ctx.JSON(http.StatusCreated, res)
}

func (sc *ScheduleController) GetAllSchedule(ctx *gin.Context) {
	result, err := sc.scheduleService.GetAllSchedule(ctx.Request.Context())
	if err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_GET_ALL_SCHEDULE, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(constants.MESSAGE_SUCCESS_GET_ALL_SCHEDULE, result)
	ctx.JSON(http.StatusOK, res)
}

func (sc *ScheduleController) UpdateSchedule(ctx *gin.Context) {
	scheduleID := ctx.Param("id")
	role := ctx.GetString("role")

	if role != constants.ENUM_ROLE_ADMIN {
		res := utils.BuildResponseFailed(constants.ErrDeniedAccess.Error(), "only admin can update schedule", nil)
		ctx.AbortWithStatusJSON(http.StatusForbidden, res)
		return
	}

	var payload dto.UpdateScheduleRequest
	payload.ScheduleID = scheduleID

	if err := ctx.ShouldBind(&payload); err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := sc.scheduleService.UpdateSchedule(ctx.Request.Context(), payload)
	if err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_UPDATE_SCHEDULE, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(constants.MESSAGE_SUCCESS_UPDATE_SCHEDULE, result)
	ctx.JSON(http.StatusOK, res)
}

func (sc *ScheduleController) DeleteSchedule(ctx *gin.Context) {
	scheduleID := ctx.Param("id")
	role := ctx.GetString("role")

	if role != constants.ENUM_ROLE_ADMIN {
		res := utils.BuildResponseFailed(constants.ErrDeniedAccess.Error(), "only admin can delete schedule", nil)
		ctx.AbortWithStatusJSON(http.StatusForbidden, res)
		return
	}

	var payload dto.DeleteScheduleRequest
	payload.ScheduleID = scheduleID

	if err := ctx.ShouldBind(&payload); err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := sc.scheduleService.DeleteScheduleByID(ctx.Request.Context(), payload)
	if err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_DELETE_SCHEDULE, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(constants.MESSAGE_SUCCESS_DELETE_SCHEDULE, result)
	ctx.JSON(http.StatusOK, res)
}

func (sc *ScheduleController) GetScheduleByID(ctx *gin.Context) {
	scheduleID := ctx.Param("id")

	if _, err := uuid.Parse(scheduleID); err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_UUID_FORMAT, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := sc.scheduleService.GetScheduleByID(ctx.Request.Context(), scheduleID)
	if err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_GET_DETAIL_SCHEDULE, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusNotFound, res)
		return
	}

	res := utils.BuildResponseSuccess(constants.MESSAGE_SUCCESS_GET_DETAIL_SCHEDULE, result)
	ctx.JSON(http.StatusOK, res)
}

func (sc *ScheduleController) GetSchedulesByFieldID(ctx *gin.Context) {
	fieldID := ctx.Param("field_id")

	if _, err := uuid.Parse(fieldID); err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_UUID_FORMAT, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := sc.scheduleService.GetSchedulesByFieldID(ctx.Request.Context(), fieldID)
	if err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_GET_ALL_SCHEDULE, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusNotFound, res)
		return
	}

	res := utils.BuildResponseSuccess(constants.MESSAGE_SUCCESS_GET_ALL_SCHEDULE, result)
	ctx.JSON(http.StatusOK, res)
}

func (sc *ScheduleController) GetScheduleByFieldIDAndDay(ctx *gin.Context) {
	fieldID := ctx.Param("field_id")
	dayStr := ctx.Param("day")

	if _, err := uuid.Parse(fieldID); err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_UUID_FORMAT, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	day, err := strconv.Atoi(dayStr)
	if err != nil || day < 0 || day > 6 {
		res := utils.BuildResponseFailed("invalid day parameter", "day must be between 0 and 6", nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := sc.scheduleService.GetScheduleByFieldIDAndDay(ctx.Request.Context(), fieldID, day)
	if err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_GET_DETAIL_SCHEDULE, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusNotFound, res)
		return
	}

	res := utils.BuildResponseSuccess(constants.MESSAGE_SUCCESS_GET_DETAIL_SCHEDULE, result)
	ctx.JSON(http.StatusOK, res)
}