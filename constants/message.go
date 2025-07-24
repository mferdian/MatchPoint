package constants

import "errors"

const (
	// failed
	MESSAGE_FAILED_PROSES_REQUEST      = "failed proses request"
	MESSAGE_FAILED_ACCESS_DENIED       = "failed access denied"
	MESSAGE_FAILED_TOKEN_NOT_FOUND     = "failed token not found"
	MESSAGE_FAILED_TOKEN_NOT_VALID     = "failed token not valid"
	MESSAGE_FAILED_TOKEN_DENIED_ACCESS = "failed token denied access"
	MESSAGE_FAILED_GET_DATA_FROM_BODY  = "failed get data from body"
	MESSAGE_FAILED_CREATE_USER         = "failed create user"
	MESSAGE_FAILED_GET_DETAIL_USER     = "failed get detail user"
	MESSAGE_FAILED_GET_LIST_USER       = "failed get list user"
	MESSAGE_FAILED_UPDATE_USER         = "failed update user"
	MESSAGE_FAILED_DELETE_USER         = "failed delete user"
	MESSAGE_FAILED_LOGIN_USER          = "failed login user"
	MESSAGE_FAILED_CREATE_CATEGORY     = "failed create category"
	MESSAGE_FAILED_GET_ALL_CATEGORY    = "failed get all category"
	MESSAGE_FAILED_UUID_FORMAT         = "failed uuid format"
	MESSAGE_FAILED_GET_DETAIL_CATEGORY = "failed get detail category"
	MESSAGE_FAILED_UPDATE_CATEGORY     = "failed update category"
	MESSAGE_FAILED_DELETE_CATEGORY     = "failed delete category"

	// success
	MESSAGE_SUCCESS_CREATE_USER         = "success create user"
	MESSAGE_SUCCESS_GET_DETAIL_USER     = "success get detail user"
	MESSAGE_SUCCESS_GET_LIST_USER       = "success get list user"
	MESSAGE_SUCCESS_UPDATE_USER         = "success update user"
	MESSAGE_SUCCESS_DELETE_USER         = "success delete user"
	MESSAGE_SUCCESS_CREATE_CATEGORY     = "success create category"
	MESSAGE_SUCCESS_GET_ALL_CATEGORY    = "success get all category"
	MESSAGE_SUCCESS_GET_DETAIL_CATEGORY = "success get detail category"
	MESSAGE_SUCCESS_UPDATE_CATEGORY     = "success update category"
	MESSAGE_SUCCESS_DELETE_CATEGORY     = "success delete category"
)

var (
	ErrGenerateAccessToken      = errors.New("failed to generate access token")
	ErrGenerateRefreshToken     = errors.New("failed to generate refresh token")
	ErrUnexpectedSigningMethod  = errors.New("unexpected signing method")
	ErrDecryptToken             = errors.New("failed to decrypt token")
	ErrTokenInvalid             = errors.New("token invalid")
	ErrValidateToken            = errors.New("failed to validate token")
	ErrInvalidName              = errors.New("failed invalid name")
	ErrInvalidEmail             = errors.New("failed invalid email")
	ErrInvalidPassword          = errors.New("failed invalid password")
	ErrEmailAlreadyExists       = errors.New("email already exists")
	ErrRegisterUser             = errors.New("failed to register user")
	ErrGetAllUserWithPagination = errors.New("failed get list user with pagination")
	ErrGetUserByID              = errors.New("failed get user by id")
	ErrUpdateUser               = errors.New("failed to update user")
	ErrPasswordSame             = errors.New("failed new password same as old password")
	ErrHashPassword             = errors.New("failed hash password")
	ErrDeleteUserByID           = errors.New("failed delete user by id")
	ErrEmailNotFound            = errors.New("email not found")
	ErrPasswordNotMatch         = errors.New("password not match")
	ErrDeniedAccess             = errors.New("denied access")
	ErrGetPermissionsByRoleID   = errors.New("failed get all permission by role id")
	ErrInvalidPhoneNumber       = errors.New("invalid phone number")
	ErrCreateCategory           = errors.New("failed created category")
	ErrGetAllCategory           = errors.New("failed get all category")
	ErrInvalidUUID              = errors.New("uuid is invalid")
	ErrGetCategoryByID          = errors.New("failed get category by id")
	ErrUpdateCategory           = errors.New("failed update category")
	ErrDeleteCategoryByID       = errors.New("failed deleted category")
)
