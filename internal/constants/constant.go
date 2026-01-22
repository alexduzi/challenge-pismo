package constants

// ContextKey is a custom type for context keys to avoid collisions
type ContextKey string

const RequestIDKey ContextKey = "request_id"
const XRequestID = "X-Request-ID"
