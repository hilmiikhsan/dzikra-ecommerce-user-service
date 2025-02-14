package constants

var (
	ErrRoleNotFound                    = "Role not found"
	ErrEmailAlreadyRegistered          = "Email already registered"
	ErrUsernameAlreadyRegistered       = "Username already registered"
	ErrInternalServerError             = "Internal server error"
	ErrEmailOrOTPNumberIsIncorrect     = "Email or otp number is incorrect"
	ErrEmailNotRegistered              = "Email is not registered"
	ErrUserAlreadyVerified             = "User already verified"
	ErrOTPNumberIsAlreadyExpired       = "OTP number is already expired"
	ErrTooManyReuqestOTPNumber         = "Too many requests for the same email, please wait until the current OTP expires"
	ErrPhoneNumberAlreadyRegistered    = "Phone number already registered"
	ErrIdentifierOrPasswordIsIncorrect = "Identifier or password is incorrect"
	ErrTokenAlreadyExpired             = "Token already expired"
	ErrUserRegistrationInProgress      = "Your registration is already in progress. Please verify your OTP to complete the registration process"
	ErrAccessTokenIsRequired           = "Access token is required"
	ErrInvalidAccessToken              = "Invalid access token"
	ErrInvalidOtpNumber                = "Invalid otp number"
)
