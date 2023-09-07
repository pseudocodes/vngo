package strategy

import (
	"fmt"
	"log"

	"time"

	"github.com/pseudocodes/goctp"
	"github.com/pseudocodes/goctp/thost"

	"github.com/pseudocodes/vngo/pkg/utils"
)

func CreateFtdcMdSpi(c *CtpClient) *FtdcMdSpi {
	p := &FtdcMdSpi{
		CtpClient: c,
		MdSpiLite: &goctp.MdSpiLite{},
	}
	p.SetOnFrontConnected(p.OnFrontConnected)
	p.SetOnFrontDisconnected(p.OnFrontDisconnected)
	p.SetOnHeartBeatWarning(p.OnHeartBeatWarning)
	p.SetOnRspSubMarketData(p.OnRspSubMarketData)
	p.SetOnRspUserLogin(p.OnRspUserLogin)
	p.SetOnRtnDepthMarketData(p.OnRtnDepthMarketData)

	return p
}

// GetMdRequestId 获得行情请求编号
func (p *FtdcMdSpi) GetMdRequestId() int {
	return int(p.MdRequestId.Add(1))
}

// OnFrontDisconnected 当客户端与交易后台通信连接断开时，该方法被调用。当发生这个情况后，API会自动重新连接，客户端可不做处理。
// 服务器已断线，该函数也会被调用。【api 会自动初始化程序，并重新登陆】
func (p *FtdcMdSpi) OnFrontDisconnected(nReason int) {
	log.Println("行情服务器已断线，尝试重新连接中...")
}

// OnFrontConnected 当客户端与交易后台建立起通信连接时（还未登录前），该方法被调用。
func (p *FtdcMdSpi) OnFrontConnected() {

	MdStr := "=================================================================================================\n" +
		"= 行情模块初始化成功, API 版本：" + p.CtpClient.MdApi.GetApiVersion() + "\n" +
		"================================================================================================="
	fmt.Println(MdStr)

	// 登录（如果行情模块在交易模块后初始化则直接登录行情）
	//if p.CtpClient.IsTraderInit {
	go func() {
		for !p.IsTraderInitFinish.Load() {
			time.Sleep(time.Duration(100) * time.Millisecond)
			// fmt.Println("here")
		}
		p.ReqUserLogin()
	}()
	//}
}

// ReqUserLogin 行情用户登录
func (p *FtdcMdSpi) ReqUserLogin() {
	log.Println("行情系统账号登陆中...")
	// req := lib.NewCThostFtdcReqUserLoginField()
	req := &goctp.ReqUserLoginField{
		BrokerID: p.Config.BrokerID,
		UserID:   p.Config.InvestorID,
		Password: p.Config.Password,
	}

	iResult := p.MdApi.ReqUserLogin(req, p.GetMdRequestId())

	if iResult != 0 {
		utils.ReqFailMsg("发送用户登录请求失败！", iResult)
	}
}

// OnRspUserLogin 登录请求响应
func (p *FtdcMdSpi) OnRspUserLogin(pRspUserLogin *goctp.RspUserLoginField, pRspInfo *goctp.RspInfoField, nRequestID int, bIsLast bool) {

	if bIsLast && !p.IsErrorRspInfo(pRspInfo) {
		log.Printf("行情系统登陆成功，当前交易日： %v\n", p.MdApi.GetTradingDay())
		go p.Strategy.OnStart(p)
	}
}

// SubscribeMarketData 订阅行情
func (p *FtdcMdSpi) SubscribeMarketData(InstrumentID []string) int {

	if len(InstrumentID) == 0 {
		log.Println("没有指定需要订阅的行情数据")
		return 0
	}

	fmt.Println("")
	log.Printf("订阅行情数据中... <%+v>", InstrumentID)

	iResult := p.MdApi.SubscribeMarketData(InstrumentID...)

	if iResult != 0 {
		utils.ReqFailMsg("发送订阅行情请求失败！", iResult)
	}

	return iResult
}

