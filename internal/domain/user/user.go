package user

import (
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	Active   Status = "active"
	Inactive Status = "inactive"
)

type User struct {
	ID        string
	Name      string
	Username  string
	Email     string
	Password  string
	Role      string
	Status    Status
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(name, username, email, password, role string) *User {
	return &User{
		ID:        uuid.New().String(),
		Name:      name,
		Username:  username,
		Email:     email,
		Password:  password,
		Role:      role,
		Status:    Active,
		CreatedAt: time.Now(),
	}
}
