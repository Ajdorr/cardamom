package migrations

import (
	"embed"

	"cardamom/core/source/ext"

	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/peterldowns/pgtestdb/migrators/golangmigrator"
)

//go:embed versions/*.sql
var migrations embed.FS

func GetDatabaseMigration() *golangmigrator.GolangMigrator {
	return golangmigrator.New("versions", golangmigrator.WithFS(migrations))
}

func GetSourceDriver() source.Driver {
	d, err := iofs.New(migrations, "versions")
	ext.PanicIfError(err)
	return d
}
