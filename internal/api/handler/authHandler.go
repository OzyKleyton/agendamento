package handler

import (
	"github.com/OzyKleyton/agendamento-api/internal/api/router"
	"github.com/OzyKleyton/agendamento-api/internal/model"
	"github.com/OzyKleyton/agendamento-api/internal/model/user"
	"github.com/OzyKleyton/agendamento-api/internal/service"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	service service.AuthService
}

func NewAuthHandler(service service.AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

func (ah *AuthHandler) Routes() router.Router {
	return func(route fiber.Router) {

		auth := route.Group("auth")
		auth.Post("/login", ah.LoginHandler)
		auth.Post("/refresh", ah.RefreshHandler)
		auth.Post("/register", ah.RegisterHandler)
	}
}

func (ah *AuthHandler) LoginHandler(c *fiber.Ctx) error {
	loginReq := new(user.LoginRequest)
	if err := c.BodyParser(loginReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.NewErrorResponse(err, fiber.ErrBadRequest))
	}

	res := ah.service.Login(loginReq)

	return c.Status(res.Status).JSON(res)
}

func (ah *AuthHandler) RefreshHandler(c *fiber.Ctx) error {
	refreshReq := new(user.RefreshRequest)
	if err := c.BodyParser(refreshReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.NewErrorResponse(err, fiber.ErrBadRequest))
	}

	res := ah.service.RefreshToken(refreshReq)

	return c.Status(res.Status).JSON(res)
}

func (ah *AuthHandler) RegisterHandler(c *fiber.Ctx) error {
	// Esta rota pode ser implementada se quiser registro público
	return c.Status(fiber.StatusNotImplemented).JSON(model.NewErrorResponse(
		fiber.ErrNotImplemented,
		fiber.ErrNotImplemented,
	))
}
