package blogs

import (
	"context"
	"errors"
	"time"
)

type Blog struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	WebsiteURL   string    `json:"websiteUrl"`
	CreatedAt    time.Time `json:"createdAt"`
	IsMembership bool      `json:"isMembership"`
}

var ErrNotFound = errors.New("blog not found")

type Repository interface {
	Create(ctx context.Context, blog Blog) (Blog, error)
	GetByID(ctx context.Context, id string) (Blog, error)
	GetAll(ctx context.Context) ([]Blog, error)
	Update(ctx context.Context, id string, blog Blog) error
	Delete(ctx context.Context, id string) error
}
