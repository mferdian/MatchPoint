package controller

import (
	"fieldreserve/constants"
	"fieldreserve/dto"
	"fieldreserve/service"
	"fieldreserve/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type (
	IBookingController interface {
		CreateBooking(ctx *gin.Context)
		GetAllBooking(ctx *gin.Context)
		GetBookingByID(ctx *gin.Context)
		GetUserBookingHistory(ctx *gin.Context)
		UpdateStatusBooking(ctx *gin.Context)
		DeleteBooking(ctx *gin.Context)
		DownloadInvoice(ctx *gin.Context)
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

func (bc *BookingController) GetUserBookingHistory(ctx *gin.Context) {
	var payload dto.BookingPaginationRequest
	if err := ctx.ShouldBindQuery(&payload); err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	userID := ctx.GetString("user_id")
	payload.UserID = userID

	result, err := bc.bookingService.GetUserBookingHistory(ctx.Request.Context(), payload)
	if err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_GET_ALL_BOOKING, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.Response{
		Status:   true,
		Messsage: constants.MESSAGE_SUCCESS_GET_ALL_BOOKING,
		Data:     result.Data,
		Meta:     result.PaginationResponse,
	}
	ctx.JSON(http.StatusOK, res)
}


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


func (bc *BookingController) UpdateStatusBooking(ctx *gin.Context) {
	bookingID := ctx.Param("id")

	if _, err := uuid.Parse(bookingID); err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_UUID_FORMAT, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var payload dto.UpdateBookingStatusRequest
	payload.BookingID = bookingID

	if err := ctx.ShouldBind(&payload); err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := bc.bookingService.UpdateBookingStatus(ctx.Request.Context(), payload)
	if err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_UPDATE_BOOKING, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(constants.MESSAGE_SUCCESS_UPDATE_BOOKING, result)
	ctx.JSON(http.StatusOK, res)
}

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


func (bc *BookingController) DownloadInvoice(ctx *gin.Context) {
	bookingID := ctx.Param("id")

	if _, err := uuid.Parse(bookingID); err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_UUID_FORMAT, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	booking, err := bc.bookingService.GetBookingByID(ctx.Request.Context(), bookingID)
	if err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_GET_BOOKING, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusNotFound, res)
		return
	}

	pdfBytes, err := utils.GenerateInvoicePDF(booking)
	if err != nil {
		res := utils.BuildResponseFailed("Failed to generate invoice", err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	ctx.Header("Content-Type", "application/pdf")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=invoice-%s.pdf", booking.BookingID))
	ctx.Data(http.StatusOK, "application/pdf", pdfBytes)
}
