package jwt

import (
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
	"github.com/humanbelnik/backy-auth-ms/internal/domain/models"
	"github.com/ilyakaznacheev/cleanenv"
)

type JWT struct {
	token string `env:"JWT"`
}

// New creates new JWT instance with secret from config file.
func New() *JWT {
	var jwt JWT
	if err := cleanenv.ReadEnv(&jwt); err != nil {
		panic("error while parsing jwt secret from .env:" + err.Error())
	}

	return &jwt
}

// NewToken generates new token based on User info, JWT secret and token's duration.
func (jwt *JWT) NewToken(user models.User, dur time.Duration) (string, error) {
	token := jwtlib.New(jwtlib.SigningMethodHS256)
	claims := token.Claims.(jwtlib.MapClaims)
	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(dur).Unix()
	claims["uid"] = user.ID

	tokenString, err := token.SignedString([]byte(jwt.token))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
