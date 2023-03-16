package trader

import . "github.com/pseudocodes/vngo/core/event"

type VtModule interface {
	Configure(name string, configRoot string)
	Setup(engine VtEngine, bus *Eventbus) error
	Start() error
	Stop() error

	Name() string
	Description() interface{}
}
