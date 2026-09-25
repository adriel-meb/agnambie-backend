// File: internal/biblebrain/errors.go
// Purpose: Typed error types for Bible Brain upstream failures.
// Author: Backend Team
// Created: 2026-09-25
// Last Modified: 2026-09-25

package biblebrain

import "fmt"

// UpstreamError represents an unexpected HTTP status from the Bible Brain API.
type UpstreamError struct {
	StatusCode int
	Endpoint   string
	Body       string // truncated response body for debugging
}

// Error implements the error interface.
func (e *UpstreamError) Error() string {
	return fmt.Sprintf("bible brain returned status %d for %s", e.StatusCode, e.Endpoint)
}

// NotAllowedError is returned when a requested resource is not in the Gabon allowlist.
type NotAllowedError struct {
	Resource string // e.g. "fileset_id", "language"
	Value    string
}

// Error implements the error interface.
func (e *NotAllowedError) Error() string {
	return fmt.Sprintf("%s %q is not in the Gabon allowlist", e.Resource, e.Value)
}
