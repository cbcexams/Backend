package main

import (
	_ "cbc-backend/routers"
	"fmt"

	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	fmt.Println("Starting server...")
	fmt.Printf("RunMode: %s\n", beego.BConfig.RunMode)

	beego.Run()
}
