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
	MESSAGE_FAILED_CREATE_FIELD        = "failed create field"
	MESSAGE_FAILED_GET_ALL_FIELD       = "failed get all field"
	MESSAGE_FAILED_GET_DETAIL_FIELD    = "failed get detail field"
	MESSAGE_FAILED_UPDATE_FIELD        = "failed update field"
	MESSAGE_FAILED_DELETE_FIELD        = "failed delete field"
	MESSAGE_FAILED_CREATE_SCHEDULE     = "failed create schedule"
	MESSAGE_FAILED_GET_ALL_SCHEDULE    = "failed get all schedule"
	MESSAGE_FAILED_UPDATE_SCHEDULE     = "failed update schedule"
	MESSAGE_FAILED_DELETE_SCHEDULE     = "failed delete schedule"
	MESSAGE_FAILED_GET_DETAIL_SCHEDULE = "failed get detail schedule"
	MESSAGE_FAILED_CREATE_BOOKING      = "failed create booking"
	MESSAGE_FAILED_GET_ALL_BOOKING     = "failed get all bookings"
	MESSAGE_FAILED_GET_DETAIL_BOOKING  = "failed get detail booking"
	MESSAGE_FAILED_UPDATE_BOOKING      = "failed update booking"
	MESSAGE_FAILED_DELETE_BOOKING      = "failed delete booking"

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
	MESSAGE_SUCCESS_CREATE_FIELD        = "success create field"
	MESSAGE_SUCCESS_GET_ALL_FIELD       = "success get all field"
	MESSAGE_SUCCESS_GET_DETAIL_FIELD    = "success get detail field"
	MESSAGE_SUCCESS_UPDATE_FIELD        = "success update field"
	MESSAGE_SUCCESS_DELETE_FIELD        = "success delete field"
	MESSAGE_SUCCESS_CREATE_SCHEDULE     = "success create schedule"
	MESSAGE_SUCCESS_GET_ALL_SCHEDULE    = "success get all schedule"
	MESSAGE_SUCCESS_UPDATE_SCHEDULE     = "success update schedule"
	MESSAGE_SUCCESS_DELETE_SCHEDULE     = "success delete schedule"
	MESSAGE_SUCCESS_GET_DETAIL_SCHEDULE = "success get detail schedule"
	MESSAGE_SUCCESS_CREATE_BOOKING      = "success create booking"
	MESSAGE_SUCCESS_GET_ALL_BOOKING     = "success get all bookings"
	MESSAGE_SUCCESS_GET_DETAIL_BOOKING  = "success get detail booking"
	MESSAGE_SUCCESS_UPDATE_BOOKING      = "success update booking"
	MESSAGE_SUCCESS_DELETE_BOOKING      = "success delete booking"
)

var (
	// Token-related errors
	ErrGenerateAccessToken     = errors.New("unable to generate access token")
	ErrGenerateRefreshToken    = errors.New("unable to generate refresh token")
	ErrUnexpectedSigningMethod = errors.New("invalid token signing method")
	ErrDecryptToken            = errors.New("unable to decrypt token")
	ErrTokenInvalid            = errors.New("token is invalid")
	ErrValidateToken           = errors.New("unable to validate token")

	// User-related errors
	ErrInvalidName              = errors.New("invalid name provided")
	ErrInvalidEmail             = errors.New("invalid email address")
	ErrInvalidPassword          = errors.New("invalid password provided")
	ErrEmailAlreadyExists       = errors.New("email address already registered")
	ErrRegisterUser             = errors.New("unable to register user")
	ErrGetAllUserWithPagination = errors.New("unable to retrieve paginated user list")
	ErrGetUserByID              = errors.New("unable to retrieve user by ID")
	ErrUpdateUser               = errors.New("unable to update user")
	ErrPasswordSame             = errors.New("new password cannot be the same as the old password")
	ErrHashPassword             = errors.New("unable to hash password")
	ErrDeleteUserByID           = errors.New("unable to delete user by ID")
	ErrEmailNotFound            = errors.New("email address not found")
	ErrPasswordNotMatch         = errors.New("incorrect password")
	ErrDeniedAccess             = errors.New("access denied")
	ErrGetPermissionsByRoleID   = errors.New("unable to retrieve permissions for role ID")
	ErrInvalidPhoneNumber       = errors.New("invalid phone number provided")

	// Category-related errors
	ErrCreateCategory     = errors.New("unable to create category")
	ErrGetAllCategory     = errors.New("unable to retrieve all categories")
	ErrInvalidUUID        = errors.New("invalid UUID provided")
	ErrGetCategoryByID    = errors.New("unable to retrieve category by ID")
	ErrUpdateCategory     = errors.New("unable to update category")
	ErrDeleteCategoryByID = errors.New("unable to delete category by ID")

	// Field-related errors
	ErrCreateField       = errors.New("unable to create field")
	ErrGetFieldByID      = errors.New("unable to retrieve field by ID")
	ErrGetAllField       = errors.New("unable to retrieve all fields")
	ErrUpdateField       = errors.New("unable to update field")
	ErrDeleteFieldByID   = errors.New("unable to delete field by ID")
	ErrInvalidFieldPrice = errors.New("field price cannot be negative")
	ErrSaveImages        = errors.New("unable to save image")
	ErrFieldNotFound     = errors.New("field not found")

	// Schedule-related errors
	ErrCreateSchedule         = errors.New("unable to create schedule")
	ErrGetAllSchedule         = errors.New("unable to retrieve all schedules")
	ErrScheduleNotFound       = errors.New("schedule not found")
	ErrUpdateSchedule         = errors.New("unable to update schedule")
	ErrDeleteSchedule         = errors.New("unable to delete schedule")
	ErrGetScheduleByID        = errors.New("unable to retrieve schedule by ID")
	ErrInvalidDayOfWeek       = errors.New("invalid day of the week")
	ErrCloseTimeMustAfterOpen = errors.New("closing time must be after opening time")
	ErrInvalidTimeFormat      = errors.New("invalid time format, expected HH:MM")

	// Booking-related errors
	ErrBookingTooSoon          = errors.New("booking must be at least 2 hours in advance")
	ErrOutsideOperatingHours   = errors.New("booking time is outside operating hours")
	ErrCreateBooking           = errors.New("unable to create booking")
	ErrUnauthorized            = errors.New("unauthorized access")
	ErrGetBookingByID          = errors.New("unable to retrieve booking by ID")
	ErrUpdateBooking           = errors.New("unable to update booking")
	ErrCannotCancelLate        = errors.New("cannot cancel booking; too close to booking time")
	ErrDeleteBooking           = errors.New("unable to delete booking")
	ErrInvalidBookingDate      = errors.New("invalid booking date format, expected YYYY-MM-DD")
	ErrInvalidStartTime        = errors.New("invalid start time format, expected HH:MM")
	ErrInvalidEndTime          = errors.New("invalid end time format, expected HH:MM")
	ErrCheckOverlap            = errors.New("unable to check for overlapping bookings")
	ErrBookingOverlap          = errors.New("booking time conflicts with an existing booking")
	ErrInvalidTotalPayment     = errors.New("total payment does not match field price and duration")
	ErrInvalidTimeRange        = errors.New("invalid time range provided")
	ErrInvalidStatusTransition = errors.New("invalid status transition for booking")
	ErrBookingAlreadyFinal     = errors.New("booking has already been finalized and cannot be updated")
	ErrInvalidStatusUpdate     = errors.New("invalid status update, only 'booked' or 'cancelled' allowed")
	ErrBookingNotFound = errors.New("")

	// General errors
	ErrInternalServer = errors.New("internal server error")
)
