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

	// Register all models in one place
	orm.RegisterModel(
		new(Resource),
		new(User),
		new(Job),
		new(Session),
	)

	// Register default database
	connStr := "user=postgres password=0000 dbname=cbcexams sslmode=disable"
	err = orm.RegisterDataBase("default", "postgres", connStr)
	if err != nil {
		fmt.Printf("Failed to register database: %v\n", err)
		return
	}

	// Test database connection
	o := orm.NewOrm()
	var result int
	err = o.Raw("SELECT 1").QueryRow(&result)
	if err != nil {
		fmt.Printf("Database connection test failed: %v\n", err)
		return
	}
	fmt.Println("Database connection successful!")
}

func TestDatabaseConnection() error {
	o := orm.NewOrm()

	var count int64
	err := o.Raw("SELECT COUNT(*) FROM web_crawler_resources").QueryRow(&count)
	if err != nil {
		return fmt.Errorf("database test query failed: %v", err)
	}

	fmt.Printf("Database connection test: found %d resources\n", count)
	return nil
}
