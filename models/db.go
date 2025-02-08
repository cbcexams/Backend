package models

import (
	"fmt"
	"cbc-backend/config"

	"github.com/beego/beego/v2/client/orm"
	_ "github.com/lib/pq" // PostgreSQL driver
)



// InitDB initializes the database connection and creates tables
func InitDB(connStr string) error {
	if connStr == "" {
		connStr = config.GetDBConnString()
	}

	// Register database driver
	orm.RegisterDriver("postgres", orm.DRPostgres)

	// Register default database
	err := orm.RegisterDataBase("default", "postgres", connStr)
	if err != nil {
		return err
	}

	return nil
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
