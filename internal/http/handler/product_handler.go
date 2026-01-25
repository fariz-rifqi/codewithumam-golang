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

type ProductHandler struct {
	svc *service.ProductService
}

func NewProductHandler(s *service.ProductService) *ProductHandler {
	return &ProductHandler{svc: s}
}

func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
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

func (h *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
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

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var in domain.Product
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

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		responder.Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	var in domain.Product
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

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		responder.Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.svc.Delete(r.Context(), id); err != nil {
		responder.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	responder.Success(w, map[string]any{"deleted": true})
}
