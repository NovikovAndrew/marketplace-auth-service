package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io/ioutil"
	"time"

	"auth-service/app/cmd/internal/config"
	"auth-service/app/cmd/internal/domain/user/model"
	"auth-service/app/cmd/package/logger"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

const (
	accessTokenType  = "access"
	refreshTokenType = "refresh"
)

type RefreshTokenCustomClaims struct {
	UserID    string
	CustomKey string
	KeyType   string
	jwt.StandardClaims
}

type AccessTokenCustomClaims struct {
	UserID  string
	KeyType string
	jwt.StandardClaims
}

type AuthenticationService struct {
	logger logger.Logger
	config config.Config
}

func NewAuthenticationService(logger logger.Logger, cfg config.Config) Authentication {
	return &AuthenticationService{
		logger: logger,
		config: cfg,
	}
}

// Authenticate check the user credentials from request and db
func (auth *AuthenticationService) Authenticate(reqUser *model.User, user *model.User) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqUser.Password)); err != nil {
		auth.logger.Debug("password hashed is not same", user.Password, reqUser.Password)
		return false
	}

	return true
}

// GenerateAccessToken generate new access token for user
func (auth *AuthenticationService) GenerateAccessToken(user *model.User) (string, error) {
	userID := user.ID

	claims := AccessTokenCustomClaims{
		userID,
		accessTokenType,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * time.Duration(auth.config.Token.JwtExpiration)).Unix(),
			Issuer:    "auth.service",
		},
	}

	signBytes, err := ioutil.ReadFile(auth.config.Token.AccessTokenPrivateKeyPath)
	if err != nil {
		auth.logger.Error("unable to read private key, error", err)
		return "", errors.New("could not generate access token. please try again")
	}

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		auth.logger.Error("unable to parse private key, error", err)
		return "", errors.New("could not parse private key. please try again")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	return token.SignedString(signKey)
}

// GenerateRefreshToken generate new refresh token for user
func (auth *AuthenticationService) GenerateRefreshToken(user *model.User) (string, error) {
	customKey := auth.GenerateCustomKey(user.ID, user.TokenHash)
	claims := RefreshTokenCustomClaims{
		user.ID,
		customKey,
		refreshTokenType,
		jwt.StandardClaims{
			Issuer: "auth.service",
		},
	}

	signBytes, err := ioutil.ReadFile(auth.config.Token.RefreshTokenPrivateKeyPath)
	if err != nil {
		auth.logger.Error("unable to read private key", err)
		return "", err
	}

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		auth.logger.Error("unable to parse private key", err)
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	return token.SignedString(signKey)
}

// GenerateCustomKey create new key for jwt payload
func (auth *AuthenticationService) GenerateCustomKey(userID string, tokenHash string) string {
	h := hmac.New(sha256.New, []byte(tokenHash))
	h.Write([]byte(userID))
	sha := hex.EncodeToString(h.Sum(nil))
	return sha
}

// ValidateAccessToken parser nad validate the given access token
func (auth *AuthenticationService) ValidateAccessToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessTokenCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			auth.logger.Error("unexpected signing method in auth token")
			return nil, errors.New("unexpected signing method in auth")
		}
		verifyBytes, err := ioutil.ReadFile(auth.config.Token.AccessTokenPublicKeyPath)
		if err != nil {
			auth.logger.Error("unable to read public access key", err)
			return nil, errors.New("unable to read public access key")
		}

		verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
		if err != nil {
			auth.logger.Error("unable to parse public access key", err)
			return nil, errors.New("unable to parse public access key")
		}

		return verifyKey, nil
	})

	if err != nil {
		auth.logger.Error("unable to parse claims, error", err)
		return "", err
	}

	claims, ok := token.Claims.(*AccessTokenCustomClaims)

	if !ok || !token.Valid || claims.UserID == "" || claims.KeyType != accessTokenType {
		return "", errors.New("invalid token: authentication failed")
	}

	return claims.UserID, nil
}

// ValidateRefreshToken parser and validate the given refresh token
func (auth *AuthenticationService) ValidateRefreshToken(tokenString string) (string, string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &RefreshTokenCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			auth.logger.Error("unexpected signing method in auth")
			return nil, errors.New("unexpected signing method in auth")
		}

		verifyBytes, err := ioutil.ReadFile(auth.config.Token.RefreshTokenPublicKeyPath)
		if err != nil {
			auth.logger.Error("unable to read public refresh key", err)
			return nil, errors.New("")
		}

		verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
		if err != nil {
			auth.logger.Error("unable to parse public refresh key", err)
			return nil, errors.New("unable to parse public refresh key")
		}

		return verifyKey, err
	})

	if err != nil {
		auth.logger.Error("unable to parse claims, error", err)
		return "", "", err
	}

	claims, ok := token.Claims.(*RefreshTokenCustomClaims)

	if !ok || !token.Valid || claims.UserID == "" || claims.KeyType != refreshTokenType {
		return "", "", errors.New("invalid token: authentication failed")
	}

	return claims.UserID, claims.CustomKey, nil
}
