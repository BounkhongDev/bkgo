package jwt_test

import (
	"testing"
	"time"

	"github.com/bounkhongdev/kbgo/adapter/jwt"
	"github.com/bounkhongdev/kbgo/config"
	"github.com/bounkhongdev/kbgo/contract"
)

func newJWT() *jwt.JWT {
	return jwt.New(config.JWT{Secret: "test-secret-key"})
}

func TestSignAndVerify(t *testing.T) {
	j := newJWT()

	claims := contract.Claims{
		"user_id": "abc-123",
		"role":    "admin",
	}

	token, err := j.Sign(claims, time.Hour)
	if err != nil {
		t.Fatalf("Sign failed: %v", err)
	}
	if token == "" {
		t.Fatal("Sign returned empty token")
	}

	got, err := j.Verify(token)
	if err != nil {
		t.Fatalf("Verify failed: %v", err)
	}
	if got["user_id"] != "abc-123" {
		t.Errorf("user_id: want abc-123 got %v", got["user_id"])
	}
	if got["role"] != "admin" {
		t.Errorf("role: want admin got %v", got["role"])
	}
}

func TestVerifyInvalidToken(t *testing.T) {
	j := newJWT()
	_, err := j.Verify("this.is.not.valid")
	if err == nil {
		t.Error("expected error for invalid token")
	}
}

func TestVerifyWrongSecret(t *testing.T) {
	signer := jwt.New(config.JWT{Secret: "secret-a"})
	verifier := jwt.New(config.JWT{Secret: "secret-b"})

	token, _ := signer.Sign(contract.Claims{"id": "1"}, time.Hour)
	_, err := verifier.Verify(token)
	if err == nil {
		t.Error("expected error when verifying with wrong secret")
	}
}

func TestExpiredToken(t *testing.T) {
	j := newJWT()
	// Sign with a negative expiry so it's already expired
	token, err := j.Sign(contract.Claims{"id": "1"}, -time.Second)
	if err != nil {
		t.Fatalf("Sign failed: %v", err)
	}
	_, err = j.Verify(token)
	if err == nil {
		t.Error("expected error for expired token")
	}
}
