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
	IBookingController interface {
		CreateBooking(ctx *gin.Context)
		GetAllBooking(ctx *gin.Context)
		GetBookingByID(ctx *gin.Context)
		UpdateBooking(ctx *gin.Context)
		DeleteBooking(ctx *gin.Context)
	}

	BookingController struct {
		bookingService service.IBookingService
	}
)

func NewBookingController(bookingService service.IBookingService) *BookingController {
	return &BookingController{
		bookingService: bookingService,
	}
}

// POST /bookings
func (bc *BookingController) CreateBooking(ctx *gin.Context) {
	var payload dto.CreateBookingRequest
	if err := ctx.ShouldBind(&payload); err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := bc.bookingService.CreateBooking(ctx.Request.Context(), payload)
	if err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_CREATE_BOOKING, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(constants.MESSAGE_SUCCESS_CREATE_BOOKING, result)
	ctx.JSON(http.StatusCreated, res)
}

// GET /bookings
func (bc *BookingController) GetAllBooking(ctx *gin.Context) {
	var payload dto.BookingPaginationRequest
	if err := ctx.ShouldBind(&payload); err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := bc.bookingService.GetAllBooking(ctx.Request.Context(), payload)
	if err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_GET_ALL_BOOKING, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(constants.MESSAGE_SUCCESS_GET_ALL_BOOKING, result)
	ctx.JSON(http.StatusOK, res)
}

// GET /bookings/:id
func (bc *BookingController) GetBookingByID(ctx *gin.Context) {
	bookingID := ctx.Param("id")

	if _, err := uuid.Parse(bookingID); err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_UUID_FORMAT, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := bc.bookingService.GetBookingByID(ctx.Request.Context(), bookingID)
	if err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_GET_DETAIL_BOOKING, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusNotFound, res)
		return
	}

	res := utils.BuildResponseSuccess(constants.MESSAGE_SUCCESS_GET_DETAIL_BOOKING, result)
	ctx.JSON(http.StatusOK, res)
}

// PUT /bookings/:id
func (bc *BookingController) UpdateBooking(ctx *gin.Context) {
	bookingID := ctx.Param("id")

	if _, err := uuid.Parse(bookingID); err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_UUID_FORMAT, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var payload dto.UpdateBookingRequest
	payload.BookingID = bookingID

	if err := ctx.ShouldBind(&payload); err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := bc.bookingService.UpdateBooking(ctx.Request.Context(), payload)
	if err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_UPDATE_BOOKING, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(constants.MESSAGE_SUCCESS_UPDATE_BOOKING, result)
	ctx.JSON(http.StatusOK, res)
}

// DELETE /bookings/:id
func (bc *BookingController) DeleteBooking(ctx *gin.Context) {
	bookingID := ctx.Param("id")

	if _, err := uuid.Parse(bookingID); err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_UUID_FORMAT, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var payload dto.DeleteBookingRequest
	payload.BookingID = bookingID

	result, err := bc.bookingService.DeleteBooking(ctx.Request.Context(), payload)
	if err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_DELETE_BOOKING, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(constants.MESSAGE_SUCCESS_DELETE_BOOKING, result)
	ctx.JSON(http.StatusOK, res)
}
