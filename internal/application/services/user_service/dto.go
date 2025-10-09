package user_service

type CreateUserInput struct {
	Name     string
	Username string
	Email    string
	Password string
}

type UpdateUserInput struct {
	Name     string
	Username string
	Email    string
	Password string
}
