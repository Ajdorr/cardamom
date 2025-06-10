package events

import (
	cfg "cardamom/core/source/config"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const EVENT_BUFFER_SIZE = 32

type FileEventStream struct {
	UnpublishedEvents []*Event
}

func (stream *FileEventStream) Publish(event *Event) {

	stream.UnpublishedEvents = append(stream.UnpublishedEvents, event)
	if len(stream.UnpublishedEvents) < EVENT_BUFFER_SIZE {
		return
	}

	stream.Flush()
}

func (stream *FileEventStream) Flush() {

	if err := os.MkdirAll(cfg.C.Events.EventFileStreamDirectory, os.ModePerm); err != nil {
		panic(err)
	}
	// Open the file
	filename := filepath.Join(
		cfg.C.Events.EventFileStreamDirectory,
		fmt.Sprintf("events-%s.txt", time.Now().UTC().Format("2006-01-02")))
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

	stream.UnpublishedEvents = make([]*Event, 0, EVENT_BUFFER_SIZE)
}

func NewFileEventStream() *FileEventStream {
	return &FileEventStream{
		UnpublishedEvents: make([]*Event, 0, EVENT_BUFFER_SIZE),
	}
}
