package strategy

import (
	"reflect"

	"github.com/pseudocodes/vngo/pkg/utils"

	"github.com/pseudocodes/goctp"
)

var _ IStrategy = &Strategy{}

type IStrategy interface {
	SetConfig(cfg *Config, ctp *CtpClient, spi *FtdcTraderSpi)
	SetInstruments(mapKey string, symInfos InstrumentInfoStruct) //初始化合约品种
	GetInstrumentInfo(InstrumentID string) (InstrumentInfoStruct, bool)
	GetFuturesList() []string                //获取合约列表
	EmptyPosition()                          //查询仓位信息前，清空仓位信息
	UpdatePosition(p InvestorPositionStruct) //更新仓位信息
	OnStart(spi *FtdcMdSpi)                  // *交易账户初始化完毕，行情登录成功后，启动策略（FtdcMdSpi.OnRspUserLogin）
	OnQuote(sym string, ticker Ticker)       // tick 报单事件
	OnBar(sym string, kdata Kline)           // K线发生变化事件
	OnOrderChange(pOrder OrderListStruct)    // 订单发生变化事件
	OnTradeDeal(pTrade goctp.TradeField)     // 成交事件
}

type Strategy struct {
	oCtp      *CtpClient               // main CtpClient
	TraderSpi *FtdcTraderSpi           // 交易相关
	Symbol    []string                 `json:"symbol"` // 跟踪的品种列表，跨期或是跨品种需要配置2个品种
	Params    []int64                  `json:"params"`
	Period    int64                    `json:"period"`   // 主要周期级别，60 代表1分钟
	MaxKLen   int64                    `json:"maxklen"`  // 存储的K线最大长度
	Position  []InvestorPositionStruct `json:"position"` // 仓位列表
	AM        map[string][]Kline       `json:"am"`       // k线序列数据，长度为MaxKLen, string 为symbol
	TK        map[string]Ticker        `json:"tk"`       // Tick 数据，保存上一个，用于计算现手和仓位差

	MapInstrumentInfos utils.Map `json:"map_instrument_infos"` // 交易所合约详情列表 InstrumentInfoStruct

}

func (s *Strategy) SetInstruments(mapKey string, symInfos InstrumentInfoStruct) {
	s.MapInstrumentInfos.Set(mapKey, symInfos)
}

func (s *Strategy) OnStart(mdspi *FtdcMdSpi) { //策略启动时调用，可被实际策略事先调用
	mdspi.SubscribeMarketData(s.Symbol)
	// 初始化K线
	//log.Printf("%v", s.GetFuturesList())
	for _, sym := range s.Symbol {
		_, ok := s.GetInstrumentInfo(sym)
		if ok {
			kl := len(s.AM[sym])
			if kl > 0 {
				s.AM[sym] = s.AM[sym][:0]
			}
		}
	}
}

func (s *Strategy) OnQuote(sym string, ticker Ticker) { //触发Ticker事件

	ticker.Volume = ticker.DayVolume - s.TK[sym].DayVolume
	ticker.Interest = ticker.OpenInterest - s.TK[sym].OpenInterest
	/*
		fmt.Println(utils.TimeToStr(int64(ticker.UpdateTime), ""), "\t合约："+sym,
			"\t价格："+utils.Float64ToString(ticker.Price),
			"\t现手："+utils.IntToString(ticker.Volume),
			"\t仓差："+utils.Float64ToString(ticker.Interest),
			"\t买一："+utils.Float64ToString(ticker.BidPrice1)+"\t"+utils.IntToString(ticker.BidVolume1),
			"\t卖一："+utils.Float64ToString(ticker.AskPrice1)+"\t"+utils.IntToString(ticker.AskVolume1),
		)
	*/

	if s.TK[sym].UpdateTime == ticker.UpdateTime { //1秒内多条推送取消
		currKline[sym] = UpdateKLine(currKline[sym], ticker)
		//更新PreTicker数据
	} else { //下一秒的新数据
		//计算是否要更新K线
		if ticker.UpdateTime%s.Period == 0 {
			if Kl, ok := currKline[sym]; ok && Kl.Datetime > 0 {
				go s.oCtp.Strategy.OnBar(sym, currKline[sym]) //此处需要调用最终实现接口的OnBar，否则策略最终策略无法执行
			} else { //重新初始化K线
				//log.Println(currKline[sym])
				go func() {
					syminfo, ok := s.GetInstrumentInfo(sym)
					if ok {
						symall := syminfo.ExchangeID + "." + sym
						symall, s.AM[sym] = InitKline(symall, s.Period, s.MaxKLen)
						kl := len(s.AM[sym])
						if int64(kl) != s.MaxKLen || s.AM[sym][0].Datetime < 946656000 {
							return
						}
						currK := s.AM[sym][kl-1]

						s.AM[sym] = append(s.AM[sym], currK)
						copy(s.AM[sym][1:], s.AM[sym][0:kl-1])
						s.AM[sym][0] = currK
						s.AM[sym] = s.AM[sym][:kl]
						s.oCtp.Strategy.OnBar(sym, currK)
					}
				}()
			}
			currKline[sym] = MakeKLine(ticker)
		} else {
			currKline[sym] = UpdateKLine(currKline[sym], ticker)
		}
	}
	s.TK[sym] = ticker
}

