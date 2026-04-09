package user

import (
	"encoding/json"
	"net/http"

	http_dto "github.com/AlexandreJSimon/hexagonal-golang-project/internal/infra/http/dto"
	"github.com/AlexandreJSimon/hexagonal-golang-project/internal/infra/http/response"
)

func (h *Handler) Login(tokenService TokenService) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var login http_dto.LoginUserRequest

		if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
			response.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		user, err := h.userService.LoginUser(r.Context(), login.Email, login.Password)
		if err != nil {
			response.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		token, err := tokenService.Generate(user.ID)
		if err != nil {
			response.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		response.Success(w, "Login successful", token)
	}
}
