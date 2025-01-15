package middleware

import (
	"cbc-backend/models"
	"strings"

	"github.com/beego/beego/v2/server/web/context"
)

func JWT(ctx *context.Context) {
	if shouldSkipAuth(ctx) {
		// For public endpoints, still create/validate session
		handleSession(ctx)
		return
	}

	authHeader := ctx.Input.Header("Authorization")
	if authHeader == "" {
		ctx.Output.SetStatus(401)
		ctx.Output.JSON(map[string]string{"error": "No authorization header"}, false, false)
		return
	}

	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
	token, err := models.ValidateToken(tokenString)
	if err != nil || !token.Valid {
		ctx.Output.SetStatus(401)
		ctx.Output.JSON(map[string]string{"error": "Invalid token"}, false, false)
		return
	}

	// Set session for authenticated requests
	handleSession(ctx)
}

func shouldSkipAuth(ctx *context.Context) bool {
	path := ctx.Input.URL()
	method := ctx.Input.Method()

	// List of paths that don't require authentication
	publicPaths := map[string][]string{
		"/v1/user/login":  {"POST"},
		"/v1/user/signup": {"POST"},
		"/v1/resources":   {"GET"}, // Make resources GET public
		"/v1/resources/*": {"GET"}, // Include any sub-paths
	}

	if methods, exists := publicPaths[path]; exists {
		for _, m := range methods {
			if m == method {
				return true
			}
		}
	}

	// Check if path matches /v1/resources/* pattern
	if strings.HasPrefix(path, "/v1/resources/") && method == "GET" {
		return true
	}

	return false
}

func handleSession(ctx *context.Context) {
	// Get existing session ID from cookie or header
	sessionID := ctx.Input.Cookie("session_id")
	if sessionID == "" {
		sessionID = ctx.Input.Header("X-Session-ID")
	}

	// If no session exists, create new one
	if sessionID == "" {
		sessionID = models.CreateSession()
		// Set session cookie
		ctx.Output.Cookie("session_id", sessionID)
	} else {
		// Validate and refresh existing session
		if err := models.ValidateSession(sessionID); err != nil {
			sessionID = models.CreateSession()
			ctx.Output.Cookie("session_id", sessionID)
		}
	}

	// Add session ID to response headers
	ctx.Output.Header("X-Session-ID", sessionID)
}
