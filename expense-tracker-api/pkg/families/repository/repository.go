package repository

import (
	"context"
	"expense-tracker-api/internal/db/dbutil"
	"expense-tracker-api/internal/db/generated"
	"expense-tracker-api/pkg/families/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type FamilyRepository interface {
	GetByID(ctx context.Context, id string) (*model.Family, error)
	GetByUniqueCode(ctx context.Context, code string) (*model.Family, error)
	ListMembers(ctx context.Context, familyID string) ([]model.FamilyMember, error)
	Create(ctx context.Context, f *model.Family) error
	AddMember(ctx context.Context, member *model.FamilyMember) error
	UpdateMemberStatus(ctx context.Context, familyID, userID, status string) error
}

type familyRepository struct {
	queries *generated.Queries
}

func NewFamilyRepository(pool *pgxpool.Pool) *familyRepository {
	return &familyRepository{queries: generated.New(pool)}
}

func (r *familyRepository) GetByID(ctx context.Context, id string) (*model.Family, error) {
	row, err := r.queries.GetFamilyByID(ctx, dbutil.ParseUUID(id))
	if err != nil {
		return nil, err
	}

	return &model.Family{
		ID:         dbutil.UUIDToString(row.ID),
		Name:       row.Name,
		UniqueCode: row.UniqueCode,
		CreatedBy:  dbutil.UUIDToString(row.CreatedBy),
		Status:     dbutil.TextToString(row.Status),
		CreatedAt:  dbutil.TimestampToString(row.CreatedAt),
		UpdatedAt:  dbutil.TimestampToString(row.UpdatedAt),
	}, nil
}

func (r *familyRepository) GetByUniqueCode(ctx context.Context, code string) (*model.Family, error) {
	row, err := r.queries.GetFamilyByUniqueCode(ctx, code)
	if err != nil {
		return nil, err
	}

	return &model.Family{
		ID:         dbutil.UUIDToString(row.ID),
		Name:       row.Name,
		UniqueCode: row.UniqueCode,
		CreatedBy:  dbutil.UUIDToString(row.CreatedBy),
		Status:     dbutil.TextToString(row.Status),
		CreatedAt:  dbutil.TimestampToString(row.CreatedAt),
		UpdatedAt:  dbutil.TimestampToString(row.UpdatedAt),
	}, nil
}

func (r *familyRepository) ListMembers(ctx context.Context, familyID string) ([]model.FamilyMember, error) {
	rows, err := r.queries.ListFamilyMembers(ctx, dbutil.ParseUUID(familyID))
	if err != nil {
		return nil, err
	}

	members := make([]model.FamilyMember, len(rows))
	for i, row := range rows {
		members[i] = model.FamilyMember{
			ID:         dbutil.UUIDToString(row.ID),
			FamilyID:   dbutil.UUIDToString(row.FamilyID),
			UserID:     dbutil.UUIDToString(row.UserID),
			FamilyRole: row.FamilyRole,
			Status:     dbutil.TextToString(row.Status),
			JoinedAt:   dbutil.TimestampToString(row.JoinedAt),
			CreatedAt:  dbutil.TimestampToString(row.CreatedAt),
			UpdatedAt:  dbutil.TimestampToString(row.UpdatedAt),
		}
	}
	return members, nil
}

func (r *familyRepository) Create(ctx context.Context, f *model.Family) error {
	return r.queries.CreateFamily(ctx, generated.CreateFamilyParams{
		Name:       f.Name,
		UniqueCode: f.UniqueCode,
		CreatedBy:  dbutil.ParseUUID(f.CreatedBy),
		Status:     dbutil.StringToText(f.Status),
	})
}

func (r *familyRepository) AddMember(ctx context.Context, member *model.FamilyMember) error {
	return r.queries.AddFamilyMember(ctx, generated.AddFamilyMemberParams{
		FamilyID:   dbutil.ParseUUID(member.FamilyID),
		UserID:     dbutil.ParseUUID(member.UserID),
		FamilyRole: member.FamilyRole,
		Status:     dbutil.StringToText(member.Status),
		JoinedAt:   dbutil.ParseTimestamp(member.JoinedAt),
	})
}

func (r *familyRepository) UpdateMemberStatus(ctx context.Context, familyID, userID, status string) error {
	return r.queries.UpdateFamilyMemberStatus(ctx, generated.UpdateFamilyMemberStatusParams{
		Status:   dbutil.StringToText(status),
		FamilyID: dbutil.ParseUUID(familyID),
		UserID:   dbutil.ParseUUID(userID),
	})
}
