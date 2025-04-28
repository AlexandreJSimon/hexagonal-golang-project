package user

import (
	"errors"
	"slices"

	"github.com/AlexandreJSimon/hexagonal-golang-project/internal/domain/user"
)

type MemoryRepository struct{}

var users []*user.User

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{}
}

func (mr *MemoryRepository) Save(user *user.User) error {
	users = append(users, user)
	return nil
}

func (mr *MemoryRepository) GetByID(id string) (*user.User, error) {
	for _, user := range users {
		if user.ID == id {
			return user, nil
		}
	}

	return nil, errors.New("user not found")
}

func (mr *MemoryRepository) Update(user *user.User) error {
	for i, u := range users {
		if u.ID == user.ID {
			users[i] = user
		}
	}

	return nil
}

func (mr *MemoryRepository) Delete(id string) error {
	for i, user := range users {
		if user.ID == id {
			users = slices.Delete(users, i, i+1)
			return nil
		}
	}

	return errors.New("user not found")
}

func (mr *MemoryRepository) List(limit, offset int) ([]*user.User, error) {
	return users, nil
}

func (mr *MemoryRepository) Count() (int, error) {
	return len(users), nil
}
