package handler

import (
	"net/http"
	"strconv"
	"strings"

	"pos-api/internal/domain"
	"pos-api/internal/http/httputil"
	"pos-api/internal/http/responder"
	"pos-api/internal/service"
)

type CategoryHandler struct {
	svc *service.CategoryService
}

func NewCategoryHandler(s *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{svc: s}
}

func (h *CategoryHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	limit := httputil.QueryInt(r, "limit", 50)
	offset := httputil.QueryInt(r, "offset", 0)

	items, err := h.svc.List(r.Context(), limit, offset)
	if err != nil {
		responder.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	responder.Success(w, map[string]any{
		"items":  items,
		"limit":  limit,
		"offset": offset,
	})
}

func (h *CategoryHandler) GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		responder.Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	p, err := h.svc.Get(r.Context(), id)
	if err != nil {
		responder.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	responder.Success(w, p)
}

func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var in domain.Category
	if err := httputil.DecodeJSON(w, r, &in); err != nil {
		responder.Error(w, http.StatusBadRequest, "invalid json: "+err.Error())
		return
	}

	created, err := h.svc.Create(r.Context(), in)
	if err != nil {
		responder.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	responder.Success(w, created)
}

func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		responder.Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	var in domain.Category
	if err := httputil.DecodeJSON(w, r, &in); err != nil {
		responder.Error(w, http.StatusBadRequest, "invalid json: "+err.Error())
		return
	}

	updated, err := h.svc.Update(r.Context(), id, in)
	if err != nil {
		responder.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	responder.Success(w, updated)
}

func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		responder.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.svc.Delete(r.Context(), id); err != nil {
		responder.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	responder.Success(w, map[string]any{"deleted": true})
}
