package dto

import (
	"fieldreserve/model"

	"github.com/google/uuid"
)

type (
	UserResponse struct {
		ID      uuid.UUID `json:"user_id"`
		Name    string    `json:"user_name"`
		Email   string    `json:"user_email"`
		Address string    `json:"address"`
		NoTelp  string    `json:"no_telp"`
	}

	CreateUserRequest struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	LoginUserRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	LoginResponse struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	UpdateUserRequest struct {
		ID       string `json:"-"`
		Name     string `json:"name,omitempty"`
		Email    string `json:"email,omitempty"`
		Address  string `json:"address,omitempty"`
		NoTelp   string `json:"no_telp,omitempty"`
		Password string `json:"password,omitempty"`
	}

	DeleteUserRequest struct {
		UserID string `json:"-"`
	}

	UserPaginationRequest struct {
		PaginationRequest
		UserID string `form:"user_id"`
	}

	UserPaginationResponse struct {
		PaginationResponse
		Data []UserResponse `json:"data"`
	}

	UserPaginationRepositoryResponse struct {
		PaginationResponse
		Users []model.User
	}
)
