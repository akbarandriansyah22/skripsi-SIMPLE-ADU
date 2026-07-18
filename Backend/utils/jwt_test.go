package utils

import (
	"strings"
	"testing"
)

func TestJWTRequiresStrongConfiguredSecret(t *testing.T) {
	t.Setenv("JWT_SECRET", "replace_with_secure_random_secret")
	if _, err := GenerateJWT(1, RoleMahasiswa); err == nil {
		t.Fatal("placeholder JWT secret must be rejected")
	}

	secret := strings.Repeat("x", 32)
	t.Setenv("JWT_SECRET", secret)
	token, err := GenerateJWT(42, "petugas")
	if err != nil {
		t.Fatalf("GenerateJWT returned error: %v", err)
	}
	userID, role, err := ParseJWT(token)
	if err != nil {
		t.Fatalf("ParseJWT returned error: %v", err)
	}
	if userID != 42 || role != RoleAdmin {
		t.Fatalf("unexpected claims: userID=%d role=%q", userID, role)
	}
}
