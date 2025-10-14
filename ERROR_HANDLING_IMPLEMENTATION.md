# Error Handling Implementation

## Overview

Implemented comprehensive error handling across handlers, usecases, and repositories with proper error categorization, logging, and client responses.

## Components Added

### 1. Custom Error Types (`internal/errors/errors.go`)

- **AppError struct** with type, message, code, and details
- **Error types**: ValidationError, AuthenticationError, NotFoundError, ConflictError, DatabaseError, ExternalAPIError, InternalError
- **Helper functions**: NewValidationError, NewAuthenticationError, etc.
- **Logging utilities**: LogError, WrapDatabaseError, WrapExternalAPIError

### 2. Error Response Handler (`internal/errors/response.go`)

- **HandleError function** for consistent API error responses
- Automatic GORM error detection and conversion
- Proper HTTP status codes for different error types

### 3. Middleware (`internal/middleware/error.go`)

- **ErrorHandler**: Global panic recovery with logging
- **RequestLogger**: Formatted request logging with timestamps

### 4. Validation Utilities (`internal/utils/validation.go`)

- **ValidateEmail**: Email format validation
- **ValidatePassword**: Password strength validation
- **ValidateRequired**: Required field validation

## Updated Components

### 1. User Repository (`internal/repository/user.go`)

- Wrapped all database operations with proper error handling
- Convert GORM errors to custom AppError types
- Added detailed error logging for database operations

### 2. User Usecase (`internal/usecase/user.go`)

- Added input validation for all methods
- Proper error handling for Firebase operations
- Enhanced error messages for authentication failures

### 3. Agent Repository (`internal/repository/agent.go`)

- Added error handling for all CRUD operations
- Proper NotFoundError handling for missing records
- Database error wrapping with operation context

### 4. Agent Usecase (`internal/usecase/agent.go`)

- Input validation for required fields
- Proper error propagation from repository layer
- Enhanced validation messages

### 5. Auth Handler (`internal/handler/auth.go`)

- Replaced manual error handling with HandleError function
- Consistent error responses across all endpoints
- Proper error logging with context

### 6. Agent Handler (`internal/handler/agent.go`)

- Unified error handling approach
- Proper HTTP status codes based on error types
- Enhanced error logging

## Error Flow

```log
Client Request → Handler → Usecase → Repository
                   ↓         ↓         ↓
              HandleError ← AppError ← Database Error
                   ↓
            JSON Response + Logging
```

## Error Response Format

```json
{
  "error": "User-friendly error message",
  "type": "ERROR_TYPE"
}
```

## Logging Format

```
[ERROR] Context: Error message (Type: ERROR_TYPE, Code: HTTP_CODE)
[ERROR] Details: Detailed error information
```

## Usage Examples

### In Repository

```go
if err := r.db.Create(user).Error; err != nil {
    return nil, appErrors.WrapDatabaseError(err, "create user")
}
```

### In Usecase

```go
if err := utils.ValidateEmail(req.Email); err != nil {
    return nil, err
}
```

### In Handler

```go
if err != nil {
    appErrors.HandleError(c, err, "SignupWithEmail")
    return
}
```

## Benefits

1. **Consistent Error Responses**: All endpoints return standardized error format
2. **Proper HTTP Status Codes**: Automatic mapping of error types to status codes
3. **Comprehensive Logging**: All errors logged with context and details
4. **Input Validation**: Prevents invalid data from reaching business logic
5. **Error Categorization**: Clear distinction between different error types
6. **Client-Friendly Messages**: User-friendly error messages without exposing internals
7. **Debugging Support**: Detailed logging for development and troubleshooting

## Next Steps

To use this error handling system in your application:

1. Import the middleware in your main application file
2. Add the error handler and request logger to your Gin router
3. Use the HandleError function in all your handlers
4. Follow the established patterns for new handlers, usecases, and repositories
