package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Nezent/auth-service/internal/application/dto"
	"github.com/Nezent/auth-service/internal/application/services"
	"github.com/Nezent/auth-service/pkg/response"
)

type UserHandler struct {
	service services.UserService
}

func NewUserHandler(service services.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}
	res, err := h.service.CreateUser(req.Email, req.Password)
	if err != nil {
		response.WriteError(w, err.Error(), err.StatusCode)
		return
	}

	response.WriteSuccess(w, res, http.StatusCreated)
}
