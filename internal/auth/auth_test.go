package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestHashPassword(t *testing.T) {
	password := "testpassword123"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}
	if hash == "" {
		t.Fatal("HashPassword returned empty hash")
	}
	if hash == password {
		t.Fatal("HashPassword returned the original password")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "testpassword123"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	err = CheckPasswordHash(password, hash)
	if err != nil {
		t.Fatalf("CheckPasswordHash failed for correct password: %v", err)
	}

	err = CheckPasswordHash("wrongpassword", hash)
	if err == nil {
		t.Fatal("CheckPasswordHash should have failed for wrong password")
	}
}

func TestMakeJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "testsecret"
	expiresIn := time.Hour

	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}
	if token == "" {
		t.Fatal("MakeJWT returned empty token")
	}
}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "testsecret"
	expiresIn := time.Hour

	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	validatedUserID, err := ValidateJWT(token, tokenSecret)
	if err != nil {
		t.Fatalf("ValidateJWT failed: %v", err)
	}
	if validatedUserID != userID {
		t.Fatalf("ValidateJWT returned wrong user ID. Expected %v, got %v", userID, validatedUserID)
	}
}

func TestValidateJWTWithWrongSecret(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "testsecret"
	wrongSecret := "wrongsecret"
	expiresIn := time.Hour

	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	_, err = ValidateJWT(token, wrongSecret)
	if err == nil {
		t.Fatal("ValidateJWT should have failed with wrong secret")
	}
}

func TestValidateJWTExpired(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "testsecret"
	expiresIn := -time.Hour // Expired token

	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	_, err = ValidateJWT(token, tokenSecret)
	if err == nil {
		t.Fatal("ValidateJWT should have failed for expired token")
	}
}

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name        string
		authHeader  string
		expectToken string
		expectError bool
	}{
		{
			name:        "Valid bearer token",
			authHeader:  "Bearer abc123",
			expectToken: "abc123",
			expectError: false,
		},
		{
			name:        "Valid bearer token with extra spaces",
			authHeader:  "Bearer   abc123   ",
			expectToken: "abc123",
			expectError: false, // TrimSpace should handle this
		},
		{
			name:        "Missing Authorization header",
			authHeader:  "",
			expectToken: "",
			expectError: true,
		},
		{
			name:        "Wrong prefix",
			authHeader:  "Basic abc123",
			expectToken: "",
			expectError: true,
		},
		{
			name:        "Bearer without token",
			authHeader:  "Bearer",
			expectToken: "",
			expectError: true,
		},
		{
			name:        "Bearer with only spaces",
			authHeader:  "Bearer   ",
			expectToken: "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			headers := make(http.Header)
			if tt.authHeader != "" {
				headers.Set("Authorization", tt.authHeader)
			}

			token, err := GetBearerToken(headers)

			if tt.expectError {
				if err == nil {
					t.Fatalf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}
				if token != tt.expectToken {
					t.Fatalf("Expected token %q, got %q", tt.expectToken, token)
				}
			}
		})
	}
}

func TestMakeRefreshToken(t *testing.T) {
	token1, err := MakeRefreshToken()
	if err != nil {
		t.Fatalf("MakeRefreshToken failed: %v", err)
	}
	if token1 == "" {
		t.Fatal("MakeRefreshToken returned empty token")
	}
	if len(token1) != 64 { // 32 bytes * 2 (hex encoding)
		t.Fatalf("Expected token length 64, got %d", len(token1))
	}

	// Test that tokens are unique
	token2, err := MakeRefreshToken()
	if err != nil {
		t.Fatalf("MakeRefreshToken failed: %v", err)
	}
	if token1 == token2 {
		t.Fatal("MakeRefreshToken returned duplicate tokens")
	}
}

func TestGetAPIKey(t *testing.T) {
	tests := []struct {
		name        string
		authHeader  string
		expectKey   string
		expectError bool
	}{
		{
			name:        "Valid API key",
			authHeader:  "ApiKey abc123",
			expectKey:   "abc123",
			expectError: false,
		},
		{
			name:        "Valid API key with extra spaces",
			authHeader:  "ApiKey   def456   ",
			expectKey:   "def456",
			expectError: false,
		},
		{
			name:        "Missing Authorization header",
			authHeader:  "",
			expectKey:   "",
			expectError: true,
		},
		{
			name:        "Wrong prefix",
			authHeader:  "Bearer abc123",
			expectKey:   "",
			expectError: true,
		},
		{
			name:        "ApiKey without key",
			authHeader:  "ApiKey",
			expectKey:   "",
			expectError: true,
		},
		{
			name:        "ApiKey with only spaces",
			authHeader:  "ApiKey   ",
			expectKey:   "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			headers := make(http.Header)
			if tt.authHeader != "" {
				headers.Set("Authorization", tt.authHeader)
			}

			key, err := GetAPIKey(headers)

			if tt.expectError {
				if err == nil {
					t.Fatalf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}
				if key != tt.expectKey {
					t.Fatalf("Expected key %q, got %q", tt.expectKey, key)
				}
			}
		})
	}
}