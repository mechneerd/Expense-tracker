package repository

import (
	"context"
	"expense-tracker-api/internal/db/generated"
	"expense-tracker-api/pkg/paymentmethods/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PaymentMethodRepository interface {
	ListAll(ctx context.Context) ([]model.PaymentMethod, error)
	GetByName(ctx context.Context, name string) (*model.PaymentMethod, error)
}

type paymentMethodRepository struct {
	queries *generated.Queries
}

func NewPaymentMethodRepository(pool *pgxpool.Pool) *paymentMethodRepository {
	return &paymentMethodRepository{queries: generated.New(pool)}
}

func (r *paymentMethodRepository) ListAll(ctx context.Context) ([]model.PaymentMethod, error) {
	rows, err := r.queries.ListPaymentMethods(ctx)
	if err != nil {
		return nil, err
	}

	methods := make([]model.PaymentMethod, len(rows))
	for i, row := range rows {
		methods[i] = model.PaymentMethod{
			ID:   row.ID,
			Name: row.Name,
		}
	}
	return methods, nil
}

func (r *paymentMethodRepository) GetByName(ctx context.Context, name string) (*model.PaymentMethod, error) {
	row, err := r.queries.GetPaymentMethodByName(ctx, name)
	if err != nil {
		return nil, err
	}

	return &model.PaymentMethod{
		ID:   row.ID,
		Name: row.Name,
	}, nil
}
