package repository_postgres

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"pos-api/internal/domain"
	"pos-api/internal/repository"
)

type ProductRepo struct {
	db *sql.DB
}

func NewProductRepo(db *sql.DB) *ProductRepo {
	return &ProductRepo{db: db}
}

func (r *ProductRepo) Create(ctx context.Context, p domain.Product) (domain.Product, error) {
	p.Name = strings.TrimSpace(p.Name)

	var out domain.Product
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO products (name, price, quantity, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id, name, price, quantity, created_at, updated_at
	`, p.Name, p.Price, p.Quantity).Scan(
		&out.ID,
		&out.Name,
		&out.Price,
		&out.Quantity,
		&out.CreatedAt,
		&out.UpdatedAt,
	)
	if err != nil {
		return domain.Product{}, err
	}
	return out, nil
}

func (r *ProductRepo) GetByID(ctx context.Context, id int) (domain.Product, error) {
	var out domain.Product
	err := r.db.QueryRowContext(ctx, `
		SELECT id, name, price, quantity, created_at, updated_at
		FROM products
		WHERE id = $1
	`, id).Scan(
		&out.ID,
		&out.Name,
		&out.Price,
		&out.Quantity,
		&out.CreatedAt,
		&out.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Product{}, errors.New("not found")
		}
		return domain.Product{}, err
	}
	return out, nil
}

func (r *ProductRepo) List(ctx context.Context, lp repository.ListParams) ([]domain.Product, error) {
	limit := lp.Limit
	offset := lp.Offset
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, price, quantity, created_at, updated_at
		FROM products
		ORDER BY id DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]domain.Product, 0)
	for rows.Next() {
		var p domain.Product
		if err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Price,
			&p.Quantity,
			&p.CreatedAt,
			&p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *ProductRepo) Update(ctx context.Context, id int, patch domain.Product) (domain.Product, error) {
	patch.Name = strings.TrimSpace(patch.Name)

	var out domain.Product
	err := r.db.QueryRowContext(ctx, `
		UPDATE products
		SET name = $1, price = $2, quantity = $3, updated_at = NOW()
		WHERE id = $4
		RETURNING id, name, price, quantity, created_at, updated_at
	`, patch.Name, patch.Price, patch.Quantity, id).Scan(
		&out.ID,
		&out.Name,
		&out.Price,
		&out.Quantity,
		&out.CreatedAt,
		&out.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Product{}, errors.New("not found")
		}
		return domain.Product{}, err
	}
	return out, nil
}

func (r *ProductRepo) Delete(ctx context.Context, id int) error {
	res, err := r.db.ExecContext(ctx, `
		DELETE FROM products
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
