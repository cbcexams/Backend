package models

import (
	"fmt"
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
			initErr = fmt.Errorf("failed to register driver: %v", err)
			return
		}

		// Register database
		err := orm.RegisterDataBase("default", "postgres",
			"postgres://postgres:postgres@localhost:5432/cbcexams?sslmode=disable")
		if err != nil {
			initErr = fmt.Errorf("failed to register database: %v", err)
			return
		}

		// Register models
		orm.RegisterModel(new(Resource))
	})

	return initErr
}
