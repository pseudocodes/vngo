package eventbus

import (
	"log"
	"sync"
	"time"
)

type EventType string

const (
	EventTimer    EventType = "eTimer"
	EventLog      EventType = "eLog"
	EventTick     EventType = "eTick."
	EventTrade    EventType = "eTrade."
	EventOrder    EventType = "eOrder."
	EventPosition EventType = "ePostion."
	EventAccount  EventType = "eAccount."
	EventContract EventType = "eContract."
	EventError    EventType = "eError."

	EventCTALog      EventType = "eCTALog."
	EventCTAStrategy EventType = "eCTAStrategy."

	EventDataRecorderLog EventType = "eDataRecorderLog."
)

func (et EventType) String() string {
	return string(et)
}

type Event struct {
	Type EventType
	Dict map[EventType]interface{}
}

func NewEvent(eventType EventType) *Event {
	return &Event{
		Type: eventType,
		Dict: make(map[EventType]interface{}),
	}
}

type Handler func(event *Event) error

func (h Handler) Handle(event *Event) error {
	return h(event)
}

type Eventbus struct {
	eventChan       chan *Event
	handlers        map[EventType]map[*Handler]bool
	generalHandlers map[*Handler]bool
	active          bool
	sync.Mutex
}

func NewEventbus() *Eventbus {
	engine := &Eventbus{
		eventChan:       make(chan *Event, 2048),
		handlers:        make(map[EventType]map[*Handler]bool),
		generalHandlers: make(map[*Handler]bool),
		active:          false,
	}
	return engine
}

func (e *Eventbus) run() {
	ticker := time.NewTicker(time.Second)
	for e.active {
		select {
		case evt := <-e.eventChan:
			e.process(evt)
		case <-ticker.C:
			log.Println("ticker Now")
			e.eventChan <- NewEvent(EventTimer)
		}
	}
}

func (e *Eventbus) process(event *Event) {
	if handlers, ok := e.handlers[event.Type]; ok {
		for handler := range handlers {
			handler.Handle(event)
		}
	}
	for handler := range e.generalHandlers {
		go handler.Handle(event)
	}
}

func (e *Eventbus) Start() {
	e.active = true
	go e.run()
}

func (e *Eventbus) Stop() {
	e.active = false
}

func (e *Eventbus) Register(type_ EventType, handler *Handler) {
	if handlers, ok := e.handlers[type_]; ok {
		if _, exists := handlers[handler]; !exists {
			handlers[handler] = true
		}
	}

}

func (e *Eventbus) Unregister(type_ EventType, handler *Handler) {
	if handlers, ok := e.handlers[type_]; ok {
		delete(handlers, handler)
	}
}

func (e *Eventbus) Put(event *Event) {
	e.eventChan <- event

}

func (e *Eventbus) RegisterGeneralHandler(handler *Handler) {
	e.Lock()
	if _, ok := e.generalHandlers[handler]; !ok {
		e.generalHandlers[handler] = true
	}
	e.Unlock()

}

func (e *Eventbus) UnregisterGeneralHandler(handler *Handler) {
	e.Lock()
	delete(e.generalHandlers, handler)
	e.Unlock()
}
