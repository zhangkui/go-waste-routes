package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"go-waste-routes/internal/domain"
	"go-waste-routes/internal/transport/contextx"
	"go-waste-routes/internal/transport/http/response"
	"go-waste-routes/internal/service"
)

type UserHandler struct {
	Users  *service.UserService
	Auth   *service.AuthService
}

func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	claims := contextx.ClaimsFromContext(r.Context())
	if claims == nil {
		response.Error(w, http.StatusUnauthorized, response.CodeUnauthorized, "unauthorized", contextx.RequestID(r.Context()))
		return
	}
	user, err := h.Auth.Me(r.Context(), claims.UserID)
	if err != nil {
		response.Error(w, http.StatusNotFound, response.CodeNotFound, err.Error(), contextx.RequestID(r.Context()))
		return
	}
	response.OK(w, user, contextx.RequestID(r.Context()))
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	page, pageSize := parsePaging(r)
	items, total, err := h.Users.List(r.Context(), page, pageSize)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, response.CodeSystemError, err.Error(), contextx.RequestID(r.Context()))
		return
	}
	response.OK(w, map[string]any{"items": items, "pagination": map[string]any{"page": page, "page_size": pageSize, "total": total, "total_pages": int((total+int64(pageSize)-1)/int64(pageSize))}}, contextx.RequestID(r.Context()))
}

func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		response.Error(w, http.StatusBadRequest, response.CodeValidation, "invalid id", contextx.RequestID(r.Context()))
		return
	}
	item, ok, err := h.Users.GetByID(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, response.CodeSystemError, err.Error(), contextx.RequestID(r.Context()))
		return
	}
	if !ok {
		response.Error(w, http.StatusNotFound, response.CodeNotFound, "not found", contextx.RequestID(r.Context()))
		return
	}
	response.OK(w, item, contextx.RequestID(r.Context()))
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
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
	user, err := h.Users.Create(r.Context(), domain.User{Username: req.Username, DisplayName: req.DisplayName, Phone: req.Phone, Email: req.Email}, req.Password, []domain.Permission{})
	if err != nil {
		response.Error(w, http.StatusBadRequest, response.CodeBusiness, err.Error(), contextx.RequestID(r.Context()))
		return
	}
	response.Created(w, user, contextx.RequestID(r.Context()))
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		response.Error(w, http.StatusBadRequest, response.CodeValidation, "invalid id", contextx.RequestID(r.Context()))
		return
	}
	existing, ok, err := h.Users.GetByID(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, response.CodeSystemError, err.Error(), contextx.RequestID(r.Context()))
		return
	}
	if !ok {
		response.Error(w, http.StatusNotFound, response.CodeNotFound, "not found", contextx.RequestID(r.Context()))
		return
	}
	var payload domain.User
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		response.Error(w, http.StatusBadRequest, response.CodeValidation, "invalid json", contextx.RequestID(r.Context()))
		return
	}
	payload.PasswordHash = existing.PasswordHash
	payload.Roles = existing.Roles
	payload.Permissions = existing.Permissions
	updated, err := h.Users.Update(r.Context(), id, payload)
	if err != nil {
		response.Error(w, http.StatusBadRequest, response.CodeBusiness, err.Error(), contextx.RequestID(r.Context()))
		return
	}
	response.OK(w, updated, contextx.RequestID(r.Context()))
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		response.Error(w, http.StatusBadRequest, response.CodeValidation, "invalid id", contextx.RequestID(r.Context()))
		return
	}
	if err := h.Users.Delete(r.Context(), id); err != nil {
		response.Error(w, http.StatusBadRequest, response.CodeBusiness, err.Error(), contextx.RequestID(r.Context()))
		return
	}
	response.OK(w, map[string]any{"deleted": true}, contextx.RequestID(r.Context()))
}

func (h *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	claims := contextx.ClaimsFromContext(r.Context())
	if claims == nil {
		response.Error(w, http.StatusUnauthorized, response.CodeUnauthorized, "unauthorized", contextx.RequestID(r.Context()))
		return
	}
	var req struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, response.CodeValidation, "invalid json", contextx.RequestID(r.Context()))
		return
	}
	if err := h.Users.ChangePassword(r.Context(), claims.UserID, req.OldPassword, req.NewPassword); err != nil {
		response.Error(w, http.StatusBadRequest, response.CodeBusiness, err.Error(), contextx.RequestID(r.Context()))
		return
	}
	response.OK(w, map[string]any{"changed": true}, contextx.RequestID(r.Context()))
}

func (h *UserHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, response.CodeValidation, "invalid id", contextx.RequestID(r.Context()))
		return
	}
	var req struct { NewPassword string `json:"new_password"` }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, response.CodeValidation, "invalid json", contextx.RequestID(r.Context()))
		return
	}
	if err := h.Users.ResetPassword(r.Context(), id, req.NewPassword); err != nil {
		response.Error(w, http.StatusBadRequest, response.CodeBusiness, err.Error(), contextx.RequestID(r.Context()))
		return
	}
	response.OK(w, map[string]any{"reset": true}, contextx.RequestID(r.Context()))
}

func (h *UserHandler) Status(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, response.CodeValidation, "invalid id", contextx.RequestID(r.Context()))
		return
	}
	var req struct { Status string `json:"status"` }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, response.CodeValidation, "invalid json", contextx.RequestID(r.Context()))
		return
	}
	if err := h.Users.UpdateStatus(r.Context(), id, req.Status); err != nil {
		response.Error(w, http.StatusBadRequest, response.CodeBusiness, err.Error(), contextx.RequestID(r.Context()))
		return
	}
	response.OK(w, map[string]any{"status": req.Status}, contextx.RequestID(r.Context()))
}

func parsePaging(r *http.Request) (int, int) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page <= 0 { page = 1 }
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize <= 0 { pageSize = 20 }
	return page, pageSize
}

func parseID(r *http.Request) (int64, error) {
	return strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
}
