package protocol

// 默认空值
const (
	EMPTY_STRING  = ""
	EMPTY_UNICODE = ""
	EMPTY_INT     = 0
	EMPTY_FLOAT   = 0.0

	// 方向常量
	DIRECTION_NONE         = "无方向"
	DIRECTION_LONG         = "多"
	DIRECTION_SHORT        = "空"
	DIRECTION_UNKNOWN      = "未知"
	DIRECTION_NET          = "净"
	DIRECTION_SELL         = "卖出"  // IB接口
	DIRECTION_COVEREDSHORT = "备兑空" // 证券期权

	// 开平常量
	OFFSET_NONE           = "无开平"
	OFFSET_OPEN           = "开仓"
	OFFSET_CLOSE          = "平仓"
	OFFSET_CLOSETODAY     = "平今"
	OFFSET_CLOSEYESTERDAY = "平昨"
	OFFSET_UNKNOWN        = "未知"

	// 状态常量
	STATUS_NOTTRADED  = "未成交"
	STATUS_PARTTRADED = "部分成交"
	STATUS_ALLTRADED  = "全部成交"
	STATUS_CANCELLED  = "已撤销"
	STATUS_REJECTED   = "拒单"
	STATUS_UNKNOWN    = "未知"

	// 合约类型常量
	PRODUCT_EQUITY      = "股票"
	PRODUCT_FUTURES     = "期货"
	PRODUCT_OPTION      = "期权"
	PRODUCT_INDEX       = "指数"
	PRODUCT_COMBINATION = "组合"
	PRODUCT_FOREX       = "外汇"
	PRODUCT_UNKNOWN     = "未知"
	PRODUCT_SPOT        = "现货"
	PRODUCT_DEFER       = "延期"
	PRODUCT_ETF         = "ETF"
	PRODUCT_WARRANT     = "权证"
	PRODUCT_BOND        = "债券"
	PRODUCT_NONE        = ""

	// 价格类型常量
	PRICETYPE_LIMITPRICE  = "限价"
	PRICETYPE_MARKETPRICE = "市价"
	PRICETYPE_FAK         = "FAK"
	PRICETYPE_FOK         = "FOK"

	// 期权类型
	OPTION_CALL = "看涨期权"
	OPTION_PUT  = "看跌期权"

	// 交易所类型
	EXCHANGE_SSE     = "SSE"     // 上交所
	EXCHANGE_SZSE    = "SZSE"    // 深交所
	EXCHANGE_CFFEX   = "CFFEX"   // 中金所
	EXCHANGE_SHFE    = "SHFE"    // 上期所
	EXCHANGE_CZCE    = "CZCE"    // 郑商所
	EXCHANGE_DCE     = "DCE"     // 大商所
	EXCHANGE_SGE     = "SGE"     // 上金所
	EXCHANGE_INE     = "INE"     // 国际能源交易中心
	EXCHANGE_UNKNOWN = "UNKNOWN" // 未知交易所
	EXCHANGE_NONE    = ""        // 空交易所
	EXCHANGE_HKEX    = "HKEX"    // 港交所
	EXCHANGE_HKFE    = "HKFE"    // 香港期货交易所

	EXCHANGE_SMART    = "SMART"    // IB智能路由（股票、期权）
	EXCHANGE_NYMEX    = "NYMEX"    // IB 期货
	EXCHANGE_GLOBEX   = "GLOBEX"   // CME电子交易平台
	EXCHANGE_IDEALPRO = "IDEALPRO" // IB外汇ECN

	EXCHANGE_CME = "CME" // CME交易所
	EXCHANGE_ICE = "ICE" // ICE交易所
	EXCHANGE_LME = "LME" // LME交易所

	EXCHANGE_OANDA  = "OANDA"  // OANDA外汇做市商
	EXCHANGE_OKCOIN = "OKCOIN" // OKCOIN比特币交易所
	EXCHANGE_HUOBI  = "HUOBI"  // 火币比特币交易所
	EXCHANGE_LHANG  = "LHANG"  // 链行比特币交易所

	// 货币类型
	CURRENCY_USD     = "USD"     // 美元
	CURRENCY_CNY     = "CNY"     // 人民币
	CURRENCY_HKD     = "HKD"     // 港币
	CURRENCY_UNKNOWN = "UNKNOWN" // 未知货币
	CURRENCY_NONE    = ""        // 空货币

	// 数据库
	LOG_DB_NAME = "VnTrader_Log_Db"

	// 接口类型
	GATEWAYTYPE_EQUITY        = "equity"        // 股票、ETF、债券
	GATEWAYTYPE_FUTURES       = "futures"       // 期货、期权、贵金属
	GATEWAYTYPE_INTERNATIONAL = "international" // 外盘
	GATEWAYTYPE_BTC           = "btc"           // 比特币
	GATEWAYTYPE_DATA          = "data"          // 数据（非交易）

)

