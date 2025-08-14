// Package dhttp provides HTTP utilities for microservices with integrated logging, tracing, and error handling.
//
// This package is designed for StreamingFast products and provides common HTTP handling patterns
// used across various microservices. It includes structured error responses, request validation,
// response writing utilities, and handler wrappers that automatically handle logging and tracing.
//
// # Handler Wrappers
//
// The package provides three main handler wrappers that simplify common HTTP patterns:
//
//   - JSONHandler: Wraps functions that return JSON responses
//   - RawHandler: Streams raw content from io.ReadCloser
//   - DirectHandler: Provides direct response writer control with error handling
//
// Example using JSONHandler:
//
//	http.Handle("/api/users", dhttp.JSONHandler(func(r *http.Request) (any, error) {
//		// Your logic here
//		return users, nil
//	}))
//
// # Request Processing
//
// Extract and validate request data using ExtractRequest for URL/query parameters
// or ExtractJSONRequest for JSON payloads:
//
//	type UserRequest struct {
//		ID   int    `schema:"id"`
//		Name string `schema:"name"`
//	}
//
//	var req UserRequest
//	if err := dhttp.ExtractRequest(ctx, r, &req, dhttp.NoValidation); err != nil {
//		// Handle validation error
//	}
//
// # Response Writing
//
// Write responses with proper headers, logging, and tracing:
//
//	dhttp.WriteJSON(ctx, w, responseData)
//	dhttp.WriteError(ctx, w, err)
//	dhttp.WriteText(ctx, w, "Hello World")
//
// # Error Handling
//
// Create structured HTTP errors with proper status codes:
//
//	return dhttp.BadRequestError(ctx, nil, "invalid_input", "The request is invalid")
//	return dhttp.NotFoundError(ctx, nil, "user_not_found", "User not found")
//	return dhttp.InternalServerError(ctx, err, "database_error", "Database connection failed")
//
// The package provides pre-made error functions for all standard HTTP status codes:
//
//	return dhttp.UnauthorizedError(ctx, nil, "auth_required", "Authentication required")
//	return dhttp.ForbiddenError(ctx, nil, "access_denied", "Access denied")
//	return dhttp.TooManyRequestsError(ctx, nil, "rate_limit", "Rate limit exceeded")
//
// See errors.go for the complete list of available HTTP error functions (4xx and 5xx).
//
// All HTTP errors are wrapped and can be easily used throughout your application.
// Error responses include trace IDs for debugging and follow a consistent JSON structure.
//
// # Utilities
//
// Additional utilities include:
//
//   - RealIP: Extract the real client IP from requests (Google Cloud Load Balancer aware)
//   - ForwardResponse: Forward HTTP responses from upstream services
//   - EmptyBody: Return empty response bodies from JSON handlers
//
// The package integrates with OpenTelemetry for tracing and uses structured logging
// throughout all operations.
package dhttp
