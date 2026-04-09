package user

import (
	"context"
	"errors"

	"github.com/AlexandreJSimon/hexagonal-golang-project/internal/domain/user"
)

type UserServiceInput struct {
	UserRepository user.UserRepository
	PasswordHasher user.PasswordHasher
}

type UserService struct {
	userRepository user.UserRepository
	passwordHasher user.PasswordHasher
}

func NewUserService(userServiceInput UserServiceInput) *UserService {
	return &UserService{
		userRepository: userServiceInput.UserRepository,
		passwordHasher: userServiceInput.PasswordHasher,
	}
}

func (s *UserService) CreateUser(ctx context.Context, input CreateUserInput) (string, error) {
	pwd, err := s.passwordHasher.Hash(input.Password)
	if err != nil {
		return "", err
	}

	user := user.NewUser(
		input.Name,
		input.Username,
		input.Email,
		pwd,
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
	pwd, err := s.passwordHasher.Hash(input.Password)
	if err != nil {
		return err
	}

	user, err := s.userRepository.GetByID(id)
	if err != nil {
		return err
	}

	//TODO: Validate input fields before updating and improve the way of assigning values ​​to the user
	user.Name = input.Name
	user.Username = input.Username
	user.Email = input.Email
	user.Password = pwd

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

func (s *UserService) LoginUser(ctx context.Context, email, password string) (*user.User, error) {
	user, err := s.userRepository.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	if !s.passwordHasher.Compare(password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}
