package events

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBytes(t *testing.T) {
	e := Event{}
	assert.Equal(t, `{"domain":"","type":"","id":"","data":null}`, string(e.getBytes()))
}

func TestShutdown(t *testing.T) {
	Shutdown()
}
