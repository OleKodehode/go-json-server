package app

import (
	"encoding/json"
	"net/http"

	"github.com/OleKodehode/go-json-server/internal/service"
)

// These endpoints needs to communicate with the service layer - Needs a pointer to it
type Handler struct {
	Service *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{Service: s}
}

// GET /:name (collection)
func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	collection := r.PathValue("name")

	filters := map[string]string{}
	for key, values := range r.URL.Query() {
		if len(values) > 0 {
			filters[key] = values[0]
		}
	}
	items := h.Service.GetAll(collection, filters)
	RespondJSON(w, http.StatusOK, items)
}

// get /:name/:id (collection/entry)
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	collection := r.PathValue("name")
	id := r.PathValue("id")
	item := h.Service.GetByID(collection, id)
	if item == nil {
		RespondError(w, http.StatusNotFound, "Entry not found")
		return
	}

	RespondJSON(w, http.StatusOK, item)
}

// POST /:name (collection)
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	collection := r.PathValue("name")

	body := map[string]any{}
	json.NewDecoder(r.Body).Decode(&body)

	item, err := h.Service.Create(collection, body)
	if err != nil {
		RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	RespondJSON(w, http.StatusCreated, item)
}

// PUT /:name/:id (collection/entry)
func (h *Handler) Replace(w http.ResponseWriter, r *http.Request) {
	collection := r.PathValue("name")
	id := r.PathValue("id")

	body := map[string]any{}
	json.NewDecoder(r.Body).Decode(&body)

	item, err := h.Service.Replace(collection, id, body)
	if err != nil {
		RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, item)
}

// PATCH /:name/:id (collection/entry)
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	collection := r.PathValue("name")
	id := r.PathValue("id")

	body := map[string]any{}
	json.NewDecoder(r.Body).Decode(&body)

	item, err := h.Service.Update(collection, id, body)
	if err != nil {
		RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, item)
}

// DELETE /:name/:id (collection/entry)
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	collection := r.PathValue("name")
	id := r.PathValue("id")

	err := h.Service.Delete(collection, id)
	if err != nil {
		RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	RespondJSON(w, http.StatusNoContent, nil)
}