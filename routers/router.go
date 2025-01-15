// @APIVersion 1.0.0
// @Title Teaching App API
// @Description Teaching resources and jobs API
package routers

import (
	"cbc-backend/controllers"
	"cbc-backend/middleware"

	"fmt"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	fmt.Println("Initializing routes...")

	// Add JWT middleware
	beego.InsertFilter("/v1/*", beego.BeforeRouter, middleware.JWT)

	// Register direct routes
	beego.Router("/v1/resources", &controllers.ResourceController{}, "get:Get;post:Post")
	beego.Router("/v1/jobs", &controllers.JobController{}, "get:Get;post:Post")

	// Register namespace routes
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/resources",
			beego.NSInclude(
				&controllers.ResourceController{},
			),
		),
		beego.NSNamespace("/jobs",
			beego.NSInclude(
				&controllers.JobController{},
			),
		),
		beego.NSNamespace("/user",
			beego.NSInclude(
				&controllers.UserController{},
			),
		),
	)
	beego.AddNamespace(ns)

	fmt.Println("Routes initialized")
}
