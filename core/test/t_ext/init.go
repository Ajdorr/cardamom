package t_ext

import (
	cfg "cardamom/core/source/config"
	"cardamom/core/source/db"
	"cardamom/core/source/db/migrations"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/lib/pq"
	"github.com/peterldowns/pgtestdb"
)

func Init(t *testing.T) {
	newTestDB(t)
	EnsureTestUser()
}

func newTestDB(t pgtestdb.TB) {
	conf := pgtestdb.Config{
		DriverName: "postgres",
		Host:       cfg.C.DB.Host,
		Port:       cfg.C.DB.Port,
		User:       cfg.C.DB.Username,
		Password:   cfg.C.DB.Password,
		Database:   cfg.C.DB.Database,
		Options:    "sslmode=disable",
	}
	conn := pgtestdb.New(t, conf, migrations.GetDatabaseMigration())

	db.ConnectWith(conn)

}

func init() {
	if err := cfg.LoadConfig("testing.yaml"); err != nil {
		panic(err)
	}
}
