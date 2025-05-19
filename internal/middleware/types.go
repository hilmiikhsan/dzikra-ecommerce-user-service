package middleware

import (
	role "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/role/dto"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type Locals struct {
	UserID      string
	Email       string
	FullName    string
	SessionID   string
	DeviceID    string
	DeviceType  string
	FcmToken    string
	UserRoles   []role.UserRoleDetail
	PhoneNumber string
}

func GetLocals(c *fiber.Ctx) *Locals {
	var l = Locals{}
	UserID, ok := c.Locals("user_id").(string)
	if ok {
		l.UserID = UserID
	} else {
		log.Warn().Msg("middleware::Locals-GetLocals failed to get user_id from locals")
	}

	email, ok := c.Locals("email").(string)
	if ok {
		l.Email = email
	} else {
		log.Warn().Msg("middleware::Locals-GetLocals failed to get email from locals")
	}

	fullName, ok := c.Locals("full_name").(string)
	if ok {
		l.FullName = fullName
	} else {
		log.Warn().Msg("middleware::Locals-GetLocals failed to get full_name from locals")
	}

	sessionID, ok := c.Locals("session_id").(string)
	if ok {
		l.SessionID = sessionID
	} else {
		log.Warn().Msg("middleware::Locals-GetLocals failed to get session_id from locals")
	}

	deviceID, ok := c.Locals("device_id").(string)
	if ok {
		l.DeviceID = deviceID
	} else {
		log.Warn().Msg("middleware::Locals-GetLocals failed to get device_id from locals")
	}

	deviceType, ok := c.Locals("device_type").(string)
	if ok {
		l.DeviceType = deviceType
	} else {
		log.Warn().Msg("middleware::Locals-GetLocals failed to get device_type from locals")
	}

	fcmToken, ok := c.Locals("fcm_token").(string)
	if ok {
		l.FcmToken = fcmToken
	} else {
		log.Warn().Msg("middleware::Locals-GetLocals failed to get fcm_token from locals")
	}

	userRoles, ok := c.Locals("user_roles").([]role.UserRoleDetail)
	if ok {
		l.UserRoles = userRoles
	} else {
		log.Warn().Msg("middleware::Locals-GetLocals failed to get user_roles from locals")
	}

	phoneNumber, ok := c.Locals("phone_number").(string)
	if ok {
		l.PhoneNumber = phoneNumber
	} else {
		log.Warn().Msg("middleware::Locals-GetLocals failed to get phone_number from locals")
	}

	return &l
}

func (l *Locals) GetUserID() string {
	return l.UserID
}

func (l *Locals) GetEmail() string {
	return l.Email
}

func (l *Locals) GetFullName() string {
	return l.FullName
}

func (l *Locals) GetSessionID() string {
	return l.SessionID
}

func (l *Locals) GetDeviceID() string {
	return l.DeviceID
}

func (l *Locals) GetDeviceType() string {
	return l.DeviceType
}

func (l *Locals) GetFcmToken() string {
	return l.FcmToken
}

func (l *Locals) GetUserRoles() []role.UserRoleDetail {
	return l.UserRoles
}

func (l *Locals) GetPhoneNumber() string {
	return l.PhoneNumber
}
