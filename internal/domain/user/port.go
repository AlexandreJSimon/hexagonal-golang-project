package user

type UserRepository interface {
	Save(user *User) error
	GetByID(id string) (*User, error)
	Update(user *User) error
	Delete(id string) error
	List(limit, offset int) ([]*User, error)
	Count() (int, error)
}
