package event

type TickEvent struct {
	Exchange string
}

type OrderEvent struct {
	Exchange string
}

type TradeEvent struct {
	Exchange string
}

type FillEvent struct {
	Exchange string
}
