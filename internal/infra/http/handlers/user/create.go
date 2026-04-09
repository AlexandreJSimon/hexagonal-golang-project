package user

import (
	"encoding/json"
	"net/http"

	userApp "github.com/AlexandreJSimon/hexagonal-golang-project/internal/application/user"
	http_dto "github.com/AlexandreJSimon/hexagonal-golang-project/internal/infra/http/dto"
	"github.com/AlexandreJSimon/hexagonal-golang-project/internal/infra/http/response"
)

// CreateUser Creates a new user
// @Summary Create a new user
// @Description Creates a new user with the provided information
// @Tags users
// @Accept json
// @Produce json
// @Param user body http_dto.CreateUserRequest true "User to create"
// @Success 201 {object} http_dto.UserResponse
// @Failure 400 {object} map[string]string
// @Router /users [post]
// @Security BearerAuth
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var user http_dto.CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		response.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	userId, err := h.userService.CreateUser(r.Context(), userApp.CreateUserInput{
		Name:     user.Name,
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
	})
	if err != nil {
		response.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	response.Success(w, "User created successfully", http_dto.UserResponse{
		ID:       userId,
		Name:     user.Name,
		Username: user.Username,
		Email:    user.Email,
	})
}
