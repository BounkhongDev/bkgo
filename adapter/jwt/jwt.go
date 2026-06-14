package jwt

import (
	"fmt"
	"time"

	"github.com/bounkhongdev/kbgo/config"
	"github.com/bounkhongdev/kbgo/contract"
	gojwt "github.com/golang-jwt/jwt/v5"
)

// JWT is the JWT adapter implementing contract.Token.
type JWT struct {
	secret []byte
}

// New creates a JWT token signer/verifier from config.
func New(cfg config.JWT) *JWT {
	return &JWT{secret: []byte(cfg.Secret)}
}

func (j *JWT) Sign(claims contract.Claims, expiry time.Duration) (string, error) {
	mc := gojwt.MapClaims{}
	for k, v := range claims {
		mc[k] = v
	}
	mc["iat"] = time.Now().Unix()
	mc["exp"] = time.Now().Add(expiry).Unix()

	token := gojwt.NewWithClaims(gojwt.SigningMethodHS256, mc)
	signed, err := token.SignedString(j.secret)
	if err != nil {
		return "", fmt.Errorf("jwt: sign: %w", err)
	}
	return signed, nil
}

func (j *JWT) Verify(tokenStr string) (contract.Claims, error) {
	token, err := gojwt.Parse(tokenStr, func(t *gojwt.Token) (any, error) {
		if _, ok := t.Method.(*gojwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("jwt: unexpected signing method: %v", t.Header["alg"])
		}
		return j.secret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("jwt: verify: %w", err)
	}

	mc, ok := token.Claims.(gojwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("jwt: invalid token")
	}

	claims := make(contract.Claims)
	for k, v := range mc {
		claims[k] = v
	}
	return claims, nil
}
