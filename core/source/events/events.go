package events

import (
	"encoding/json"
	"sync"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

type EventStream interface {
	Publish(event *Event)
	Flush()
}

type Event struct {
	Domain string `json:"domain"`
	Type   string `json:"type"`
	ID     string `json:"id"`
	Data   any    `json:"data"`
}

func (event *Event) GenerateID() string {
	event.ID = gonanoid.Must(32)
	return event.ID
}

func Publish(event *Event) {

	// Automatically generate id
	if len(event.ID) == 0 {
		event.GenerateID()
	}

	eventChannel <- event
}

func (event *Event) getBytes() []byte {
	if payload, err := json.Marshal(event); err != nil {
		panic(err)
	} else {
		return payload
	}
}

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

	eventStreams = append(eventStreams, NewFileEventStream())
	doneSync.Add(1)

	go eventStreamLoop()
}
