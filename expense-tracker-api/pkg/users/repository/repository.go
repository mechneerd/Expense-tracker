package repository

import (
	"context"
	"expense-tracker-api/internal/db/dbutil"
	"expense-tracker-api/internal/db/generated"
	"expense-tracker-api/pkg/users/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByID(ctx context.Context, id string) (*model.User, error)
	Create(ctx context.Context, u *model.User) error
	Update(ctx context.Context, u *model.User) error
	VerifyEmail(ctx context.Context, id string) error
}

type userRepository struct {
	queries *generated.Queries
}

func NewUserRepository(pool *pgxpool.Pool) *userRepository {
	return &userRepository{queries: generated.New(pool)}
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	row, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return &model.User{
		ID:              dbutil.UUIDToString(row.ID),
		Email:           row.Email,
		FirstName:       dbutil.TextToString(row.FirstName),
		LastName:        dbutil.TextToString(row.LastName),
		Phone:           dbutil.TextToString(row.Phone),
		AvatarURL:       dbutil.TextPtrToString(row.AvatarUrl),
		Status:          dbutil.TextToString(row.Status),
		CreatedAt:       dbutil.TimestampToTime(row.CreatedAt),
		UpdatedAt:       dbutil.TimestampToTime(row.UpdatedAt),
		EmailVerifiedAt: dbutil.TimestampToTime(row.EmailVerifiedAt),
	}, nil
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
	row, err := r.queries.GetUserByID(ctx, dbutil.ParseUUID(id))
	if err != nil {
		return nil, err
	}

	return &model.User{
		ID:              dbutil.UUIDToString(row.ID),
		Email:           row.Email,
		FirstName:       dbutil.TextToString(row.FirstName),
		LastName:        dbutil.TextToString(row.LastName),
		Phone:           dbutil.TextToString(row.Phone),
		AvatarURL:       dbutil.TextPtrToString(row.AvatarUrl),
		Status:          dbutil.TextToString(row.Status),
		CreatedAt:       dbutil.TimestampToTime(row.CreatedAt),
		UpdatedAt:       dbutil.TimestampToTime(row.UpdatedAt),
		EmailVerifiedAt: dbutil.TimestampToTime(row.EmailVerifiedAt),
	}, nil
}

func (r *userRepository) Create(ctx context.Context, u *model.User) error {
	return r.queries.CreateUser(ctx, generated.CreateUserParams{
		Email:     u.Email,
		FirstName: dbutil.StringToText(u.FirstName),
		LastName:  dbutil.StringToText(u.LastName),
		AvatarUrl: dbutil.StringToText(derefString(u.AvatarURL)),
		Status:    dbutil.StringToText(u.Status),
	})
}

func (r *userRepository) Update(ctx context.Context, u *model.User) error {
	return r.queries.UpdateUser(ctx, generated.UpdateUserParams{
		FirstName: dbutil.StringToText(u.FirstName),
		LastName:  dbutil.StringToText(u.LastName),
		Phone:     dbutil.StringToText(u.Phone),
		AvatarUrl: dbutil.StringToText(derefString(u.AvatarURL)),
		Status:    dbutil.StringToText(u.Status),
		ID:        dbutil.ParseUUID(u.ID),
	})
}

func (r *userRepository) VerifyEmail(ctx context.Context, id string) error {
	return r.queries.VerifyUserEmail(ctx, dbutil.ParseUUID(id))
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
