package trader

import . "vngo/event"

//VtGateway 对接交易接口
type IVtGateway interface {
	Init(eventbus *Eventbus, name string)

	Connect() error
	Subscribe(subscribeReq *VtSubscribeReq) error
	SendOrder(orderReq *VtOrderReq) error
	CancelOrder(cancelOrderReq *VtCancelOrderReq) error
	QueryAccount() error
	QueryPosition() error
	Close() error

	OnTick(tick *VtTickData)
	OnTrade(trade *VtTradeData)
	OnOrder(order *VtOrderData)
	OnPosition(position *VtPositionData)
	OnAccount(account *VtAccountData)
	OnError(err *VtErrorData)
	OnLog(log *VtLogData)
	OnContract(contract *VtContractData)

	Decription() interface{}
}

type VtGatewayBase struct {
	IVtGateway
	EventBus    *Eventbus
	GatewayName string
}

func (g *VtGatewayBase) Init(eventbus *Eventbus, name string) {
	g.EventBus = eventbus
	g.GatewayName = name
}

//OnTick 市场行情推送
func (g *VtGatewayBase) OnTick(tick *VtTickData) {

	event1 := NewEvent(EventTick)
	event1.Data = tick
	g.EventBus.Put(event1)

	// 特定合约代码的事件
	event2 := NewEvent(EventType(EventTick.String() + tick.VtSymbol))
	event2.Data = tick
	g.EventBus.Put(event2)
}

//OnTrade 成交信息推送
func (g *VtGatewayBase) OnTrade(trade *VtTradeData) {
	// 通用事件
	event1 := NewEvent(EventTrade)
	event1.Data = trade
	g.EventBus.Put(event1)

	// 特定合约的成交事件
	event2 := NewEvent(EventType(EventTrade.String() + trade.VtSymbol))
	event1.Data = trade
	g.EventBus.Put(event2)
}

//OnOrder 订单变化推送
func (g *VtGatewayBase) OnOrder(order *VtOrderData) {
	// 通用事件
	event1 := NewEvent(EventOrder)
	event1.Data = order
	g.EventBus.Put(event1)

	// 特定订单编号的事件
	event2 := NewEvent(EventType(EventOrder.String() + order.VtOrderID))
	event1.Data = order
	g.EventBus.Put(event2)
}

//OnPosition 持仓信息推送
func (g *VtGatewayBase) OnPosition(position *VtPositionData) {

	// 通用事件
	event1 := NewEvent(EventPosition)
	event1.Data = position
	g.EventBus.Put(event1)

	// 特定合约代码的事件
	event2 := NewEvent(EventType(EventPosition.String() + position.VtSymbol))
	event1.Data = position
	g.EventBus.Put(event2)
}

//OnAccount 账户信息推送
func (g *VtGatewayBase) OnAccount(account *VtAccountData) {
	// 通用事件
	event1 := NewEvent(EventAccount)
	event1.Data = account
	g.EventBus.Put(event1)

	// 特定合约代码的事件
	event2 := NewEvent(EventType(EventAccount.String() + account.VtAccountID))
	event1.Data = account
	g.EventBus.Put(event2)
}

//OnError 错误信息推送
func (g *VtGatewayBase) OnError(err *VtErrorData) {
	event1 := NewEvent(EventError)
	event1.Data = err
	g.EventBus.Put(event1)
}

//OnLog 日志推送
func (g *VtGatewayBase) OnLog(log *VtLogData) {
	event1 := NewEvent(EventLog)
	event1.Data = log
	g.EventBus.Put(event1)
}

//OnContract 合约基础信息推送
func (g *VtGatewayBase) OnContract(contract *VtContractData) {
	event1 := NewEvent(EventContract)
	event1.Data = contract
	g.EventBus.Put(event1)
}

func (g *VtGatewayBase) Connect() error {
	return nil
}

func (g *VtGatewayBase) Subscribe(subscribeReq *VtSubscribeReq) error {
	return nil
}

func (g *VtGatewayBase) SendOrder(orderReq *VtOrderReq) error {
	return nil
}

func (g *VtGatewayBase) CancelOrder(cancelOrderReq *VtCancelOrderReq) error {
	return nil
}

func (g *VtGatewayBase) QueryAccount() error {
	return nil
}

func (g *VtGatewayBase) QueryPosition() error {
	return nil
}

func (g *VtGatewayBase) Close() error {
	return nil
}
