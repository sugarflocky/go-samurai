package blogs

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

var _ Repository = (*PostgresRepository)(nil)

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) Create(ctx context.Context, blog Blog) (Blog, error) {
	err := r.pool.QueryRow(ctx, `
	INSERT INTO blogs 
	( name, description, website_url ) 
	VALUES ( $1, $2, $3 )
	RETURNING id, created_at, is_membership
	`,
		blog.Name, blog.Description, blog.WebsiteURL).Scan(&blog.ID, &blog.CreatedAt, &blog.IsMembership)
	if err != nil {
		return Blog{}, err
	}

	return blog, nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (Blog, error) {
	var blog Blog
	err := r.pool.QueryRow(ctx,
		`SELECT id, name, description, website_url, created_at, is_membership FROM blogs WHERE id = $1`,
		id).Scan(&blog.ID, &blog.Name, &blog.Description, &blog.WebsiteURL, &blog.CreatedAt, &blog.IsMembership)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Blog{}, ErrNotFound
		}

		return Blog{}, err
	}

	return blog, nil
}

func (r *PostgresRepository) GetAll(ctx context.Context) ([]Blog, error) {
	blogs := make([]Blog, 0)
	rows, err := r.pool.Query(ctx, `SELECT id, name, description, website_url, created_at, is_membership FROM blogs`)
	if err != nil {
		return blogs, err
	}
	defer rows.Close()

	for rows.Next() {
		var blog Blog
		if err := rows.Scan(&blog.ID, &blog.Name, &blog.Description, &blog.WebsiteURL, &blog.CreatedAt, &blog.IsMembership); err != nil {
			return blogs, err
		}
		blogs = append(blogs, blog)
	}

	return blogs, nil
}

func (r *PostgresRepository) Update(ctx context.Context, id string, blog Blog) error {
	tag, err := r.pool.Exec(ctx, `
	UPDATE blogs 
	SET name = $2, description = $3, website_url = $4
	WHERE id = $1`, id, blog.Name, blog.Description, blog.WebsiteURL)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM blogs WHERE id = $1`, id)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}
