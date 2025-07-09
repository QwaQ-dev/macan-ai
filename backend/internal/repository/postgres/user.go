package postgres

import (
	"database/sql"
	"errors"
	"log/slog"

	"github.com/qwaq-dev/macan-ai/internal/structures"
	"github.com/qwaq-dev/macan-ai/pkg/sl"
)

type UserRepo struct {
	log *slog.Logger
	db  *sql.DB
}

func NewUserRepo(log *slog.Logger, db *sql.DB) *UserRepo {
	return &UserRepo{
		log: log,
		db:  db,
	}
}

func (r *UserRepo) CreateUser(user *structures.UserResponse) (int, error) {
	const op = "postgres.user.CreateUser"
	log := r.log.With("op", op)
	var id int

	query := "INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id"
	err := r.db.QueryRow(query, user.Username, user.Password).Scan(&id)
	if err != nil {
		log.Error("Error with inserting user data", sl.Err(err))
		return id, nil
	}

	return id, nil
}

func (r *UserRepo) GetUserByUsername(username string) (*structures.UserResponse, error) {
	const op = "postgres.user.GetUserByUsername"
	log := r.log.With("op", op)

	query := "SELECT id, username, password FROM users WHERE username=$1"

	user := new(structures.UserResponse)

	err := r.db.QueryRow(query, username).Scan(&user.Id, &user.Username, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		log.Debug("User with this username already exists", sl.Err(err))
		return nil, err
	}

	return user, err
}

func (r *UserRepo) GetUserById(id int) (*structures.UserResponse, error) {
	const op = "postgres.user.GetUserById"
	log := r.log.With("op", op)

	user := new(structures.UserResponse)

	query := "SELECT username, password FROM users WHERE id=$1"

	err := r.db.QueryRow(query, id).Scan(&user.Username, &user.Password)
	if err != nil {
		log.Error("Error with scanning user data", sl.Err(err))
		return user, err
	}

	user.Id = id

	return user, nil
}
