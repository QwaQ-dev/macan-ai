package handlers

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/qwaq-dev/macan-ai/internal/services"
	"github.com/qwaq-dev/macan-ai/pkg/sl"
)

type UserResumeHandler struct {
	log               *slog.Logger
	userResumeService *services.UserResumeService
}

func NewUserResumeHandler(log *slog.Logger, userResumeService *services.UserResumeService) *UserResumeHandler {
	return &UserResumeHandler{
		log:               log,
		userResumeService: userResumeService,
	}
}

func (u *UserResumeHandler) UploadUserResume(c *fiber.Ctx) error {
	const op = "handlers.user_resume.UploadUserResume"
	log := u.log.With("op", op)

	resume, err := c.FormFile("resume")
	if err != nil {
		log.Error("Error with uploading resume", sl.Err(err))
		return c.Status(500).JSON(fiber.Map{
			"error": "Error with uploading resume",
		})
	}

	tempDir := os.TempDir()
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		os.Mkdir(tempDir, 0755)
	}

	uniqueFilename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), strings.ReplaceAll(resume.Filename, " ", "_"))

	filepath := filepath.Join(tempDir, uniqueFilename)

	if err := c.SaveFile(resume, filepath); err != nil {
		log.Error("Error saving uploaded resume file", sl.Err(err))
		return c.Status(500).JSON(fiber.Map{
			"error": "Error saving resume file",
		})
	}
	log.Debug("Resume file saved temporarily", slog.String("filepath", filepath))

	err = u.userResumeService.UploadResume(filepath)
	if err != nil {
		log.Error("Error with reading resume", sl.Err(err))
		return c.Status(500).JSON(fiber.Map{
			"error": "Error with reading resume",
		})
	}

	log.Debug("Resume was uploaded successfully", slog.String("filename", resume.Filename))
	return c.Status(200).JSON(fiber.Map{
		"message": "Resume was uploaded successfully",
	})
}
