-- Migration: Add per-user OTP columns
ALTER TABLE users
    ADD COLUMN otp_secret VARCHAR(64),
    ADD COLUMN otp_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN otp_last_verified_at TIMESTAMP NULL;

CREATE INDEX IF NOT EXISTS idx_users_otp_secret ON users(otp_secret);
