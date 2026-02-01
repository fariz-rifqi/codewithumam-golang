package repository_postgres

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"pos-api/internal/domain"
	"pos-api/internal/repository"
)

type CategoryRepo struct {
	db *sql.DB
}

func NewCategoryRepo(db *sql.DB) *CategoryRepo {
	return &CategoryRepo{db: db}
}

func (r *CategoryRepo) Create(ctx context.Context, c domain.Category) (domain.Category, error) {
	c.Name = strings.TrimSpace(c.Name)
	c.Description = strings.TrimSpace(c.Description)

	var out domain.Category
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO categories (name, description, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
		RETURNING id, name, description, created_at, updated_at
	`, c.Name, c.Description).Scan(
		&out.ID,
		&out.Name,
		&out.Description,
		&out.CreatedAt,
		&out.UpdatedAt,
	)
	if err != nil {
		return domain.Category{}, err
	}
	return out, nil
}

func (r *CategoryRepo) GetByID(ctx context.Context, id int) (domain.Category, error) {
	var out domain.Category
	err := r.db.QueryRowContext(ctx, `
		SELECT id, name, description, created_at, updated_at
		FROM categories
		WHERE id = $1
	`, id).Scan(
		&out.ID,
		&out.Name,
		&out.Description,
		&out.CreatedAt,
		&out.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Category{}, errors.New("not found")
		}
		return domain.Category{}, err
	}
	return out, nil
}

func (r *CategoryRepo) List(ctx context.Context, lp repository.ListParams) ([]domain.Category, error) {
	limit := lp.Limit
	offset := lp.Offset
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, description, created_at, updated_at
		FROM categories
		ORDER BY id DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]domain.Category, 0)
	for rows.Next() {
		var c domain.Category
		if err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.Description,
			&c.CreatedAt,
			&c.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *CategoryRepo) Update(ctx context.Context, id int, patch domain.Category) (domain.Category, error) {
	patch.Name = strings.TrimSpace(patch.Name)
	patch.Description = strings.TrimSpace(patch.Description)

	var out domain.Category
	err := r.db.QueryRowContext(ctx, `
		UPDATE categories
		SET name = $1, description = $2, updated_at = NOW()
		WHERE id = $3
		RETURNING id, name, description, created_at, updated_at
	`, patch.Name, patch.Description, id).Scan(
		&out.ID,
		&out.Name,
		&out.Description,
		&out.CreatedAt,
		&out.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Category{}, errors.New("not found")
		}
		return domain.Category{}, err
	}
	return out, nil
}

func (r *CategoryRepo) Delete(ctx context.Context, id int) error {
	res, err := r.db.ExecContext(ctx, `
		DELETE FROM categories
		WHERE id = $1
	`, id)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("not found")
	}
	return nil
}
