package trader

import . "vngo/event"

type VtModule interface {
	Configure(name string, configRoot string)
	Setup(engine VtEngine, bus *Eventbus) error
	Start() error
	Stop() error

	Description() interface{}
}
