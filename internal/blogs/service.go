package blogs

import "context"

type service struct {
	repo Repository
}

func NewService(repo Repository) *service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, blog Blog) (Blog, error) {
	return s.repo.Create(ctx, blog)
}

func (s *service) GetByID(ctx context.Context, id string) (Blog, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) GetAll(ctx context.Context) ([]Blog, error) {
	return s.repo.GetAll(ctx)
}

func (s *service) Update(ctx context.Context, id string, blog Blog) error {
	return s.repo.Update(ctx, id, blog)
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
