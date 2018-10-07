package protocol

//VtGateway 对接交易接口
type VtGateway interface {
	Init(name string)
	Name() string
	Start() error
	Stop()
	Connect() error

	Close() error
	/*
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
	*/
}
