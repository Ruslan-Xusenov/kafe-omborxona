package handler

import (
	"net/http"

	"kafe-omborxona/internal/domain"
	"kafe-omborxona/internal/middleware"
	"kafe-omborxona/internal/service"
)

type AuthHandler struct {
	authSvc *service.AuthService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{authSvc: svc}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req domain.LoginRequest
	if err := Decode(r, &req); err != nil {
		Error(w, http.StatusBadRequest, "noto'g'ri so'rov formati")
		return
	}
	if req.Username == "" || req.Password == "" {
		Error(w, http.StatusBadRequest, "login va parol kiritilishi shart")
		return
	}

	resp, err := h.authSvc.Login(req)
	if err != nil {
		Error(w, http.StatusUnauthorized, err.Error())
		return
	}
	JSON(w, http.StatusOK, resp)
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	user, err := h.authSvc.GetUserByID(userID)
	if err != nil {
		Error(w, http.StatusNotFound, "foydalanuvchi topilmadi")
		return
	}
	JSON(w, http.StatusOK, user)
}

func (h *AuthHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.authSvc.GetAllUsers()
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, http.StatusOK, users)
}

func (h *AuthHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateUserRequest
	if err := Decode(r, &req); err != nil {
		Error(w, http.StatusBadRequest, "noto'g'ri format")
		return
	}
	if req.Username == "" || req.Password == "" || req.FullName == "" {
		Error(w, http.StatusBadRequest, "barcha maydonlar to'ldirilishi shart")
		return
	}
	if req.Role == "" {
		req.Role = domain.RoleWarehouseManager
	}

	user, err := h.authSvc.CreateUser(req)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, http.StatusCreated, user)
}

func (h *AuthHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := ParseID(r)
	if err != nil {
		Error(w, http.StatusBadRequest, "noto'g'ri ID")
		return
	}
	var req domain.UpdateUserRequest
	if err := Decode(r, &req); err != nil {
		Error(w, http.StatusBadRequest, "noto'g'ri format")
		return
	}
	if err := h.authSvc.UpdateUser(id, req); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, http.StatusOK, map[string]string{"message": "yangilandi"})
}

func (h *AuthHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := ParseID(r)
	if err != nil {
		Error(w, http.StatusBadRequest, "noto'g'ri ID")
		return
	}
	if err := h.authSvc.DeleteUser(id); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, http.StatusOK, map[string]string{"message": "o'chirildi"})
}
