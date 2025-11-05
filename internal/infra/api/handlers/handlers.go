package handlers

import "github.com/AlexandreJSimon/hexagonal-golang-project/internal/application/services/user_service"

type HandlerInput struct {
	UserService user_service.UserServiceProvider
}

type Handler struct {
	userService user_service.UserServiceProvider
}

func NewHandler(HandlerInput HandlerInput) *Handler {
	return &Handler{
		userService: HandlerInput.UserService,
	}
}
