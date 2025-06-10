package main

import (
	_ "cardamom/core/source/config"
	"cardamom/core/source/db"
	"cardamom/core/source/db/migrations"
	"cardamom/core/source/ext"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
)

func main() {
	print("Running migration\n")
	db.Connect()
	dbDriver, err := postgres.WithInstance(db.Conn(), &postgres.Config{})
	ext.PanicIfError(err)

	sourceDriver := migrations.GetSourceDriver()
	m, err := migrate.NewWithInstance("iofs", sourceDriver, "postgres", dbDriver)
	ext.PanicIfError(err)

	ext.PanicIfError(m.Up())
	print("Migration successful\n")
}
