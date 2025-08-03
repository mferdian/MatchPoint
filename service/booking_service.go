package service

import (
	"context"
	"fieldreserve/constants"
	"fieldreserve/dto"
	"fieldreserve/helpers"
	"fieldreserve/model"
	"fieldreserve/repository"
	"fieldreserve/utils"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
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
	utils.Log.Info("Initializing new BookingService")
	return &BookingService{
		bookingRepo:  bookingRepo,
		jwtService:   jwtService,
		scheduleRepo: scheduleRepo,
		fieldRepo:    fieldRepo,
	}
}

func (bs *BookingService) CreateBooking(ctx context.Context, req dto.CreateBookingRequest) (dto.BookingResponse, error) {
	utils.Log.WithFields(logrus.Fields{
		"fieldID":       req.FieldID,
		"bookingDate":   req.BookingDate,
		"startTime":     req.StartTime,
		"endTime":       req.EndTime,
		"paymentMethod": req.PaymentMethod,
		"totalPayment":  req.TotalPayment,
	}).Info("Starting booking creation process")

	loc := helpers.GetAppLocation()

	// === [1] Extract Token & User ID ===
	utils.Log.Debug("Extracting token and user ID from context")
	tokenStr, ok := ctx.Value("token").(string)
	if !ok || tokenStr == "" {
		utils.Log.Error("Token not found or empty in context")
		return dto.BookingResponse{}, constants.ErrUnauthorized
	}

	userIDStr, err := bs.jwtService.GetUserIDByToken(tokenStr)
	if err != nil {
		utils.Log.WithError(err).Error("Failed to get user ID from token")
		return dto.BookingResponse{}, constants.ErrUnauthorized
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		utils.Log.WithError(err).WithField("userIDStr", userIDStr).Error("Failed to parse user ID")
		return dto.BookingResponse{}, constants.ErrInvalidUUID
	}

	utils.Log.WithField("userID", userID).Info("Successfully extracted user ID from token")

	// === [2] Parse Field ID ===
	utils.Log.WithField("fieldID", req.FieldID).Debug("Parsing field ID")
	fieldID, err := uuid.Parse(req.FieldID)
	if err != nil {
		utils.Log.WithError(err).WithField("fieldID", req.FieldID).Error("Failed to parse field ID")
		return dto.BookingResponse{}, constants.ErrInvalidUUID
	}

	// === [3] Parse Booking Date & Time ===
	utils.Log.WithFields(logrus.Fields{
		"bookingDate": req.BookingDate,
		"startTime":   req.StartTime,
		"endTime":     req.EndTime,
	}).Debug("Parsing booking date and time")

	bookingDate, err := time.ParseInLocation("2006-01-02", req.BookingDate, loc)
	if err != nil {
		utils.Log.WithError(err).WithField("bookingDate", req.BookingDate).Error("Failed to parse booking date")
		return dto.BookingResponse{}, constants.ErrInvalidBookingDate
	}

	startTimeParsed, err := time.ParseInLocation("15:04", req.StartTime, loc)
	if err != nil {
		utils.Log.WithError(err).WithField("startTime", req.StartTime).Error("Failed to parse start time")
		return dto.BookingResponse{}, constants.ErrInvalidTimeFormat
	}
	endTimeParsed, err := time.ParseInLocation("15:04", req.EndTime, loc)
	if err != nil {
		utils.Log.WithError(err).WithField("endTime", req.EndTime).Error("Failed to parse end time")
		return dto.BookingResponse{}, constants.ErrInvalidTimeFormat
	}

	startTime := time.Date(bookingDate.Year(), bookingDate.Month(), bookingDate.Day(), startTimeParsed.Hour(), startTimeParsed.Minute(), 0, 0, loc)
	endTime := time.Date(bookingDate.Year(), bookingDate.Month(), bookingDate.Day(), endTimeParsed.Hour(), endTimeParsed.Minute(), 0, 0, loc)

	utils.Log.WithFields(logrus.Fields{
		"parsedStartTime": startTime,
		"parsedEndTime":   endTime,
	}).Debug("Successfully parsed booking times")

	// === [4] Validasi Waktu ===
	utils.Log.Debug("Validating booking time constraints")
	now := time.Now().In(loc)
	if startTime.Before(now.Add(2 * time.Hour)) {
		utils.Log.WithFields(logrus.Fields{
			"startTime":   startTime,
			"currentTime": now,
			"minTime":     now.Add(2 * time.Hour),
		}).Warn("Booking attempt too soon - less than 2 hours notice")
		return dto.BookingResponse{}, constants.ErrBookingTooSoon
	}
	if !endTime.After(startTime) {
		utils.Log.WithFields(logrus.Fields{
			"startTime": startTime,
			"endTime":   endTime,
		}).Warn("Invalid time range - end time not after start time")
		return dto.BookingResponse{}, constants.ErrInvalidTimeRange
	}

	// === [5] Validasi Field ===
	utils.Log.WithField("fieldID", req.FieldID).Debug("Validating field existence and calculating payment")
	field, _, err := bs.fieldRepo.GetFieldByID(ctx, nil, req.FieldID)
	if err != nil {
		utils.Log.WithError(err).WithField("fieldID", req.FieldID).Error("Field not found")
		return dto.BookingResponse{}, constants.ErrFieldNotFound
	}

	durationHours := endTime.Sub(startTime).Hours()
	expectedTotal := float64(field.FieldPrice) * durationHours
	if math.Abs(req.TotalPayment-expectedTotal) > 1 {
		utils.Log.WithFields(logrus.Fields{
			"expectedTotal":   expectedTotal,
			"providedTotal":   req.TotalPayment,
			"fieldPrice":      field.FieldPrice,
			"durationHours":   durationHours,
		}).Warn("Invalid total payment amount")
		return dto.BookingResponse{}, constants.ErrInvalidTotalPayment
	}

	utils.Log.WithFields(logrus.Fields{
		"fieldName":     field.FieldName,
		"fieldPrice":    field.FieldPrice,
		"duration":      durationHours,
		"totalPayment":  expectedTotal,
	}).Info("Field validation successful")

	// === [6] Validasi Jadwal Field ===
	dayOfWeek := int(bookingDate.Weekday())
	utils.Log.WithFields(logrus.Fields{
		"fieldID":   req.FieldID,
		"dayOfWeek": dayOfWeek,
	}).Debug("Checking field schedule")

	schedule, err := bs.scheduleRepo.GetScheduleByFieldIDAndDay(ctx, nil, req.FieldID, dayOfWeek)
	if err != nil {
		utils.Log.WithError(err).WithFields(logrus.Fields{
			"fieldID":   req.FieldID,
			"dayOfWeek": dayOfWeek,
		}).Error("Schedule not found for field and day")
		return dto.BookingResponse{}, constants.ErrScheduleNotFound
	}

	openTime := time.Date(bookingDate.Year(), bookingDate.Month(), bookingDate.Day(), schedule.OpenTime.Hour(), schedule.OpenTime.Minute(), 0, 0, loc)
	closeTime := time.Date(bookingDate.Year(), bookingDate.Month(), bookingDate.Day(), schedule.CloseTime.Hour(), schedule.CloseTime.Minute(), 0, 0, loc)
	if startTime.Before(openTime) || endTime.After(closeTime) {
		utils.Log.WithFields(logrus.Fields{
			"requestedStart": startTime,
			"requestedEnd":   endTime,
			"openTime":       openTime,
			"closeTime":      closeTime,
		}).Warn("Booking time outside operating hours")
		return dto.BookingResponse{}, constants.ErrOutsideOperatingHours
	}

	utils.Log.WithFields(logrus.Fields{
		"openTime":  openTime,
		"closeTime": closeTime,
	}).Debug("Schedule validation successful")

	// === [7] Validasi Overlap Booking ===
	utils.Log.Debug("Checking for booking overlaps")
	overlap, err := bs.bookingRepo.CheckBookingOverlap(ctx, nil, fieldID, bookingDate, startTime, endTime)
	if err != nil {
		utils.Log.WithError(err).WithFields(logrus.Fields{
			"fieldID":     fieldID,
			"bookingDate": bookingDate,
			"startTime":   startTime,
			"endTime":     endTime,
		}).Error("Failed to check booking overlap")
		return dto.BookingResponse{}, constants.ErrCheckOverlap
	}
	if overlap {
		utils.Log.WithFields(logrus.Fields{
			"fieldID":     fieldID,
			"bookingDate": bookingDate,
			"startTime":   startTime,
			"endTime":     endTime,
		}).Warn("Booking time slot already occupied")
		return dto.BookingResponse{}, constants.ErrBookingOverlap
	}

	utils.Log.Debug("No booking overlap found")

	// === [8] Handle Bukti Pembayaran (Opsional) ===
	var proofPath string
	var paymentUploadedAt *time.Time
	status := constants.ENUM_STATUS_BOOKING_PENDING

	if req.ProofPayment != nil {
		utils.Log.Debug("Processing payment proof upload")
		imageName, err := helpers.SaveImage(req.ProofPayment, "./assets/proof", "proof")
		if err != nil {
			utils.Log.WithError(err).Error("Failed to save payment proof image")
			return dto.BookingResponse{}, constants.ErrSaveImages
		}
		proofPath = imageName
		status = constants.ENUM_STATUS_BOOKING_WAITING
		now := time.Now().In(loc)
		paymentUploadedAt = &now
		
		utils.Log.WithFields(logrus.Fields{
			"proofPath":         proofPath,
			"paymentUploadedAt": paymentUploadedAt,
			"newStatus":         status,
		}).Info("Payment proof uploaded successfully")
	} else {
		utils.Log.Debug("No payment proof provided, booking set to pending status")
	}

	// === [9] Simpan Booking ===
	bookingID := uuid.New()
	booking := model.Booking{
		BookingID:         bookingID,
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

	utils.Log.WithFields(logrus.Fields{
		"bookingID":     bookingID,
		"userID":        userID,
		"fieldID":       fieldID,
		"paymentMethod": req.PaymentMethod,
		"status":        status,
	}).Info("Attempting to save booking to database")

	if err := bs.bookingRepo.CreateBooking(ctx, nil, booking); err != nil {
		utils.Log.WithError(err).WithField("bookingID", bookingID).Error("Failed to create booking in database")
		return dto.BookingResponse{}, constants.ErrCreateBooking
	}

	utils.Log.WithField("bookingID", bookingID).Info("Booking created successfully")

	// === [10] Return DTO Response ===
	response := dto.BookingResponse{
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
	}

	utils.Log.WithField("bookingID", bookingID).Info("Booking creation process completed successfully")
	return response, nil
}

func (bs *BookingService) GetAllBooking(ctx context.Context, req dto.BookingPaginationRequest) (dto.BookingPaginationResponse, error) {
	utils.Log.WithFields(logrus.Fields{
		"page":    req.Page,
		"perPage": req.PerPage,
	}).Info("Fetching all bookings with pagination")

	dataWithPaginate, err := bs.bookingRepo.GetAllBooking(ctx, nil, req)
	if err != nil {
		utils.Log.WithError(err).WithFields(logrus.Fields{
			"page":    req.Page,
			"perPage": req.PerPage,
		}).Error("Failed to fetch all bookings")
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

	utils.Log.WithFields(logrus.Fields{
		"totalBookings": len(datas),
		"page":          dataWithPaginate.Page,
		"maxPage":       dataWithPaginate.MaxPage,
		"count":         dataWithPaginate.Count,
	}).Info("Successfully fetched all bookings")

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
	utils.Log.WithFields(logrus.Fields{
		"page":    req.Page,
		"perPage": req.PerPage,
	}).Info("Fetching user booking history")

	dataWithPaginate, err := bs.bookingRepo.GetAllBooking(ctx, nil, req)
	if err != nil {
		utils.Log.WithError(err).WithFields(logrus.Fields{
			"page":    req.Page,
			"perPage": req.PerPage,
		}).Error("Failed to fetch user booking history")
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

	utils.Log.WithFields(logrus.Fields{
		"totalBookings": len(datas),
		"page":          dataWithPaginate.Page,
		"maxPage":       dataWithPaginate.MaxPage,
		"count":         dataWithPaginate.Count,
	}).Info("Successfully fetched user booking history")

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
	utils.Log.WithField("bookingID", bookingID).Info("Fetching booking by ID")

	if _, err := uuid.Parse(bookingID); err != nil {
		utils.Log.WithError(err).WithField("bookingID", bookingID).Error("Invalid booking ID format")
		return dto.BookingFullResponse{}, constants.ErrInvalidUUID
	}

	booking, _, err := bs.bookingRepo.GetBookingByID(ctx, nil, bookingID)
	if err != nil {
		utils.Log.WithError(err).WithField("bookingID", bookingID).Error("Failed to fetch booking by ID")
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

	utils.Log.WithFields(logrus.Fields{
		"bookingID":   bookingID,
		"fieldName":   field.FieldName,
		"userName":    user.Name,
		"status":      booking.Status,
	}).Info("Successfully fetched booking details")

	return res, nil
}

func (bs *BookingService) UpdateBookingStatus(ctx context.Context, req dto.UpdateBookingStatusRequest) (dto.BookingResponse, error) {
	utils.Log.WithFields(logrus.Fields{
		"bookingID": req.BookingID,
		"newStatus": *req.Status,
	}).Info("Starting booking status update")

	loc := helpers.GetAppLocation()

	// ====== 1. Validasi UUID Booking ID ======
	if _, err := uuid.Parse(req.BookingID); err != nil {
		utils.Log.WithError(err).WithField("bookingID", req.BookingID).Error("Invalid booking ID format")
		return dto.BookingResponse{}, constants.ErrInvalidUUID
	}

	// ====== 2. Ambil Booking dari Database ======
	utils.Log.WithField("bookingID", req.BookingID).Debug("Fetching booking from database")
	booking, _, err := bs.bookingRepo.GetBookingByID(ctx, nil, req.BookingID)
	if err != nil {
		utils.Log.WithError(err).WithField("bookingID", req.BookingID).Error("Booking not found")
		return dto.BookingResponse{}, constants.ErrBookingNotFound
	}

	utils.Log.WithFields(logrus.Fields{
		"bookingID":    req.BookingID,
		"currentStatus": booking.Status,
	}).Debug("Current booking status retrieved")

	// ====== 3. Cek Status Saat Ini ======
	if booking.Status == constants.ENUM_STATUS_BOOKING_CALCEL || booking.Status == constants.ENUM_STATUS_BOOKING_BOOKED {
		utils.Log.WithFields(logrus.Fields{
			"bookingID":     req.BookingID,
			"currentStatus": booking.Status,
			"requestedStatus": *req.Status,
		}).Warn("Cannot update status - booking already in final state")
		return dto.BookingResponse{}, constants.ErrBookingAlreadyFinal
	}

	// ====== 4. Validasi Status Baru ======
	newStatus := strings.ToLower(*req.Status)
	if newStatus != constants.ENUM_STATUS_BOOKING_CALCEL && newStatus != constants.ENUM_STATUS_BOOKING_BOOKED {
		utils.Log.WithFields(logrus.Fields{
			"bookingID":       req.BookingID,
			"requestedStatus": newStatus,
		}).Error("Invalid status update requested")
		return dto.BookingResponse{}, constants.ErrInvalidStatusUpdate
	}

	// ====== 5. Update Status & Timestamp ======
	now := time.Now().In(loc)
	oldStatus := booking.Status
	booking.Status = newStatus

	if newStatus == constants.ENUM_STATUS_BOOKING_BOOKED {
		booking.PaymentVerifiedAt = &now
		utils.Log.WithFields(logrus.Fields{
			"bookingID":          req.BookingID,
			"paymentVerifiedAt": now,
		}).Info("Booking approved and payment verified")
	} else if newStatus == constants.ENUM_STATUS_BOOKING_CALCEL {
		booking.CancelledAt = &now
		utils.Log.WithFields(logrus.Fields{
			"bookingID":    req.BookingID,
			"cancelledAt": now,
		}).Info("Booking cancelled")
	}

	// ====== 6. Simpan ke Database ======
	utils.Log.WithField("bookingID", req.BookingID).Debug("Saving updated booking to database")
	if err := bs.bookingRepo.UpdateBooking(ctx, nil, booking); err != nil {
		utils.Log.WithError(err).WithField("bookingID", req.BookingID).Error("Failed to update booking in database")
		return dto.BookingResponse{}, constants.ErrUpdateBooking
	}

	utils.Log.WithFields(logrus.Fields{
		"bookingID": req.BookingID,
		"oldStatus": oldStatus,
		"newStatus": newStatus,
	}).Info("Booking status updated successfully")

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
	utils.Log.WithField("bookingID", req.BookingID).Info("Starting booking deletion process")

	if _, err := uuid.Parse(req.BookingID); err != nil {
		utils.Log.WithError(err).WithField("bookingID", req.BookingID).Error("Invalid booking ID format")
		return dto.BookingResponse{}, constants.ErrInvalidUUID
	}

	// Ambil booking terlebih dahulu
	utils.Log.WithField("bookingID", req.BookingID).Debug("Fetching booking for deletion")
	booking, _, err := bs.bookingRepo.GetBookingByID(ctx, nil, req.BookingID)
	if err != nil {
		utils.Log.WithError(err).WithField("bookingID", req.BookingID).Error("Booking not found for deletion")
		return dto.BookingResponse{}, constants.ErrGetBookingByID
	}

	// Validasi: hanya bisa dibatalkan maksimal 3 jam sebelum startTime
	timeUntilStart := time.Until(booking.StartTime)
	if timeUntilStart < 3*time.Hour {
		utils.Log.WithFields(logrus.Fields{
			"bookingID":      req.BookingID,
			"startTime":      booking.StartTime,
			"timeUntilStart": timeUntilStart,
		}).Warn("Cannot cancel booking - less than 3 hours before start time")
		return dto.BookingResponse{}, constants.ErrCannotCancelLate
	}

	utils.Log.WithFields(logrus.Fields{
		"bookingID":      req.BookingID,
		"startTime":      booking.StartTime,
		"timeUntilStart": timeUntilStart,
	}).Debug("Cancellation time validation passed")

	// Lakukan soft delete atau ubah status
	utils.Log.WithField("bookingID", req.BookingID).Debug("Performing booking deletion")
	if err := bs.bookingRepo.DeleteBooking(ctx, nil, req.BookingID); err != nil {
		utils.Log.WithError(err).WithField("bookingID", req.BookingID).Error("Failed to delete booking")
		return dto.BookingResponse{}, constants.ErrDeleteBooking
	}

	utils.Log.WithField("bookingID", req.BookingID).Info("Booking deleted successfully")

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