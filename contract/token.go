package contract

import "time"

// Claims is a generic map of token payload data.
type Claims map[string]any

// Token is the port for any token/auth adapter.
// Implement this interface to swap between JWT, PASETO, opaque tokens, etc.
type Token interface {
	Sign(claims Claims, expiry time.Duration) (string, error)
	Verify(token string) (Claims, error)
}
