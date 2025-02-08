package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

// Upload represents a file upload in the system
type Upload struct {
	Id          string    `orm:"pk;column(id)" json:"id"`
	FileName    string    `orm:"column(file_name)" json:"file_name"`
	FilePath    string    `orm:"column(file_path);unique" json:"file_path"`
	FileSize    int64     `orm:"column(file_size)" json:"file_size"`
	ContentType string    `orm:"column(content_type)" json:"content_type"`
	UserID      string    `orm:"column(user_id);size(36)" json:"user_id"`
	CreatedAt   time.Time `orm:"auto_now_add;type(timestamp with time zone);column(created_at)" json:"created_at"`
}

// TableName specifies the database table name
func (u *Upload) TableName() string {
	return "uploads"
}

// EnsureUploadsTable creates the uploads table if it doesn't exist
func EnsureUploadsTable() error {
	sql := `
	CREATE TABLE IF NOT EXISTS uploads (
		id VARCHAR(36) PRIMARY KEY,
		file_name VARCHAR(255) NOT NULL,
		file_path VARCHAR(255) UNIQUE NOT NULL,
		file_size BIGINT NOT NULL,
		content_type VARCHAR(100),
		user_id VARCHAR(36) NOT NULL REFERENCES users(id),
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	)`

	// Drop existing table if it exists
	dropSQL := `DROP TABLE IF EXISTS uploads`
	o := orm.NewOrm()

	// Drop and recreate
	_, err := o.Raw(dropSQL).Exec()
	if err != nil {
		return err
	}

	_, err = o.Raw(sql).Exec()
	return err
}

// CreateUpload creates a new upload record
func CreateUpload(upload *Upload) error {
	o := orm.NewOrm()
	// Use Insert instead of InsertOrUpdate since we're providing the UUID
	_, err := o.Raw(`INSERT INTO uploads (id, file_name, file_path, file_size, content_type, user_id, created_at) 
					 VALUES (?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)`,
		upload.Id, upload.FileName, upload.FilePath, upload.FileSize, upload.ContentType, upload.UserID).Exec()
	return err
}
