package auth

import (
	"auth-service/app/cmd/internal/config"
	"auth-service/app/cmd/internal/service/auth"
	"auth-service/app/cmd/package/logger"
)

type AuthenticationHandler struct {
	logger      logger.Logger
	cfg         config.Config
	authService auth.Authentication
}

func NewAuthenticationHandler(logger logger.Logger, cfg config.Config, authService auth.Authentication) *AuthenticationHandler {
	return &AuthenticationHandler{
		logger:      logger,
		cfg:         cfg,
		authService: authService,
	}
}
