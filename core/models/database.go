package models

import (
	cfg "cardamom/core/config"
	"cardamom/core/ext/log_ext"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {

	if cfg.IsLocal() {
		if db, err := gorm.Open(sqlite.Open(cfg.C.DB_Sqlite), &gorm.Config{}); err != nil {
			panic(log_ext.Errorf("failed to connect to database -- %w", err))
		} else {
			DB = db
		}
	} else if db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			cfg.C.DB_Host, cfg.C.DB_Port, cfg.C.DB_Username, cfg.C.DB_Password, cfg.C.DB_Name),
		PreferSimpleProtocol: true,
	}), &gorm.Config{}); err != nil {
		panic(log_ext.Errorf("failed to connect to database -- %w", err))
	} else {
		DB = db
	}

}