func (s *Strategy) OnBar(sym string, kdata Kline) { //触发K线更新事件
	kl := int64(len(s.AM[sym]))
	if kl == s.MaxKLen {
		s.AM[sym] = s.AM[sym][1:]
		s.AM[sym] = append(s.AM[sym], kdata)
	}
	//fmt.Println(sym, s.AM[sym])
}

func (s *Strategy) OnOrderChange(pOrder OrderListStruct) {

}

func (s *Strategy) OnTradeDeal(pTrade goctp.TradeField) {

}

func (s *Strategy) EmptyPosition() {
	s.Position = s.Position[:0]
}

func (s *Strategy) UpdatePosition(p InvestorPositionStruct) {
	s.Position = append(s.Position, p)
}

// GetInstrumentInfo 获得合约详情信息
func (s *Strategy) GetInstrumentInfo(InstrumentID string) (InstrumentInfoStruct, bool) {
	if v, ok := s.MapInstrumentInfos.Get(InstrumentID); ok {
		return v.(InstrumentInfoStruct), true
	} else {
		var mInstrumentInfo InstrumentInfoStruct
		return mInstrumentInfo, false
	}
}

// GetFuturesList 获得期货合约列表【只有期货，不含期权】
func (s *Strategy) GetFuturesList() []string {
	var InstrumentList []string

	mInstrumentInfos := s.MapInstrumentInfos.GetAll()
	for _, v := range mInstrumentInfos {
		val := v.(InstrumentInfoStruct)

		// 类型为期货的合约
		if val.ProductClass == "1" {
			InstrumentList = append(InstrumentList, val.InstrumentID)
		}
	}

	return InstrumentList
}

func (s *Strategy) OpenOrder(sym string, direction byte, price float64, volume int) {
	// var Input InputOrderStruct
	// Input.InstrumentID = sym
	// Input.Direction = direction
	// Input.Price = price
	// Input.Volume = volume
	// s.TraderSpi.OrderOpen(Input)
}

func (s *Strategy) OrderClose(sym string, direction byte, price float64, volume int) {
	// var Input InputOrderStruct
	// Input.InstrumentID = sym
	// Input.Direction = direction
	// Input.Price = price
	// Input.Volume = volume
	// s.TraderSpi.OrderClose(Input)
}
func (s *Strategy) OrderCancel(sym string, orderId string) {
	// s.TraderSpi.OrderCancel(sym, orderId)
}

/**
 *   计算盈亏
 *
 * @param   InstrumentID  string  合约
 * @param   OpenPrice     float64 开仓价格
 * @param   LastPrice     float64 最新价|平仓价格
 * @param   Number        int     数量
 * @param   PosiDirection string  持仓方向[2：买，3：卖]
 */

func (s *Strategy) GetPositionProfit(InstrumentID string, OpenPrice float64, LastPrice float64, Number int, PosiDirection string) float64 {

	InstrumentInfo, _ := s.GetInstrumentInfo(InstrumentID)
	if PosiDirection == "2" {
		return ((LastPrice - OpenPrice) * float64(InstrumentInfo.VolumeMultiple)) * float64(Number)
	} else {
		return ((OpenPrice - LastPrice) * float64(InstrumentInfo.VolumeMultiple)) * float64(Number)
	}
}

func (s *Strategy) SetConfig(cfg *Config, ctp *CtpClient, spi *FtdcTraderSpi) {
	//fmt.Printf("Strategy MaxKlen %v", cfg.MaxKlen)
	s.Symbol = cfg.Symbol
	s.Period = cfg.Period
	s.Params = cfg.Params[:]
	s.MaxKLen = cfg.MaxKlen
	s.oCtp = ctp
	s.TraderSpi = spi
	for _, sym := range s.Symbol {
		if len(s.AM[sym]) > 0 {
			s.AM[sym] = s.AM[sym][:0]
		}
		if utils.IsNil(s.AM) {
			s.AM = make(map[string][]Kline)
		}
		if utils.IsNil(s.AM[sym]) {
			s.AM[sym] = make([]Kline, s.MaxKLen)
		}

		if utils.IsNil(s.TK) {
			s.TK = make(map[string]Ticker)
		}
		if utils.IsNil(s.TK[sym]) {
			s.TK[sym] = Ticker{
				DayVolume:    0,
				OpenInterest: 0,
			}
		}
	}
}

var (
	RegisterStrategy = make(map[string]interface{})
	currKline        = make(map[string]Kline)
)

func init() { //所有的策略编写完毕后，需要在这地方进行注册登记
	RegisterStrategy["StrategyEMA"] = &SuperDoubleEMA{}
}

func (c *CtpClient) Register(cfg *Config, spi *FtdcTraderSpi) bool {
	//println(cfg.Class, cfg.Symbol)
	if obj, OK := RegisterStrategy[cfg.Class]; OK {
		t := reflect.TypeOf(obj).Elem()
		c.Strategy = reflect.New(t).Interface().(IStrategy)
		c.Strategy.SetConfig(cfg, c, spi)
		//reflect.TypeOf(obj)
		//_s := c.Strategy.()

		return true
	} else {
		return false
	}
}
