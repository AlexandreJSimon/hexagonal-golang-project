package env

import (
	"os"
	"strconv"
)

type Env struct {
	JWTSecretKey string
	Port         int
	AllawOrigins string
}

func Load() *Env {
	return &Env{
		JWTSecretKey: MustGetString("JWT_SECRET_KEY"),
		Port:         GetInt("PORT", 8080),
		AllawOrigins: MustGetString("ALLOWED_ORIGINS"),
	}
}

func GetString(key string, def string) string {
	if str, ok := os.LookupEnv(key); ok {
		return str
	}

	return def
}

func MustGetString(key string) string {
	if str, ok := os.LookupEnv(key); ok {
		return str
	}

	panic("missing required environment variable: " + key)
}

func GetInt(key string, def int) int {
	str, ok := os.LookupEnv(key)
	if !ok {
		return def
	}

	i, err := strconv.Atoi(str)
	if err != nil {
		panic("invalid integer environment variable: " + key)
	}

	return i
}

func MustGetInt(key string) int {
	str, ok := os.LookupEnv(key)
	if !ok {
		panic("missing required integer environment variable: " + key)
	}

	i, err := strconv.Atoi(str)
	if err != nil {
		panic("invalid integer environment variable: " + key)
	}

	return i
}
