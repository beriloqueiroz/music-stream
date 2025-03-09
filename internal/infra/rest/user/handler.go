package rest_server_user

import (
	"encoding/json"
	"net/http"

	"github.com/beriloqueiroz/music-stream/internal/application"
)

type UserHandler struct {
	service *application.UserService
}

func NewUserHandler(service *application.UserService) *UserHandler {
	return &UserHandler{service: service}
}

type registerRequest struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	InviteCode string `json:"invite_code"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token   string `json:"token"`
	ID      string `json:"id"`
	Email   string `json:"email"`
	IsAdmin bool   `json:"is_admin"`
}

type createInviteRequest struct {
	Email string `json:"email"`
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.service.Register(r.Context(), req.Email, req.Password, req.InviteCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, user, err := h.service.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(loginResponse{Token: *token, ID: user.ID.Hex(), Email: user.Email, IsAdmin: user.IsAdmin})
}

func (h *UserHandler) CreateInvite(w http.ResponseWriter, r *http.Request) {
	// TODO: Implementar middleware de autenticação
	var req createInviteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	invite, err := h.service.CreateInvite(r.Context(), req.Email, r.Header.Get("user_id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(invite)
}
