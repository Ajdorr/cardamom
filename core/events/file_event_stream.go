package events

import (
	cfg "cardamom/core/config"
	"fmt"
	"os"
	"time"
)

const event_buffer_size = 32

type FileEventStream struct {
	EventDirectory    string
	UnpublishedEvents []*Event
}

func (stream *FileEventStream) Publish(event *Event) {

	stream.UnpublishedEvents = append(stream.UnpublishedEvents, event)
	if len(stream.UnpublishedEvents) < event_buffer_size {
		return
	}

	stream.Flush()
}

func (stream *FileEventStream) Flush() {

	// Open the file
	filename := fmt.Sprintf(
		"%s/events-%s.txt", cfg.C.EventFileStreamDirectory, time.Now().UTC().Format("2006-01-02"))
	file, err := os.OpenFile(
		filename,
		os.O_APPEND|os.O_WRONLY|os.O_CREATE,
		0666,
	)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	for _, evt := range stream.UnpublishedEvents {
		if _, err := file.Write(evt.getBytes()); err != nil {
			panic(err)
		}
		if _, err := file.WriteString("\n"); err != nil {
			panic(err)
		}
	}

	stream.UnpublishedEvents = make([]*Event, 0, event_buffer_size)
}
