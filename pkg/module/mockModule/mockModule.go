package mockmodule

import (
	"log"
	"sync"
	. "vngo/core/event"
	"vngo/core/protocol"

	"go.uber.org/zap"
)

type MockModule struct {
	Name string
	Ctx  *protocol.ApplicationContext
	Sub  *TypeMuxSubscription

	Log *zap.Logger

	quitChannel chan struct{}
	running     sync.WaitGroup
}

func NewMockModule(ctx *protocol.ApplicationContext) *MockModule {
	return &MockModule{
		Ctx: ctx,
	}
}

func (m *MockModule) Configure(name string, configRoot string) {
	m.Name = name
}

func (m *MockModule) Start() error {
	m.registerEvent()
	go m.eventloop()
	return nil
}

func (m *MockModule) Stop() error {
	m.Sub.Unsubscribe()
	return nil
}

func (m *MockModule) eventloop() {

	for evt := range m.Sub.Chan() {
		switch evt.Data.(type) {
		case (*TickEvent):
			m.processTickEvent(evt.Data.(*TickEvent))
		case (*OrderEvent):
			m.processOrderEvent(evt.Data.(*OrderEvent))
		case (*TradeEvent):
			m.processTradeEvent(evt.Data.(*TradeEvent))
		}
	}
}
func (m *MockModule) registerEvent() {
	sub := m.Ctx.EventQueue.Subscribe(&TickEvent{}, &OrderEvent{}, &TradeEvent{})
	m.Sub = sub
}

func (m *MockModule) processTickEvent(event *TickEvent) error {
	log.Printf("Process Tick Event %+v\n", event)
	return nil
}

func (m *MockModule) processOrderEvent(event *OrderEvent) error {
	log.Printf("Process Order Event %+v\n", event)
	return nil
}

func (m *MockModule) processTradeEvent(event *TradeEvent) error {
	log.Printf("Process Trade Event %+v\n", event)
	return nil
}
