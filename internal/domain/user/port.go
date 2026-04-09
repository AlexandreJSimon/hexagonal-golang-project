package user

type UserRepository interface {
	Save(user *User) error
	GetByID(id string) (*User, error)
	GetByEmail(email string) (*User, error)
	Update(user *User) error
	Delete(id string) error
	List(limit, offset int) ([]*User, error)
	Count() (int, error)
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(password, hash string) bool
}
