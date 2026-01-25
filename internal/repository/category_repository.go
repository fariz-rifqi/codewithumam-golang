package repository

import (
	"context"
	"pos-api/internal/domain"
)

type CategoryRepository interface {
	Create(ctx context.Context, c domain.Category) (domain.Category, error)
	GetByID(ctx context.Context, id int) (domain.Category, error)
	List(ctx context.Context, p ListParams) ([]domain.Category, error)
	Update(ctx context.Context, id int, c domain.Category) (domain.Category, error)
	Delete(ctx context.Context, id int) error
}
