package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/humanbelnik/backy-auth-ms/internal/config"
	"github.com/humanbelnik/backy-auth-ms/internal/domain/models"
	"github.com/humanbelnik/backy-auth-ms/internal/storage"
	_ "github.com/lib/pq"
	"golang.org/x/exp/slog"
)

const (
	queryRegisterUser        = "INSERT INTO users(id, email, nickname, pass_hash) VALUES($1, $2, $3, $4)"
	queryProvideWithEmail    = "SELECT id, email, nickname, pass_hash FROM users WHERE email = $1"
	queryProvideWithNickname = "SELECT id, email, nickname, pass_hash FROM users WHERE nickname = $1"
)

type Storage struct {
	db  *sql.DB
	log *slog.Logger
}

// New opens new Postgres connection.
// Returns storage instance and possible error.
func New(log *slog.Logger, config config.DatabaseConfig) (*Storage, error) {
	const fn = "storage.postgres.New"
	log = log.With(
		slog.String("fn", fn),
	)
	log.Info("starting db connection")

	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.Host, config.User, config.Password, config.Name, config.Port)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("%s : %w", fn, err)
	}
	log.Info("connection opened")

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s : %w", fn, err)
	}
	log.Info("db pinged")

	return &Storage{
		db:  db,
		log: log,
	}, nil
}

// RegisterUser inserts new User into database respectfully to UNIQUE constraints on Email and Nickname.
// Returns new User's ID and possible error.
func (s *Storage) RegisterUser(ctx context.Context, email string, nickname string, passHash []byte) (ID int64, err error) {
	const fn = "storage.postgres.RegisterUser"
	log := s.log.With(
		slog.String("fn", fn),
	)
	log.Info("starting registration")

	q, err := s.db.PrepareContext(ctx, queryRegisterUser)
	if err != nil {
		return 0, fmt.Errorf("%s : %w", fn, err)
	}
	log.Info("query prepared")

	id := int64(uuid.New().ID())
	_, err = q.ExecContext(ctx, id, email, nickname, passHash)
	if err != nil {
		return 0, fmt.Errorf("%s : %w", fn, err)
	}
	log.Info("query executed")

	return id, nil
}

// ProvideUser returns User's info based on loginStr (Check comment below).
func (s *Storage) ProvideUser(ctx context.Context, loginStr string) (user models.User, err error) {
	const fn = "storage.postgres.ProvideUser"
	log := s.log.With(
		slog.String("fn", fn),
	)
	log.Info("starting user providing")

	// Handle what User specified in login field (Email/Nickname).
	// Decision based on if there's a '@' symbol as it's restricted for Nicknames to contain '@'.
	var (
		q *sql.Stmt
	)

	byEmail := provideByEmail(loginStr)
	if byEmail {
		q, err = s.db.PrepareContext(ctx, queryProvideWithEmail)
	} else {
		q, err = s.db.PrepareContext(ctx, queryProvideWithNickname)
	}
	if err != nil {
		return models.User{}, fmt.Errorf("%s : %w", fn, err)
	}
	log.Info("query prepared; provideByEmail: %t", byEmail)

	var u models.User
	row := q.QueryRowContext(ctx, loginStr)
	err = row.Scan(&u.ID, &u.Email, &u.Nickname, &u.PasswordHashed)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s : %w", fn, storage.ErrUserNotFound)
		}
		return models.User{}, fmt.Errorf("%s : %w", fn, err)
	}
	log.Info("query executed")

	return u, nil
}

func provideByEmail(s string) bool {
	return strings.Contains(s, "@")
}
