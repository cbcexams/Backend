package models

import (
	"fmt"
	"strings"
	"sync"

	"cbc-backend/config"

	"github.com/beego/beego/v2/client/orm"
	_ "github.com/lib/pq" // PostgreSQL driver
)

var (
	initOnce sync.Once
	initErr  error
)

// InitDB initializes the database connection and creates tables
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
		connStr := fmt.Sprintf(
			"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
			config.DBUser,
			config.DBPassword,
			config.DBName,
			config.DBHost,
			config.DBPort,
		)
		err := orm.RegisterDataBase("default", "postgres", connStr)
		if err != nil {
			// Ignore if already registered
			if !strings.Contains(err.Error(), "already registered") {
				initErr = fmt.Errorf("failed to register database: %v", err)
				return
			}
		}

		// Ensure tables exist
		if err := EnsureUsersTable(); err != nil {
			initErr = fmt.Errorf("failed to create users table: %v", err)
			return
		}

		if err := EnsureUploadsTable(); err != nil {
			initErr = fmt.Errorf("failed to create uploads table: %v", err)
			return
		}

		fmt.Println("Database connection successful!")
	})

	return initErr
}

func TestDatabaseConnection() error {
	o := orm.NewOrm()
	var count int64
	err := o.Raw("SELECT COUNT(*) FROM uploads").QueryRow(&count)
	if err != nil {
		return fmt.Errorf("database test query failed: %v", err)
	}
	fmt.Printf("✅ Database connected successfully\n")
	fmt.Printf("✅ Found %d uploads in uploads table\n", count)
	return nil
}
