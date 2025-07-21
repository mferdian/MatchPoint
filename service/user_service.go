package service

import (
	"context"
	"fieldreserve/constants"
	"fieldreserve/dto"
	"fieldreserve/helpers"
	"fieldreserve/model"
	"fieldreserve/repository"

	"github.com/google/uuid"
)

type (
	IUserService interface {
		CreateUser(ctx context.Context, req dto.CreateUserRequest) (dto.UserResponse, error)
		ReadUserByEmail(ctx context.Context, req dto.LoginUserRequest) (dto.LoginResponse, error)
		ReadAllUserWithPagination(ctx context.Context, req dto.UserPaginationRequest) (dto.UserPaginationResponse, error)
		UpdateUser(ctx context.Context, req dto.UpdateUserRequest) (dto.UserResponse, error)
		DeleteUser(ctx context.Context, req dto.DeleteUserRequest) (dto.UserResponse, error)
	}

	UserService struct {
		userRepo   repository.IUserRepository
		jwtService InterfaceJWTService
	}
)

func NewUserService(userRepo repository.IUserRepository, jwtService InterfaceJWTService) *UserService {
	return &UserService{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

func (us *UserService) CreateUser(ctx context.Context, req dto.CreateUserRequest) (dto.UserResponse, error) {
	if len(req.Name) < 3 {
		return dto.UserResponse{}, constants.ErrInvalidName
	}

	if !helpers.IsValidEmail(req.Email) {
		return dto.UserResponse{}, constants.ErrInvalidEmail
	}

	_, flag, err := us.userRepo.GetUserByEmail(ctx, nil, req.Email)
	if flag || err == nil {
		return dto.UserResponse{}, constants.ErrEmailAlreadyExists
	}

	if len(req.Password) < 8 {
		return dto.UserResponse{}, constants.ErrInvalidPassword
	}

	user := model.User{
		UserID:   uuid.New(),
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Role: constants.ENUM_ROLE_USER,
	}

	err = us.userRepo.CreateUser(ctx, nil, user)
	if err != nil {
		return dto.UserResponse{}, constants.ErrRegisterUser
	}

	res := dto.UserResponse{
		ID:      user.UserID,
		Name:    user.Name,
		Email:   user.Email,
	}

	return res, nil
}


func (us *UserService) ReadUserByEmail(ctx context.Context, req dto.LoginUserRequest) (dto.LoginResponse, error) {
	if !helpers.IsValidEmail(req.Email) {
		return dto.LoginResponse{}, constants.ErrInvalidEmail
	}

	if len(req.Password) < 8 {
		return dto.LoginResponse{}, constants.ErrInvalidPassword
	}

	user, flag, err := us.userRepo.GetUserByEmail(ctx, nil, req.Email)
	if !flag || err != nil {
		return dto.LoginResponse{}, constants.ErrEmailNotFound
	}

	checkPassword, err := helpers.CheckPassword(user.Password, []byte(req.Password))
	if err != nil || !checkPassword {
		return dto.LoginResponse{}, constants.ErrPasswordNotMatch
	}

	accessToken, refreshToken, err := us.jwtService.GenerateToken(user.UserID.String(), user.Role)
	if err != nil {
		return dto.LoginResponse{}, err
	}

	return dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}


func (us *UserService) ReadAllUserWithPagination(ctx context.Context, req dto.UserPaginationRequest) (dto.UserPaginationResponse, error) {
	dataWithPaginate, err := us.userRepo.GetAllUserWithPagination(ctx, nil, req)
	if err != nil {
		return dto.UserPaginationResponse{}, constants.ErrGetAllUserWithPagination
	}

	var datas []dto.UserResponse
	for _, user := range dataWithPaginate.Users {
		data := dto.UserResponse{
			ID:      user.UserID,
			Name:    user.Name,
			Email:   user.Email,
			Address: user.Address,
			NoTelp:  user.NoTelp,
		}

		datas = append(datas, data)
	}

	return dto.UserPaginationResponse{
		Data: datas,
		PaginationResponse: dto.PaginationResponse{
			Page:    dataWithPaginate.Page,
			PerPage: dataWithPaginate.PerPage,
			MaxPage: dataWithPaginate.MaxPage,
			Count:   dataWithPaginate.Count,
		},
	}, nil
}

func (us *UserService) UpdateUser(ctx context.Context, req dto.UpdateUserRequest) (dto.UserResponse, error) {
	user, _, err := us.userRepo.GetUserByID(ctx, nil, req.ID)
	if err != nil {
		return dto.UserResponse{}, constants.ErrGetUserByID
	}


	if req.Name != "" {
		if len(req.Name) < 5 {
			return dto.UserResponse{}, constants.ErrInvalidName
		}
		user.Name = req.Name
	}

	if req.Email != "" {
		if !helpers.IsValidEmail(req.Email) {
			return dto.UserResponse{}, constants.ErrInvalidEmail
		}

		_, exists, err := us.userRepo.GetUserByEmail(ctx, nil, req.Email)
		if exists || err == nil {
			return dto.UserResponse{}, constants.ErrEmailAlreadyExists
		}

		user.Email = req.Email
	}

	if req.Password != "" {
		if checkPassword, err := helpers.CheckPassword(user.Password, []byte(req.Password)); checkPassword || err == nil {
			return dto.UserResponse{}, constants.ErrPasswordSame
		}

		hashP, err := helpers.HashPassword(req.Password)
		if err != nil {
			return dto.UserResponse{}, constants.ErrHashPassword
		}

		user.Password = hashP
	}

	if req.Address != "" {
		user.Address = req.Address
	}

	if req.NoTelp != "" {
		if len(req.NoTelp) < 10 {
			return dto.UserResponse{}, constants.ErrInvalidPhoneNumber
		}
		user.NoTelp = req.NoTelp
	}

	err = us.userRepo.UpdateUser(ctx, nil, user)
	if err != nil {
		return dto.UserResponse{}, constants.ErrUpdateUser
	}

	res := dto.UserResponse{
		ID:      user.UserID,
		Name:    user.Name,
		Email:   user.Email,
		Address: user.Address,
		NoTelp:  user.NoTelp,
	}

	return res, nil
}


func (us *UserService) DeleteUser(ctx context.Context, req dto.DeleteUserRequest) (dto.UserResponse, error) {
	deletedUser, _, err := us.userRepo.GetUserByID(ctx, nil, req.UserID)
	if err != nil {
		return dto.UserResponse{}, constants.ErrGetUserByID
	}

	err = us.userRepo.DeleteUserByID(ctx, nil, req.UserID)
	if err != nil {
		return dto.UserResponse{}, constants.ErrDeleteUserByID
	}

	res := dto.UserResponse{
		ID:      deletedUser.UserID,
		Name:    deletedUser.Name,
		Email:   deletedUser.Email,
		Address: deletedUser.Address,
		NoTelp:  deletedUser.NoTelp,
	}

	return res, nil
}
