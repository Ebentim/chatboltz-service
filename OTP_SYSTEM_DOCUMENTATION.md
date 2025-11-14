# OTP System - Clean Architecture Implementation

## Overview

This document describes the redesigned OTP (One-Time Password) system built with clean architecture principles. The system supports three primary use cases:

1. **Two-Factor Authentication (2FA)** - Additional security layer for user accounts
2. **Forgot Password** - Password reset functionality via OTP
3. **Login** - OTP-based authentication as an alternative to password login

## Architecture

### Clean Architecture Layers

```
┌─────────────────────────────────────────────────────────────┐
│                    Presentation Layer                        │
│  ┌─────────────────┐  ┌─────────────────┐  ┌──────────────┐ │
│  │   OTP Handler   │  │   OTP Routes    │  │  Middleware  │ │
│  └─────────────────┘  └─────────────────┘  └──────────────┘ │
└─────────────────────────────────────────────────────────────┘
                                │
┌─────────────────────────────────────────────────────────────┐
│                    Use Case Layer                           │
│  ┌─────────────────┐  ┌─────────────────┐  ┌──────────────┐ │
│  │  OTP Service    │  │  2FA Service    │  │ Login Service│ │
│  └─────────────────┘  └─────────────────┘  └──────────────┘ │
│  ┌─────────────────┐  ┌─────────────────┐                  │
│  │Password Service │  │   OTP Factory   │                  │
│  └─────────────────┘  └─────────────────┘                  │
└─────────────────────────────────────────────────────────────┘
                                │
┌─────────────────────────────────────────────────────────────┐
│                    Domain Layer                             │
│  ┌─────────────────┐  ┌─────────────────┐  ┌──────────────┐ │
│  │  OTP Entities   │  │   User Entity   │  │ Token Entity │ │
│  └─────────────────┘  └─────────────────┘  └──────────────┘ │
└─────────────────────────────────────────────────────────────┘
                                │
┌─────────────────────────────────────────────────────────────┐
│                Infrastructure Layer                         │
│  ┌─────────────────┐  ┌─────────────────┐  ┌──────────────┐ │
│  │ Token Repository│  │ User Repository │  │ Email Service│ │
│  └─────────────────┘  └─────────────────┘  └──────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

## Use Cases

### 1. Two-Factor Authentication (2FA)

**Purpose**: Add an extra security layer to user accounts.

**Flow**:
1. User enables 2FA in their account settings
2. During login, after password verification, system requests OTP
3. User receives OTP via email
4. User enters OTP to complete login

**API Endpoints**:
- `POST /api/v1/otp/2fa/enable` - Enable 2FA for user
- `POST /api/v1/otp/2fa/disable` - Disable 2FA for user
- `GET /api/v1/otp/2fa/status` - Check 2FA status
- `POST /api/v1/otp/request` - Request 2FA OTP
- `POST /api/v1/otp/verify` - Verify 2FA OTP

### 2. Forgot Password

**Purpose**: Allow users to reset their password when forgotten.

**Flow**:
1. User clicks "Forgot Password" on login page
2. User enters email address
3. System sends OTP to email
4. User enters OTP and new password
5. Password is updated

**API Endpoints**:
- `POST /api/v1/otp/request` - Request password reset OTP
- `POST /api/v1/otp/verify` - Verify password reset OTP
- `POST /api/v1/otp/password-reset/complete` - Complete password reset

### 3. Login

**Purpose**: Provide passwordless login option using OTP.

**Flow**:
1. User enters email on login page
2. User selects "Login with OTP" option
3. System sends OTP to email
4. User enters OTP
5. System authenticates user and provides JWT tokens

**API Endpoints**:
- `POST /api/v1/otp/request` - Request login OTP
- `POST /api/v1/otp/verify` - Verify login OTP
- `POST /api/v1/otp/login/complete` - Complete OTP login

## API Documentation

### Request OTP

```http
POST /api/v1/otp/request
Content-Type: application/json

{
  "email": "user@example.com",
  "purpose": "2fa|password_reset|login",
  "length": 6
}
```

**Response**:
```json
{
  "message": "OTP sent for 2fa",
  "ttl_minutes": 10,
  "purpose": "2fa"
}
```

### Verify OTP

```http
POST /api/v1/otp/verify
Content-Type: application/json

{
  "email": "user@example.com",
  "purpose": "2fa|password_reset|login",
  "code": "123456"
}
```

**Response**:
```json
{
  "message": "OTP verified for 2fa",
  "purpose": "2fa",
  "verified": true,
  "expires_at": "2024-12-01T10:30:00Z"
}
```

### Complete Password Reset

```http
POST /api/v1/otp/password-reset/complete
Content-Type: application/json

