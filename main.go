package main

import (
	_ "cbc-backend/routers"
	"fmt"
	"os"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	// Create uploads directory if it doesn't exist
	if err := os.MkdirAll("uploads", 0755); err != nil {
		logs.Error("Failed to create uploads directory:", err)
	}
}

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	// Configure CORS
	beego.InsertFilter("*", beego.BeforeRouter, func(ctx *beego.Context) {
		origin := ctx.Input.Header("Origin")
		if origin != "" {
			ctx.Output.Header("Access-Control-Allow-Origin", origin)
			ctx.Output.Header("Access-Control-Allow-Methods",
				beego.AppConfig.String("CORSAllowMethods"))
			ctx.Output.Header("Access-Control-Allow-Headers",
				beego.AppConfig.String("CORSAllowHeaders"))
			ctx.Output.Header("Access-Control-Allow-Credentials", "true")
		}

		// Handle preflight requests
		if ctx.Input.Method() == "OPTIONS" {
			ctx.Output.SetStatus(200)
			ctx.ResponseWriter.WriteHeader(200)
			return
		}
		ctx.ResponseWriter.WriteHeader(200)
	})

	fmt.Println("Starting server...")
	fmt.Printf("RunMode: %s\n", beego.BConfig.RunMode)
	beego.Run()
}
