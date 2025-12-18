package utils

// HTTP status messages
const (
	// Success messages
	MsgSuccess          = "Operation successful"
	MsgCreated          = "Resource created successfully"
	MsgUpdated          = "Resource updated successfully"
	MsgDeleted          = "Resource deleted successfully"
	MsgLoginSuccess     = "Login successful"
	MsgRegisterSuccess  = "Registration successful"

	// Error messages
	MsgBadRequest       = "Invalid request"
	MsgUnauthorized     = "Unauthorized access"
	MsgForbidden        = "Access forbidden"
	MsgNotFound         = "Resource not found"
	MsgConflict         = "Resource already exists"
	MsgInternalError    = "Internal server error"
	
	// Validation messages
	MsgInvalidEmail     = "Invalid email format"
	MsgInvalidPassword  = "Password must be at least 6 characters"
	MsgEmailExists      = "Email already registered"
	MsgInvalidCredentials = "Invalid email or password"
	MsgAccountInactive  = "Account is inactive"
	MsgMissingToken     = "Missing authorization token"
	MsgInvalidToken     = "Invalid or expired token"
	MsgAdminRequired    = "Admin access required"
)

// Database error messages
const (
	MsgDBConnectionFailed = "Failed to connect to database"
	MsgDBQueryFailed      = "Database query failed"
	MsgDBCreateFailed     = "Failed to create record"
	MsgDBUpdateFailed     = "Failed to update record"
	MsgDBDeleteFailed     = "Failed to delete record"
)
