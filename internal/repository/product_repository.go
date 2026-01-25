package repository

import (
	"context"
	"pos-api/internal/domain"
)

type ProductRepository interface {
	Create(ctx context.Context, p domain.Product) (domain.Product, error)
	GetByID(ctx context.Context, id int) (domain.Product, error)
	List(ctx context.Context, p ListParams) ([]domain.Product, error)
	Update(ctx context.Context, id int, p domain.Product) (domain.Product, error)
	Delete(ctx context.Context, id int) error
}
