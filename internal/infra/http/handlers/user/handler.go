package user

import (
	"context"

	userApp "github.com/AlexandreJSimon/hexagonal-golang-project/internal/application/user"
	userDomain "github.com/AlexandreJSimon/hexagonal-golang-project/internal/domain/user"
)

type UserServiceProvider interface {
	CreateUser(context.Context, userApp.CreateUserInput) (string, error)
	ListUsers(context.Context, int, int) ([]*userDomain.User, error)
	GetUserByID(context.Context, string) (*userDomain.User, error)
	UpdateUser(context.Context, string, userApp.UpdateUserInput) error
	DeleteUser(context.Context, string) error
	CountUsers(context.Context) (int, error)
}

type Handler struct {
	userService UserServiceProvider
}

func NewHandler(UserService UserServiceProvider) *Handler {
	return &Handler{
		userService: UserService,
	}
}
