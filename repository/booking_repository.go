package repository

import (
	"context"
	"fieldreserve/dto"
	"fieldreserve/model"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	IBookingRepository interface {
		CreateBooking(ctx context.Context, tx *gorm.DB, booking model.Booking) error
		GetAllBooking(ctx context.Context, tx *gorm.DB, req dto.BookingPaginationRequest) (dto.BookingPaginationRepositoryResponse, error)
		GetBookingByID(ctx context.Context, tx *gorm.DB, bookingID string) (model.Booking, bool, error)
		UpdateBooking(ctx context.Context, tx *gorm.DB, booking model.Booking) error
		DeleteBooking(ctx context.Context, tx *gorm.DB, bookingID string) error
		CheckBookingOverlap(ctx context.Context, tx *gorm.DB, fieldID uuid.UUID, bookingDate time.Time, startTime, endTime time.Time) (bool, error)

	}

	BookingRepository struct {
		db *gorm.DB
	}
)

func NewBookingRepository(db *gorm.DB) *BookingRepository {
	return &BookingRepository{
		db: db,
	}
}

func (br *BookingRepository) CreateBooking(ctx context.Context, tx *gorm.DB, booking model.Booking) error {
	if tx == nil {
		tx = br.db
	}

	return tx.WithContext(ctx).Create(&booking).Error
}

func (br *BookingRepository) GetAllBooking(ctx context.Context, tx *gorm.DB, req dto.BookingPaginationRequest) (dto.BookingPaginationRepositoryResponse, error) {
	if tx == nil {
		tx = br.db
	}

	var bookings []model.Booking
	var count int64

	// Default pagination
	if req.PaginationRequest.PerPage == 0 {
		req.PaginationRequest.PerPage = 10
	}
	if req.PaginationRequest.Page == 0 {
		req.PaginationRequest.Page = 1
	}

	query := tx.WithContext(ctx).
		Model(&model.Booking{}).
		Joins("Field").
		Preload("Field")

	// Search logic (booking_date, status, field name, field address)
	if search := strings.TrimSpace(req.PaginationRequest.Search); search != "" {
		searchValue := "%" + strings.ToLower(search) + "%"
		query = query.Where(`
			CAST(bookings.booking_date AS TEXT) ILIKE ? OR
			LOWER(bookings.status) ILIKE ? OR
			LOWER(fields.field_name) ILIKE ? OR
			LOWER(fields.field_address) ILIKE ?`,
			searchValue, searchValue, searchValue, searchValue,
		)
	}

	// Filtering
	if req.BookingID != "" {
		query = query.Where("bookings.booking_id = ?", req.BookingID)
	}
	if req.FieldID != "" {
		query = query.Where("bookings.field_id = ?", req.FieldID)
	}
	if req.UserID != "" {
		query = query.Where("bookings.user_id = ?", req.UserID)
	}
	if req.Status != "" {
		query = query.Where("bookings.status = ?", req.Status)
	}

	// Count total
	if err := query.Count(&count).Error; err != nil {
		return dto.BookingPaginationRepositoryResponse{}, err
	}

	// Ambil data dengan pagination
	if err := query.
		Order("bookings.created_at DESC").
		Scopes(Paginate(req.PaginationRequest.Page, req.PaginationRequest.PerPage)).
		Find(&bookings).Error; err != nil {
		return dto.BookingPaginationRepositoryResponse{}, err
	}

	totalPage := int64(math.Ceil(float64(count) / float64(req.PaginationRequest.PerPage)))

	return dto.BookingPaginationRepositoryResponse{
		Bookings: bookings,
		PaginationResponse: dto.PaginationResponse{
			Page:    req.PaginationRequest.Page,
			PerPage: req.PaginationRequest.PerPage,
			MaxPage: totalPage,
			Count:   count,
		},
	}, nil
}

func (br *BookingRepository) GetBookingByID(ctx context.Context, tx *gorm.DB, bookingID string) (model.Booking, bool, error) {
	if tx == nil {
		tx = br.db
	}

	var booking model.Booking
	if err := tx.WithContext(ctx).
		Preload("Field").
		Preload("Field.Category").
		Where("booking_id = ?", bookingID).
		Take(&booking).Error; err != nil {
		return model.Booking{}, false, err
	}

	return booking, true, nil
}

func (br *BookingRepository) UpdateBooking(ctx context.Context, tx *gorm.DB, booking model.Booking) error {
	if tx == nil {
		tx = br.db
	}

	return tx.WithContext(ctx).Where("booking_id = ?", booking.BookingID).Updates(&booking).Error
}

func (br *BookingRepository) DeleteBooking(ctx context.Context, tx *gorm.DB, bookingID string) error {
	if tx == nil {
		tx = br.db
	}

	return tx.WithContext(ctx).Where("booking_id = ?", bookingID).Delete(&model.Booking{}).Error
}

func (br *BookingRepository) CheckBookingOverlap(ctx context.Context, tx *gorm.DB, fieldID uuid.UUID, bookingDate time.Time, startTime, endTime time.Time) (bool, error) {
    if tx == nil {
        tx = br.db
    }

    var count int64
    err := tx.WithContext(ctx).
        Model(&model.Booking{}).
        Where("field_id = ? AND booking_date = ? AND status != ?", fieldID, bookingDate, "cancelled").
        Where("? < end_time AND ? > start_time", startTime, endTime).
        Count(&count).Error

    if err != nil {
        return false, err
    }

    return count > 0, nil
}

