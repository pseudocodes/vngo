package main

import (
	. "vngo/common"
	. "vngo/eventbus"
	. "vngo/gateway"
)

type Broker struct {
}

func (b *Broker) InitGateway() {

}

func (b *Broker) AddGateway(name string, gateway Gateway) {

}

func (b *Broker) Connect(name string) {

}

func (b *Broker) Subscribe(req *SubscribeReq, name string) {

}

func (b *Broker) SendOrder(req *OrderReq, name string) {

}

func (b *Broker) CancelOrder(req *CancelOrderReq, name string) {

}

func (b *Broker) QueryAccount(name string) {

}

func (b *Broker) QueryPosition(name string) {

}

func (b *Broker) Exit(name string) {

}

func (b *Broker) WriteLog(content string) {

}

func (b *Broker) StoreConnect() error {
	return nil
}

func (b *Broker) StoreInsert(dbName, collectionName string, data []byte) error {
	return nil
}

func (b *Broker) StoreQuery(dbName, collectionName string, data []byte) {

}

func (b *Broker) StoreUpdate(dbName, collectionName string, data []byte, flt string, upsert bool) {

}

func (b *Broker) StoreLogging(event *EventType) {

}

func (b *Broker) GetContract(vtSymbol string) {

}

func (b *Broker) GetAllContracts() {

}

func (b *Broker) GetOrder(vtOrderID string) {

}

func (b *Broker) GetAllWorkingOrders() {

}

func (b *Broker) GetAllGatewayNames() {

}
