package blogs

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type Handler struct {
	svc *service
}

func NewHandler(svc *service) *Handler {
	return &Handler{svc: svc}
}

type createBlogDto struct {
	Name        string `json:"name" validate:"required,max=255"`
	Description string `json:"description" validate:"max=1000"`
	WebsiteURL  string `json:"websiteUrl" validate:"required,url,max=255"`
}

type updateBlogDto struct {
	Name        string `json:"name" validate:"required,max=255"`
	Description string `json:"description" validate:"max=1000"`
	WebsiteURL  string `json:"websiteUrl" validate:"required,url,max=255"`
}

type GetAllBlogsInputModel struct {
	SearchNameTerm string `json:"searchNameTerm"`
	SortBy         string `json:"sortBy"`
	SortDirection  string `json:"sortDirection"`
	PageNumber     int    `json:"pageNumber"`
	PageSize       int    `json:"pageSize"`
}

func validateIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if err := validate.Var(id, "required,uuid"); err != nil {
			http.Error(w, "Invalid ID format", http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) Routes(r chi.Router) {
	r.Get("/", h.getAll)
	r.Post("/", h.create)

	r.Route("/{id}", func(r chi.Router) {
		r.Use(validateIDMiddleware)
		r.Get("/", h.getByID)
		r.Put("/", h.update)
		r.Delete("/", h.delete)
	})
}

func (h *Handler) getAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	blogs, err := h.svc.GetAll(ctx)
	if err != nil {
		log.Printf("unexpected error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(blogs); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
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
		log.Printf("unexpected error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(blog); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var dto createBlogDto
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := validate.Struct(dto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	blog, err := h.svc.Create(ctx, Blog{
		Name:        dto.Name,
		Description: dto.Description,
		WebsiteURL:  dto.WebsiteURL,
	})
	if err != nil {
		log.Printf("unexpected error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(blog); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	var dto updateBlogDto

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := validate.Struct(dto); err != nil {
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
		log.Printf("unexpected error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
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
		log.Printf("unexpected error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
