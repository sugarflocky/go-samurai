package blogs

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	svc *service
}

func NewHandler(svc *service) *Handler {
	return &Handler{svc: svc}
}

type createBlogDto struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	WebsiteURL  string `json:"websiteUrl"`
}

type updateBlogDto struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	WebsiteURL  string `json:"websiteUrl"`
}

func (h *Handler) Routes(r chi.Router) {
	r.Get("/", h.getAll)
	r.Get("/{id}", h.getByID)
	r.Post("/", h.create)
	r.Put("/{id}", h.update)
	r.Delete("/{id}", h.delete)
}

func (h *Handler) getAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	blogs, err := h.svc.GetAll(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blogs)
}

func (h *Handler) getByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")
	blog, err := h.svc.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			http.Error(w, "Blog not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blog)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var dto createBlogDto
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	blog, err := h.svc.Create(ctx, Blog{
		Name:        dto.Name,
		Description: dto.Description,
		WebsiteURL:  dto.WebsiteURL,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(blog)
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	var dto updateBlogDto

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.svc.Update(ctx, id, Blog{
		Name:        dto.Name,
		Description: dto.Description,
		WebsiteURL:  dto.WebsiteURL,
	})
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			http.Error(w, "Blog not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	err := h.svc.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			http.Error(w, "Blog not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
