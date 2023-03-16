package protocol

import (
	"fmt"
	"sync"

	"github.com/pseudocodes/vngo/core/event"

	"go.uber.org/zap"
)

type ApplicationContext struct {
	ConfigurationValid bool
	EventQueue         *event.TypeMux

	Logger   *zap.Logger
	Gateways map[string]VtGateway
	Modules  map[string]VtModule

	gatewaySlice []string
	moduleSlice  []string
	sync.RWMutex
}

func (ctx *ApplicationContext) Configurate() error {

	ctx.EventQueue = &event.TypeMux{}

	ctx.Gateways = make(map[string]VtGateway)
	ctx.Modules = make(map[string]VtModule)

	ctx.gatewaySlice = make([]string, 0)
	ctx.moduleSlice = make([]string, 0)

	return nil
}

func (ctx *ApplicationContext) Start() error {
	ctx.Lock()
	defer ctx.Unlock()
	for i, gw := range ctx.gatewaySlice {
		if err := ctx.Gateways[gw].Start(); err != nil {
			//TODO: log error
			fmt.Println(err)
			for j := i - 1; j >= 0; j-- {
				ctx.Gateways[ctx.gatewaySlice[j]].Stop()
			}
		}
	}
	for i, mod := range ctx.moduleSlice {
		if err := ctx.Modules[mod].Start(); err != nil {
			//TODO: log error
			fmt.Println(err)
			for j := i - 1; j >= 0; j-- {
				ctx.Modules[ctx.moduleSlice[j]].Stop()
			}
		}
	}
	return nil
}

func (ctx *ApplicationContext) Stop() {
	ctx.Lock()
	defer ctx.Unlock()
	for i := len(ctx.moduleSlice); i >= 0; i-- {
		ctx.Modules[ctx.moduleSlice[i]].Stop()
	}

	for g := len(ctx.Gateways); g >= 0; g-- {
		ctx.Gateways[ctx.gatewaySlice[g]].Stop()
	}

}

func (ctx *ApplicationContext) AddModule(name string, module VtModule) {
	ctx.Lock()
	defer ctx.Unlock()

	ctx.Modules[name] = module
	ctx.moduleSlice = append(ctx.moduleSlice, name)
}

func (ctx *ApplicationContext) AddGateway(name string, gateway VtGateway) {
	ctx.Lock()
	defer ctx.Unlock()

	ctx.Gateways[name] = gateway
	ctx.gatewaySlice = append(ctx.gatewaySlice, name)
}

func (ctx *ApplicationContext) RemoveModule(name string) {
	ctx.Lock()
	defer ctx.Unlock()

	module, ok := ctx.Modules[name]
	if !ok {
		return
	}
	module.Stop()

	delete(ctx.Modules, name)
	var idx = 0
	for i, key := range ctx.moduleSlice {
		if key == name {
			idx = i
			break
		}
	}
	ctx.moduleSlice = append(ctx.moduleSlice[:idx-1], ctx.moduleSlice[idx+1:]...)
}

func (ctx *ApplicationContext) RemoveGateway(name string) {
	ctx.Lock()
	defer ctx.Unlock()

	gateway, ok := ctx.Gateways[name]
	if !ok {
		return
	}
	gateway.Stop()

	delete(ctx.Gateways, name)
	var idx = 0
	for i, key := range ctx.gatewaySlice {
		if key == name {
			idx = i
			break
		}
	}
	ctx.gatewaySlice = append(ctx.gatewaySlice[:idx-1], ctx.gatewaySlice[idx+1:]...)
}
