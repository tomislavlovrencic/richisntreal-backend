package auth

import (
	"net/http"
)

// Authenticator knows how to extract & validate a user ID from an HTTP request.
type Authenticator interface {
	// Authenticate returns the userID or an error if unauthenticated.
	Authenticate(r *http.Request) (int64, error)
}
