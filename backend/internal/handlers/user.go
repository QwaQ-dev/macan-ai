package handlers

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/qwaq-dev/macan-ai/internal/services"
	"github.com/qwaq-dev/macan-ai/internal/structures"
	"github.com/qwaq-dev/macan-ai/pkg/sl"
)

type UserHandler struct {
	log         *slog.Logger
	userService *services.UserService
}

type UserServiceInterface interface {
	CreateUser(user *structures.UserResponse) (int, error)
}

func NewUserHandler(log *slog.Logger, userService *services.UserService) *UserHandler {
	return &UserHandler{
		log:         log,
		userService: userService,
	}
}

func (u *UserHandler) SignIn(c *fiber.Ctx) error {
	const op = "handlers.user.CreateUser"
	log := u.log.With("op", op)

	user := new(structures.UserResponse)

	if err := c.BodyParser(user); err != nil {
		log.Error("Invalid user format", sl.Err(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Invalid user format",
		})
	}

	if user.Password == "" || user.Username == "" {
		return c.Status(404).JSON(fiber.Map{
			"error": "Username and password are required",
		})
	}

	id, err := u.userService.CreateUser(user)
	if err != nil {
		log.Error("Error with creating user", sl.Err(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "User already exists",
		})
	}

	log.Info("User was created successfully", slog.Int("userId", id))
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success",
		"userId":  id,
	})
}
