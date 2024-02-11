package auth_service

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/humanbelnik/backy-auth-ms/internal/domain/models"
	"github.com/humanbelnik/backy-auth-ms/internal/lib/jwt"
	"github.com/humanbelnik/backy-auth-ms/internal/storage"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/slog"
)

type Auth struct {
	log             *slog.Logger
	userRegistrator UserRegistrator
	userProvider    UserProvider
	tokenTTL        time.Duration
}

type UserRegistrator interface {
	RegisterUser(ctx context.Context, email string, nickname string, passHash []byte) (ID int64, err error)
}

type UserProvider interface {
	ProvideUser(ctx context.Context, loginStr string) (user models.User, err error)
}

// New returns new instance of Auth service.
func New(log *slog.Logger, userRegistrator UserRegistrator, userProvider UserProvider, tokenTTL time.Duration) *Auth {
	return &Auth{
		userRegistrator: userRegistrator,
		userProvider:    userProvider,
		log:             log,
		tokenTTL:        tokenTTL,
	}
}

// Login calls for Data layer, gets User's info and than validate it's password.
// Returns JWT session token and possible error.
func (a *Auth) Login(ctx context.Context, loginString string, password string) (string, error) {
	const fn = "auth_service.Login"
	log := a.log.With(
		slog.String("fn", fn),
	)
	log.Info("start logging nuw user")

	// Validate Email or Nickname.
	user, err := a.userProvider.ProvideUser(ctx, loginString)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("invalid login or email", err)

			return "", fmt.Errorf("%s : %w", fn, err)
		}
	}
	log.Info("user found")

	// Validate password
	if err := bcrypt.CompareHashAndPassword(user.PasswordHashed, []byte(password)); err != nil {
		a.log.Warn("invalid password", err)

		return "", fmt.Errorf("%s : %w", fn, err)
	}
	log.Info("password validated")

	// Create jwt token
	jwtInstance := jwt.New()
	token, err := jwtInstance.NewToken(user, a.tokenTTL)
	if err != nil {
		a.log.Warn("failed to generate jwt token", err)

		return "", fmt.Errorf("%s : %w", fn, err)
	}
	log.Info("login finished")

	return token, nil
}

var (
	ErrInvalidEmail    = errors.New("invalid email")
	ErrInvalidNickname = errors.New("invalid nickname")
)

// Register validates credentials, encrypt password and calls Data layer to register new User.
// Returns new User's ID and possib;e error.
func (a *Auth) Register(ctx context.Context, email string, nickname string, password string) (int64, error) {
	const fn = "auth_service.Register"
	log := a.log.With(
		slog.String("fn", fn),
	)
	log.Info("start user registration")

	if !isValidEmail(email) {
		err := ErrInvalidEmail
		log.Error("invalid email", err)
		return 0, fmt.Errorf("%s : %w", fn, err)
	}
	log.Info("email validated")

	if !isValidNickname(nickname) {
		err := ErrInvalidNickname
		log.Error("invalid nickname", err)
		return 0, fmt.Errorf("%s : %w", fn, err)
	}
	log.Info("nickname validated")

	// Hash password.
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate hash from password", err)
		return 0, fmt.Errorf("%s : %w", fn, err)
	}
	log.Info("password hashed")

	// Save user.
	id, err := a.userRegistrator.RegisterUser(ctx, email, nickname, passwordHash)
	if err != nil {
		log.Error("failed to register user", err)
		return 0, fmt.Errorf("%s : %w", fn, err)
	}
	log.Info("user registration finished")

	return id, nil
}

func (a *Auth) Logout(ctx context.Context, token string) (ok bool, err error) {
	panic("implement me!")
}

func (a *Auth) Unregister(ctx context.Context, unregisterString string, password string, passowordConfirmed string) (userID int64, err error) {
	panic("implement me!")
}

func (a *Auth) IsAdmin(ctx context.Context, userID int64) (isAdmin bool, err error) {
	panic("implement me!")
}

////////////////////////////////////////////////////////////////////////////////////////////////

var (
	emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
)

func isValidEmail(e string) bool {
	re := regexp.MustCompile(emailRegex)

	return re.MatchString(e)
}

func isValidNickname(nn string) bool {
	return !strings.Contains(nn, "@")
}
