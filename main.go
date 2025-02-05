package main

import (
	_ "cbc-backend/routers"
	"fmt"
	"os"

	"cbc-backend/models"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
)

// init initializes the application by setting up required directories
func init() {
	// Create uploads directory if it doesn't exist
	// This directory is used to store uploaded resource files
	if err := os.MkdirAll("uploads", 0755); err != nil {
		logs.Error("Failed to create uploads directory:", err)
	}

	// Initialize database connection
	if err := models.InitDB(); err != nil {
		logs.Error("Failed to initialize database:", err)
		os.Exit(1)
	}
}

func main() {
	// Print application startup banner
	fmt.Println("\n==================================================")
	fmt.Println("                 Application Start                  ")
	fmt.Println("==================================================")

	// Test database connection
	if err := testDatabaseConnection(); err != nil {
		fmt.Printf("❌ Database connection error: %v\n", err)
		os.Exit(1)
	}

	// Configure development mode settings
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"

		// Configure admin server only in dev mode
		beego.BConfig.Listen.EnableAdmin = true
		beego.BConfig.Listen.AdminAddr = "localhost"
		beego.BConfig.Listen.AdminPort = 8088
	} else {
		// Disable admin server in non-dev modes
		beego.BConfig.Listen.EnableAdmin = false
	}

	// Configure main server
	beego.BConfig.Listen.HTTPAddr = "localhost"
	beego.BConfig.Listen.HTTPPort = 8081

	// Configure CORS middleware
	// This allows cross-origin requests from frontend applications
	beego.InsertFilter("*", beego.BeforeRouter, func(ctx *context.Context) {
		origin := ctx.Input.Header("Origin")
		if origin != "" {
			// Set CORS headers to allow cross-origin requests
			ctx.Output.Header("Access-Control-Allow-Origin", origin)
			ctx.Output.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
			ctx.Output.Header("Access-Control-Allow-Headers", "Origin,Content-Type,Accept,Authorization")
			ctx.Output.Header("Access-Control-Allow-Credentials", "true")
		}

		// Handle preflight OPTIONS requests
		if ctx.Input.Method() == "OPTIONS" {
			ctx.Output.SetStatus(200)
			return
		}
	})

	// Start the server
	fmt.Println("\nStarting server...")
	fmt.Printf("RunMode: %s\n", beego.BConfig.RunMode)
	fmt.Printf("HTTP Server: http://%s:%d\n", beego.BConfig.Listen.HTTPAddr, beego.BConfig.Listen.HTTPPort)
	if beego.BConfig.Listen.EnableAdmin {
		fmt.Printf("Admin Server: http://%s:%d\n", beego.BConfig.Listen.AdminAddr, beego.BConfig.Listen.AdminPort)
	}

	beego.Run()
}

func testDatabaseConnection() error {
	o := orm.NewOrm()
	var count int64
	err := o.Raw("SELECT COUNT(*) FROM web_crawler_resources").QueryRow(&count)
	if err != nil {
		return fmt.Errorf("database test failed: %v", err)
	}
	fmt.Printf("✅ Database connected successfully\n")
	fmt.Printf("✅ Found %d resources in web_crawler_resources table\n", count)
	return nil
}
