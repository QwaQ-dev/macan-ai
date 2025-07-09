package services

import (
	"fmt"
	"log/slog"

	"github.com/qwaq-dev/macan-ai/internal/repository/postgres"
	"github.com/qwaq-dev/macan-ai/internal/structures"
	"github.com/qwaq-dev/macan-ai/pkg/sl"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	log  *slog.Logger
	repo *postgres.UserRepo
}

type UserRepoInterface interface {
	CreateUser(user *structures.UserResponse) (int, error)
	GetUserByUsername(username string) (*structures.UserResponse, error)
	GetUserById(id int) (*structures.UserResponse, error)
}

func NewUserService(log *slog.Logger, repo *postgres.UserRepo) *UserService {
	return &UserService{
		log:  log,
		repo: repo,
	}
}

func (s *UserService) CreateUser(user *structures.UserResponse) (int, error) {
	const op = "services.user.CreateUser"
	log := s.log.With("op", op)

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("Error with generating hash from password", sl.Err(err))
		return 0, err
	}
	user.Password = string(hash)

	userExists, err := s.repo.GetUserByUsername(user.Username)
	if err != nil {
		log.Error("Error with db", sl.Err(err))
		return 0, err
	}
	if userExists != nil {
		log.Info("User is already exists")
		return 0, fmt.Errorf("User is already exists")
	}

	id, err := s.repo.CreateUser(user)
	if err != nil {
		log.Error("Error with creating user", sl.Err(err))
		return 0, err
	}

	return id, nil
}
