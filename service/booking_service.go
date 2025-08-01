package service

import (
	"context"
	"fieldreserve/constants"
	"fieldreserve/dto"
	"fieldreserve/helpers"
	"fieldreserve/model"
	"fieldreserve/repository"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
)

type (
	IBookingService interface {
		CreateBooking(ctx context.Context, req dto.CreateBookingRequest) (dto.BookingResponse, error)
		GetAllBooking(ctx context.Context, req dto.BookingPaginationRequest) (dto.BookingPaginationResponse, error)
		GetUserBookingHistory(ctx context.Context, req dto.BookingPaginationRequest) (dto.BookingPaginationResponse, error)
		GetBookingByID(ctx context.Context, bookingID string) (dto.BookingFullResponse, error)
		UpdateBookingStatus(ctx context.Context, req dto.UpdateBookingStatusRequest) (dto.BookingResponse, error)
		DeleteBooking(ctx context.Context, req dto.DeleteBookingRequest) (dto.BookingResponse, error)
	}

	BookingService struct {
		bookingRepo  repository.IBookingRepository
		jwtService   InterfaceJWTService
		scheduleRepo repository.IScheduleRepository
		fieldRepo    repository.IFieldRepository
	}
)

func NewBookingService(
	bookingRepo repository.IBookingRepository,
	jwtService InterfaceJWTService,
	scheduleRepo repository.IScheduleRepository,
	fieldRepo repository.IFieldRepository,
) *BookingService {
	return &BookingService{
		bookingRepo:  bookingRepo,
		jwtService:   jwtService,
		scheduleRepo: scheduleRepo,
		fieldRepo:    fieldRepo,
	}
}

