package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Config struct {
	SecretKey string
	ExpiresIn time.Duration
}

type JWT struct {
	secret []byte
	exp    time.Duration
}

func New(cfg Config) (*JWT, error) {
	if cfg.SecretKey == "" {
		return nil, errors.New("secret key is required")
	}
	if cfg.ExpiresIn == 0 {
		cfg.ExpiresIn = 24 * time.Hour
	}

	return &JWT{
		secret: []byte(cfg.SecretKey),
		exp:    cfg.ExpiresIn,
	}, nil
}

func (j *JWT) Generate(sub string) (string, error) {
	claims := jwt.MapClaims{
		"sub": sub,
		"exp": time.Now().Add(j.exp).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

func (j *JWT) Validate(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return j.secret, nil
	}, jwt.WithExpirationRequired())

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		sub, ok := claims["sub"].(string)
		if !ok {
			return "", errors.New("invalid sub claim")
		}
		return sub, nil
	}

	return "", errors.New("invalid token")
}
