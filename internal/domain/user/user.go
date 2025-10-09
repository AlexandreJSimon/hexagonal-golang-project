package domain

import (
	"time"

	"github.com/google/uuid"
)

type Status string

type Role string

const (
	Active   Status = "active"
	Inactive Status = "inactive"
	Admin    Role   = "admin"
	Editor   Role   = "editor"
	Viewer   Role   = "viewer"
)

type User struct {
	ID        string
	Name      string
	Username  string
	Email     string
	Password  string
	Role      Role
	Status    Status
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(name, username, email, password string) *User {
	return &User{
		ID:        uuid.New().String(),
		Name:      name,
		Username:  username,
		Email:     email,
		Password:  password,
		Role:      Viewer,
		Status:    Active,
		CreatedAt: time.Now(),
	}
}

func (u *User) SetUserEditor() {
	u.Role = Editor
}

func (u *User) SetUserAdmin() {
	u.Role = Admin
}
