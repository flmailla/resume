package auth

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"reflect"
	"github.com/golang-jwt/jwt/v5"
	"math/big"
	"crypto/rand"
)

func TestJWTValidatorCreation(t *testing.T) {
	url := "dummy"
	client := &http.Client{
			Timeout: 10 * time.Second,
		}
	cache := make(map[string]*rsa.PublicKey)
	v:= NewJWTValidator(url)
	
	t.Run("check initialization", func(t *testing.T) {
		if v.jwksURL != url {
			t.Errorf("expected url %s, got %s", url, v.jwksURL)
		}

		if v.httpClient.Timeout != client.Timeout {
			t.Errorf("expected client %s, got %s", client.Timeout, v.httpClient.Timeout)
		}

		if reflect.TypeOf(v.keyCache) != reflect.TypeOf(cache) {
			t.Errorf("expected key cache %v, got %v", reflect.TypeOf(cache), reflect.TypeOf(v.keyCache))
		}

		if v.cacheTTL != (1 * time.Hour) {
			t.Errorf("expected cache %d, got %d", 1 * time.Hour, v.cacheTTL)
		}
	})
}


func TestFetchJWKS(t *testing.T) {
	sampleJWKS := JWKS{
		Keys: []JWK{
			{
				Kty: "RSA",
				Kid: "key1",
				Use: "sig",
				N:   base64.RawURLEncoding.EncodeToString([]byte("modulus")),
				E:   base64.RawURLEncoding.EncodeToString([]byte{0x01, 0x00, 0x01}), // Exponent 65537
				Alg: "RS256",
			},
		},
	}

	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(sampleJWKS)
	}))
	defer server.Close()

	validator := NewJWTValidator(server.URL)
	jwks, err := validator.fetchJWKS()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(jwks.Keys) != 1 {
		t.Errorf("expected 1 key, got %d", len(jwks.Keys))
	}
	if jwks.Keys[0].Kid != "key1" {
		t.Errorf("expected key ID 'key1', got %s", jwks.Keys[0].Kid)
	}

	validator = NewJWTValidator("http://nonexistent")
	_, err = validator.fetchJWKS()
	if err == nil {
		t.Error("expected error for invalid URL, got none")
	}

	// Test error case: non-200 status
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	validator = NewJWTValidator(server.URL)
	_, err = validator.fetchJWKS()
	if err == nil {
		t.Error("expected error for non-200 status, got none")
	}
}

func TestJwkToRSAPublicKey(t *testing.T) {
	validator := NewJWTValidator("")

	// Valid JWK
	jwk := JWK{
		Kty: "RSA",
		N:   base64.RawURLEncoding.EncodeToString([]byte("modulus")),
		E:   base64.RawURLEncoding.EncodeToString([]byte{0x01, 0x00, 0x01}), // Exponent 65537
	}
	key, err := validator.jwkToRSAPublicKey(jwk)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if key.E != 65537 {
		t.Errorf("expected exponent 65537, got %d", key.E)
	}

	jwk.Kty = "EC"
	_, err = validator.jwkToRSAPublicKey(jwk)
	if err == nil {
		t.Error("expected error for invalid key type, got none")
	}

	jwk.Kty = "RSA"
	jwk.N = "@"
	_, err = validator.jwkToRSAPublicKey(jwk)
	if err == nil {
		t.Error("expected error for invalid modulus, got none")
	}

	jwk.N = base64.RawURLEncoding.EncodeToString([]byte("modulus"))
	jwk.E = "@"
	_, err = validator.jwkToRSAPublicKey(jwk)
	if err == nil {
		t.Error("expected error for invalid exponent, got none")
	}
}

func TestGetPublicKey(t *testing.T) {
	sampleJWKS := JWKS{
		Keys: []JWK{
			{
				Kty: "RSA",
				Kid: "key1",
				N:   base64.RawURLEncoding.EncodeToString([]byte("modulus")),
				E:   base64.RawURLEncoding.EncodeToString([]byte{0x01, 0x00, 0x01}),
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(sampleJWKS)
	}))
	defer server.Close()

	validator := NewJWTValidator(server.URL)
	key, err := validator.getPublicKey("key1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if key == nil {
		t.Fatal("expected non-nil public key")
	}

	key, err = validator.getPublicKey("key1")
	if err != nil {
		t.Fatalf("expected no error from cache, got %v", err)
	}
	if key == nil {
		t.Fatal("expected non-nil public key from cache")
	}

	_, err = validator.getPublicKey("key2")
	if err == nil {
		t.Error("expected error for missing key, got none")
	}
}

func TestVerifyToken(t *testing.T) {
	rsaKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate RSA key: %v", err)
	}

	nBytes := rsaKey.PublicKey.N.Bytes()
	eBytes := big.NewInt(int64(rsaKey.PublicKey.E)).Bytes()
	sampleJWKS := JWKS{
		Keys: []JWK{
			{
				Kty: "RSA",
				Kid: "key1",
				N:   base64.RawURLEncoding.EncodeToString(nBytes),
				E:   base64.RawURLEncoding.EncodeToString(eBytes),
				Alg: "RS256",
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(sampleJWKS)
	}))
	defer server.Close()

	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{aud},
			Issuer:    iss,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = "key1"
	tokenString, err := token.SignedString(rsaKey)
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}

	validator := NewJWTValidator(server.URL)
	err = validator.verifyToken(tokenString)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = validator.verifyToken("invalid.token.string")
	if err == nil {
		t.Error("expected error for invalid token, got none")
	}

	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(-time.Hour))
	token = jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = "key1"
	tokenString, err = token.SignedString(rsaKey)
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}
	err = validator.verifyToken(tokenString)
	if err == nil {
		t.Error("expected error for expired token, got none")
	}
}