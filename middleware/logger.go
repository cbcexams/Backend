package middleware

import (
	"fmt"
	"time"

	"github.com/beego/beego/v2/server/web/context"
)

// LoggerMiddleware logs information about each request
func LoggerMiddleware(ctx *context.Context) {
	// Record start time
	start := time.Now()

	// Get request details
	method := ctx.Input.Method()
	path := ctx.Input.URL()
	ip := ctx.Input.IP()

	// Print request start
	fmt.Printf("\n[%s] %s %s - Started\n", method, path, ip)

	// Continue processing
	ctx.Input.SetData("request_start", start)

	// Calculate duration after the request is processed
	defer func() {
		duration := time.Since(start)

		// Get response status
		status := ctx.ResponseWriter.Status

		// Print complete log
		fmt.Printf("[%s] %s %d - %v\n", method, path, status, duration)

		// Log additional details in debug mode
		if ctx.Input.Query("debug") == "true" {
			// Get request headers
			headers := ctx.Request.Header
			fmt.Printf("Headers: %+v\n", headers)

			// Get query parameters
			queryValues := ctx.Request.URL.Query()
			fmt.Printf("Query: %+v\n", queryValues)
		}
	}()
}
