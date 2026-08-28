package repository

import (
	"context"
	"expense-tracker-api/internal/db/dbutil"
	"expense-tracker-api/internal/db/generated"
	"expense-tracker-api/pkg/transactions/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionRepository interface {
	Create(ctx context.Context, t *model.Transaction) error
	GetByID(ctx context.Context, id string) (*model.Transaction, error)
	ListByFamily(ctx context.Context, familyID string) ([]model.Transaction, error)
	ListByUser(ctx context.Context, userID string, familyID *string) ([]model.Transaction, error)
	Delete(ctx context.Context, id string) error
}

type transactionRepository struct {
	queries *generated.Queries
}

func NewTransactionRepository(pool *pgxpool.Pool) *transactionRepository {
	return &transactionRepository{queries: generated.New(pool)}
}

func (r *transactionRepository) Create(ctx context.Context, t *model.Transaction) error {
	return r.queries.CreateTransaction(ctx, generated.CreateTransactionParams{
		FamilyID:        dbutil.ParseUUID(t.FamilyID),
		UserID:          dbutil.ParseUUID(derefString(t.UserID)),
		TransactionType: t.Type,
		CategoryID:      dbutil.ParseUUID(derefString(t.Category)),
		PaymentMethodID: dbutil.StringToText(derefString(t.PaymentMethod)),
		UpiAppID:        dbutil.ParseUUID(derefString(t.UPIApp)),
		Amount:          dbutil.FloatToNumeric(t.Amount),
		TransactionDate: dbutil.ParseDate(t.Date),
		Description:     dbutil.StringToText(t.Description),
	})
}

func (r *transactionRepository) GetByID(ctx context.Context, id string) (*model.Transaction, error) {
	row, err := r.queries.GetTransactionByID(ctx, dbutil.ParseUUID(id))
	if err != nil {
		return nil, err
	}
	return convertTransactionRow(
		dbutil.UUIDToString(row.ID),
		dbutil.UUIDToString(row.FamilyID),
		dbutil.UUIDToString(row.UserID),
		row.TransactionType.(string),
		dbutil.TextToString(row.Category),
		dbutil.TextToString(row.PaymentMethod),
		dbutil.UUIDToString(row.UpiAppID),
		dbutil.NumericToFloat64(row.Amount),
		dbutil.DateToString(row.TransactionDate),
		dbutil.TextToString(row.Description),
		dbutil.TimestampToString(row.CreatedAt),
		dbutil.TimestampToString(row.UpdatedAt),
		dbutil.TimestampToString(row.DeletedAt),
	), nil
}

func (r *transactionRepository) ListByFamily(ctx context.Context, familyID string) ([]model.Transaction, error) {
	rows, err := r.queries.ListTransactionsByFamily(ctx, dbutil.ParseUUID(familyID))
	if err != nil {
		return nil, err
	}

	transactions := make([]model.Transaction, len(rows))
	for i, row := range rows {
		transactions[i] = *convertTransactionRow(
			dbutil.UUIDToString(row.ID),
			dbutil.UUIDToString(row.FamilyID),
			dbutil.UUIDToString(row.UserID),
			row.TransactionType.(string),
			dbutil.TextToString(row.Category),
			dbutil.TextToString(row.PaymentMethod),
			dbutil.UUIDToString(row.UpiAppID),
			dbutil.NumericToFloat64(row.Amount),
			dbutil.DateToString(row.TransactionDate),
			dbutil.TextToString(row.Description),
			dbutil.TimestampToString(row.CreatedAt),
			dbutil.TimestampToString(row.UpdatedAt),
			dbutil.TimestampToString(row.DeletedAt),
		)
	}
	return transactions, nil
}

func (r *transactionRepository) ListByUser(ctx context.Context, userID string, familyID *string) ([]model.Transaction, error) {
	famID := ""
	if familyID != nil {
		famID = *familyID
	}

	rows, err := r.queries.ListTransactionsByUser(ctx, generated.ListTransactionsByUserParams{
		UserID:  dbutil.ParseUUID(userID),
		Column2: dbutil.ParseUUID(famID),
	})
	if err != nil {
		return nil, err
	}

	transactions := make([]model.Transaction, len(rows))
	for i, row := range rows {
		transactions[i] = *convertTransactionRow(
			dbutil.UUIDToString(row.ID),
			dbutil.UUIDToString(row.FamilyID),
			dbutil.UUIDToString(row.UserID),
			row.TransactionType.(string),
			dbutil.TextToString(row.Category),
			dbutil.TextToString(row.PaymentMethod),
			dbutil.UUIDToString(row.UpiAppID),
			dbutil.NumericToFloat64(row.Amount),
			dbutil.DateToString(row.TransactionDate),
			dbutil.TextToString(row.Description),
			dbutil.TimestampToString(row.CreatedAt),
			dbutil.TimestampToString(row.UpdatedAt),
			dbutil.TimestampToString(row.DeletedAt),
		)
	}
	return transactions, nil
}

func (r *transactionRepository) Delete(ctx context.Context, id string) error {
	return r.queries.SoftDeleteTransaction(ctx, dbutil.ParseUUID(id))
}

func convertTransactionRow(id, familyID, userID, txType, category, paymentMethod, upiAppID string, amount float64, date, description, createdAt, updatedAt, deletedAt string) *model.Transaction {
	t := &model.Transaction{
		ID:          id,
		FamilyID:    familyID,
		Type:        txType,
		Category:    &category,
		PaymentMethod: &paymentMethod,
		UPIApp:      &upiAppID,
		Amount:      amount,
		Date:        date,
		Description: description,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}

	if userID != "" {
		t.UserID = &userID
	}
	if deletedAt != "" {
		t.DeletedAt = &deletedAt
	}

	return t
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
