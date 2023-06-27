package storage

import (
	"auth-service/app/cmd/internal/domain/user/model"
	"context"
)

type Repository interface {
	Create(ctx context.Context, user *model.User) error
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetUserByID(ctx context.Context, userID string) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) error
	UpdateVerificationStatus(ctx context.Context, email string, status bool) error
	StoreVerificationData(ctx context.Context, verificationData *model.VerificationData) error
	GetVerificationData(ctx context.Context, email string, verificationType model.VerificationType) (*model.VerificationData, error)
	DeleteVerificationData(ctx context.Context, data model.VerificationData) error
	UpdatePassword(ctx context.Context, userID string, password string, tokenHash string) error
}
