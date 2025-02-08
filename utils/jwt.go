package utils

import (
	"cbc-backend/config"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Claims represents the JWT claims
type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

// GenerateJWT generates a new JWT token
func GenerateJWT(userID string, username string, role string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(config.JWTSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}

	return tokenString, nil
}

// ValidateJWT validates a JWT token and returns the claims
func ValidateJWT(tokenString string) (*Claims, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return config.JWTSecret, nil
	})

	if err != nil {
		// Add more detailed error logging
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, fmt.Errorf("token expired")
			} else if ve.Errors&jwt.ValidationErrorSignatureInvalid != 0 {
				return nil, fmt.Errorf("invalid signature")
			}
		}
		return nil, fmt.Errorf("token validation error: %v", err)
	}

	// Get the claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		// Add claims debugging
		fmt.Printf("Token claims: UserID=%s, Username=%s, ExpiresAt=%v\n",
			claims.UserID, claims.Username, time.Unix(claims.ExpiresAt, 0))
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// GetUserIDFromToken extracts the user ID from a token string
func GetUserIDFromToken(tokenString string) (string, error) {
	claims, err := ValidateJWT(tokenString)
	if err != nil {
		return "", err
	}
	return claims.UserID, nil
}

// GetUsernameFromToken extracts the username from a token string
func GetUsernameFromToken(tokenString string) (string, error) {
	claims, err := ValidateJWT(tokenString)
	if err != nil {
		return "", err
	}
	return claims.Username, nil
}
