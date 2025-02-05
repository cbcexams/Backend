package middleware

import (
	"cbc-backend/config"
	"fmt"
	"strings"

	"github.com/beego/beego/v2/server/web/context"
	"github.com/dgrijalva/jwt-go"
)

// JWTMiddleware handles JWT authentication
func JWTMiddleware(ctx *context.Context) {
	// Get the Authorization header
	authHeader := ctx.Input.Header("Authorization")
	if authHeader == "" {
		ctx.Output.SetStatus(401)
		ctx.Output.JSON(map[string]interface{}{
			"success": false,
			"error":   "Authorization header is required",
		}, true, false)
		return
	}

	// Check Bearer prefix
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		ctx.Output.SetStatus(401)
		ctx.Output.JSON(map[string]interface{}{
			"success": false,
			"error":   "Invalid authorization format. Use: Bearer <token>",
		}, true, false)
		return
	}

	tokenString := parts[1]

	// Parse and validate the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Return the secret key used for signing
		return config.JWTSecret, nil
	})

	if err != nil {
		ctx.Output.SetStatus(401)
		ctx.Output.JSON(map[string]interface{}{
			"success": false,
			"error":   "Invalid or expired token",
		}, true, false)
		return
	}

	if !token.Valid {
		ctx.Output.SetStatus(401)
		ctx.Output.JSON(map[string]interface{}{
			"success": false,
			"error":   "Invalid token",
		}, true, false)
		return
	}

	// Token is valid - extract claims if needed
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		// Store user info in context for later use
		ctx.Input.SetData("user_id", claims["user_id"])
		ctx.Input.SetData("username", claims["username"])
	}
}
