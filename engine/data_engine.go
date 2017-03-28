package engine

import . "vngo/eventbus"

type DataEngine struct {
	eventbus         *Eventbus
	ContractDict     map[string]interface{}
	OrderDict        map[string]interface{}
	WorkingOrderDict map[string]interface{}
}

func NewDataEngine(event *Eventbus) *DataEngine {
	return &DataEngine{
		eventbus:         event,
		ContractDict:     make(map[string]interface{}),
		OrderDict:        make(map[string]interface{}),
		WorkingOrderDict: make(map[string]interface{}),
	}
}

func (self *DataEngine) UpdateContract(event *Event) error {
	return nil

}

func (self *DataEngine) GetContract(vtSymbol string) {

}

func (self *DataEngine) GetAllContracts() {

}

func (self *DataEngine) SaveContracts() {

}

func (self *DataEngine) LoadContracts() {

}

func (self *DataEngine) UpdateOrder(event *Event) error {
	return nil

}

func (self *DataEngine) GetOrder(vtOrderID string) {

}

func (self *DataEngine) GetAllWorkingOrders() {

}

func (self *DataEngine) registerEvent() {
	var updateContract Handler = self.UpdateContract
	var updateOrder Handler = self.UpdateOrder
	self.eventbus.Register(EventContract, &updateContract)
	self.eventbus.Register(EventOrder, &updateOrder)
}
