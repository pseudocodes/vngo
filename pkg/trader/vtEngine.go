package trader

type VtEngine interface {
	Start() error
	Stop() error
	Connect(gatewayName string) error
	Subscribe(req *VtSubscribeReq, gateway string) error
	SendOrder(req *VtOrderReq, gateway string) error
	CancelOrder(req *VtCancelOrderReq, gateway string) error
	QueryAccount(gateway string) error
	QueryPosition(gateway string) error
	Close() error
}
