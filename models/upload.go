package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

// Upload represents a file upload in the system
type Upload struct {
	Id          string    `orm:"pk;column(id)" json:"id"`
	FileName    string    `orm:"size(255)" json:"file_name"`
	FilePath    string    `orm:"size(255);unique" json:"file_path"`
	FileSize    int64     `orm:"" json:"file_size"`
	ContentType string    `orm:"size(100)" json:"content_type"`
	UserID      int       `orm:"column(user_id)" json:"user_id"`
	CreatedAt   time.Time `orm:"auto_now_add;type(timestamp)" json:"created_at"`
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
		user_id INTEGER NOT NULL REFERENCES users(id),
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	)`

	o := orm.NewOrm()
	_, err := o.Raw(sql).Exec()
	return err
}

// CreateUpload creates a new upload record
func CreateUpload(upload *Upload) error {
	o := orm.NewOrm()
	_, err := o.Insert(upload)
	return err
}
