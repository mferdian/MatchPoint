package service

import (
	"context"
	"fieldreserve/constants"
	"fieldreserve/dto"
	"fieldreserve/helpers"
	"fieldreserve/model"
	"fieldreserve/repository"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type (
	IBookingService interface {
		CreateBooking(ctx context.Context, req dto.CreateBookingRequest) (dto.BookingResponse, error)
		GetAllBooking(ctx context.Context, req dto.BookingPaginationRequest) (dto.BookingPaginationResponse, error)
		GetBookingByID(ctx context.Context, bookingID string) (dto.BookingFullResponse, error)
		UpdateBooking(ctx context.Context, req dto.UpdateBookingRequest) (dto.BookingResponse, error)
		DeleteBooking(ctx context.Context, req dto.DeleteBookingRequest) (dto.BookingResponse, error)
	}

	BookingService struct {
		bookingRepo  repository.IBookingRepository
		jwtService   InterfaceJWTService
		scheduleRepo repository.IScheduleRepository
	}
)

func NewBookingService(
	bookingRepo repository.IBookingRepository,
	jwtService InterfaceJWTService,
	scheduleRepo repository.IScheduleRepository,
) *BookingService {
	return &BookingService{
		bookingRepo:  bookingRepo,
		jwtService:   jwtService,
		scheduleRepo: scheduleRepo,
	}
}

func (bs *BookingService) CreateBooking(ctx context.Context, req dto.CreateBookingRequest) (dto.BookingResponse, error) {
	loc := helpers.GetAppLocation()

	// Ambil token dari context
	tokenString, ok := ctx.Value("token").(string)
	if !ok || tokenString == "" {
		return dto.BookingResponse{}, constants.ErrUnauthorized
	}

	userIDStr, err := bs.jwtService.GetUserIDByToken(tokenString)
	if err != nil {
		return dto.BookingResponse{}, constants.ErrUnauthorized
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return dto.BookingResponse{}, constants.ErrInvalidUUID
	}

	fieldID, err := uuid.Parse(req.FieldID)
	if err != nil {
		return dto.BookingResponse{}, constants.ErrInvalidUUID
	}

	// Parsing tanggal dan waktu booking (string)
	bookingDate, err := time.ParseInLocation("2006-01-02", req.BookingDate, loc)
	if err != nil {
		return dto.BookingResponse{}, constants.ErrInvalidBookingDate
	}

	startParsed, err := time.ParseInLocation("15:04", req.StartTime, loc)
	if err != nil {
		return dto.BookingResponse{}, constants.ErrInvalidTimeFormat
	}
	endParsed, err := time.ParseInLocation("15:04", req.EndTime, loc)
	if err != nil {
		return dto.BookingResponse{}, constants.ErrInvalidTimeFormat
	}

	// Gabungkan tanggal booking + jam
	startTime := time.Date(bookingDate.Year(), bookingDate.Month(), bookingDate.Day(),
		startParsed.Hour(), startParsed.Minute(), 0, 0, loc)

	endTime := time.Date(bookingDate.Year(), bookingDate.Month(), bookingDate.Day(),
		endParsed.Hour(), endParsed.Minute(), 0, 0, loc)

	// Validasi minimal 2 jam dari sekarang
	now := time.Now().In(loc)
	if startTime.Before(now.Add(2 * time.Hour)) {
		return dto.BookingResponse{}, constants.ErrBookingTooSoon
	}

	// Ambil schedule field di hari tersebut
	dayOfWeek := int(bookingDate.Weekday()) // 0 = Sunday, 1 = Monday, ..., 6 = Saturday
	schedule, err := bs.scheduleRepo.GetScheduleByFieldIDAndDay(ctx, nil, req.FieldID, dayOfWeek)
	if err != nil {
		return dto.BookingResponse{}, constants.ErrScheduleNotFound
	}

	// Buat open-close time dari schedule (dalam zona waktu Asia/Jakarta)
	openTime := time.Date(bookingDate.Year(), bookingDate.Month(), bookingDate.Day(),
		schedule.OpenTime.Hour(), schedule.OpenTime.Minute(), 0, 0, loc)
	closeTime := time.Date(bookingDate.Year(), bookingDate.Month(), bookingDate.Day(),
		schedule.CloseTime.Hour(), schedule.CloseTime.Minute(), 0, 0, loc)

	if startTime.Before(openTime) || endTime.After(closeTime) {
		return dto.BookingResponse{}, constants.ErrOutsideOperatingHours
	}

	// Validasi tabrakan waktu booking
	overlap, err := bs.bookingRepo.CheckBookingOverlap(ctx, nil, fieldID, bookingDate, startTime, endTime)
	if err != nil {
		return dto.BookingResponse{}, constants.ErrCheckOverlap
	}

	if overlap {
		return dto.BookingResponse{}, constants.ErrBookingOverlap
	}

	// Buat model booking
	booking := model.Booking{
		BookingID:     uuid.New(),
		UserID:        userID,
		FieldID:       fieldID,
		PaymentMethod: req.PaymentMethod,
		BookingDate:   bookingDate,
		StartTime:     startTime,
		EndTime:       endTime,
		Status:        "pending",
	}

	if err := bs.bookingRepo.CreateBooking(ctx, nil, booking); err != nil {
		return dto.BookingResponse{}, constants.ErrCreateBooking
	}

	// Build response
	return dto.BookingResponse{
		BookingID:     booking.BookingID,
		UserID:        booking.UserID,
		FieldID:       booking.FieldID,
		PaymentMethod: booking.PaymentMethod,
		BookingDate:   booking.BookingDate,
		StartTime:     booking.StartTime,
		EndTime:       booking.EndTime,
		Status:        booking.Status,
	}, nil
}

func (bs *BookingService) GetAllBooking(ctx context.Context, req dto.BookingPaginationRequest) (dto.BookingPaginationResponse, error) {
	dataWithPaginate, err := bs.bookingRepo.GetAllBooking(ctx, nil, req)
	if err != nil {
		return dto.BookingPaginationResponse{}, constants.ErrGetAllField
	}

	var datas []dto.BookingResponse
	for _, booking := range dataWithPaginate.Bookings {
		data := dto.BookingResponse{
			BookingID:     booking.BookingID,
			UserID:        booking.UserID,
			FieldID:       booking.FieldID,
			PaymentMethod: booking.PaymentMethod,
			BookingDate:   booking.BookingDate,
			StartTime:     booking.StartTime,
			EndTime:       booking.EndTime,
			Status:        booking.Status,
		}

		datas = append(datas, data)
	}

	return dto.BookingPaginationResponse{
		Data: datas,
		PaginationResponse: dto.PaginationResponse{
			Page:    dataWithPaginate.Page,
			PerPage: dataWithPaginate.PerPage,
			MaxPage: dataWithPaginate.MaxPage,
			Count:   dataWithPaginate.Count,
		},
	}, nil
}

func (bs *BookingService) GetBookingByID(ctx context.Context, bookingID string) (dto.BookingFullResponse, error) {
	if _, err := uuid.Parse(bookingID); err != nil {
		return dto.BookingFullResponse{}, constants.ErrInvalidUUID
	}

	booking, _, err := bs.bookingRepo.GetBookingByID(ctx, nil, bookingID)
	if err != nil {
		return dto.BookingFullResponse{}, constants.ErrGetFieldByID
	}

	field := booking.Field
	category := field.Category

	fieldDTO := dto.FieldCompactResponse{
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
	}

	res := dto.BookingFullResponse{
		BookingID:     booking.BookingID,
		UserID:        booking.UserID,
		PaymentMethod: booking.PaymentMethod,
		BookingDate:   booking.BookingDate,
		StartTime:     booking.StartTime,
		EndTime:       booking.EndTime,
		Status:        booking.Status,
		Field:         fieldDTO,
	}

	return res, nil
}

func (bs *BookingService) UpdateBooking(ctx context.Context, req dto.UpdateBookingRequest) (dto.BookingResponse, error) {
	if _, err := uuid.Parse(req.BookingID); err != nil {
		return dto.BookingResponse{}, constants.ErrInvalidUUID
	}

	booking, _, err := bs.bookingRepo.GetBookingByID(ctx, nil, req.BookingID)
	if err != nil {
		return dto.BookingResponse{}, constants.ErrGetBookingByID
	}

	loc, _ := time.LoadLocation("Asia/Jakarta")

	if req.FieldID != nil {
		fieldUUID, err := uuid.Parse(*req.FieldID)
		if err != nil {
			return dto.BookingResponse{}, constants.ErrInvalidUUID
		}
		booking.FieldID = fieldUUID
	}

	if req.PaymentMethod != nil {
		booking.PaymentMethod = *req.PaymentMethod
	}

	if req.BookingDate != nil {
		date, err := time.ParseInLocation("2006-01-02", *req.BookingDate, loc)
		if err != nil {
			return dto.BookingResponse{}, constants.ErrInvalidBookingDate
		}
		booking.BookingDate = date
	}

	if req.StartTime != nil {
		startDateTimeStr := fmt.Sprintf("%sT%s:00", booking.BookingDate.Format("2006-01-02"), *req.StartTime)
		startDateTime, err := time.ParseInLocation("2006-01-02T15:04:05", startDateTimeStr, loc)
		if err != nil {
			return dto.BookingResponse{}, constants.ErrInvalidStartTime
		}
		booking.StartTime = startDateTime
	}

	if req.EndTime != nil {
		endDateTimeStr := fmt.Sprintf("%sT%s:00", booking.BookingDate.Format("2006-01-02"), *req.EndTime)
		endDateTime, err := time.ParseInLocation("2006-01-02T15:04:05", endDateTimeStr, loc)
		if err != nil {
			return dto.BookingResponse{}, fmt.Errorf("invalid end_time format: %w", err)
		}
		booking.EndTime = endDateTime
	}

	if req.Status != nil {
		booking.Status = *req.Status
	}

	if err := bs.bookingRepo.UpdateBooking(ctx, nil, booking); err != nil {
		return dto.BookingResponse{}, constants.ErrUpdateBooking
	}

	return dto.BookingResponse{
		BookingID:     booking.BookingID,
		UserID:        booking.UserID,
		FieldID:       booking.FieldID,
		PaymentMethod: booking.PaymentMethod,
		BookingDate:   booking.BookingDate,
		StartTime:     booking.StartTime,
		EndTime:       booking.EndTime,
		Status:        booking.Status,
	}, nil
}

func (bs *BookingService) DeleteBooking(ctx context.Context, req dto.DeleteBookingRequest) (dto.BookingResponse, error) {
	if _, err := uuid.Parse(req.BookingID); err != nil {
		return dto.BookingResponse{}, constants.ErrInvalidUUID
	}

	// Ambil booking terlebih dahulu
	booking, _, err := bs.bookingRepo.GetBookingByID(ctx, nil, req.BookingID)
	if err != nil {
		return dto.BookingResponse{}, constants.ErrGetBookingByID
	}

	// Validasi: hanya bisa dibatalkan maksimal 3 jam sebelum startTime
	if time.Until(booking.StartTime) < 3*time.Hour {
		return dto.BookingResponse{}, constants.ErrCannotCancelLate
	}

	// Lakukan soft delete atau ubah status
	if err := bs.bookingRepo.DeleteBooking(ctx, nil, req.BookingID); err != nil {
		return dto.BookingResponse{}, constants.ErrDeleteBooking
	}

	return dto.BookingResponse{
		BookingID:     booking.BookingID,
		UserID:        booking.UserID,
		FieldID:       booking.FieldID,
		PaymentMethod: booking.PaymentMethod,
		BookingDate:   booking.BookingDate,
		StartTime:     booking.StartTime,
		EndTime:       booking.EndTime,
		Status:        booking.Status,
	}, nil
}
