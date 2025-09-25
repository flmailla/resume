package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"fmt"
	"net/http"
	"crypto/rsa"
	"time"
	"math/big"
	"encoding/json"
	"encoding/base64"
)

// List of Json Web Key Sets
type JWKS struct {
	Keys []JWK `json:"keys"`
}

// Parts of a Json Web Key set
type JWK struct {
	Kty string `json:"kty"` // Key type
	Kid string `json:"kid"` // Key ID
	Use string `json:"use"` // Public key use
	N   string `json:"n"`   // RSA modulus
	E   string `json:"e"`   // RSA exponent
	Alg string `json:"alg"` // Algorithm
}

// Standard claims
type Claims struct {
	jwt.RegisteredClaims
}


// JWT validator components
type JWTValidator struct {
	jwksURL    string
	httpClient *http.Client
	keyCache   map[string]*rsa.PublicKey
	cacheTime  time.Time
	cacheTTL   time.Duration
}

// Instantiate a new JWT Validator
func NewJWTValidator(jwksURL string) *JWTValidator {
	return &JWTValidator{
		jwksURL: jwksURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		keyCache: make(map[string]*rsa.PublicKey),
		cacheTTL: 1 * time.Hour,
	}
}

// Get The key sets from an URL defined in config.go
func (v *JWTValidator) fetchJWKS() (*JWKS, error) {
	resp, err := v.httpClient.Get(v.jwksURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JWKS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("JWKS endpoint returned status: %d", resp.StatusCode)
	}

	var jwks JWKS
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return nil, fmt.Errorf("failed to decode JWKS: %w", err)
	}

	return &jwks, nil
}

// Return a RSA public Key from a JWK attributes
func (v *JWTValidator) jwkToRSAPublicKey(jwk JWK) (*rsa.PublicKey, error) {
	if jwk.Kty != "RSA" {
		return nil, fmt.Errorf("unsupported key type: %s", jwk.Kty)
	}

	nBytes, err := base64.RawURLEncoding.DecodeString(jwk.N)
	if err != nil {
		return nil, fmt.Errorf("failed to decode modulus: %w", err)
	}

	eBytes, err := base64.RawURLEncoding.DecodeString(jwk.E)
	if err != nil {
		return nil, fmt.Errorf("failed to decode exponent: %w", err)
	}

	n := new(big.Int).SetBytes(nBytes)
	e := new(big.Int).SetBytes(eBytes)

	return &rsa.PublicKey{
		N: n,
		E: int(e.Int64()),
	}, nil
}

// Extract the RSA public Key for a given kid in the JWKS
func (v *JWTValidator) getPublicKey(kid string) (*rsa.PublicKey, error) {
	if time.Since(v.cacheTime) < v.cacheTTL {
		if key, exists := v.keyCache[kid]; exists {
			return key, nil
		}
	}

	jwks, err := v.fetchJWKS()
	if err != nil {
		return nil, err
	}

	v.keyCache = make(map[string]*rsa.PublicKey)
	v.cacheTime = time.Now()

	for _, jwk := range jwks.Keys {
		publicKey, err := v.jwkToRSAPublicKey(jwk)
		if err != nil {
			continue
		}
		v.keyCache[jwk.Kid] = publicKey
	}

	if key, exists := v.keyCache[kid]; exists {
		return key, nil
	}

	return nil, fmt.Errorf("key with ID %s not found", kid)
}


// Validate the token signature
func (v *JWTValidator) keyFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}

	kid, ok := token.Header["kid"].(string)
	if !ok {
		return nil, fmt.Errorf("token header missing kid")
	}

	return v.getPublicKey(kid)
}


// Validate custom claims in the JWT token
func (v *JWTValidator) validateCustomClaims(claims *Claims) error {
	return nil
}


// Globally check the validity of a JWT token 
func (v *JWTValidator) verifyToken(tokenString string) error {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, 
		claims, 
		v.keyFunc, 
		jwt.WithAudience(aud), 
		jwt.WithIssuer(iss),
		jwt.WithExpirationRequired() )
	if err != nil {
		return fmt.Errorf("token validation failed: %w", err)
	}

	if !token.Valid {
		return fmt.Errorf("token is not valid")
	}

	err = v.validateCustomClaims(claims)
	if err != nil {
		return fmt.Errorf("token validation failed: %w", err)
	}

	return nil
}