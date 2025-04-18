package constants

var (
	ErrRoleNotFound                        = "Role not found"
	ErrUserNotFound                        = "User not found"
	ErrEmailAlreadyRegistered              = "Email already registered"
	ErrInternalServerError                 = "Internal server error"
	ErrEmailOrOTPNumberIsIncorrect         = "Email or otp number is incorrect"
	ErrEmailNotRegistered                  = "Email is not registered"
	ErrUserAlreadyVerified                 = "User already verified"
	ErrOTPNumberIsAlreadyExpired           = "OTP not found, please request a new OTP"
	ErrTooManyReuqestOTPNumber             = "Too many requests for the same email, please wait until the current OTP expires"
	ErrPhoneNumberAlreadyRegistered        = "Phone number already registered"
	ErrEmailOrPasswordIsIncorrect          = "Email or password is incorrect"
	ErrTokenAlreadyExpired                 = "Token already expired"
	ErrUserRegistrationInProgress          = "Your registration is already in progress. Please verify your OTP to complete the registration process"
	ErrAccessTokenIsRequired               = "Access token is required"
	ErrInvalidAccessToken                  = "Invalid access token"
	ErrInvalidOtpNumber                    = "Invalid otp number"
	ErrUserProfileNotFound                 = "User profile not found"
	ErrApplicationPermissionNotFound       = "Application permission is not found"
	ErrRoleAlreadyExist                    = "Role already exist"
	ErrApplicationDoenstHaveAccess         = "Application does not have access to this action"
	ErrApplicationForbidden                = "You do not have permission for this action"
	ErrUserNotVerified                     = "User not verified"
	ErrStaticRoleCannotBeDeleted           = "Static role cannot be deleted"
	ErrProductCategoryAlreadyRegistered    = "Product category already registered"
	ErrProductCategoryNotFound             = "Product category not found"
	ErrProductSubCategoryAlreadyRegistered = "Product sub category already registered"
	ErrProductSubCategoryNotFound          = "Product sub category not found"
	ErrInvalidProductData                  = "Invalid product_data JSON format"
	ErrProductNotFound                     = "Product not found"
	ErrProductVariantsNotFound             = "Product variants not found"
	ErrProductImagesNotFound               = "Product images not found"
	ErrProductGroceriesNotFound            = "Product groceries not found"
)
