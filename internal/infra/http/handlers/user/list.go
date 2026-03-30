package user

import (
	"net/http"
	"strconv"

	http_dto "github.com/AlexandreJSimon/hexagonal-golang-project/internal/infra/http/dto"
	"github.com/AlexandreJSimon/hexagonal-golang-project/internal/infra/http/response"
)

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
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
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

	var listUsers []http_dto.UserResponse

	for _, user := range users {
		listUsers = append(listUsers, http_dto.UserResponse{
			ID:       user.ID,
			Name:     user.Name,
			Username: user.Username,
			Email:    user.Email,
		})
	}

	response.Success(w, "Users retrieved successfully", listUsers)
}
