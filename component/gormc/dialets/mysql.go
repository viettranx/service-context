package dialets

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// MySqlDB Get MySQL DB connection
// dsn string
// Ex: user:password@/db_name?charset=utf8&parseTime=True&loc=Local
func MySqlDB(dsn string) (db *gorm.DB, err error) {
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}
