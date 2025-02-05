package models

import (
	"github.com/beego/beego/v2/client/orm"
)

func init() {
	// Register models
	orm.RegisterModel(
		new(User),
		new(Upload),
		new(Resource), // For reading only
	)
}
