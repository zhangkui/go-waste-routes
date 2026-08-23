package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go-waste-routes/internal/domain"
	"go-waste-routes/internal/service"
	"go-waste-routes/internal/transport/contextx"
	"go-waste-routes/internal/transport/http/response"
)

type CRUDService[T any] interface {
	Create(context.Context, T) (T, error)
	Get(context.Context, int64) (T, error)
	List(context.Context, int, int) ([]T, int64, error)
	Update(context.Context, int64, T) (T, error)
	Delete(context.Context, int64) error
}

type CRUDHandler[T any] struct {
	Service CRUDService[T]
}

type RouteStopPresentation struct {
	Preview   []domain.RouteStop `json:"preview"`
	Execution []domain.TaskStop  `json:"execution"`
}

type RoutePresentationHandler struct {
	Routes *service.RouteService
	Tasks  *service.TaskService
}

func NewRoutePresentationHandler(routes *service.RouteService, tasks *service.TaskService) *RoutePresentationHandler {
	return &RoutePresentationHandler{Routes: routes, Tasks: tasks}
}

func (h *RoutePresentationHandler) Build(stops []domain.RouteStop) RouteStopPresentation {
	return RouteStopPresentation{
		Preview:   h.Routes.PreviewStops(stops),
		Execution: h.Tasks.BuildStopsFromRoute(stops),
	}
}

func NewCRUDHandler[T any](service CRUDService[T]) *CRUDHandler[T] {
	return &CRUDHandler[T]{Service: service}
}

func (h *CRUDHandler[T]) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page <= 0 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize <= 0 {
		pageSize = 20
	}
	items, total, err := h.Service.List(r.Context(), page, pageSize)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, response.CodeSystemError, err.Error(), contextx.RequestID(r.Context()))
		return
	}
	response.OK(w, map[string]any{"items": items, "pagination": map[string]any{"page": page, "page_size": pageSize, "total": total, "total_pages": int((total + int64(pageSize) - 1) / int64(pageSize))}}, contextx.RequestID(r.Context()))
}

func (h *CRUDHandler[T]) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, response.CodeValidation, "invalid id", contextx.RequestID(r.Context()))
		return
	}
	item, err := h.Service.Get(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusNotFound, response.CodeNotFound, "not found", contextx.RequestID(r.Context()))
		return
	}
	response.OK(w, item, contextx.RequestID(r.Context()))
}

func (h *CRUDHandler[T]) Create(w http.ResponseWriter, r *http.Request) {
	var payload T
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		response.Error(w, http.StatusBadRequest, response.CodeValidation, "invalid json", contextx.RequestID(r.Context()))
		return
	}
	item, err := h.Service.Create(r.Context(), payload)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, response.CodeSystemError, err.Error(), contextx.RequestID(r.Context()))
		return
	}
	response.Created(w, item, contextx.RequestID(r.Context()))
}

func (h *CRUDHandler[T]) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, response.CodeValidation, "invalid id", contextx.RequestID(r.Context()))
		return
	}
	var payload T
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		response.Error(w, http.StatusBadRequest, response.CodeValidation, "invalid json", contextx.RequestID(r.Context()))
		return
	}
	item, err := h.Service.Update(r.Context(), id, payload)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, response.CodeSystemError, err.Error(), contextx.RequestID(r.Context()))
		return
	}
	response.OK(w, item, contextx.RequestID(r.Context()))
}

func (h *CRUDHandler[T]) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, response.CodeValidation, "invalid id", contextx.RequestID(r.Context()))
		return
	}
	if err := h.Service.Delete(r.Context(), id); err != nil {
		response.Error(w, http.StatusInternalServerError, response.CodeSystemError, err.Error(), contextx.RequestID(r.Context()))
		return
	}
	response.OK(w, map[string]any{"deleted": true}, contextx.RequestID(r.Context()))
}