const (
	SAVE_DATA = "保存数据"

	CONTRACT_SYMBOL = "合约代码"
	CONTRACT_NAME   = "名称"
	LAST_PRICE      = "最新价"
	PRE_CLOSE_PRICE = "昨收盘"
	VOLUME          = "成交量"
	OPEN_INTEREST   = "持仓量"
	OPEN_PRICE      = "开盘价"
	HIGH_PRICE      = "最高价"
	LOW_PRICE       = "最低价"
	TIME            = "时间"
	GATEWAY         = "接口"
	CONTENT         = "内容"

	ERROR_CODE    = "错误代码"
	ERROR_MESSAGE = "错误信息"

	TRADE_ID   = "成交编号"
	ORDER_ID   = "委托编号"
	DIRECTION  = "方向"
	OFFSET     = "开平"
	PRICE      = "价格"
	TRADE_TIME = "成交时间"

	ORDER_VOLUME    = "委托数量"
	TRADED_VOLUME   = "成交数量"
	ORDER_STATUS    = "委托状态"
	ORDER_TIME      = "委托时间"
	CANCEL_TIME     = "撤销时间"
	FRONT_ID        = "前置编号"
	SESSION_ID      = "会话编号"
	POSITION        = "持仓量"
	YD_POSITION     = "昨持仓"
	FROZEN          = "冻结量"
	POSITION_PROFIT = "持仓盈亏"

	ACCOUNT_ID   = "账户编号"
	PRE_BALANCE  = "昨净值"
	BALANCE      = "净值"
	AVAILABLE    = "可用"
	COMMISSION   = "手续费"
	MARGIN       = "保证金"
	CLOSE_PROFIT = "平仓盈亏"

	TRADING           = "交易"
	PRICE_TYPE        = "价格类型"
	EXCHANGE          = "交易所"
	CURRENCY          = "货币"
	PRODUCT_CLASS     = "产品类型"
	LAST              = "最新"
	SEND_ORDER        = "发单"
	CANCEL_ALL        = "全撤"
	VT_SYMBOL         = "vt系统代码"
	CONTRACT_SIZE     = "合约大小"
	PRICE_TICK        = "最小价格变动"
	STRIKE_PRICE      = "行权价"
	UNDERLYING_SYMBOL = "标的代码"
	OPTION_TYPE       = "期权类型"
	EXPIRY_DATE       = "到期日"

	REFRESH         = "刷新"
	SEARCH          = "查询"
	CONTRACT_SEARCH = "合约查询"

	BID_1 = "买一"
	BID_2 = "买二"
	BID_3 = "买三"
	BID_4 = "买四"
	BID_5 = "买五"
	ASK_1 = "卖一"
	ASK_2 = "卖二"
	ASK_3 = "卖三"
	ASK_4 = "卖四"
	ASK_5 = "卖五"

	BID_PRICE_1 = "买一价"
	BID_PRICE_2 = "买二价"
	BID_PRICE_3 = "买三价"
	BID_PRICE_4 = "买四价"
	BID_PRICE_5 = "买五价"
	ASK_PRICE_1 = "卖一价"
	ASK_PRICE_2 = "卖二价"
	ASK_PRICE_3 = "卖三价"
	ASK_PRICE_4 = "卖四价"
	ASK_PRICE_5 = "卖五价"

	BID_VOLUME_1 = "买一量"
	BID_VOLUME_2 = "买二量"
	BID_VOLUME_3 = "买三量"
	BID_VOLUME_4 = "买四量"
	BID_VOLUME_5 = "买五量"
	ASK_VOLUME_1 = "卖一量"
	ASK_VOLUME_2 = "卖二量"
	ASK_VOLUME_3 = "卖三量"
	ASK_VOLUME_4 = "卖四量"
	ASK_VOLUME_5 = "卖五量"

	MARKET_DATA   = "行情"
	LOG           = "日志"
	ERROR         = "错误"
	TRADE         = "成交"
	ORDER         = "委托"
	ACCOUNT       = "账户"
	WORKING_ORDER = "可撤"

	SYSTEM           = "系统"
	CONNECT_DATABASE = "连接数据库"
	EXIT             = "退出"
	APPLICATION      = "功能"
	DATA_RECORDER    = "行情记录"
	RISK_MANAGER     = "风控管理"

	STRATEGY     = "策略"
	CTA_STRATEGY = "CTA策略"

	HELP         = "帮助"
	RESTORE      = "还原窗口"
	ABOUT        = "关于"
	TEST         = "测试"
	CONNECT      = "连接"
	EDIT_SETTING = "编辑配置"
	LOAD         = "读取"
	SAVE         = "保存"

	CPU_MEMORY_INFO = "CPU使用率：{cpu}%   内存使用率：{memory}%"
	CONFIRM_EXIT    = "确认退出？"

	GATEWAY_NOT_EXIST             = "接口不存在：{gateway}"
	DATABASE_CONNECTING_COMPLETED = "MongoDB连接成功"
	DATABASE_CONNECTING_FAILED    = "MongoDB连接失败"
	DATA_INSERT_FAILED            = "数据插入失败，MongoDB没有连接"
	DATA_QUERY_FAILED             = "数据查询失败，MongoDB没有连接"
	DATA_UPDATE_FAILED            = "数据更新失败，MongoDB没有连接"
)
