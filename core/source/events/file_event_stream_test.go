package events_test

import (
	"cardamom/core/source/config"
	"cardamom/core/source/events"
	_ "cardamom/core/test/t_ext"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileEventStream(t *testing.T) {

	es := events.NewFileEventStream()

	os.MkdirAll(config.C.Events.EventFileStreamDirectory, os.ModePerm)
	files, err := os.ReadDir(config.C.Events.EventFileStreamDirectory)
	if err != nil {
		t.Fail()
	}
	for _, f := range files {
		if f.Type().IsRegular() && strings.HasPrefix(f.Name(), "events-") && strings.HasPrefix(f.Name(), "events-") {
			os.Remove(filepath.Join(config.C.Events.EventFileStreamDirectory, f.Name()))
		}
	}

	evt := &events.Event{
		Domain: "test",
		Type:   "test",
		Data:   nil,
	}

	es.Publish(evt)
	files, err = os.ReadDir(config.C.Events.EventFileStreamDirectory)
	if err != nil {
		t.Fail()
	}
	assert.Zero(t, len(files))

	for range events.EVENT_BUFFER_SIZE {
		es.Publish(evt)
	}

	files, err = os.ReadDir(config.C.Events.EventFileStreamDirectory)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 1, len(files))

	data, err := os.ReadFile(filepath.Join(config.C.Events.EventFileStreamDirectory, files[0].Name()))
	if err != nil {
		t.Fail()
	}

	evts := strings.Split(string(data), "\n")
	for _, evtRaw := range evts {
		if len(evtRaw) == 0 {
			continue
		}

		if json.Unmarshal([]byte(evtRaw), &events.Event{}) != nil {
			t.Fail()
		}
	}
}
