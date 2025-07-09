package services

import (
	"fmt"
	"log/slog"

	"github.com/qwaq-dev/macan-ai/pkg/sl"
)

type UserResumeService struct {
	log                 *slog.Logger
	ResumeParsingClient *ResumeParsingClient
}

func NewUserResumeService(log *slog.Logger, ResumeParsingClient *ResumeParsingClient) *UserResumeService {
	return &UserResumeService{
		log:                 log,
		ResumeParsingClient: ResumeParsingClient,
	}
}

func (r *UserResumeService) UploadResume(filepath string) error {
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

	log.Info("Resume successfully parsed by gRPC service",
		slog.String("filename", filepath),
		slog.String("fullName", fmt.Sprintf("%s %s", parsedData.GetFullName().GetFirstName(), parsedData.GetFullName().GetLastName())),
		slog.Any("skills", parsedData.GetSkills()),
	)

	return nil
}
