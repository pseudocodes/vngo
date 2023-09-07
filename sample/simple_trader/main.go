// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"log"
	"os"
	"sync/atomic"
	"time"

	"github.com/gookit/goutil/dump"
	"github.com/pseudocodes/goctp"
	"github.com/pseudocodes/goctp/thost"
)

var SimnowEnv map[string]map[string]string = map[string]map[string]string{
	"td": {
		"7x24":      "tcp://180.168.146.187:10130",
		"telesim1":  "tcp://180.168.146.187:10201",
		"telesim2":  "tcp://180.168.146.187:10202",
		"moblesim3": "tcp://218.202.237.33:10203",
	},
	"md": {
		"7x24":      "tcp://180.168.146.187:10131",
		"telesim1":  "tcp://180.168.146.187:10211",
		"telesim2":  "tcp://180.168.146.187:10212",
		"moblesim3": "tcp://218.202.237.33:10213",
	},
}

var (
	InvestorID = os.Getenv("SIMNOW_USER_ID")       // <- 环境变量设置
	SimnowPass = os.Getenv("SIMNOW_USER_PASSWORD") // <- 环境变量设置
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

type baseSpi struct {
	tdspi *goctp.TraderSpiLite
	tdapi *goctp.TraderApiLite

	brokerID   string
	investorID string
	password   string
	appid      string
	authCode   string

	requestID atomic.Int32
}

func CreateBaseSpi() *baseSpi {
	s := &baseSpi{
		// tdapi: tdapi,
		tdspi: &goctp.TraderSpiLite{},

		brokerID:   "9999",
		investorID: InvestorID, // <- 环境变量设置
		password:   SimnowPass, // <- 环境变量设置
		appid:      "simnow_client_test",
		authCode:   "0000000000000000",
	}

	s.tdspi.SetOnFrontConnected(s.OnFrontConnected)

	s.tdspi.SetOnFrontDisconnected(s.OnFrontDisconnected)

	s.tdspi.SetOnHeartBeatWarning(s.OnHeartBeatWarning)

	s.tdspi.SetOnRspAuthenticate(s.OnRspAuthenticate)

	s.tdspi.SetOnRspUserLogin(s.OnRspUserLogin)

	s.tdspi.SetOnRspOrderInsert(s.OnRspOrderInsert)

	s.tdspi.SetOnRspOrderAction(s.OnRspOrderAction)

	s.tdspi.SetOnRspSettlementInfoConfirm(s.OnRspSettlementInfoConfirm)

	s.tdspi.SetOnRspQryOrder(s.OnRspQryOrder)

	s.tdspi.SetOnRspQryInvestorPosition(s.OnRspQryInvestorPosition)

	s.tdspi.SetOnRspQryTradingAccount(s.OnRspQryTradingAccount)

	// s.tdspi.SetOnRspQryInstrumentMarginRate(s.OnRspQryInstrumentMarginRate)

	// s.tdspi.SetOnRspQryInstrumentCommissionRate(s.OnRspQryInstrumentCommissionRate)

	s.tdspi.SetOnRspQryInstrument(s.OnRspQryInstrument)

	// s.tdspi.SetOnRspQrySettlementInfo(s.OnRspQrySettlementInfo)

	s.tdspi.SetOnRspError(s.OnRspError)

	s.tdspi.SetOnRtnOrder(s.OnRtnOrder)

	s.tdspi.SetOnRtnTrade(s.OnRtnTrade)

	// s.tdspi.SetOnErrRtnOrderInsert(s.OnErrRtnOrderInsert)

	s.tdspi.SetOnErrRtnOrderAction(s.OnErrRtnOrderAction)

	s.tdspi.SetOnRtnInstrumentStatus(s.OnRtnInstrumentStatus)
	return s
}

func (s *baseSpi) OnFrontDisconnected(nReason int) {
	log.Printf("OnFrontDissconnected: %v\n", nReason)
}

func (p *baseSpi) OnHeartBeatWarning(nTimeLapse int) {
	log.Println("(OnHeartBeatWarning) nTimerLapse=", nTimeLapse)
}

func (s *baseSpi) OnFrontConnected() {
	var ret int
	log.Printf("OnFrontConnected\n")

	ret = s.tdapi.ReqAuthenticate(&goctp.ReqAuthenticateField{
		BrokerID: s.brokerID,
		UserID:   s.investorID,
		AuthCode: s.authCode,
		AppID:    s.appid,
	}, int(s.requestID.Add(1)))

	log.Printf("user auth: %v\n", ret)
}

func (s *baseSpi) OnRspAuthenticate(f *goctp.RspAuthenticateField, r *goctp.RspInfoField, nRequestID int, bIsLast bool) {
	dump.Println(r)
	dump.Println(f)
	req := &goctp.ReqUserLoginField{
		BrokerID: "9999",
		UserID:   os.Getenv("SIMNOW_USER_ID"), // <- 环境变量设置
		Password: os.Getenv("SIMNOW_USER_PASSWORD"),
	}
	dump.Println(req)
	ret := s.tdapi.ReqUserLogin(req, int(s.requestID.Add(1)))
	log.Printf("user login: %v\n", ret)
}

func (s *baseSpi) OnRspUserLogin(f *goctp.RspUserLoginField, r *goctp.RspInfoField, nRequestID int, bIsLast bool) {
	dump.Println(r)
	dump.Println(f)

	req := &goctp.SettlementInfoConfirmField{
		BrokerID:   s.brokerID,
		InvestorID: s.investorID,
	}
	ret := s.tdapi.ReqSettlementInfoConfirm(req, int(s.requestID.Add(1)))
	log.Printf("req_settlement_info_confirm : %v\n", ret)

}

func (s *baseSpi) OnRspSettlementInfoConfirm(f *goctp.SettlementInfoConfirmField, r *goctp.RspInfoField, nRequestID int, bIsLast bool) {
	dump.Println(r)
	dump.Println(f)

	req := &goctp.QryInstrumentField{}
	ret := s.tdapi.ReqQryInstrument(req, 3)
	log.Printf("user qry ins: %v\n", ret)
}

// OnRspSettlementInfoConfirm 发送投资者结算单确认响应

func (s *baseSpi) OnRspQryInstrument(pInstrument *goctp.InstrumentField, pRspInfo *goctp.RspInfoField, nRequestID int, bIsLast bool) {
	// dump.Print(pRspInfo, nRequestID, bIsLast)
	// dump.Println(pInstrument.InstrumentName)

	if bIsLast {
		log.Printf("qry ins finished\n")

		req := &goctp.QryTradingAccountField{
			BrokerID:   s.brokerID,
			InvestorID: s.investorID,
		}
		time.Sleep(1500 * time.Millisecond)
		ret := s.tdapi.ReqQryTradingAccount(req, int(s.requestID.Add(1)))
		if ret != 0 {
			log.Printf("req_qry_trading_account failed %v\n", ret)
		}
	}
}

// OnRspQryTradingAccount 请求查询资金账户响应
func (s *baseSpi) OnRspQryTradingAccount(pTradingAccount *goctp.TradingAccountField, pRspInfo *goctp.RspInfoField, nRequestID int, bIsLast bool) {

	if bIsLast && !s.isErrorRspInfo(pRspInfo) {
		dump.P(pTradingAccount)

		req := &goctp.QryOrderField{
			BrokerID:   s.brokerID,
			InvestorID: s.investorID,
		}
		time.Sleep(1500 * time.Millisecond)
		ret := s.tdapi.ReqQryOrder(req, int(s.requestID.Add(1)))
		if ret != 0 {
			log.Printf("req_qry_order failed: %v\n", ret)
		}
	}
}

// 合约交易状态通知
func (s *baseSpi) OnRtnInstrumentStatus(pInstrumentStatus *goctp.InstrumentStatusField) {
	// dump.P(pInstrumentStatus)
}

func (s *baseSpi) OnRspQryOrder(pOrder *goctp.OrderField, pRspInfo *goctp.RspInfoField, nRequestID int, bIsLast bool) {
	if s.isErrorRspInfo(pRspInfo) {
		return
	}
	if pOrder == nil {
		log.Printf("OnRspQryOrder: %v\n", pOrder)
	}
	if bIsLast {
		req := &goctp.QryInvestorPositionField{
			BrokerID:   s.brokerID,
			InvestorID: s.investorID,
		}
		time.Sleep(1500 * time.Millisecond)
		ret := s.tdapi.ReqQryInvestorPosition(req, int(s.requestID.Add(1)))
		if ret != 0 {
			log.Printf("req_qry_investor_position failed %v\n", ret)
		}
	}
}

func (s *baseSpi) OnRspQryInvestorPosition(pInvestorPosition *goctp.InvestorPositionField, pRspInfo *goctp.RspInfoField, nRequestID int, bIsLast bool) {
	if s.isErrorRspInfo(pRspInfo) {
		return
	}
	if bIsLast {
		log.Printf("finish rsp_qry_investor_position\n")
		s.requestID.Add(1)
		order := &goctp.InputOrderField{
			BrokerID:   "9999",
			InvestorID: s.investorID,
			// UserID:     s.investorID,
			// ExchangeID:          "SHFE",
			InstrumentID:        "ag2306",
			OrderPriceType:      thost.THOST_FTDC_OPT_LimitPrice,
			Direction:           thost.THOST_FTDC_D_Buy,
			CombOffsetFlag:      string(thost.THOST_FTDC_OF_Open),
			CombHedgeFlag:       string(thost.THOST_FTDC_HF_Speculation),
			LimitPrice:          5670,
			VolumeTotalOriginal: 1,
			TimeCondition:       thost.THOST_FTDC_TC_GFD,
			VolumeCondition:     thost.THOST_FTDC_VC_AV,
			ContingentCondition: thost.THOST_FTDC_CC_Immediately,
			ForceCloseReason:    thost.THOST_FTDC_FCC_NotForceClose,
		}
		ret := s.tdapi.ReqOrderInsert(order, int(s.requestID.Load()))
		log.Printf("input order ret %v\n", ret)
	}
}

func (s *baseSpi) OnRspOrderInsert(pOrder *goctp.InputOrderField, pRspInfo *goctp.RspInfoField, requestID int, isLast bool) {
	s.isErrorRspInfo(pRspInfo)
}

// 错误应答
func (s *baseSpi) OnRspError(pRspInfo *goctp.RspInfoField, nRequestID int, bIsLast bool) {
	s.isErrorRspInfo(pRspInfo)
}

// 报单操作错误回报
func (s *baseSpi) OnErrRtnOrderAction(pOrderAction *goctp.OrderActionField, pRspInfo *goctp.RspInfoField) {
	s.isErrorRspInfo(pRspInfo)
}

// 报单操作请求响应（撤单失败会触发）
func (s *baseSpi) OnRspOrderAction(pInputOrderAction *goctp.InputOrderActionField, pRspInfo *goctp.RspInfoField, nRequestID int, bIsLast bool) {
	s.isErrorRspInfo(pRspInfo)
}

// OnRtnTrade 成交通知（委托单在交易所成交了）
func (s *baseSpi) OnRtnTrade(pTrade *goctp.TradeField) {
	dump.V(pTrade)
}

// OnRtnOrder 报单通知（委托单）
func (s *baseSpi) OnRtnOrder(pOrder *goctp.OrderField) {
	dump.V(pOrder)

}

func (s *baseSpi) isErrorRspInfo(pRspInfo *goctp.RspInfoField) bool {

	// 容错处理 pRspInfo ，部分响应函数中，pRspInfo 为 0
	if pRspInfo == nil {
		return false
	}
	// 如果ErrorID != 0, 说明收到了错误的响应
	bResult := (pRspInfo.ErrorID != 0)
	if bResult {
		log.Printf("ErrorID=%v ErrorMsg=%v\n", pRspInfo.ErrorID, pRspInfo.ErrorMsg)
	}
	return bResult

}

func sample1() {

	tdapi := goctp.CreateTraderApiLite(goctp.TraderFlowPath("./data/"))
	baseSpi := CreateBaseSpi()
	baseSpi.tdapi = tdapi
	log.Printf("baseSpi %+v\n", baseSpi)
	tdapi.RegisterSpi(baseSpi.tdspi)
	// tdapi.RegisterFront(SimnowEnv["td"]["7x24"])
	tdapi.RegisterFront(SimnowEnv["td"]["telesim1"])

	tdapi.Init()

	println(tdapi.GetTradingDay())
	println(tdapi.GetApiVersion())

	tdapi.Join()
}

func main() {
	sample1()
	// sample2()
}
