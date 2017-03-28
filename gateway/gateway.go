package gateway

type Gateway interface {
	OnTick()
	OnTrade()
	OnOrder()
	OnPosition()
	OnAccount()
	OnError()
	OnLog()
	OnContract()

	Connect()
	Subscribe()
	SendOrder()
	CancelOrder()
	QueryAccount()
	QueryPostion()
	Close()
}
