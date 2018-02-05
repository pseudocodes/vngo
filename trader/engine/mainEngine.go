package engine

import (
	"fmt"
	"sync"

	. "vngo/event"
	. "vngo/trader"
)

type MainEngine struct {
	VtEngine

	Gateways map[string]IVtGateway
	Modules  map[string]VtModule

	Eventbus *Eventbus
	sync.RWMutex
}

func NewMainEngine() *MainEngine {
	me := &MainEngine{
		Gateways: make(map[string]IVtGateway),
		Modules:  make(map[string]VtModule),
		Eventbus: NewEventbus(),
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

func (me *MainEngine) AddModule(name string, module VtModule) {
	me.Lock()
	defer me.Unlock()
	me.Modules[name] = module
}

func (me *MainEngine) AddGateway(name string, gateway IVtGateway) {
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

func (me *MainEngine) Subscribe(req *VtSubscribeReq, gatewayName string) error {
	gateway, err := me.getGateway(gatewayName)
	if err != nil {
		return err
	}
	return gateway.Subscribe(req)
}

func (me *MainEngine) SendOrder(req *VtOrderReq, gatewayName string) error {
	gateway, err := me.getGateway(gatewayName)
	if err != nil {
		return err
	}
	return gateway.SendOrder(req)
}

func (me *MainEngine) CancelOrder(req *VtCancelOrderReq, gatewayName string) error {
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

func (me *MainEngine) getGateway(gatewayName string) (IVtGateway, error) {
	me.RLock()
	defer me.RUnlock()
	gateway, ok := me.Gateways[gatewayName]
	if !ok {
		err := fmt.Errorf("gateway %v not exists \n", gatewayName)
		return nil, err
	}
	return gateway, nil
}
