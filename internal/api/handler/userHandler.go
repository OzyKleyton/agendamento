package handler

import (
	"strconv"

	"github.com/OzyKleyton/agendamento-api/internal/api/router"
	"github.com/OzyKleyton/agendamento-api/internal/model"
	"github.com/OzyKleyton/agendamento-api/internal/model/user"
	"github.com/OzyKleyton/agendamento-api/internal/service"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (uh *UserHandler) Routes() router.Router {
	return func(route fiber.Router) {
		user := route.Group("users")

		// Todas as rotas de usuário agora são protegidas pelo middleware JWT
		user.Post("/", uh.CreateUserHandler)
		user.Get("/", uh.FindAllUsersHandler)
		user.Get("/profile", uh.GetProfileHandler) // Nova rota para pegar perfil do usuário logado
		user.Get("/:email", uh.FindUserByEmailHandler)
		user.Put("/:id", uh.UpdateUserHandler)
		user.Put("/profile", uh.UpdateProfileHandler) // Nova rota para atualizar perfil do usuário logado
		user.Delete("/:id", uh.DeleteUserHandler)
	}
}

func (uh *UserHandler) CreateUserHandler(c *fiber.Ctx) error {
	userReq := new(user.UserReq)
	if err := c.BodyParser(userReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.NewErrorResponse(err, fiber.ErrBadRequest))
	}

	res := uh.service.CreateUser(userReq)

	return c.Status(res.Status).JSON(res)
}

func (uh *UserHandler) FindAllUsersHandler(c *fiber.Ctx) error {
	res := uh.service.FindAllUsers()

	return c.Status(res.Status).JSON(res)
}

func (uh *UserHandler) FindUserByEmailHandler(c *fiber.Ctx) error {
	userEmail := c.Params("email")

	res := uh.service.FindUserByEmail(userEmail)

	return c.Status(res.Status).JSON(res)
}

func (uh *UserHandler) UpdateUserHandler(c *fiber.Ctx) error {
	userReq := new(user.UserReq)
	if err := c.BodyParser(userReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.NewErrorResponse(err, fiber.ErrBadRequest))
	}

	userID, _ := strconv.Atoi(c.Params("id", "0"))

	res := uh.service.UpdateUser(uint(userID), userReq)

	return c.Status(res.Status).JSON(res)
}

func (uh *UserHandler) DeleteUserHandler(c *fiber.Ctx) error {
	userID, _ := strconv.Atoi(c.Params("id", "0"))

	res := uh.service.DeleteUser(uint(userID))

	return c.Status(res.Status).JSON(res)
}

// NOVOS HANDLERS PARA USUÁRIO LOGADO

// GetProfileHandler - Retorna o perfil do usuário logado
func (uh *UserHandler) GetProfileHandler(c *fiber.Ctx) error {
	// Pega o email do usuário logado do contexto (setado pelo middleware)
	userEmail := c.Locals("email").(string)

	if userEmail == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(model.NewErrorResponse(
			fiber.ErrUnauthorized,
			fiber.ErrUnauthorized,
		))
	}

	res := uh.service.FindUserByEmail(userEmail)

	return c.Status(res.Status).JSON(res)
}

// UpdateProfileHandler - Atualiza o perfil do usuário logado
func (uh *UserHandler) UpdateProfileHandler(c *fiber.Ctx) error {
	// Pega o ID do usuário logado do contexto
	userID := c.Locals("user_id").(uint)

	if userID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(model.NewErrorResponse(
			fiber.ErrUnauthorized,
			fiber.ErrUnauthorized,
		))
	}

	userReq := new(user.UserReq)
	if err := c.BodyParser(userReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.NewErrorResponse(err, fiber.ErrBadRequest))
	}

	res := uh.service.UpdateUser(userID, userReq)

	return c.Status(res.Status).JSON(res)
}
