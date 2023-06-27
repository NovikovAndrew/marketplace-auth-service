package storage

import (
	"auth-service/app/cmd/internal/domain/user/model"
	"auth-service/app/cmd/package/logger"
	"context"
	"github.com/hashicorp/go-hclog"
	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
	"time"
)

type PostgresRepository struct {
	db     *sqlx.DB
	logger logger.Logger
}

func NewUserRepository(db *sqlx.DB, logger logger.Logger) Repository {
	return &PostgresRepository{
		db:     db,
		logger: logger,
	}
}

func (pr *PostgresRepository) Create(ctx context.Context, user *model.User) error {
	user.ID = uuid.NewV4().String()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	query := "INSERT INTO users (id, email, password, username, token_hash, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	pr.logger.Debug("create user", hclog.Fmt("%#v", user))

	_, err := pr.db.ExecContext(ctx, query, user.ID, user.Email, user.Password, user.Username, user.TokenHash, user.CreatedAt, user.UpdatedAt)
	return err
}

func (pr *PostgresRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	pr.logger.Debug("get user by email", email)
	query := "SELECT * FROM users WHERE email = $1"
	var user model.User
	if err := pr.db.GetContext(ctx, &user, query, email); err != nil {
		pr.logger.Error("get user by email error: %v, query %s", err, query)
		return nil, err
	}

	return &user, nil
}

func (pr *PostgresRepository) GetUserByID(ctx context.Context, userID string) (*model.User, error) {
	pr.logger.Debug("get user by id", userID)
	query := "SELECT * FROM users WHERE id = $1"
	var user model.User
	if err := pr.db.GetContext(ctx, &user, query, userID); err != nil {
		pr.logger.Error("get user by id error: %v, query %s", err, query)
		return nil, err
	}

	return &user, nil
}

func (pr *PostgresRepository) UpdateUser(ctx context.Context, user *model.User) error {
	user.UpdatedAt = time.Now()
	query := "UPDATE users SET username = $1, updated_at = $2 WHERE id = $3"
	pr.logger.Debug("update user", user)

	_, err := pr.db.ExecContext(ctx, query, user.Username, user.UpdatedAt, user.ID)
	return err
}

func (pr *PostgresRepository) UpdateVerificationStatus(ctx context.Context, email string, status bool) error {
	pr.logger.Debug("update verification status by email", status, email)
	query := "UPDATE users SET is_verified = $1 WHERE email = $2"

	_, err := pr.db.ExecContext(ctx, query, status, email)
	return err
}

func (pr *PostgresRepository) StoreVerificationData(ctx context.Context, verificationData *model.VerificationData) error {
	pr.logger.Debug("store verification data", verificationData)
	query := "INSERT INTO verifications(code, email, expires_at, type) VALUES($1, $2, $3, $4)"

	_, err := pr.db.ExecContext(ctx, query, verificationData.Code, verificationData.Email, verificationData.ExpiresAt, verificationData.Type)
	return err
}

func (pr *PostgresRepository) GetVerificationData(ctx context.Context, email string, verificationType model.VerificationType) (*model.VerificationData, error) {
	pr.logger.Debug("get verification data by email and verification type", email, verificationType)
	query := "SELECT * FROM verifications WHERE email = $1 AND type = $2"
	var verificationData model.VerificationData

	if err := pr.db.GetContext(ctx, &verificationData, query, email, verificationType); err != nil {
		pr.logger.Error("get verification by email error: %v, query %s", err, query)
		return nil, err
	}

	return &verificationData, nil
}

func (pr *PostgresRepository) DeleteVerificationData(ctx context.Context, data model.VerificationData) error {
	pr.logger.Debug("delete verification data", data)
	query := "DELETE FROM verifications WHERE email = $1 AND type = $2"

	_, err := pr.db.ExecContext(ctx, query, data.Email, data.Type)
	return err
}

func (pr *PostgresRepository) UpdatePassword(ctx context.Context, userID string, password string, tokenHash string) error {
	pr.logger.Debug("update password", userID, password, tokenHash)
	query := "UPDATE users SET password = $1, token_hash = $2 WHERE id = $3"

	_, err := pr.db.ExecContext(ctx, query, password, tokenHash, userID)
	return err
}
