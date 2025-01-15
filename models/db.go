package models

import (
	"fmt"

	"github.com/beego/beego/v2/client/orm"
	_ "github.com/lib/pq"
)

func init() {
	// Register driver
	err := orm.RegisterDriver("postgres", orm.DRPostgres)
	if err != nil {
		fmt.Printf("Failed to register driver: %v\n", err)
		return
	}

	// Register default database
	connStr := "user=postgres password=0000 dbname=cbcexams sslmode=disable"
	err = orm.RegisterDataBase("default", "postgres", connStr)
	if err != nil {
		fmt.Printf("Failed to register database: %v\n", err)
		return
	}

	// Register models
	orm.RegisterModel(
		new(User),
		new(Resource),
		new(Job),
	)

	// Test database connection
	o := orm.NewOrm()
	var result int
	err = o.Raw("SELECT 1").QueryRow(&result)
	if err != nil {
		fmt.Printf("Database connection test failed: %v\n", err)
		return
	}
	fmt.Println("Database connection successful!")

	// Create tables
	err = orm.RunSyncdb("default", false, true)
	if err != nil {
		fmt.Printf("Failed to sync database: %v\n", err)
		return
	}
}
