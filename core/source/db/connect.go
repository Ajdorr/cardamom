package db

import (
	cfg "cardamom/core/source/config"
	"database/sql"

	"cardamom/core/source/ext/log_ext"

	_ "ariga.io/atlas-provider-gorm/gormschema"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB
var conn *sql.DB

func Conn() *sql.DB {
	return conn
}

func DB() *gorm.DB {
	return db
}

func ConnectWith(conn gorm.ConnPool) {
	if _db, err := gorm.Open(postgres.New(postgres.Config{Conn: conn}), &gorm.Config{}); err != nil {
		panic(log_ext.Errorf("failed to open gorm database connection -- %w", err))
	} else {
		db = _db
	}
}

func Connect() {
	if _conn, err := sql.Open("pgx", cfg.C.GetDSN()); err != nil {
		panic(log_ext.Errorf("failed to connect to database -- %w", err))
	} else {
		conn = _conn
		ConnectWith(_conn)
	}
}
