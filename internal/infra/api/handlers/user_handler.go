package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/AlexandreJSimon/hexagonal-golang-project/internal/application/services/user_service"
	api_dto "github.com/AlexandreJSimon/hexagonal-golang-project/internal/infra/api/dto"
	"github.com/AlexandreJSimon/hexagonal-golang-project/internal/infra/api/response"
)

// CreateUser Creates a new user
// @Summary Create a new user
// @Description Creates a new user with the provided information
// @Tags users
// @Accept json
// @Produce json
// @Param user body api_dto.CreateUserRequest true "User to create"
// @Success 201 {object} api_dto.UserResponse
// @Failure 400 {object} map[string]string
// @Router /users [post]
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user api_dto.CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		response.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	userId, err := h.userService.CreateUser(r.Context(), user_service.CreateUserInput{
		Name:     user.Name,
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
	})
	if err != nil {
		response.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	response.Success(w, "User created successfully", api_dto.UserResponse{
		ID:       userId,
		Name:     user.Name,
		Username: user.Username,
		Email:    user.Email,
	})
}

// ListUsers Lists users with pagination
// @Summary List users with pagination
// @Description Retrieves a list of users with pagination support
// @Tags users
// @Accept json
// @Produce json
// @Param limit query int true "Number of users to return"
// @Param offset query int true "Number of users to skip"
// @Success 201 {object} []api_dto.UserResponse
// @Failure 400 {object} map[string]string
// @Router /users [get]
func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	var limit, offset int

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		response.Error(w, "Invalid offset parameter", http.StatusBadRequest)
		return
	}

	limit, err = strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		response.Error(w, "Invalid limit parameter", http.StatusBadRequest)
		return
	}

	users, err := h.userService.ListUsers(r.Context(), limit, offset)
	if err != nil {
		response.Error(w, "Failed to list users", http.StatusInternalServerError)
		return
	}

	var listUsers []api_dto.UserResponse

	for _, user := range users {
		listUsers = append(listUsers, api_dto.UserResponse{
			ID:       user.ID,
			Name:     user.Name,
			Username: user.Username,
			Email:    user.Email,
		})
	}

	response.Success(w, "Users retrieved successfully", listUsers)
}
