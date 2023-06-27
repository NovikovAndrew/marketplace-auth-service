package auth

import "auth-service/app/cmd/internal/domain/user/model"

type Authentication interface {
	Authenticate(reqUser *model.User, user *model.User) bool
	GenerateAccessToken(user *model.User) (string, error)
	GenerateRefreshToken(user *model.User) (string, error)
	GenerateCustomKey(userID string, tokenHash string) string
	ValidateAccessToken(tokenString string) (string, error)
	ValidateRefreshToken(tokenString string) (string, string, error)
}
