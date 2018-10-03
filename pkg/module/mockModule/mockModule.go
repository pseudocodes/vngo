package mockModule

import (
	"log"
	"vngo/event"
	"vngo/trader"
)

type MockModule struct {
	trader.VtModule
	engine   trader.VtEngine
	eventbus *event.Eventbus
	name     string
}

func NewMockModule() *MockModule {
	return &MockModule{}
}

func (m *MockModule) Configure(name string, configRoot string) {
	m.name = name
}

func (m *MockModule) Setup(engine trader.VtEngine, bus *event.Eventbus) error {
	m.engine = engine
	m.eventbus = bus
	return nil
}

func (m *MockModule) Start() error {
	m.registerEvent()
	return nil
}

func (m *MockModule) Stop() error {
	return nil
}

func (m *MockModule) Description() interface{} {
	return nil
}

func (m *MockModule) registerEvent() {
	m.eventbus.Register(event.EventTick, event.Handler(m.processTickEvent))
	m.eventbus.Register(event.EventOrder, event.Handler(m.processOrderEvent))
	m.eventbus.Register(event.EventTrade, event.Handler(m.processTradeEvent))
}

func (m *MockModule) processTickEvent(event *event.Event) error {
	log.Printf("Process Tick Event %+v\n", event)
	return nil
}

func (m *MockModule) processOrderEvent(event *event.Event) error {
	log.Printf("Process Order Event %+v\n", event)
	return nil
}

func (m *MockModule) processTradeEvent(event *event.Event) error {
	log.Printf("Process Trade Event %+v\n", event)
	return nil
}
