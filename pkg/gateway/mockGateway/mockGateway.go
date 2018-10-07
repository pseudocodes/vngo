package mockGateway

import (
	"sync"
	"vngo/core/event"
	"vngo/core/protocol"
	. "vngo/pkg/trader"

	"go.uber.org/zap"
)

type MockGateway struct {
	Base VtGatewayBase
	Name string

	Ctx *protocol.ApplicationContext
	Log *zap.Logger

	quitChannel chan struct{}
	running     sync.WaitGroup
}

func NewMockGateway(name string) *MockGateway {
	return &MockGateway{
		Name: name,
	}
}
func (g *MockGateway) Init(bus *event.Eventbus, name string) {
	g.Base.Init(bus, name)
}

func (g *MockGateway) Connect() error {
	return nil
}

func (g *MockGateway) Subscribe(subscribeReq *VtSubscribeReq) error {
	return nil
}

func (g *MockGateway) SendOrder(orderReq *VtOrderReq) error {
	return nil
}

func (g *MockGateway) CancelOrder(cancelOrderReq *VtCancelOrderReq) error {
	return nil
}

func (g *MockGateway) QueryAccount() error {
	return nil
}

func (g *MockGateway) QueryPosition() error {
	return nil
}

func (g *MockGateway) Close() error {
	return nil
}
