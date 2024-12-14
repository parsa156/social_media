package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"social_media/internal/domain"
)

// JWTManager manages JWT operations.
type JWTManager struct {
	secretKey     string
	tokenDuration time.Duration
}

// Claims defines the custom claims structure.
// Note that UserID is now a string to match the UUID type in the domain model.
type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// NewJWTManager creates a new JWTManager.
func NewJWTManager(secretKey string, tokenDuration time.Duration) *JWTManager {
	return &JWTManager{
		secretKey:     secretKey,
		tokenDuration: tokenDuration,
	}
}

// Generate creates a new JWT token for a user.
func (manager *JWTManager) Generate(user *domain.User) (string, error) {
	username := ""
	if user.Username != nil {
		username = *user.Username
	}

	claims := Claims{
		UserID:   user.ID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(manager.tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(manager.secretKey))
}

// Verify validates the token and returns the claims.
func (manager *JWTManager) Verify(accessToken string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(manager.secretKey), nil
		},
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, jwt.ErrTokenInvalidClaims
	}
	return claims, nil
}
