package user

import (
	"net/http"

	"github.com/prodanov17/znk/internal/services/auth"
	"github.com/prodanov17/znk/internal/types"
	"github.com/prodanov17/znk/internal/utils"
)

type handler struct {
	service types.UserService
}

func NewHandler(service types.UserService) *handler {
	return &handler{service: service}
}

func (h *handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /user", auth.WithJWTAuth(h.handleGetUser))
	router.HandleFunc("POST /login", h.handleLogin)
	router.HandleFunc("POST /register", h.handleRegister)
}

func (h *handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var user types.RegisterUserPayload

	if err := utils.ParseJSON(r, &user); err != nil {
		utils.WriteError(w, r, http.StatusBadRequest, err)
		return
	}

	token, err := h.service.RegisterUser(&user)

	if err != nil {
		utils.WriteError(w, r, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (h *handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var user types.LoginUserPayload

	if err := utils.ParseJSON(r, &user); err != nil {
		utils.WriteError(w, r, http.StatusBadRequest, err)
		return
	}

	token, err := h.service.LoginUser(&user)
	if err != nil {
		utils.WriteError(w, r, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (h *handler) handleGetUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserKey).(int)

	user, err := h.service.GetUserByID(userID)
	if err != nil {
		utils.WriteError(w, r, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, user)
}
