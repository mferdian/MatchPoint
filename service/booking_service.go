package service

import (
	"context"
	"fieldreserve/constants"
	"fieldreserve/dto"
	"fieldreserve/model"
	"fieldreserve/repository"
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
	tokenString, ok := ctx.Value("token").(string)
	if !ok || tokenString == "" {
		return dto.BookingResponse{}, constants.ErrUnauthorized
	}

	userIDstr, err := bs.jwtService.GetUserIDByToken(tokenString)
	if err != nil {
		return dto.BookingResponse{}, constants.ErrUnauthorized
	}

	userID, _ := uuid.Parse(userIDstr)

	// Parse fieldID
	fieldUUID, err := uuid.Parse(req.FieldID)
	if err != nil {
		return dto.BookingResponse{}, constants.ErrInvalidUUID
	}

	// Parse tanggal dan jam
	bookingDate, err := time.Parse("2006-01-02", req.BookingDate)
	if err != nil {
		return dto.BookingResponse{}, constants.ErrInvalidBookingDate
	}
	startParsed, err := time.Parse("15:04", req.StartTime)
	if err != nil {
		return dto.BookingResponse{}, constants.ErrInvalidTimeFormat
	}
	endParsed, err := time.Parse("15:04", req.EndTime)
	if err != nil {
		return dto.BookingResponse{}, constants.ErrInvalidTimeFormat
	}

	// Gabungkan bookingDate + jam
	startTime := time.Date(bookingDate.Year(), bookingDate.Month(), bookingDate.Day(), startParsed.Hour(), startParsed.Minute(), 0, 0, time.UTC)
	endTime := time.Date(bookingDate.Year(), bookingDate.Month(), bookingDate.Day(), endParsed.Hour(), endParsed.Minute(), 0, 0, time.UTC)

	// Validasi: minimal 2 jam dari sekarang
	if startTime.Before(time.Now().Add(2 * time.Hour)) {
		return dto.BookingResponse{}, constants.ErrBookingTooSoon
	}

	// Ambil schedule dari field & day
	dayOfWeek := int(bookingDate.Weekday()) // 0: Sunday, 1: Monday, ..., 6: Saturday
	schedule, err := bs.scheduleRepo.GetScheduleByFieldIDAndDay(ctx, nil, req.FieldID, dayOfWeek)
	if err != nil {
		return dto.BookingResponse{}, constants.ErrScheduleNotFound
	}

	// Validasi dalam jam operasional
	open := time.Date(startTime.Year(), startTime.Month(), startTime.Day(), schedule.OpenTime.Hour(), schedule.OpenTime.Minute(), 0, 0, time.UTC)
	close := time.Date(startTime.Year(), startTime.Month(), startTime.Day(), schedule.CloseTime.Hour(), schedule.CloseTime.Minute(), 0, 0, time.UTC)

	if startTime.Before(open) || endTime.After(close) {
		return dto.BookingResponse{}, constants.ErrOutsideOperatingHours
	}

	// Simpan booking
	booking := model.Booking{
		BookingID:     uuid.New(),
		UserID:        userID,
		FieldID:       fieldUUID,
		PaymentMethod: req.PaymentMethod,
		BookingDate:   bookingDate,
		StartTime:     startTime,
		EndTime:       endTime,
		Status:        "pending",
	}

	if err := bs.bookingRepo.CreateBooking(ctx, nil, booking); err != nil {
		return dto.BookingResponse{}, constants.ErrCreateBooking
	}

	res := dto.BookingResponse{
		BookingID:     booking.BookingID,
		UserID:        booking.UserID,
		FieldID:       booking.FieldID,
		PaymentMethod: booking.PaymentMethod,
		BookingDate:   booking.BookingDate,
		StartTime:     booking.StartTime,
		EndTime:       booking.EndTime,
		Status:        booking.Status,
	}

	return res, nil
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

	// Ambil data booking lama
	booking, _, err := bs.bookingRepo.GetBookingByID(ctx, nil, req.BookingID)
	if err != nil {
		return dto.BookingResponse{}, constants.ErrGetBookingByID
	}

	// Update field jika ada perubahan
	if req.UserID != nil {
		booking.UserID = *req.UserID
	}
	if req.FieldID != nil {
		booking.FieldID = *req.FieldID
	}
	if req.PaymentMethod != nil {
		booking.PaymentMethod = *req.PaymentMethod
	}
	if req.BookingDate != nil {
		booking.BookingDate = *req.BookingDate
	}
	if req.StartTime != nil {
		booking.StartTime = *req.StartTime
	}
	if req.EndTime != nil {
		booking.EndTime = *req.EndTime
	}
	if req.Status != nil {
		booking.Status = *req.Status
	}

	// Simpan perubahan
	if err := bs.bookingRepo.UpdateBooking(ctx, nil, booking); err != nil {
		return dto.BookingResponse{}, constants.ErrUpdateBooking
	}

	// Kembalikan response
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
