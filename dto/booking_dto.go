package dto

import (
	"fieldreserve/model"
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type (
	// Create Booking
	CreateBookingRequest struct {
		FieldID       string                `form:"field_id" binding:"required"`
		BookingDate   string                `form:"booking_date" binding:"required"`
		StartTime     string                `form:"start_time" binding:"required"`
		EndTime       string                `form:"end_time" binding:"required"`
		PaymentMethod string                `form:"payment_method" binding:"required"`
		ProofPayment  *multipart.FileHeader `form:"proof_payment" binding:"required"`
		TotalPayment  float64               `form:"total_payment" binding:"required"`
	}

	// Response Ringkas
	BookingResponse struct {
		BookingID         uuid.UUID  `json:"booking_id"`
		UserID            uuid.UUID  `json:"user_id"`
		FieldID           uuid.UUID  `json:"field_id"`
		PaymentMethod     string     `json:"payment_method"`
		BookingDate       time.Time  `json:"booking_date"`
		StartTime         time.Time  `json:"start_time"`
		EndTime           time.Time  `json:"end_time"`
		Status            string     `json:"status"`
		TotalPayment      float64    `json:"total_payment"`
		ProofPayment      string     `json:"proof_payment"`
		PaymentUploadedAt *time.Time `json:"payment_uploaded_at,omitempty"`
		PaymentVerifiedAt *time.Time `json:"payment_verified_at,omitempty"`
		CancelledAt       *time.Time `json:"cancelled_at,omitempty"`
	}

	// Response Lengkap (untuk history / admin)
	BookingFullResponse struct {
		BookingID         uuid.UUID            `json:"booking_id"`
		UserID            uuid.UUID            `json:"user_id"`
		PaymentMethod     string               `json:"payment_method"`
		BookingDate       time.Time            `json:"booking_date"`
		StartTime         time.Time            `json:"start_time"`
		EndTime           time.Time            `json:"end_time"`
		TotalPayment      float64              `json:"total_payment"`
		ProofPayment      string               `json:"proof_payment"`
		Status            string               `json:"status"`
		Field             FieldCompactResponse `json:"field"`
		PaymentVerifiedAt *time.Time           `json:"payment_verified_at,omitempty"`
		CancelledAt       *time.Time           `json:"cancelled_at,omitempty"`
		PaymentUploadedAt *time.Time           `json:"payment_uploaded_at,omitempty"`
		VerifiedAt        *time.Time           `json:"verified_at,omitempty"`
	}

	// Update Booking (khusus admin, user tidak bisa update booking)
	UpdateBookingStatusRequest struct {
		BookingID string  `json:"-"`
		Status    *string `json:"status,omitempty"`
	}

	// Delete Booking
	DeleteBookingRequest struct {
		BookingID string `json:"-"`
	}

	// Pagination Request
	BookingPaginationRequest struct {
		PaginationRequest
		BookingID string `form:"booking_id"`
		FieldID   string `form:"field_id"`
		UserID    string `form:"user_id"`
		Status    string `form:"status"`
	}

	// Pagination Response
	BookingPaginationResponse struct {
		PaginationResponse
		Data []BookingResponse `json:"data"`
	}

	// Digunakan di Repository
	BookingPaginationRepositoryResponse struct {
		PaginationResponse
		Bookings []model.Booking
	}
)
