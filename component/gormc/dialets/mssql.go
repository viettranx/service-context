package dialets

import (
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

// MSSqlDB Get MS SQL DB connection
// dsn string
// Ex: sqlserver://username:password@localhost:1433?database=dbname
func MSSqlDB(dsn string) (db *gorm.DB, err error) {
	return gorm.Open(sqlserver.Open(dsn))
}
