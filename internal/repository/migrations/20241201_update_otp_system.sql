-- Migration: Update OTP system for clean architecture
-- Date: 2024-12-01
-- Description: Ensures the tokens table supports the new OTP system with proper enum values

-- Update the tokens table enum to ensure all required purposes are available
ALTER TABLE tokens MODIFY COLUMN purpose ENUM('password_reset', '2fa', 'login') NOT NULL;

-- Add index for better performance on email + purpose lookups
CREATE INDEX IF NOT EXISTS idx_tokens_email_purpose ON tokens(email, purpose);

-- Add index for expiration cleanup
CREATE INDEX IF NOT EXISTS idx_tokens_expires_at ON tokens(expires_at);

-- Ensure users table has the required OTP columns (they should already exist from previous migration)
-- This is just to ensure consistency

-- Add OTP secret column if it doesn't exist
ALTER TABLE users ADD COLUMN IF NOT EXISTS otp_secret VARCHAR(64) NULL;

-- Add OTP enabled column if it doesn't exist  
ALTER TABLE users ADD COLUMN IF NOT EXISTS otp_enabled BOOLEAN NOT NULL DEFAULT FALSE;

-- Add OTP last verified timestamp if it doesn't exist
ALTER TABLE users ADD COLUMN IF NOT EXISTS otp_last_verified_at TIMESTAMP NULL;

-- Add indexes for better performance
CREATE INDEX IF NOT EXISTS idx_users_otp_enabled ON users(otp_enabled);
CREATE INDEX IF NOT EXISTS idx_users_otp_secret ON users(otp_secret);

-- Clean up any expired tokens (optional maintenance)
DELETE FROM tokens WHERE expires_at < NOW();