package model

import "time"

type VerificationType int

const (
	MainVerification VerificationType = iota + 1
	PassReset
)

type VerificationData struct {
	Email     string           `json:"email"`
	Code      string           `json:"code"`
	ExpiresAt time.Time        `json:"expires_at"`
	Type      VerificationType `json:"type"`
}
