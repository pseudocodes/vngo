package trader

import "time"

//VtBaseData 基本数据
type VtBaseData struct {
	Gateway string
	RawData []byte
}

//VtTickData Tick行情数据类
type VtTickData struct {
	VtBaseData
	// 代码相关
	Symbol   string // 合约代码
	Exchange string // 交易所代码
	VtSymbol string // 合约在vt系统中的唯一代码，通常是 合约代码.交易所代码

	// 成交数据
	LastPrice    float64   //最新成交价
	LastVolume   int64     // 最新成交量
	Volume       int64     // 今天总成交量
	OpenInterest int64     // 持仓量
	Time         string    // 时间 11:20:56.5
	Date         string    // 日期 20151009
	Datetime     time.Time // golang 的datetime时间对象

	// 常规行情
	OpenPrice     float64 // 今日开盘价
	HighPrice     float64 // 今日最高价
	HowPrice      float64 // 今日最低价
	PreClosePrice float64

	UpperLimit float64 // 涨停价
	LowerLimit float64 // 跌停价

	// 五档行情
	BidPrice1 float64
	BidPrice2 float64
	BidPrice3 float64
	BidPrice4 float64
	BidPrice5 float64

	AskPrice1 float64
	AskPrice2 float64
	AskPrice3 float64
	AskPrice4 float64
	AskPrice5 float64

	BidVolume1 int64
	BidVolume2 int64
	BidVolume3 int64
	BidVolume4 int64
	BidVolume5 int64

	AskVolume1 int64
	AskVolume2 int64
	AskVolume3 int64
	AskVolume4 int64
	AskVolume5 int64
}

// VtBarData: K线数据
type VtBarData struct {
	VtBaseData

	VtSymbol string // vt系统代码
	Symbol   string // 代码
	Exchange string // 交易所

	Open  float64 // OHLC
	High  float64
	Low   float64
	Close float64

	Date     string    // bar开始的时间，日期
	Time     string    // 时间
	Datetime time.Time // golang 的 datetime 时间对象

	Volume       int64 // 成交量
	OpenInterest int64 // 持仓量

}

// VtTradeData 成交数据类
type VtTradeData struct {
	VtBaseData

	// 代码编号相关
	Symbol   string // 合约代码
	Exchange string // 交易所代码
	VtSymbol string // 合约在vt系统中的唯一代码，通常是 合约代码.交易所代码

	TradeID   string // 成交编号
	VtTradeID string // 成交在vt系统中的唯一编号，通常是 Gateway 名.成交编号

	OrderID   string // 订单编号
	VtOrderID string // 订单在vt系统中的唯一编号，通常是 Gateway 名.订单编号

	// 成交相关
	Direction string  // 成交方向
	Offset    string  // 成交开平仓
	Price     float64 // 成交价格
	Volume    int64   // 成交数量
	TradeTime string  // 成交时间
}

// VtOrderData 订单数据类
type VtOrderData struct {
	VtBaseData

	// 代码编号相关
	Symbol   string // 合约代码
	Exchange string // 交易所代码
	VtSymbol string // 合约在vt系统中的唯一代码，通常是 合约代码.交易所代码

	OrderID   string // 订单编号
	VtOrderID string // 订单在vt系统中的唯一编号，通常是 Gateway名.订单编号

	// 报单相关
	Direction    string  // 报单方向
	Offset       string  // 报单开平仓
	Price        float64 // 报单价格
	TotalVolume  int64   // 报单总数量
	TradedVolume int64   // 报单成交数量
	Status       string  // 报单状态

	OrderTime  string // 发单时间
	CancelTime string // 撤单时间

	// CTP/LTS相关
	FrontID   int64 // 前置机编号
	SessionID int64 // 连接编号
}

// VtPositionData 持仓数据类
type VtPositionData struct {
	VtBaseData

	// 代码编号相关
	Symbol   string // 合约代码
	Exchange string // 交易所代码
	VtSymbol string // 合约在vt系统中的唯一代码，合约代码.交易所代码

	// 持仓相关
	Direction      string  // 持仓方向
	Position       int64   // 持仓量
	Frozen         int64   // 冻结数量
	Price          float64 // 持仓均价
	VtPositionName string  // 持仓在vt系统中的唯一代码，通常是vtSymbol.方向
	YdPosition     int64   // 昨持仓
	PositionProfit float64 // 持仓盈亏
}

// VtAccountData 账户数据类
type VtAccountData struct {
	VtBaseData

	// 账号代码相关
	AccountID   string // 账户代码
	VtAccountID string // 账户在vt中的唯一代码，通常是 Gateway名.账户代码

	// 数值相关
	PreBalance     float64 // 昨日账户结算净值
	Balance        float64 // 账户净值
	Available      float64 // 可用资金
	Commission     float64 // 今日手续费
	Margin         float64 // 保证金占用
	CloseProfit    float64 // 平仓盈亏
	PositionProfit float64 // 持仓盈亏
}

// VtErrorData 错误数据类
type VtErrorData struct {
	VtBaseData

	ErrorID        string // 错误代码
	ErrorMsg       string // 错误信息
	AdditionalInfo string // 补充信息

	ErrorTime string // 错误生成时间

}

func (e VtErrorData) Error() string {
	return ""
}

// VtLogData 日志数据类
type VtLogData struct {
	VtBaseData

	logTime    string // 日志生成时间
	logContent string // 日志信息
	logLevel   string // 日志级别
}

// VtContractData 合约详细信息类
type VtContractData struct {
	VtBaseData

	Symbol   string // 代码
	Exchange string // 交易所代码
	VtSymbol string // 合约在vt系统中的唯一代码，通常是 合约代码.交易所代码
	Name     string // 合约中文名

	ProductClass string  // 合约类型
	Size         int64   // 合约大小
	PriceTick    float64 // 合约最小价格TICK

	// 期权相关
	StrikePrice      float64 // 期权行权价
	UnderlyingSymbol string  // 标的物合约代码
	OptionType       string  // 期权类型
	ExpiryDate       string  // 到期日
}

// VtSubscribeReq 订阅行情时传入结构体
type VtSubscribeReq struct {
	Symbol   string // 代码
	Exchange string // 交易所

	// 以下为IB相关
	ProductClass string  // 合约类型
	Currency     string  // 合约货币
	Expiry       string  // 到期日
	StrikePrice  float64 // 行权价
	OptionType   string  // 期权类型
}

// VtOrderReq 发单时传入的结构体
type VtOrderReq struct {
	Symbol   string  // 代码
	Exchange string  // 交易所
	VtSymbol string  // VT合约代码
	Price    float64 // 价格
	Volume   int64   // 数量

	PriceType string // 价格类型
	Direction string // 买卖
	Offset    string // 开平

	// 以下为IB相关
	ProductClass                 string  // 合约类型
	Currency                     string  // 合约货币
	Expiry                       string  // 到期日
	StrikePrice                  float64 // 行权价
	OptionType                   string  // 期权类型
	LastTradeDateOrContractMonth string  // 合约月,IB专用
	Multiplier                   string  // 乘数,IB专用
}

// VtCancelOrderReq 撤单时传入结构体
type VtCancelOrderReq struct {
	Symbol   string // 代码
	Exchange string // 交易所
	VtSymbol string // VT合约代码

	// 以下字段主要和CTP、LTS类接口相关
	OrderID   string // 报单号
	FrontID   string // 前置机号
	SessionID string // 会话号

}
