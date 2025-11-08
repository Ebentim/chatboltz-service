package entity

import "time"

// OTPPurpose defines the three supported OTP use cases
type OTPPurpose string

const (
	OTPPurpose2FA            OTPPurpose = "2fa"
	OTPPurposeForgotPassword OTPPurpose = "password_reset"
	OTPPurposeLogin          OTPPurpose = "login"
)

// OTPRequest represents an OTP generation request
type OTPRequest struct {
	Email   string     `json:"email" binding:"required,email"`
	Purpose OTPPurpose `json:"purpose" binding:"required"`
	Length  int        `json:"length,omitempty"`
}

// OTPVerifyRequest represents an OTP verification request
type OTPVerifyRequest struct {
	Email   string     `json:"email" binding:"required,email"`
	Purpose OTPPurpose `json:"purpose" binding:"required"`
	Code    string     `json:"code" binding:"required"`
}

// OTPResponse represents the response after OTP generation
type OTPResponse struct {
	Message    string `json:"message"`
	TTLMinutes int    `json:"ttl_minutes"`
	Purpose    string `json:"purpose"`
}

// OTPVerifyResponse represents the response after OTP verification
type OTPVerifyResponse struct {
	Message   string     `json:"message"`
	Purpose   string     `json:"purpose"`
	Verified  bool       `json:"verified"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

// IsValid checks if the OTP purpose is valid
func (p OTPPurpose) IsValid() bool {
	switch p {
	case OTPPurpose2FA, OTPPurposeForgotPassword, OTPPurposeLogin:
		return true
	default:
		return false
	}
}

// String returns the string representation of OTPPurpose
func (p OTPPurpose) String() string {
	return string(p)
}