func (bs *BookingService) CreateBooking(ctx context.Context, req dto.CreateBookingRequest) (dto.BookingResponse, error) {
	loc := helpers.GetAppLocation()

	// === [1] Extract Token & User ID ===
	tokenStr, ok := ctx.Value("token").(string)
	if !ok || tokenStr == "" {
		return dto.BookingResponse{}, constants.ErrUnauthorized
	}

	userIDStr, err := bs.jwtService.GetUserIDByToken(tokenStr)
	if err != nil {
		return dto.BookingResponse{}, constants.ErrUnauthorized
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return dto.BookingResponse{}, constants.ErrInvalidUUID
	}

	// === [2] Parse Field ID ===
	fieldID, err := uuid.Parse(req.FieldID)
	if err != nil {
		return dto.BookingResponse{}, constants.ErrInvalidUUID
	}

	// === [3] Parse Booking Date & Time ===
	bookingDate, err := time.ParseInLocation("2006-01-02", req.BookingDate, loc)
	if err != nil {
		return dto.BookingResponse{}, constants.ErrInvalidBookingDate
	}

	startTimeParsed, err := time.ParseInLocation("15:04", req.StartTime, loc)
	if err != nil {
		return dto.BookingResponse{}, constants.ErrInvalidTimeFormat
	}
	endTimeParsed, err := time.ParseInLocation("15:04", req.EndTime, loc)
	if err != nil {
		return dto.BookingResponse{}, constants.ErrInvalidTimeFormat
	}

	startTime := time.Date(bookingDate.Year(), bookingDate.Month(), bookingDate.Day(), startTimeParsed.Hour(), startTimeParsed.Minute(), 0, 0, loc)
	endTime := time.Date(bookingDate.Year(), bookingDate.Month(), bookingDate.Day(), endTimeParsed.Hour(), endTimeParsed.Minute(), 0, 0, loc)

	// === [4] Validasi Waktu ===
	now := time.Now().In(loc)
	if startTime.Before(now.Add(2 * time.Hour)) {
		return dto.BookingResponse{}, constants.ErrBookingTooSoon
	}
	if !endTime.After(startTime) {
		return dto.BookingResponse{}, constants.ErrInvalidTimeRange
	}

	// === [5] Validasi Field ===
	field, _, err := bs.fieldRepo.GetFieldByID(ctx, nil, req.FieldID)
	if err != nil {
		return dto.BookingResponse{}, constants.ErrFieldNotFound
	}

	durationHours := endTime.Sub(startTime).Hours()
	expectedTotal := float64(field.FieldPrice) * durationHours
	if math.Abs(req.TotalPayment-expectedTotal) > 1 {
		return dto.BookingResponse{}, constants.ErrInvalidTotalPayment
	}

	// === [6] Validasi Jadwal Field ===
	dayOfWeek := int(bookingDate.Weekday())
	schedule, err := bs.scheduleRepo.GetScheduleByFieldIDAndDay(ctx, nil, req.FieldID, dayOfWeek)
	if err != nil {
		return dto.BookingResponse{}, constants.ErrScheduleNotFound
	}

	openTime := time.Date(bookingDate.Year(), bookingDate.Month(), bookingDate.Day(), schedule.OpenTime.Hour(), schedule.OpenTime.Minute(), 0, 0, loc)
	closeTime := time.Date(bookingDate.Year(), bookingDate.Month(), bookingDate.Day(), schedule.CloseTime.Hour(), schedule.CloseTime.Minute(), 0, 0, loc)
	if startTime.Before(openTime) || endTime.After(closeTime) {
		return dto.BookingResponse{}, constants.ErrOutsideOperatingHours
	}

	// === [7] Validasi Overlap Booking ===
	overlap, err := bs.bookingRepo.CheckBookingOverlap(ctx, nil, fieldID, bookingDate, startTime, endTime)
	if err != nil {
		return dto.BookingResponse{}, constants.ErrCheckOverlap
	}
	if overlap {
		return dto.BookingResponse{}, constants.ErrBookingOverlap
	}

	// === [8] Handle Bukti Pembayaran (Opsional) ===
	var proofPath string
	var paymentUploadedAt *time.Time
	status := constants.ENUM_STATUS_BOOKING_PENDING

	if req.ProofPayment != nil {
		imageName, err := helpers.SaveImage(req.ProofPayment, "./assets/proof", "proof")
		if err != nil {
			return dto.BookingResponse{}, constants.ErrSaveImages
		}
		proofPath = imageName
		status = constants.ENUM_STATUS_BOOKING_WAITING
		now := time.Now().In(loc)
		paymentUploadedAt = &now
	}

	// === [9] Simpan Booking ===
	booking := model.Booking{
		BookingID:         uuid.New(),
		UserID:            userID,
		FieldID:           fieldID,
		PaymentMethod:     req.PaymentMethod,
		BookingDate:       bookingDate,
		StartTime:         startTime,
		EndTime:           endTime,
		TotalPayment:      req.TotalPayment,
		ProofPayment:      proofPath,
		Status:            status,
		PaymentUploadedAt: paymentUploadedAt,
	}

	if err := bs.bookingRepo.CreateBooking(ctx, nil, booking); err != nil {
		return dto.BookingResponse{}, constants.ErrCreateBooking
	}

	// === [10] Return DTO Response ===
	return dto.BookingResponse{
		BookingID:     booking.BookingID,
		UserID:        booking.UserID,
		FieldID:       booking.FieldID,
		PaymentMethod: booking.PaymentMethod,
		BookingDate:   booking.BookingDate,
		StartTime:     booking.StartTime,
		EndTime:       booking.EndTime,
		TotalPayment:  booking.TotalPayment,
		Status:        booking.Status,
		ProofPayment:  booking.ProofPayment,
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
			BookingID:         booking.BookingID,
			UserID:            booking.UserID,
			FieldID:           booking.FieldID,
			PaymentMethod:     booking.PaymentMethod,
			BookingDate:       booking.BookingDate,
			StartTime:         booking.StartTime,
			EndTime:           booking.EndTime,
			Status:            booking.Status,
			TotalPayment:      booking.TotalPayment,
			ProofPayment:      booking.ProofPayment,
			PaymentUploadedAt: booking.PaymentUploadedAt,
			PaymentVerifiedAt: booking.PaymentVerifiedAt,
			CancelledAt:       booking.CancelledAt,
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

func (bs *BookingService) GetUserBookingHistory(ctx context.Context, req dto.BookingPaginationRequest) (dto.BookingPaginationResponse, error) {
	dataWithPaginate, err := bs.bookingRepo.GetAllBooking(ctx, nil, req)
	if err != nil {
		return dto.BookingPaginationResponse{}, constants.ErrGetAllField
	}

	var datas []dto.BookingResponse
	for _, booking := range dataWithPaginate.Bookings {
		data := dto.BookingResponse{
			BookingID:         booking.BookingID,
			UserID:            booking.UserID,
			FieldID:           booking.FieldID,
			PaymentMethod:     booking.PaymentMethod,
			BookingDate:       booking.BookingDate,
			StartTime:         booking.StartTime,
			EndTime:           booking.EndTime,
			Status:            booking.Status,
			TotalPayment:      booking.TotalPayment,
			ProofPayment:      booking.ProofPayment,
			PaymentUploadedAt: booking.PaymentUploadedAt,
			PaymentVerifiedAt: booking.PaymentVerifiedAt,
			CancelledAt:       booking.CancelledAt,
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
	user := booking.User
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

	userDTO := dto.UserCompactResponse{
		ID:      user.UserID,
		Name:    user.Name,
		Email:   user.Email,
		NoTelp:  user.NoTelp,
		Address: user.Address,
	}

	res := dto.BookingFullResponse{
		BookingID:         booking.BookingID,
		PaymentMethod:     booking.PaymentMethod,
		BookingDate:       booking.BookingDate,
		StartTime:         booking.StartTime,
		EndTime:           booking.EndTime,
		Status:            booking.Status,
		TotalPayment:      booking.TotalPayment,
		ProofPayment:      booking.ProofPayment,
		PaymentUploadedAt: booking.PaymentUploadedAt,
		PaymentVerifiedAt: booking.PaymentVerifiedAt,
		CancelledAt:       booking.CancelledAt,
		User:              userDTO,
		Field:             fieldDTO,
	}

	return res, nil
}

func (bs *BookingService) UpdateBookingStatus(ctx context.Context, req dto.UpdateBookingStatusRequest) (dto.BookingResponse, error) {
	loc := helpers.GetAppLocation()

	// // ====== 1. Validasi UUID Booking ID ======
	if _, err := uuid.Parse(req.BookingID); err != nil {
		return dto.BookingResponse{}, constants.ErrInvalidUUID
	}

	// ====== 2. Ambil Booking dari Database ======
	booking, _, err := bs.bookingRepo.GetBookingByID(ctx, nil, req.BookingID)
	if err != nil {
		return dto.BookingResponse{}, constants.ErrBookingNotFound
	}

	// ====== 3. Cek Status Saat Ini ======
	if booking.Status == constants.ENUM_STATUS_BOOKING_CALCEL || booking.Status == constants.ENUM_STATUS_BOOKING_BOOKED {
		return dto.BookingResponse{}, constants.ErrBookingAlreadyFinal
	}

	// ====== 4. Validasi Status Baru ======
	newStatus := strings.ToLower(*req.Status)
	if newStatus != constants.ENUM_STATUS_BOOKING_CALCEL && newStatus != constants.ENUM_STATUS_BOOKING_BOOKED {
		return dto.BookingResponse{}, constants.ErrInvalidStatusUpdate
	}

	// ====== 5. Update Status & Timestamp ======
	now := time.Now().In(loc)
	booking.Status = newStatus

	if newStatus == constants.ENUM_STATUS_BOOKING_BOOKED {
		booking.PaymentVerifiedAt = &now
	} else if newStatus == constants.ENUM_STATUS_BOOKING_CALCEL {
		booking.CancelledAt = &now
	}

	// ====== 6. Simpan ke Database ======
	if err := bs.bookingRepo.UpdateBooking(ctx, nil, booking); err != nil {
		return dto.BookingResponse{}, constants.ErrUpdateBooking
	}

	// ====== 7. Response DTO ======
	return dto.BookingResponse{
		BookingID:         booking.BookingID,
		UserID:            booking.UserID,
		FieldID:           booking.FieldID,
		BookingDate:       booking.BookingDate,
		StartTime:         booking.StartTime,
		EndTime:           booking.EndTime,
		PaymentMethod:     booking.PaymentMethod,
		TotalPayment:      booking.TotalPayment,
		ProofPayment:      booking.ProofPayment,
		Status:            booking.Status,
		PaymentUploadedAt: booking.PaymentUploadedAt,
		PaymentVerifiedAt: booking.PaymentVerifiedAt,
		CancelledAt:       booking.CancelledAt,
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
		BookingID:         booking.BookingID,
		UserID:            booking.UserID,
		FieldID:           booking.FieldID,
		BookingDate:       booking.BookingDate,
		StartTime:         booking.StartTime,
		EndTime:           booking.EndTime,
		PaymentMethod:     booking.PaymentMethod,
		TotalPayment:      booking.TotalPayment,
		ProofPayment:      booking.ProofPayment,
		Status:            booking.Status,
		PaymentUploadedAt: booking.PaymentUploadedAt,
		PaymentVerifiedAt: booking.PaymentVerifiedAt,
		CancelledAt:       booking.CancelledAt,
	}, nil
}
