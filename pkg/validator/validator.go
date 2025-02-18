package validator

import (
	"reflect"
	"regexp"
	"strings"

	// "github.com/go-playground/locales/en"
	// ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	// en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/rs/zerolog/log"
)

type Validator struct {
	// trans     ut.Translator
	validator *validator.Validate
}

func NewValidator() *Validator {
	validatorCustom := &Validator{}

	// en := en.New()
	// uni := ut.New(en, en)
	// trans, _ := uni.GetTranslator("en")

	v := validator.New()
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		var name string

		name = strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "" {
			name = strings.SplitN(fld.Tag.Get("query"), ",", 2)[0]
		}

		if name == "" {
			name = strings.SplitN(fld.Tag.Get("form"), ",", 2)[0]
		}

		if name == "" {
			name = strings.SplitN(fld.Tag.Get("params"), ",", 2)[0]
		}

		if name == "" {
			name = strings.SplitN(fld.Tag.Get("prop"), ",", 2)[0]
		}

		if name == "-" {
			return ""
		}

		return name
	})

	// en_translations.RegisterDefaultTranslations(v, trans)
	if err := v.RegisterValidation("email_blacklist", isEmailBlacklist); err != nil {
		log.Fatal().Err(err).Msg("Error while registering email_blacklist validator")
	}
	if err := v.RegisterValidation("strong_password", isStrongPassword); err != nil {
		log.Fatal().Err(err).Msg("Error while registering strong_password validator")
	}
	if err := v.RegisterValidation("unique_in_slice", isUniqueInSlice); err != nil {
		log.Fatal().Err(err).Msg("Error while registering unique validator")
	}
	if err := v.RegisterValidation("phone", phoneValidator); err != nil {
		log.Fatal().Err(err).Msg("Error while registering phone validator")
	}
	if err := v.RegisterValidation("otp_number", otpNumberValidation); err != nil {
		log.Fatal().Err(err).Msg("Error while registering otp_number validator")
	}
	if err := v.RegisterValidation("google_token", isGoogleTokenValid); err != nil {
		log.Fatal().Err(err).Msg("Error while registering google_token validator")
	}

	validatorCustom.validator = v
	// validatorCustom.trans = trans

	return validatorCustom
}

func (v *Validator) Validate(i any) error {
	return v.validator.Struct(i)
}

// blacklist email validator
func isEmailBlacklist(fl validator.FieldLevel) bool {
	email := fl.Field().String()
	disallowedDomains := []string{"outlook", "hotmail", "aol", "live", "inbox", "icloud", "mail", "gmx", "yandex"}

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	domain := strings.Split(parts[1], ".")[0]

	for _, disallowed := range disallowedDomains {
		if domain == disallowed {
			return false
		}
	}

	return true
}

func isStrongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) < 8 {
		return false
	}

	hasUppercase := false
	hasLowercase := false
	hasNumber := false

	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUppercase = true
		case char >= 'a' && char <= 'z':
			hasLowercase = true
		case char >= '0' && char <= '9':
			hasNumber = true
		}
	}

	return hasUppercase && hasLowercase && hasNumber
}

func isUniqueInSlice(fl validator.FieldLevel) bool {
	// Get the slice from the FieldLevel interface
	val := fl.Field()

	// Ensure the field is a slice
	if val.Kind() != reflect.Slice {
		return false
	}

	// Use a map to check for duplicates
	elements := make(map[interface{}]bool)
	for i := 0; i < val.Len(); i++ {
		elem := val.Index(i).Interface()
		if _, found := elements[elem]; found {
			return false // Duplicate found
		}
		elements[elem] = true
	}
	return true
}

func phoneValidator(fl validator.FieldLevel) bool {
	// +62, 62, 0
	phoneRegex := `^(?:\+62|62|0)[2-9][0-9]{7,14}$`
	re := regexp.MustCompile(phoneRegex)
	return re.MatchString(fl.Field().String())
}

func otpNumberValidation(fl validator.FieldLevel) bool {
	otp := fl.Field().String()
	matched, _ := regexp.MatchString(`^[0-9A-Z]{6}$`, otp)
	return matched
}

func isGoogleTokenValid(fl validator.FieldLevel) bool {
	token := fl.Field().String()

	if len(token) < 100 || len(token) > 1000 {
		return false
	}

	jwtRegex := `^[a-zA-Z0-9_-]+\.[a-zA-Z0-9_-]+\.[a-zA-Z0-9_-]+$`
	matched, err := regexp.MatchString(jwtRegex, token)
	if err != nil || !matched {
		return false
	}

	return true
}
