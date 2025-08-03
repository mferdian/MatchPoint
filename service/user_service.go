package service

import (
	"context"
	"fieldreserve/constants"
	"fieldreserve/dto"
	"fieldreserve/helpers"
	"fieldreserve/model"
	"fieldreserve/repository"
	"fieldreserve/utils"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type (
	IUserService interface {
		CreateUser(ctx context.Context, req dto.CreateUserRequest) (dto.UserResponse, error)
		GetUserByEmail(ctx context.Context, req dto.LoginUserRequest) (dto.LoginResponse, error)
		GetuserByID(ctx context.Context, userID string) (dto.UserResponse, error)
		GetAllUserWithPagination(ctx context.Context, req dto.UserPaginationRequest) (dto.UserPaginationResponse, error)
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
	if flag && err == nil {
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
		Role:     constants.ENUM_ROLE_USER,
	}

	err = us.userRepo.CreateUser(ctx, nil, user)
	if err != nil {
		utils.Log.WithError(err).WithField("email", req.Email).Error("Failed to create user")
		return dto.UserResponse{}, constants.ErrRegisterUser
	}

	utils.Log.WithFields(logrus.Fields{
		"user_id": user.UserID,
		"email":   user.Email,
	}).Info("User created successfully")

	res := dto.UserResponse{
		ID:    user.UserID,
		Name:  user.Name,
		Email: user.Email,
	}

	return res, nil
}

func (us *UserService) GetUserByEmail(ctx context.Context, req dto.LoginUserRequest) (dto.LoginResponse, error) {
	if !helpers.IsValidEmail(req.Email) {
		return dto.LoginResponse{}, constants.ErrInvalidEmail
	}

	if len(req.Password) < 8 {
		return dto.LoginResponse{}, constants.ErrInvalidPassword
	}

	user, flag, err := us.userRepo.GetUserByEmail(ctx, nil, req.Email)
	if !flag || err != nil {
		utils.Log.WithField("email", req.Email).Warn("Login failed: email not found")
		return dto.LoginResponse{}, constants.ErrEmailNotFound
	}

	checkPassword, err := helpers.CheckPassword(user.Password, []byte(req.Password))
	if err != nil || !checkPassword {
		utils.Log.WithField("email", req.Email).Warn("Login failed: password mismatch")
		return dto.LoginResponse{}, constants.ErrPasswordNotMatch
	}

	accessToken, refreshToken, err := us.jwtService.GenerateToken(user.UserID.String(), user.Role)
	if err != nil {
		utils.Log.WithError(err).WithField("user_id", user.UserID).Error("Failed to generate token")
		return dto.LoginResponse{}, err
	}

	utils.Log.WithFields(logrus.Fields{
		"user_id": user.UserID,
		"email":   user.Email,
	}).Info("User login successful")

	return dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (us *UserService) GetuserByID(ctx context.Context, userID string) (dto.UserResponse, error) {
	if _, err := uuid.Parse(userID); err != nil {
		return dto.UserResponse{}, constants.ErrInvalidUUID
	}

	user, _, err := us.userRepo.GetUserByID(ctx, nil, userID)
	if err != nil {
		utils.Log.WithError(err).WithField("user_id", userID).Error("Failed to get user by ID")
		return dto.UserResponse{}, constants.ErrGetUserByID
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

func (us *UserService) GetAllUserWithPagination(ctx context.Context, req dto.UserPaginationRequest) (dto.UserPaginationResponse, error) {
	dataWithPaginate, err := us.userRepo.GetAllUserWithPagination(ctx, nil, req)
	if err != nil {
		utils.Log.WithError(err).Error("Failed to get all users with pagination")
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
		utils.Log.WithError(err).WithField("user_id", req.ID).Error("Failed to get user for update")
		return dto.UserResponse{}, constants.ErrGetUserByID
	}

	if req.Email != "" && req.Email != user.Email {
		if !helpers.IsValidEmail(req.Email) {
			return dto.UserResponse{}, constants.ErrInvalidEmail
		}

		existingUser, exists, err := us.userRepo.GetUserByEmail(ctx, nil, req.Email)
		if err == nil && exists && existingUser.UserID != user.UserID {
			return dto.UserResponse{}, constants.ErrEmailAlreadyExists
		}

		user.Email = req.Email
	}

	if req.Name != "" {
		if len(req.Name) < 5 {
			return dto.UserResponse{}, constants.ErrInvalidName
		}
		user.Name = req.Name
	}

	if req.Password != "" {
		isSame, _ := helpers.CheckPassword(user.Password, []byte(req.Password))
		if isSame {
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
		utils.Log.WithError(err).WithField("user_id", user.UserID).Error("Failed to update user")
		return dto.UserResponse{}, constants.ErrUpdateUser
	}

	utils.Log.WithFields(logrus.Fields{
		"user_id": user.UserID,
		"email":   user.Email,
	}).Info("User updated successfully")

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
		utils.Log.WithError(err).WithField("user_id", req.UserID).Error("Failed to get user before delete")
		return dto.UserResponse{}, constants.ErrGetUserByID
	}

	err = us.userRepo.DeleteUserByID(ctx, nil, req.UserID)
	if err != nil {
		utils.Log.WithError(err).WithField("user_id", req.UserID).Error("Failed to delete user")
		return dto.UserResponse{}, constants.ErrDeleteUserByID
	}

	utils.Log.WithFields(logrus.Fields{
		"user_id": deletedUser.UserID,
		"email":   deletedUser.Email,
	}).Info("User deleted successfully")

	res := dto.UserResponse{
		ID:      deletedUser.UserID,
		Name:    deletedUser.Name,
		Email:   deletedUser.Email,
		Address: deletedUser.Address,
		NoTelp:  deletedUser.NoTelp,
	}

	return res, nil
}
