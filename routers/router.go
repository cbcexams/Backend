// @APIVersion 1.0.0
// @Title Teaching App API
// @Description Teaching resources and jobs API
package routers

import (
	"cbc-backend/controllers"

	"fmt"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	fmt.Println("Initializing routes...")

	// Register direct routes
	beego.Router("/v1/resources", &controllers.ResourceController{}, "get:Get;post:Post")

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

	// Print registered routes
	fmt.Println("Registered routes:")
	fmt.Println("- GET, POST /v1/resources")
	fmt.Println("- GET, POST /v1/jobs")
	fmt.Println("- GET, POST /v1/user")

	fmt.Println("Routes initialized")
}
