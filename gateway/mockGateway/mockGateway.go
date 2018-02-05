package mockGateway

import (
	"vngo/event"
	. "vngo/trader"
)

type MockGateway struct {
	Base VtGatewayBase
	Name string
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
