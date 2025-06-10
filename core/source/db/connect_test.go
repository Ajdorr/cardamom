package db_test

import (
	"cardamom/core/source/config"
	"cardamom/core/source/db"
	"cardamom/core/test/t_ext"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	t_ext.Init(t)
	db.Connect()
}

func TestConnectFailure(t *testing.T) {
	config.C.DB.Host = "invalid-host"

	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				assert.Contains(t, err.Error(), "failed to open gorm database connection")
			} else {
				t.Fail()
			}
		} else {
			t.Fail()
		}
	}()

	db.Connect()
}
