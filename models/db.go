package models

import (
	"fmt"
	"strings"
	"sync"

	"github.com/beego/beego/v2/client/orm"
	_ "github.com/lib/pq" // PostgreSQL driver
)

var (
	initOnce sync.Once
	initErr  error
)

// InitDB initializes the database connection and registers models
func InitDB() error {
	initOnce.Do(func() {
		// Register database driver
		if err := orm.RegisterDriver("postgres", orm.DRPostgres); err != nil {
			// Ignore if already registered
			if !strings.Contains(err.Error(), "already registered") {
				initErr = fmt.Errorf("failed to register driver: %v", err)
				return
			}
		}

		// Try to register database
		connStr := "user=postgres password=0000 dbname=cbcexams sslmode=disable"
		err := orm.RegisterDataBase("default", "postgres", connStr)
		if err != nil {
			// Ignore if already registered
			if !strings.Contains(err.Error(), "already registered") {
				initErr = fmt.Errorf("failed to register database: %v", err)
				return
			}
		}

		// Register all models in one place
		orm.RegisterModel(
			new(Resource),
			// new(User),
			// new(Job),
			// new(Session),
		)

		// Test database connection
		if err := TestDatabaseConnection(); err != nil {
			initErr = fmt.Errorf("database connection test failed: %v", err)
			return
		}

		fmt.Println("Database connection successful!")
	})

	return initErr
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
