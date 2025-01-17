package middleware

import (
	"fmt"
	"time"

	"github.com/beego/beego/v2/server/web/context"
)

// LoggerMiddleware logs request details
func LoggerMiddleware(ctx *context.Context) {
	start := time.Now()
	path := ctx.Request.URL.Path
	method := ctx.Request.Method

	defer func() {
		fmt.Printf("[%s] %s %d - %v\n",
			method,
			path,
			ctx.ResponseWriter.Status,
			time.Since(start),
		)
	}()
}
