package gormc

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	sctx "github.com/viettranx/service-context"
	"github.com/viettranx/service-context/component/gormc/dialets"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type GormDBType int

const (
	GormDBTypeMySQL GormDBType = iota + 1
	GormDBTypePostgres
	GormDBTypeSQLite
	GormDBTypeMSSQL
	GormDBTypeNotSupported
)

type GormOpt struct {
	dsn                   string
	dbType                string
	maxOpenConnections    int
	maxIdleConnections    int
	maxConnectionIdleTime int
}

type gormDB struct {
	id     string
	prefix string
	logger sctx.Logger
	db     *gorm.DB
	*GormOpt
}

func NewGormDB(id, prefix string) *gormDB {
	return &gormDB{
		GormOpt: new(GormOpt),
		id:      id,
		prefix:  strings.TrimSpace(prefix),
	}
}

func (gdb *gormDB) ID() string {
	return gdb.id
}

func (gdb *gormDB) InitFlags() {
	prefix := gdb.prefix
	if gdb.prefix != "" {
		prefix += "-"
	}

	flag.StringVar(
		&gdb.dsn,
		fmt.Sprintf("%sdb-dsn", prefix),
		"",
		"Database dsn",
	)

	flag.StringVar(
		&gdb.dbType,
		fmt.Sprintf("%sdb-driver", prefix),
		"mysql",
		"Database driver (mysql, postgres, sqlite, mssql) - Default mysql",
	)

	flag.IntVar(
		&gdb.maxOpenConnections,
		fmt.Sprintf("%sdb-max-conn", prefix),
		30,
		"maximum number of open connections to the database - Default 30",
	)

	flag.IntVar(
		&gdb.maxIdleConnections,
		fmt.Sprintf("%sdb-max-ide-conn", prefix),
		10,
		"maximum number of database connections in the idle - Default 10",
	)

	flag.IntVar(
		&gdb.maxConnectionIdleTime,
		fmt.Sprintf("%sdb-max-conn-ide-time", prefix),
		3600,
		"maximum amount of time a connection may be idle in seconds - Default 3600",
	)
}

func (gdb *gormDB) isDisabled() bool {
	return gdb.dsn == ""
}

func (gdb *gormDB) Activate(_ sctx.ServiceContext) error {
	gdb.logger = sctx.GlobalLogger().GetLogger(gdb.id)

	dbType := getDBType(gdb.dbType)
	if dbType == GormDBTypeNotSupported {
		return errors.WithStack(errors.New("Database type not supported."))
	}

	gdb.logger.Info("Connecting to database...")

	var err error
	gdb.db, err = gdb.getDBConn(dbType)

	if err != nil {
		gdb.logger.Error("Cannot connect to database", err.Error())
		return err
	}

	return nil
}

func (gdb *gormDB) Stop() error {
	/**
	From version 1.20, Jinzhu (GORM author) decided eliminate Close() method because GORM supports connection pooling.
	And best practice is open connection once and reuse it until the program exits.

	But some cases, we need to close connection manually. For that case, we can use this method.
	**/
	db, err := gdb.db.DB()
	if err != nil {
		return err
	}
	return db.Close()
}

func (gdb *gormDB) GetDB() *gorm.DB {
	if gdb.logger.GetLevel() == "debug" || gdb.logger.GetLevel() == "trace" {
		return gdb.db.Session(&gorm.Session{NewDB: true}).Debug()
	}

	newSessionDB := gdb.db.Session(&gorm.Session{NewDB: true, Logger: gdb.db.Logger.LogMode(logger.Silent)})

	if db, err := newSessionDB.DB(); err == nil {
		db.SetMaxOpenConns(gdb.maxOpenConnections)
		db.SetMaxIdleConns(gdb.maxIdleConnections)
		db.SetConnMaxIdleTime(time.Second * time.Duration(gdb.maxConnectionIdleTime))
	}

	return newSessionDB
}

func getDBType(dbType string) GormDBType {
	switch strings.ToLower(dbType) {
	case "mysql":
		return GormDBTypeMySQL
	case "postgres":
		return GormDBTypePostgres
	case "sqlite":
		return GormDBTypeSQLite
	case "mssql":
		return GormDBTypeMSSQL
	}

	return GormDBTypeNotSupported
}

func (gdb *gormDB) getDBConn(t GormDBType) (dbConn *gorm.DB, err error) {
	switch t {
	case GormDBTypeMySQL:
		return dialets.MySqlDB(gdb.dsn)
	case GormDBTypePostgres:
		return dialets.PostgresDB(gdb.dsn)
	case GormDBTypeSQLite:
		return dialets.SQLiteDB(gdb.dsn)
	case GormDBTypeMSSQL:
		return dialets.MSSqlDB(gdb.dsn)
	}

	return nil, nil
}