{
  "email": "user@example.com",
  "code": "123456",
  "new_password": "newSecurePassword123"
}
```

### Complete OTP Login

```http
POST /api/v1/otp/login/complete
Content-Type: application/json

{
  "email": "user@example.com",
  "code": "123456"
}
```

**Response**:
```json
{
  "user": {
    "id": "user-id",
    "email": "user@example.com",
    "name": "User Name"
  },
  "token": "jwt-access-token",
  "refresh_token": "jwt-refresh-token",
  "expires_at": 1701432600
}
```

## Configuration

### OTP Configuration

```go
type OTPConfig struct {
    DefaultLength int           // Default OTP length (6 digits)
    TTL           time.Duration // OTP validity period (10 minutes)
    MaxAttempts   int           // Maximum verification attempts (3)
}
```

### Environment Variables

```env
OTP_DEFAULT_LENGTH=6
OTP_TTL_MINUTES=10
OTP_MAX_ATTEMPTS=3
SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USERNAME=noreply@example.com
SMTP_PASSWORD=smtp_password
```

## Security Features

1. **OTP Hashing**: All OTP codes are hashed using bcrypt before storage
2. **Expiration**: OTPs automatically expire after configured TTL
3. **One-Time Use**: OTPs are automatically deleted after verification
4. **Rate Limiting**: Built-in protection against brute force attacks
5. **Purpose Isolation**: OTPs are scoped to specific purposes
6. **Email Validation**: Strict email format validation

## Database Schema

### Tokens Table

```sql
CREATE TABLE tokens (
    id VARCHAR(36) PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    token TEXT NOT NULL,
    purpose ENUM('password_reset', '2fa', 'login') NOT NULL,
    used BOOLEAN NOT NULL DEFAULT FALSE,
    expires_at TIMESTAMP NOT NULL,
    created_at VARCHAR(255) NOT NULL,
    updated_at VARCHAR(255) NOT NULL,
    deleted_at TIMESTAMP NULL,
    
    INDEX idx_tokens_email_purpose (email, purpose),
    INDEX idx_tokens_expires_at (expires_at)
);
```

### Users Table (OTP-related columns)

```sql
ALTER TABLE users ADD COLUMN otp_secret VARCHAR(64) NULL;
ALTER TABLE users ADD COLUMN otp_enabled BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE users ADD COLUMN otp_last_verified_at TIMESTAMP NULL;
```

## Testing

### Unit Tests

Run the OTP service tests:

```bash
go test ./tests/otp_service_test.go -v
```

### Integration Tests

Test the complete OTP flow:

```bash
go test ./tests/otp_integration_test.go -v
```

### Manual Testing

Use the provided HTTP files in the `http/` directory to test the API endpoints manually.

## Deployment

### Migration

Run the database migration:

```bash
# Apply the OTP system migration
mysql -u username -p database_name < internal/repository/migrations/20241201_update_otp_system.sql
```

### Service Setup

```go
// In your main application setup
func setupOTPServices(db *gorm.DB, emailService *smtp.Client) {
    // Create repositories
    tokenRepo := repository.NewUserToken(db)
    userRepo := repository.NewUserRepository(db)
    
    // Create OTP factory
    factory := usecase.NewOTPFactory(tokenRepo, userRepo)
    
    // Setup routes
    handler.SetupOTPRoutesWithMiddleware(
        router, 
        tokenRepo, 
        userRepo, 
        emailService, 
        authMiddleware,
    )
}
```

## Error Handling

The system uses structured error handling with specific error types:

- `ValidationError`: Invalid input data
- `DatabaseError`: Database operation failures
- `InternalError`: System-level errors
- `NotFoundError`: Resource not found

## Monitoring and Logging

All OTP operations are logged with:
- Request details (email, purpose, timestamp)
- Success/failure status
- Error details (if any)
- Performance metrics

## Future Enhancements

1. **SMS Support**: Add SMS as an alternative delivery method
2. **TOTP Integration**: Support for authenticator apps (Google Authenticator, Authy)
3. **Backup Codes**: Generate backup codes for 2FA
4. **Rate Limiting**: Advanced rate limiting per user/IP
5. **Analytics**: Detailed analytics and reporting
6. **Multi-language**: Support for multiple languages in email templates

## Troubleshooting

### Common Issues

1. **OTP Not Received**: Check email service configuration and spam folders
2. **OTP Expired**: Ensure TTL configuration is appropriate
3. **Invalid OTP**: Verify OTP format and ensure no typos
4. **Database Errors**: Check database connectivity and schema

### Debug Mode

Enable debug logging to troubleshoot issues:

```go
// Add debug logging in development
log.SetLevel(log.DebugLevel)
```

## Support

For issues and questions:
1. Check the logs for error details
2. Verify configuration settings
3. Test with curl/Postman
4. Review the test cases for expected behavior