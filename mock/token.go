package mock

import (
	"time"

	"github.com/bounkhongdev/kbgo/contract"
)

// Token is a test double for contract.Token.
// Override SignFn / VerifyFn to control behaviour, or use the defaults:
//   - Sign always returns "mock-token"
//   - Verify always returns an empty Claims map
type Token struct {
	SignFn   func(claims contract.Claims, expiry time.Duration) (string, error)
	VerifyFn func(token string) (contract.Claims, error)
}

func (t *Token) Sign(claims contract.Claims, expiry time.Duration) (string, error) {
	if t.SignFn != nil {
		return t.SignFn(claims, expiry)
	}
	return "mock-token", nil
}

func (t *Token) Verify(token string) (contract.Claims, error) {
	if t.VerifyFn != nil {
		return t.VerifyFn(token)
	}
	return contract.Claims{}, nil
}
