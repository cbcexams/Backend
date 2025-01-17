package middleware

import (
	"github.com/beego/beego/v2/server/web/context"
)

// JWTMiddleware handles JWT authentication
func JWTMiddleware(ctx *context.Context) {
	// For now, just pass through since we're debugging the resource endpoint
	// We'll implement proper JWT validation later
}
