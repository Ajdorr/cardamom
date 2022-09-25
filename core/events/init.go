package events

import (
	cfg "cardamom/core/config"
	"sync"
)

func Shutdown() {
	close(eventChannel)
	doneSync.Wait()
}

var eventChannel chan *Event
var eventStreams []EventStream
var doneSync sync.WaitGroup

func eventStreamLoop() {
	defer doneSync.Done()

	for evt := range eventChannel {
		for _, es := range eventStreams {
			es.Publish(evt)
		}
	}

	for _, es := range eventStreams {
		es.Flush()
	}
}

func init() {
	eventChannel = make(chan *Event, 32)

	if len(cfg.C.EventFileStreamDirectory) > 0 {
		eventStreams = append(eventStreams, &FileEventStream{
			EventDirectory:    cfg.C.EventFileStreamDirectory,
			UnpublishedEvents: make([]*Event, 0, event_buffer_size),
		})
		doneSync.Add(1)
	}

	go eventStreamLoop()
}
