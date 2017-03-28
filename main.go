package main

import (
	"fmt"
	"log"
	"time"

	. "vngo/eventbus"
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

	var handler Handler = simpletest
	ee.RegisterGeneralHandler(&handler)
	ee.Start()

	for {
		time.Sleep(time.Second)
	}
}

func main() {
	fmt.Println("vn-go")

	process1()

}
