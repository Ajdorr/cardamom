package models_test

import (
	m "cardamom/core/source/db/models"
	"cardamom/core/test/t_ext"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	t_ext.Init(t)

	u := m.User{}
	u.BeforeCreate(nil)
	assert.NotEmpty(t, u.Uid)

	u.SetPassword("abcd")
	assert.NotEmpty(t, u.Password)
	assert.NotEqual(t, "abcd", u.Password)

	missing, err := m.GetUserByEmail("invalid-email")
	assert.Nil(t, missing)
	assert.Equal(t, "record not found", err.Error())
}
