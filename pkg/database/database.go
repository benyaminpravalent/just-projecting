package database

import (
	"context"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"

	"github.com/jmoiron/sqlx"
	"github.com/mine/just-projecting/pkg/config"
	"github.com/mine/just-projecting/pkg/logger"
)

var DB *sqlx.DB

func InitMySql(ctx context.Context) {
	l := logger.GetLoggerContext(ctx, "database", "Connect")

	dsn := config.GetString("mysql_dsn")
	l.Info(dsn)

	dbConnection, err := sqlx.Open("mysql", dsn)
	if err != nil {
		panic(err.Error())
	}

	err = dbConnection.Ping()
	if err != nil {
		panic(err.Error())
	}

	l.Info("Connected to MySQL")

	DB = dbConnection
}
