package utils

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaim struct {
	UserID   uint   `json:"user_id"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	OutletID *uint  `json:"outlet_id,omitempty"`
	jwt.RegisteredClaims
}

func getJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "bazar-super-secret-key-2026"
	}
	return []byte(secret)
}

// GenerateToken creates a signed JWT token with user claims
func GenerateToken(userID uint, email string, role string, outletID *uint) (string, error) {
	expHoursStr := os.Getenv("JWT_EXP_HOURS")
	expHours, err := strconv.Atoi(expHoursStr)
	if err != nil || expHours <= 0 {
		expHours = 24 // default to 24 hours
	}

	claims := &JWTClaim{
		UserID:   userID,
		Email:    email,
		Role:     role,
		OutletID: outletID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "bazar-be",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getJWTSecret())
}

// ValidateToken parses and validates a JWT token string
func ValidateToken(tokenString string) (*JWTClaim, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaim{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return getJWTSecret(), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaim)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
