package hash_test

import (
	"testing"

	"github.com/BounkhongDev/bkgo/hash"
)

func TestPassword(t *testing.T) {
	plain := "supersecret123"

	hashed, err := hash.Password(plain)
	if err != nil {
		t.Fatalf("hash.Password failed: %v", err)
	}
	if hashed == plain {
		t.Error("hashed password must not equal plain text")
	}
	if len(hashed) == 0 {
		t.Error("hashed password must not be empty")
	}
}

func TestCheckPassword(t *testing.T) {
	plain := "supersecret123"

	hashed, err := hash.Password(plain)
	if err != nil {
		t.Fatalf("hash.Password failed: %v", err)
	}

	t.Run("correct password matches", func(t *testing.T) {
		if !hash.CheckPassword(plain, hashed) {
			t.Error("expected CheckPassword=true for correct password")
		}
	})

	t.Run("wrong password does not match", func(t *testing.T) {
		if hash.CheckPassword("wrongpassword", hashed) {
			t.Error("expected CheckPassword=false for wrong password")
		}
	})
}

func TestUniqueHashes(t *testing.T) {
	plain := "samepassword"
	h1, _ := hash.Password(plain)
	h2, _ := hash.Password(plain)

	// bcrypt uses a random salt — same input must produce different hashes
	if h1 == h2 {
		t.Error("two hashes of the same password should differ (bcrypt salting)")
	}
}
