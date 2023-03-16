package engine

import (
	"fmt"
	"sync"

	"github.com/pseudocodes/vngo/core/event"
	"github.com/pseudocodes/vngo/pkg/trader"
)

type MainEngine struct {
	trader.VtEngine

	Gateways map[string]trader.IVtGateway
	Modules  map[string]trader.VtModule

	Eventbus *event.Eventbus
	sync.RWMutex
}

func NewMainEngine() *MainEngine {
	me := &MainEngine{
		Gateways: make(map[string]trader.IVtGateway),
		Modules:  make(map[string]trader.VtModule),
		Eventbus: event.NewEventbus(),
	}
	return me
}

func (me *MainEngine) Configure(name string, configRoot string) error {
	return nil
}

func (me *MainEngine) Start() error {
	me.Eventbus.Start()
	return nil
}

func (me *MainEngine) Stop() error {
	me.Eventbus.Stop()
	return nil
}

func (me *MainEngine) AddModule(name string, module trader.VtModule) {
	me.Lock()
	defer me.Unlock()
	me.Modules[name] = module
}

func (me *MainEngine) AddGateway(name string, gateway trader.IVtGateway) {
	me.Lock()
	defer me.Unlock()
	me.Gateways[name] = gateway
}

func (me *MainEngine) Connect(gatewayName string) error {
	gateway, err := me.getGateway(gatewayName)
	if err != nil {
		return err
	}
	return gateway.Connect()
}

func (me *MainEngine) Subscribe(req *trader.VtSubscribeReq, gatewayName string) error {
	gateway, err := me.getGateway(gatewayName)
	if err != nil {
		return err
	}
	return gateway.Subscribe(req)
}

func (me *MainEngine) SendOrder(req *trader.VtOrderReq, gatewayName string) error {
	gateway, err := me.getGateway(gatewayName)
	if err != nil {
		return err
	}
	return gateway.SendOrder(req)
}

func (me *MainEngine) CancelOrder(req *trader.VtCancelOrderReq, gatewayName string) error {
	gateway, err := me.getGateway(gatewayName)
	if err != nil {
		return err
	}
	return gateway.CancelOrder(req)
}

func (me *MainEngine) QueryAccount(gatewayName string) error {
	gateway, err := me.getGateway(gatewayName)
	if err != nil {
		return err
	}
	return gateway.QueryAccount()
}

func (me *MainEngine) QueryPosition(gatewayName string) error {
	gateway, err := me.getGateway(gatewayName)
	if err != nil {
		return err
	}
	return gateway.QueryPosition()
}

func (me *MainEngine) Close() error {
	me.RLock()
	defer me.RUnlock()
	return nil
}

func (me *MainEngine) getGateway(gatewayName string) (trader.IVtGateway, error) {
	me.RLock()
	defer me.RUnlock()
	gateway, ok := me.Gateways[gatewayName]
	if !ok {
		err := fmt.Errorf("gateway %v not exists \n", gatewayName)
		return nil, err
	}
	return gateway, nil
}
