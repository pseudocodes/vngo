package protocol

type VtModule interface {
	Configure(name string, configRoot string)

	Start() error
	Stop() error

	Name() string
	Description() interface{}
}
