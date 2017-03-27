package main

import (
	"log"
	"sync"
	"time"
)

type EventType string

const (
	EventTimer    EventType = "eTimer"
	EventLog      EventType = "eLog"
	EventTick     EventType = "eTick"
	EventTrade    EventType = "eTrade"
	EventOrder    EventType = "eOrder"
	EventPosition EventType = "ePostion"
	EventAccount  EventType = "eAccount"
	EventContract EventType = "eContract"
	EventError    EventType = "eError"

	EventCTALog      EventType = "eCTALog"
	EventCTAStrategy EventType = "eCTAStrategy"

	EventDataRecorderLog EventType = "eDataRecorderLog"
)

type Event struct {
	Type EventType
	Dict map[EventType]interface{}
}

type Handler func(event *Event) error

type EventEngine struct {
	eventChan       chan *Event
	handlers        map[EventType][]*Handler
	generalHandlers map[*Handler]bool
	active          bool
	sync.Mutex
}

func NewEventEngine() *EventEngine {
	engine := &EventEngine{
		eventChan: make(chan *Event, 20480),
	}
	return engine
}

func (e *EventEngine) run() {
	ticker := time.NewTicker(time.Second)
	for e.active {
		select {
		case evt := <-e.eventChan:
			e.process(evt)
		case <-ticker.C:
			log.Println("ticker Now")
		}
	}
}

func (e *EventEngine) process(event *Event) {
	if handlers, ok := e.handlers[event.Type]; ok {
		for _, handler := range handlers {
			(*handler)(event)
		}
	}
}

func (e *EventEngine) Start() {
	e.active = true
	go e.run()
}

func (e *EventEngine) Stop() {
	e.active = false
}

func (e *EventEngine) Register(type_ EventType, handler *Handler) {

}

func (e *EventEngine) Unregister(type_ EventType, handler *Handler) {

}

func (e *EventEngine) Put(event *Event) {

}

func (e *EventEngine) RegisterGeneralHandler(handler *Handler) {

}

func (e *EventEngine) UnregisterGeneralHandler(handler *Handler) {

}
