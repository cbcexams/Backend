// @APIVersion 1.0.0
// @Title CBC-Backend API
// @Description CBC-Backend API
package routers

import (
	"cbc-backend/controllers"
	"cbc-backend/middleware"

	"fmt"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
)

func init() {
	fmt.Println("\n=== Route Configuration ===")

	// Add error handling middleware first
	beego.InsertFilter("/*", beego.BeforeRouter, middleware.ErrorMiddleware)

	// Configure routes
	configureUserRoutes()
	configureResourceRoutes()
	configureJobRoutes()

	// Print route summary
	fmt.Println("\nAvailable Endpoints:")
	fmt.Printf("Auth:\n")
	fmt.Printf("  POST /v1/user/signup\n")
	fmt.Printf("  POST /v1/user/login\n")
	fmt.Printf("  GET  /v1/user/logout\n")
	fmt.Printf("  POST /v1/user/forgot-password\n")
	fmt.Printf("  POST /v1/user/reset-password\n")
	fmt.Printf("  DELETE /v1/user/:uid [Protected]\n")
	fmt.Printf("\nResources:\n")
	fmt.Printf("  GET  /v1/resources? [params] [public]\n")
	fmt.Printf("  POST /v1/resources [Protected]\n")
	fmt.Printf("\nJobs:\n")
	fmt.Printf("  GET  /v1/jobs?[params] [Protected]\n")
	fmt.Printf("  POST /v1/jobs [Protected]\n")

	fmt.Println("\n=== Route Configuration Complete ===")
}

func configureUserRoutes() {
	// User routes
	beego.Router("/v1/user/signup", &controllers.UserController{}, "post:Post")
	beego.Router("/v1/user/login", &controllers.UserController{}, "post:Login")
	beego.Router("/v1/user/logout", &controllers.UserController{}, "get:Logout")
	beego.Router("/v1/user/reset-password", &controllers.UserController{}, "post:ResetPassword")
	beego.Router("/v1/user/forgot-password", &controllers.UserController{}, "post:ForgotPassword")
	beego.Router("/v1/user/:uid", &controllers.UserController{}, "delete:Delete")

	// Add JWT middleware only for POST /v1/resources
	beego.InsertFilter("/v1/resources", beego.BeforeRouter, func(ctx *context.Context) {
		// Skip JWT check for GET requests
		if ctx.Input.Method() == "GET" {
			return
		}
		if ctx.Input.Method() == "POST" {
			middleware.JWTMiddleware(ctx)
		}
	})

	// Add JWT middleware for POST /v1/jobs
	beego.InsertFilter("/v1/jobs", beego.BeforeRouter, func(ctx *context.Context) {
		// Protect both GET and POST methods
		if ctx.Input.Method() == "GET" || ctx.Input.Method() == "POST" {
			middleware.JWTMiddleware(ctx)
		}
	})

	// Add logger middleware for all routes
	beego.InsertFilter("/*", beego.BeforeRouter, middleware.LoggerMiddleware)

	// Add JWT middleware for delete user endpoint
	beego.InsertFilter("/v1/user/*", beego.BeforeRouter, func(ctx *context.Context) {
		if ctx.Input.Method() == "DELETE" {
			middleware.JWTMiddleware(ctx)
		}
	})
}

func configureResourceRoutes() {
	// Resource routes
	beego.Router("/v1/resources", &controllers.ResourceController{})
}

func configureJobRoutes() {
	// Jobs routes
	beego.Router("/v1/jobs", &controllers.JobController{})
}
