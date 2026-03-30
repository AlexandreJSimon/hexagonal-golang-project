package user

import (
	"context"

	"github.com/AlexandreJSimon/hexagonal-golang-project/internal/domain/user"
)

type UserServiceProvider interface {
	CreateUser(context.Context, CreateUserInput) (string, error)
	ListUsers(context.Context, int, int) ([]*user.User, error)
}

type UserServiceInput struct {
	UserRepository user.UserRepository
}

type UserService struct {
	userRepository user.UserRepository
}

func NewUserService(userServiceInput UserServiceInput) *UserService {
	return &UserService{
		userRepository: userServiceInput.UserRepository,
	}
}

func (s *UserService) CreateUser(ctx context.Context, input CreateUserInput) (string, error) {
	user := user.NewUser(
		input.Name,
		input.Username,
		input.Email,
		input.Password,
	)

	if err := s.userRepository.Save(user); err != nil {
		return "", err
	}

	return user.ID, nil
}

func (s *UserService) ListUsers(ctx context.Context, limit, offset int) ([]*user.User, error) {
	users, err := s.userRepository.List(limit, offset)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (*user.User, error) {
	user, err := s.userRepository.GetByID(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id string, input UpdateUserInput) error {

	user, err := s.userRepository.GetByID(id)
	if err != nil {
		return err
	}

	//TODO: Validate input fields before updating and improve the way of assigning values ​​to the user
	user.Name = input.Name
	user.Username = input.Username
	user.Email = input.Email
	user.Password = input.Password

	if err := s.userRepository.Update(user); err != nil {
		return err
	}

	return nil
}

func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	if err := s.userRepository.Delete(id); err != nil {
		return err
	}

	return nil
}

func (s *UserService) CountUsers(ctx context.Context) (int, error) {
	count, err := s.userRepository.Count()
	if err != nil {
		return 0, err
	}

	return count, nil
}
