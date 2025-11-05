-- Rollback: Remove per-user OTP columns
-- Assuming PostgreSQL; if MSSQL adjust syntax accordingly.
ALTER TABLE users DROP COLUMN IF EXISTS otp_secret;
ALTER TABLE users DROP COLUMN IF EXISTS otp_enabled;
ALTER TABLE users DROP COLUMN IF EXISTS otp_last_verified_at;

DROP INDEX IF EXISTS idx_users_otp_secret;
