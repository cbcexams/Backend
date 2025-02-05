// @APIVersion 1.0.0
// @Title Teaching App API
// @Description Teaching resources and jobs API
package routers

import (
	"cbc-backend/controllers"
	"cbc-backend/middleware"

	"fmt"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
)

func init() {
	fmt.Println("\n==================================================")
	fmt.Println("              Router Initialization                ")
	fmt.Println("==================================================")

	// Create a namespace for v1 API
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/resources",
			beego.NSRouter("/", &controllers.ResourceController{}, "get:Get"),
			beego.NSRouter("/", &controllers.ResourceController{}, "post:Post"),
		),
		beego.NSNamespace("/user",
			beego.NSRouter("/signup", &controllers.UserController{}, "post:Post"),
			beego.NSRouter("/login", &controllers.UserController{}, "post:Login"),
		),
	)

	// Add namespace to beego
	beego.AddNamespace(ns)

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

	// Add logger middleware for all routes
	beego.InsertFilter("/*", beego.BeforeRouter, middleware.LoggerMiddleware)

	// Print registered routes
	fmt.Println("\nRegistered Routes:")
	fmt.Printf("GET  /v1/resources?[params] (public)\n")
	fmt.Printf("POST /v1/resources (protected)\n")
	fmt.Printf("POST /v1/user/signup\n")
	fmt.Printf("POST /v1/user/login\n")

	fmt.Println("\nRouter initialization complete")
	fmt.Println("==================================================")
}
