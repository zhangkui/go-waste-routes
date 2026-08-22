package handler

import (
	"encoding/json"
	"net/http"

	"go-waste-routes/internal/domain"
	"go-waste-routes/internal/transport/contextx"
	"go-waste-routes/internal/transport/http/response"
	"go-waste-routes/internal/service"
)

type AuthHandler struct {
	Auth *service.AuthService
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, response.CodeValidation, "invalid json", contextx.RequestID(r.Context()))
		return
	}
	result, err := h.Auth.Login(r.Context(), req.Username, req.Password, service.RemoteIP(r.RemoteAddr))
	if err != nil {
		response.Error(w, http.StatusUnauthorized, response.CodeUnauthorized, err.Error(), contextx.RequestID(r.Context()))
		return
	}
	response.OK(w, result, contextx.RequestID(r.Context()))
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username    string `json:"username"`
		DisplayName string `json:"display_name"`
		Phone       string `json:"phone"`
		Email       string `json:"email"`
		Password    string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, response.CodeValidation, "invalid json", contextx.RequestID(r.Context()))
		return
	}
	user, err := h.Auth.Register(r.Context(), domain.User{Username: req.Username, DisplayName: req.DisplayName, Phone: req.Phone, Email: req.Email}, req.Password)
	if err != nil {
		response.Error(w, http.StatusBadRequest, response.CodeBusiness, err.Error(), contextx.RequestID(r.Context()))
		return
	}
	response.Created(w, user, contextx.RequestID(r.Context()))
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req refreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, response.CodeValidation, "invalid json", contextx.RequestID(r.Context()))
		return
	}
	result, err := h.Auth.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, response.CodeUnauthorized, err.Error(), contextx.RequestID(r.Context()))
		return
	}
	response.OK(w, result, contextx.RequestID(r.Context()))
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req refreshRequest
	_ = json.NewDecoder(r.Body).Decode(&req)
	h.Auth.Logout(req.RefreshToken)
	response.OK(w, map[string]any{"ok": true}, contextx.RequestID(r.Context()))
}
