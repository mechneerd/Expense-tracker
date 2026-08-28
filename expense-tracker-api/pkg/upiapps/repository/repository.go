package repository

import (
	"context"
	"expense-tracker-api/internal/db/dbutil"
	"expense-tracker-api/internal/db/generated"
	"expense-tracker-api/pkg/upiapps/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UPIAppRepository interface {
	ListAll(ctx context.Context) ([]model.UPIApp, error)
	GetByName(ctx context.Context, name string) (*model.UPIApp, error)
}

type upiAppRepository struct {
	queries *generated.Queries
}

func NewUPIAppRepository(pool *pgxpool.Pool) *upiAppRepository {
	return &upiAppRepository{queries: generated.New(pool)}
}

func (r *upiAppRepository) ListAll(ctx context.Context) ([]model.UPIApp, error) {
	rows, err := r.queries.ListUPIApps(ctx)
	if err != nil {
		return nil, err
	}

	apps := make([]model.UPIApp, len(rows))
	for i, row := range rows {
		apps[i] = model.UPIApp{
			ID:        dbutil.UUIDToString(row.ID),
			Name:      row.Name,
			CreatedAt: dbutil.TimestampToString(row.CreatedAt),
			UpdatedAt: dbutil.TimestampToString(row.UpdatedAt),
		}
	}
	return apps, nil
}

func (r *upiAppRepository) GetByName(ctx context.Context, name string) (*model.UPIApp, error) {
	row, err := r.queries.GetUPIAppByName(ctx, name)
	if err != nil {
		return nil, err
	}

	return &model.UPIApp{
		ID:        dbutil.UUIDToString(row.ID),
		Name:      row.Name,
		CreatedAt: dbutil.TimestampToString(row.CreatedAt),
		UpdatedAt: dbutil.TimestampToString(row.UpdatedAt),
	}, nil
}
