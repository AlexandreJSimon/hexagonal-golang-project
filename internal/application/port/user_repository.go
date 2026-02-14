package port

import domain "github.com/AlexandreJSimon/hexagonal-golang-project/internal/domain"

type UserRepository interface {
	Save(user *domain.User) error
	GetByID(id string) (*domain.User, error)
	Update(user *domain.User) error
	Delete(id string) error
	List(limit, offset int) ([]*domain.User, error)
	Count() (int, error)
}
