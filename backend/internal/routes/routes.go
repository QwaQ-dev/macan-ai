package routes

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/qwaq-dev/macan-ai/internal/handlers"
)

func InitRoutes(
	app *fiber.App,
	log *slog.Logger,
	userHandler *handlers.UserHandler,
	userResumeHandler *handlers.UserResumeHandler,
) {
	user := app.Group("/user")

	user.Post("/sign-in", userHandler.SignIn)
	user.Post("/resume", userResumeHandler.UploadUserResume)

	log.Debug("All routes has been initialized!")
}