// OnRspSubMarketData 订阅行情应答
func (p *FtdcMdSpi) OnRspSubMarketData(pSpecificInstrument string, pRspInfo *goctp.RspInfoField, nRequestID int, bIsLast bool) {
	log.Printf("OnRspSubMarketData \n")

	if !p.IsErrorRspInfo(pRspInfo) {
		log.Printf("订阅合约 %v 行情数据成功！\n", pSpecificInstrument)
	}
	// log.Printf("OnRspSubMarketData \n")
}

// UnSubscribeMarketData 退订行情
func (p *FtdcMdSpi) UnSubscribeMarketData(InstrumentID []string) int {

	if len(InstrumentID) == 0 {
		log.Println("没有指定需要退订的行情数据")
		return 0
	}

	fmt.Println("")
	log.Println("退订行情数据中...")

	iResult := p.MdApi.UnSubscribeMarketData(InstrumentID...)

	if iResult != 0 {
		utils.ReqFailMsg("发送退订行情请求失败！", iResult)
	}

	return iResult
}

// OnRspUnSubMarketData 退订行情应答
func (p *FtdcMdSpi) OnRspUnSubMarketData(pSpecificInstrument *thost.CThostFtdcSpecificInstrumentField, pRspInfo *goctp.RspInfoField, nRequestID int, bIsLast bool) {
	if !p.IsErrorRspInfo(pRspInfo) {
		log.Printf("取消订阅 %v 行情数据成功！\n", string(pSpecificInstrument.InstrumentID[:]))
	}
}

// OnRtnDepthMarketData 深度行情通知
func (p *FtdcMdSpi) OnRtnDepthMarketData(pDepthMarketData *goctp.DepthMarketDataField) {
	t := Ticker{
		InstrumentID: pDepthMarketData.InstrumentID,
		UpdateTime:   utils.StrToTime(pDepthMarketData.ActionDay + " " + pDepthMarketData.UpdateTime),
		Price:        pDepthMarketData.LastPrice,
		DayVolume:    pDepthMarketData.Volume,
		OpenInterest: pDepthMarketData.OpenInterest,
		Average:      pDepthMarketData.AveragePrice,
		OpenPrice:    pDepthMarketData.OpenPrice,
		HighestPrice: pDepthMarketData.HighestPrice,
		LowestPrice:  pDepthMarketData.LowestPrice,
		BidPrice1:    pDepthMarketData.BidPrice1,
		BidVolume1:   pDepthMarketData.BidVolume1,
		AskPrice1:    pDepthMarketData.AskPrice1,
		AskVolume1:   pDepthMarketData.AskVolume1,
	}
	// fmt.Printf("OnRtnDepthMarketData\n")
	//fmt.Printf("%v 合约：%v \t最新价：%v [%v\t%v] \t买一价：%v \t卖一价：%v \t买一量：%v \t卖一量：%v\n", pDepthMarketData.GetUpdateTime(),
	//	pDepthMarketData.GetInstrumentID(), pDepthMarketData.GetLastPrice(), pDepthMarketData.GetVolume(), pDepthMarketData.GetOpenInterest(), pDepthMarketData.GetBidPrice1(), pDepthMarketData.GetAskPrice1(), pDepthMarketData.GetBidVolume1(), pDepthMarketData.GetAskVolume1())
	p.Strategy.OnQuote(pDepthMarketData.InstrumentID, t)
}

// IsErrorRspInfo 行情系统错误通知
func (p *FtdcMdSpi) IsErrorRspInfo(pRspInfo *goctp.RspInfoField) bool {

	// 容错处理 pRspInfo ，部分响应函数中，pRspInfo 为 0
	if pRspInfo == nil {
		return false

	} else {

		// 如果ErrorID != 0, 说明收到了错误的响应
		bResult := pRspInfo.ErrorID != 0
		if bResult {
			// pRspInfo.GetErrorMsg 为 GBK 编码需要自行转成 utf8
			log.Printf("ErrorID=%v ErrorMsg=%v\n", pRspInfo.ErrorID, pRspInfo.ErrorMsg)
		}

		return bResult
	}
}
