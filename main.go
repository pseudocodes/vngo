package main

import (
	"fmt"
	"log"
	"time"
	"vngo/event"
	"vngo/gateway/mockGateway"

	. "vngo/event"
	"vngo/module/mockModule"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func simpletest(event *Event) error {
	fmt.Printf("process simpletest! event_type:%v\n", event.Type)
	log.Printf("simple test log\n")
	return nil
}

func process1() {
	ee := NewEventbus()

	// var handler Handler = simpletest
	ee.RegisterGeneralHandler(Handler(simpletest))
	ee.Start()

	for {
		time.Sleep(time.Second)
	}
}

func process2() {
	module := mockModule.NewMockModule()
	gateway := mockGateway.NewMockGateway("mock")
	eventBus := NewEventbus()
	eventBus.Start()
	gateway.Init(eventBus, "mock")
	module.Setup(nil, eventBus)
	module.Start()
	for {
		log.Println("put tick event")
		eventBus.Put(NewEvent(event.EventTick))
		log.Println("put order event")
		eventBus.Put(NewEvent(event.EventOrder))
		log.Println("put trade event")
		eventBus.Put(NewEvent(event.EventTrade))
		time.Sleep(2 * time.Second)
		log.Println()
	}
}

func main() {
	fmt.Println("vn.go")

	// process1()
	process2()
}
