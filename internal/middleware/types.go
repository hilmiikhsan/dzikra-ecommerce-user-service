package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type Locals struct {
	UserID   int
	Username string
	Email    string
	FullName string
}

func GetLocals(c *fiber.Ctx) *Locals {
	var l = Locals{}
	UserID, ok := c.Locals("user_id").(int)
	if ok {
		l.UserID = UserID
	} else {
		log.Warn().Msg("middleware::Locals-GetLocals failed to get user_id from locals")
	}

	username, ok := c.Locals("username").(string)
	if ok {
		l.Username = username
	} else {
		log.Warn().Msg("middleware::Locals-GetLocals failed to get username from locals")
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

	return &l
}

func (l *Locals) GetUserID() int {
	return l.UserID
}

func (l *Locals) GetUsername() string {
	return l.Username
}

func (l *Locals) GetEmail() string {
	return l.Email
}

func (l *Locals) GetFullName() string {
	return l.FullName
}
