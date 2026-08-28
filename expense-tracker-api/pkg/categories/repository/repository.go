package repository

import (
	"context"
	"expense-tracker-api/internal/db/dbutil"
	"expense-tracker-api/internal/db/generated"
	"expense-tracker-api/pkg/categories/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryRepository interface {
	ListByType(ctx context.Context, categoryType string) ([]model.ExpenseCategory, error)
	GetByName(ctx context.Context, name string) (*model.ExpenseCategory, error)
}

type categoryRepository struct {
	queries *generated.Queries
}

func NewCategoryRepository(pool *pgxpool.Pool) *categoryRepository {
	return &categoryRepository{queries: generated.New(pool)}
}

func (r *categoryRepository) ListByType(ctx context.Context, categoryType string) ([]model.ExpenseCategory, error) {
	rows, err := r.queries.ListCategoriesByType(ctx, categoryType)
	if err != nil {
		return nil, err
	}

	categories := make([]model.ExpenseCategory, len(rows))
	for i, row := range rows {
		categories[i] = model.ExpenseCategory{
			ID:        dbutil.UUIDToString(row.ID),
			Type:      row.Type,
			Name:      row.Name,
			CreatedAt: dbutil.TimestampToString(row.CreatedAt),
			UpdatedAt: dbutil.TimestampToString(row.UpdatedAt),
		}
	}
	return categories, nil
}

func (r *categoryRepository) GetByName(ctx context.Context, name string) (*model.ExpenseCategory, error) {
	row, err := r.queries.GetCategoryByName(ctx, name)
	if err != nil {
		return nil, err
	}

	return &model.ExpenseCategory{
		ID:        dbutil.UUIDToString(row.ID),
		Type:      row.Type,
		Name:      row.Name,
		CreatedAt: dbutil.TimestampToString(row.CreatedAt),
		UpdatedAt: dbutil.TimestampToString(row.UpdatedAt),
	}, nil
}
