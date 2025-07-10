package services

import (
	"fmt"
	"log/slog"

	"github.com/qwaq-dev/macan-ai/internal/repository/postgres"
	"github.com/qwaq-dev/macan-ai/pkg/sl"
)

type UserResumeService struct {
	log                 *slog.Logger
	ResumeParsingClient *ResumeParsingClient
	UserResumeRepo      *postgres.UserResumeRepo
}

func NewUserResumeService(log *slog.Logger, ResumeParsingClient *ResumeParsingClient, userResumeRepo *postgres.UserResumeRepo) *UserResumeService {
	return &UserResumeService{
		log:                 log,
		ResumeParsingClient: ResumeParsingClient,
		UserResumeRepo:      userResumeRepo,
	}
}

func (r *UserResumeService) UploadResume(userId int64, filepath string) error {
	const op = "services.user_resume.UploadResume"
	log := r.log.With("op", op)

	log.Debug("Calling Resume Parsing Service with filepath", slog.String("filepath", filepath))

	parsedData, err := r.ResumeParsingClient.ParseResume(filepath)
	if err != nil {
		log.Error("Failed parse resume via gRPC service", sl.Err(err))
		return err
	}

	if !parsedData.GetSuccess() {
		log.Warn("Resume parsing service reported failure", slog.String("filepath", filepath))
		return fmt.Errorf("resume parsing service reported failure for file: %s", filepath)
	}

	err = r.UserResumeRepo.AddResume(parsedData, int(userId))
	if err != nil {
		log.Error("Error with repository", sl.Err(err))
		return err
	}
	return nil
}
