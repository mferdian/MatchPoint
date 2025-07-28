package dto

import (
	"fieldreserve/model"
	"time"

	"github.com/google/uuid"
)

type (
	CreateBookingRequest struct {
		FieldID       string `json:"field_id" binding:"required"`
		BookingDate   string `json:"booking_date" binding:"required"`
		StartTime     string `json:"start_time" binding:"required"`
		EndTime       string `json:"end_time" binding:"required"`
		PaymentMethod string `json:"payment_method" binding:"required"`
	}

	BookingResponse struct {
		BookingID     uuid.UUID `json:"booking_id"`
		UserID        uuid.UUID `json:"user_id"`
		FieldID       uuid.UUID `json:"field_id"`
		PaymentMethod string    `json:"payment_method"`
		BookingDate   time.Time `json:"booking_date"`
		StartTime     time.Time `json:"start_time"`
		EndTime       time.Time `json:"end_time"`
		Status        string    `json:"status"`
	}

	BookingFullResponse struct {
		BookingID     uuid.UUID            `json:"booking_id"`
		UserID        uuid.UUID            `json:"user_id"`
		PaymentMethod string               `json:"payment_method"`
		BookingDate   time.Time            `json:"booking_date"`
		StartTime     time.Time            `json:"start_time"`
		EndTime       time.Time            `json:"end_time"`
		Status        string               `json:"status"`
		Field         FieldCompactResponse `json:"field"`
	}

	UpdateBookingRequest struct {
		BookingID     string     `json:"-"`
		UserID        *uuid.UUID `json:"user_id,omitempty"`
		FieldID       *uuid.UUID `json:"field_id,omitempty"`
		PaymentMethod *string    `json:"payment_method,omitempty"`
		BookingDate   *time.Time `json:"booking_date,omitempty"`
		StartTime     *time.Time `json:"start_time,omitempty"`
		EndTime       *time.Time `json:"end_time,omitempty"`
		Status        *string    `json:"status,omitempty"`
	}

	DeleteBookingRequest struct {
		BookingID string `json:"-"`
	}

	// Pagination
	BookingPaginationRequest struct {
		PaginationRequest
		BookingID string `form:"booking_id"`
		FieldID   string `form:"field_id"`
		UserID    string `form:"user_id"`
		Status    string `form:"status"`
	}

	BookingPaginationResponse struct {
		PaginationResponse
		Data []BookingResponse `json:"data"`
	}

	BookingPaginationRepositoryResponse struct {
		PaginationResponse
		Bookings []model.Booking
	}
)
