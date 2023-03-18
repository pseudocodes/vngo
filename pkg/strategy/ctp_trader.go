package strategy

import (
	"fmt"
	"log"
	"time"

	"github.com/pseudocodes/goctp"
	"github.com/pseudocodes/goctp/thost"
	"github.com/pseudocodes/vngo/pkg/utils"
)

func CreateFtdcTraderSpi(c *CtpClient) *FtdcTraderSpi {
	p := &FtdcTraderSpi{
		CtpClient:     c,
		TraderSpiLite: &goctp.TraderSpiLite{},
	}
	p.SetOnFrontConnected(p.OnFrontConnected)
	p.SetOnFrontDisconnected(p.OnFrontDisconnected)
	p.SetOnHeartBeatWarning(p.OnHeartBeatWarning)
	p.SetOnRspAuthenticate(p.OnRspAuthenticate)
	p.SetOnRspUserLogin(p.OnRspUserLogin)
	p.SetOnRspOrderAction(p.OnRspOrderAction)
	p.SetOnRspOrderInsert(p.OnRspOrderInsert)
	p.SetOnRspSettlementInfoConfirm(p.OnRspSettlementInfoConfirm)
	p.SetOnRspQryInstrument(p.OnRspQryInstrument)
	p.SetOnRspQryTradingAccount(p.OnRspQryTradingAccount)
	p.SetOnRspQryOrder(p.OnRspQryOrder)
	p.SetOnRspQryInvestorPosition(p.OnRspQryInvestorPosition)

	p.SetOnRtnOrder(p.OnRtnOrder)
	p.SetOnRtnTrade(p.OnRtnTrade)
	p.SetOnErrRtnOrderAction(p.OnErrRtnOrderAction)
	p.SetOnRspError(p.OnRspError)

	return p
}

// GetTraderRequestId 获得交易请求编号
func (p *FtdcTraderSpi) GetTraderRequestId() int {
	return int(p.TraderRequestId.Add(1))
}

// OnFrontDisconnected
//   当客户端与交易后台通信连接断开时，该方法被调用。
//   当发生这个情况后，API会自动重新连接，客户端可不做处理。

func (p *FtdcTraderSpi) OnFrontDisconnected(nReason int) {

	p.IsTraderLogin.Store(false)
	p.IsTraderInit.Store(false)
	p.IsTraderInitFinish.Store(false)
	log.Println("交易服务器已断线，尝试重新连接中...")
}

// ReqMsg 发送请求日志（仅查询类的函数需要调用）
func (p *FtdcTraderSpi) ReqMsg(Msg string) {

	// 交易程序未初始化完成时，执行查询类的函数需要有1.5秒间隔
	if !p.IsTraderInitFinish.Load() {
		time.Sleep(time.Millisecond * 1500)
	}

	fmt.Println("")
	log.Println(Msg)
}

// OnFrontConnected 当客户端与交易后台建立起通信连接时（还未登录前），该方法被调用。
func (p *FtdcTraderSpi) OnFrontConnected() {

	TraderStr := "=================================================================================================\n" +
		"= 交易模块初始化成功, API 版本：" + p.CtpClient.TraderApi.GetApiVersion() + "\n" +
		"================================================================================================="
	fmt.Println(TraderStr)

	p.IsTraderInit.Store(true)

	// 填写了 AppID 与 AuthCode 则进行客户端认证
	if p.Config.AppID != "" && p.Config.AuthCode != "" {
		p.ReqAuthenticate()
	} else {
		p.ReqUserLogin()
	}
}

// ReqAuthenticate 客户端认证
func (p *FtdcTraderSpi) ReqAuthenticate() {

	log.Println("客户端认证中...")

	req := &goctp.ReqAuthenticateField{
		BrokerID: p.Config.BrokerID,
		UserID:   p.Config.InvestorID, // <- 环境变量设置
		AuthCode: p.Config.AppID,
		AppID:    p.Config.AuthCode,
	}
	iResult := p.TraderApi.ReqAuthenticate(req, p.GetTraderRequestId())

	if iResult != 0 {
		utils.ReqFailMsg("发送客户端认证请求失败！", iResult)
	}
}

// OnRspAuthenticate 客户端认证响应
func (p *FtdcTraderSpi) OnRspAuthenticate(f *goctp.RspAuthenticateField, pRspInfo *goctp.RspInfoField, nRequestID int, bIsLast bool) {

	if bIsLast && !p.IsErrorRspInfo(pRspInfo) {

		log.Println("客户端认证成功！")

		p.ReqUserLogin()
	}
}

// ReqUserLogin 用户登录请求
func (p *FtdcTraderSpi) ReqUserLogin() {

	time.Sleep(time.Second * 1)

	log.Println("交易系统账号登陆中...")

	req := &goctp.ReqUserLoginField{
		BrokerID: p.Config.BrokerID,
		UserID:   p.Config.InvestorID, // <- 环境变量设置
		Password: p.Config.Password,
	}

	iResult := p.TraderApi.ReqUserLogin(req, p.GetTraderRequestId())

	if iResult != 0 {
		utils.ReqFailMsg("发送用户登录请求失败！", iResult)
	}
}

func (p *FtdcTraderSpi) OnRspUserLogin(f *goctp.RspUserLoginField, pRspInfo *goctp.RspInfoField, nRequestID int, bIsLast bool) {

	if bIsLast && !p.IsErrorRspInfo(pRspInfo) {

		p.IsTraderLogin.Store(true)

		log.Printf("交易账号已登录，当前交易日：%v\n", p.TraderApi.GetTradingDay())

		p.ReqSettlementInfoConfirm()
	}
}

// ReqSettlementInfoConfirm 投资者结算单确认
func (p *FtdcTraderSpi) ReqSettlementInfoConfirm() int {

	p.ReqMsg("投资者结算单确认中...")

	req := &goctp.SettlementInfoConfirmField{
		BrokerID:   p.Config.BrokerID,
		InvestorID: p.Config.InvestorID,
	}

	iResult := p.CtpClient.TraderApi.ReqSettlementInfoConfirm(req, p.GetTraderRequestId())

	if iResult != 0 {
		utils.ReqFailMsg("确认投资者结算单失败！", iResult)
	}

	return iResult
}

// OnRspSettlementInfoConfirm 发送投资者结算单确认响应
func (p *FtdcTraderSpi) OnRspSettlementInfoConfirm(f *goctp.SettlementInfoConfirmField, pRspInfo *goctp.RspInfoField, nRequestID int, bIsLast bool) {

	if bIsLast && !p.IsErrorRspInfo(pRspInfo) {
		log.Println("投资者结算单确认成功")

		p.ReqQryInstrument()
	}
}

// ReqQryInstrument 请求查询合约
func (p *FtdcTraderSpi) ReqQryInstrument() int {

	p.ReqMsg("查询合约中...")

	req := &goctp.QryInstrumentField{}

	iResult := p.TraderApi.ReqQryInstrument(req, p.GetTraderRequestId())

	if iResult != 0 {
		utils.ReqFailMsg("查询合约失败！", iResult)
	}

	return iResult
}

// OnRspQryInstrument 请求查询合约响应
func (p *FtdcTraderSpi) OnRspQryInstrument(pInstrument *goctp.InstrumentField, pRspInfo *goctp.RspInfoField, nRequestID int, bIsLast bool) {

	if !p.IsErrorRspInfo(pRspInfo) {

		var mInstrumentInfo InstrumentInfoStruct

		var mapKey string = pInstrument.InstrumentID

		mInstrumentInfo.InstrumentID = pInstrument.InstrumentID
		mInstrumentInfo.ExchangeID = pInstrument.ExchangeID
		mInstrumentInfo.InstrumentName = pInstrument.InstrumentName
		mInstrumentInfo.ExchangeInstID = pInstrument.ExchangeInstID
		mInstrumentInfo.ProductID = pInstrument.ProductID
		mInstrumentInfo.ProductClass = string(pInstrument.ProductClass)
		mInstrumentInfo.DeliveryYear = pInstrument.DeliveryYear
		mInstrumentInfo.DeliveryMonth = pInstrument.DeliveryMonth
		mInstrumentInfo.MaxMarketOrderVolume = pInstrument.MaxMarketOrderVolume
		mInstrumentInfo.MinMarketOrderVolume = pInstrument.MinMarketOrderVolume
		mInstrumentInfo.MaxLimitOrderVolume = pInstrument.MaxLimitOrderVolume
		mInstrumentInfo.MinLimitOrderVolume = pInstrument.MinLimitOrderVolume
		mInstrumentInfo.VolumeMultiple = pInstrument.VolumeMultiple
		mInstrumentInfo.PriceTick = pInstrument.PriceTick
		mInstrumentInfo.CreateDate = pInstrument.CreateDate
		mInstrumentInfo.OpenDate = pInstrument.OpenDate
		mInstrumentInfo.ExpireDate = pInstrument.ExpireDate
		mInstrumentInfo.StartDelivDate = pInstrument.StartDelivDate
		mInstrumentInfo.EndDelivDate = pInstrument.EndDelivDate
		mInstrumentInfo.InstLifePhase = string(pInstrument.InstLifePhase)
		mInstrumentInfo.IsTrading = pInstrument.IsTrading
		mInstrumentInfo.PositionType = string(pInstrument.PositionType)
		mInstrumentInfo.PositionDateType = string(pInstrument.PositionDateType)
		mInstrumentInfo.LongMarginRatio = pInstrument.LongMarginRatio
		mInstrumentInfo.ShortMarginRatio = pInstrument.ShortMarginRatio
		mInstrumentInfo.MaxMarginSideAlgorithm = string(pInstrument.MaxMarginSideAlgorithm)
		mInstrumentInfo.UnderlyingInstrID = pInstrument.UnderlyingInstrID
		mInstrumentInfo.StrikePrice = pInstrument.StrikePrice
		mInstrumentInfo.OptionsType = string(pInstrument.OptionsType)
		mInstrumentInfo.UnderlyingMultiple = pInstrument.UnderlyingMultiple
		mInstrumentInfo.CombinationType = string(pInstrument.CombinationType)

		// log.Printf("mapkey %v\n", mapKey)
		p.Strategy.SetInstruments(mapKey, mInstrumentInfo)

		if bIsLast {

			log.Printf("合约记录初始化完毕！")

			if !p.IsTraderInitFinish.Load() {
				// 请求查询资金账户
				p.ReqQryTradingAccount()
			}
		}
	}
}

// ReqQryTradingAccount 请求查询资金账户
func (p *FtdcTraderSpi) ReqQryTradingAccount() int {

	p.ReqMsg("查询资金账户中...")

	req := &goctp.QryTradingAccountField{
		BrokerID:   p.Config.BrokerID,
		InvestorID: p.Config.InvestorID,
	}
	iResult := p.TraderApi.ReqQryTradingAccount(req, p.GetTraderRequestId())

	if iResult != 0 {
		utils.ReqFailMsg("查询资金账户失败！", iResult)
	}

	return iResult
}

// OnRspQryTradingAccount 请求查询资金账户响应
func (p *FtdcTraderSpi) OnRspQryTradingAccount(pTradingAccount *goctp.TradingAccountField, pRspInfo *goctp.RspInfoField, nRequestID int, bIsLast bool) {

	if bIsLast && !p.IsErrorRspInfo(pRspInfo) {

		var mAccountInfo AccountInfoStruct

		mAccountInfo.MapKey = pTradingAccount.BrokerID + "_" + pTradingAccount.AccountID

		mAccountInfo.BrokerID = pTradingAccount.BrokerID
		mAccountInfo.AccountID = pTradingAccount.AccountID
		mAccountInfo.PreMortgage = utils.Decimal(pTradingAccount.PreMortgage, 2)
		mAccountInfo.PreCredit = utils.Decimal(pTradingAccount.PreCredit, 2)
		mAccountInfo.PreDeposit = utils.Decimal(pTradingAccount.PreDeposit, 2)
		mAccountInfo.PreBalance = utils.Decimal(pTradingAccount.PreBalance, 2)
		mAccountInfo.PreMargin = utils.Decimal(pTradingAccount.PreMargin, 2)
		mAccountInfo.InterestBase = utils.Decimal(pTradingAccount.InterestBase, 2)
		mAccountInfo.Interest = utils.Decimal(pTradingAccount.Interest, 2)
		mAccountInfo.Deposit = utils.Decimal(pTradingAccount.Deposit, 2)
		mAccountInfo.Withdraw = utils.Decimal(pTradingAccount.Withdraw, 2)
		mAccountInfo.FrozenMargin = utils.Decimal(pTradingAccount.FrozenMargin, 2)
		mAccountInfo.FrozenCash = utils.Decimal(pTradingAccount.FrozenCash, 2)
		mAccountInfo.FrozenCommission = utils.Decimal(pTradingAccount.FrozenCommission, 2)
		mAccountInfo.CurrMargin = utils.Decimal(pTradingAccount.CurrMargin, 2)
		mAccountInfo.CashIn = utils.Decimal(pTradingAccount.CashIn, 2)
		mAccountInfo.Commission = utils.Decimal(pTradingAccount.Commission, 2)
		mAccountInfo.CloseProfit = utils.Decimal(pTradingAccount.CloseProfit, 2)
		mAccountInfo.PositionProfit = utils.Decimal(pTradingAccount.PositionProfit, 2)
		mAccountInfo.Balance = utils.Decimal(pTradingAccount.Balance, 2)
		mAccountInfo.Available = utils.Decimal(pTradingAccount.Available, 2)
		mAccountInfo.WithdrawQuota = utils.Decimal(pTradingAccount.WithdrawQuota, 2)
		mAccountInfo.Reserve = utils.Decimal(pTradingAccount.Reserve, 2)
		mAccountInfo.TradingDay = pTradingAccount.TradingDay
		mAccountInfo.SettlementID = pTradingAccount.SettlementID
		mAccountInfo.Credit = utils.Decimal(pTradingAccount.Credit, 2)
		mAccountInfo.Mortgage = utils.Decimal(pTradingAccount.Mortgage, 2)
		mAccountInfo.ExchangeMargin = utils.Decimal(pTradingAccount.ExchangeMargin, 2)
		mAccountInfo.DeliveryMargin = utils.Decimal(pTradingAccount.DeliveryMargin, 2)
		mAccountInfo.ExchangeDeliveryMargin = utils.Decimal(pTradingAccount.ExchangeDeliveryMargin, 2)
		mAccountInfo.ReserveBalance = utils.Decimal(pTradingAccount.ReserveBalance, 2)
		mAccountInfo.CurrencyID = pTradingAccount.CurrencyID
		mAccountInfo.PreFundMortgageIn = utils.Decimal(pTradingAccount.PreFundMortgageIn, 2)
		mAccountInfo.PreFundMortgageOut = utils.Decimal(pTradingAccount.PreFundMortgageOut, 2)
		mAccountInfo.FundMortgageIn = utils.Decimal(pTradingAccount.FundMortgageIn, 2)
		mAccountInfo.FundMortgageOut = utils.Decimal(pTradingAccount.FundMortgageOut, 2)
		mAccountInfo.FundMortgageAvailable = utils.Decimal(pTradingAccount.FundMortgageAvailable, 2)
		mAccountInfo.MortgageableFund = utils.Decimal(pTradingAccount.MortgageableFund, 2)
		mAccountInfo.SpecProductMargin = utils.Decimal(pTradingAccount.SpecProductMargin, 2)
		mAccountInfo.SpecProductFrozenMargin = utils.Decimal(pTradingAccount.SpecProductFrozenMargin, 2)
		mAccountInfo.SpecProductCommission = utils.Decimal(pTradingAccount.SpecProductCommission, 2)
		mAccountInfo.SpecProductFrozenCommission = utils.Decimal(pTradingAccount.SpecProductFrozenCommission, 2)
		mAccountInfo.SpecProductPositionProfit = utils.Decimal(pTradingAccount.SpecProductPositionProfit, 2)
		mAccountInfo.SpecProductCloseProfit = utils.Decimal(pTradingAccount.SpecProductCloseProfit, 2)
		mAccountInfo.SpecProductPositionProfitByAlg = utils.Decimal(pTradingAccount.SpecProductPositionProfitByAlg, 2)
		mAccountInfo.SpecProductExchangeMargin = utils.Decimal(pTradingAccount.SpecProductExchangeMargin, 2)
		mAccountInfo.BizType = string(pTradingAccount.BizType)
		mAccountInfo.FrozenSwap = utils.Decimal(pTradingAccount.FrozenSwap, 2)
		mAccountInfo.RemainSwap = utils.Decimal(pTradingAccount.RemainSwap, 2)

		AccountInfoStr := "-------------------------------------------------------------------------------------------------\n" +
			"- 公司代码：" + pTradingAccount.BrokerID + "\n" +
			"- 资金账号：" + pTradingAccount.AccountID + "\n" +
			"- 期初资金：" + utils.Float64ToString(mAccountInfo.PreBalance) + "\n" +
			"- 动态权益：" + utils.Float64ToString(mAccountInfo.Balance) + "\n" +
			"- 可用资金：" + utils.Float64ToString(mAccountInfo.Available) + "\n" +
			"- 持仓盈亏：" + utils.Float64ToString(mAccountInfo.PositionProfit) + "\n" +
			"- 平仓盈亏：" + utils.Float64ToString(mAccountInfo.CloseProfit) + "\n" +
			"- 手续费  ：" + utils.Float64ToString(mAccountInfo.Commission) + "\n" +
			"-------------------------------------------------------------------------------------------------"
		fmt.Println(AccountInfoStr)

		if !p.IsTraderInitFinish.Load() {
			// 请求查询投资者报单（委托单）
			p.ReqQryOrder()
		}
	}
}

// ReqQryOrder 请求查询投资者报单（委托单）
func (p *FtdcTraderSpi) ReqQryOrder() int {

	p.ReqMsg("查询投资者报单中...")

	// req := lib.NewCThostFtdcQryOrderField()
	req := &goctp.QryOrderField{
		BrokerID:   p.Config.BrokerID,
		InvestorID: p.Config.InvestorID,
	}

	iResult := p.TraderApi.ReqQryOrder(req, p.GetTraderRequestId())

	if iResult != 0 {
		utils.ReqFailMsg("查询投资者报单失败！", iResult)
	}

	return iResult
}

// OnRspQryOrder 请求查询投资者报单响应
func (p *FtdcTraderSpi) OnRspQryOrder(pOrder *goctp.OrderField, pRspInfo *goctp.RspInfoField, nRequestID int, bIsLast bool) {

	if !p.IsErrorRspInfo(pRspInfo) {

		// 如果 没有数据 pOrder 等于0 // nil
		// pOrderCode := fmt.Sprintf("%v", pOrder)

		// 只记录有报单编号的报单数据
		if pOrder != nil && pOrder.OrderSysID != "" {
			// 获得报单结构体数据
			mOrder := GetOrderListStruct(pOrder)

			// 报单列表数据 key 键
			mOrder.MapKey = pOrder.InstrumentID + "_" + utils.TrimSpace(pOrder.OrderSysID)

			// 记录报单数据
			p.MapOrderList.Set(mOrder.MapKey, mOrder)
		}

		if bIsLast {

			fmt.Println("-------------------------------------------------------------------------------------------------")

			MapOrderNoTradeSize := 0

			mOrderList := p.MapOrderList.GetAll()
			for _, v := range mOrderList {
				val := v.(OrderListStruct)

				// 输出 未成交、部分成交 的报单
				if val.OrderStatus == string(thost.THOST_FTDC_OST_NoTradeQueueing) || val.OrderStatus == string(thost.THOST_FTDC_OST_PartTradedQueueing) {
					MapOrderNoTradeSize += 1
					fmt.Printf("- 合约：%v   \t%v:%v   \t数量: %v   \t价格: %v   \t报单编号: %v (%v)\n", val.InstrumentID, val.DirectionTitle, val.CombOffsetFlagTitle, val.Volume, val.LimitPrice, utils.TrimSpace(val.OrderSysID), val.OrderStatusTitle)
				}
			}

			fmt.Printf("- 共有报单记录 %v 条，未成交 %v 条（不含错单）\n", p.MapOrderList.Size(), MapOrderNoTradeSize)
			fmt.Println("-------------------------------------------------------------------------------------------------")

			if !p.IsTraderInitFinish.Load() {
				// 请求查询投资者持仓（汇总）
				p.ReqQryInvestorPosition()
			}
		}
	}
}

// ReqQryInvestorPosition 请求查询投资者持仓（汇总）
func (p *FtdcTraderSpi) ReqQryInvestorPosition() int {

	p.ReqMsg("查询投资者持仓中...")

	req := &goctp.QryInvestorPositionField{
		BrokerID:   p.Config.BrokerID,
		InvestorID: p.Config.InvestorID,
	}

	iResult := p.TraderApi.ReqQryInvestorPosition(req, p.GetTraderRequestId())

	if iResult != 0 {
		utils.ReqFailMsg("查询投资者持仓失败！", iResult)
	}

	fmt.Println("-------------------------------------------------------------------------------------------------")
	p.Strategy.EmptyPosition()
	return iResult
}

// OnRspQryInvestorPosition 请求查询投资者持仓（汇总）响应
func (p *FtdcTraderSpi) OnRspQryInvestorPosition(pInvestorPosition *goctp.InvestorPositionField, pRspInfo *goctp.RspInfoField, nRequestID int, bIsLast bool) {

	if !p.IsErrorRspInfo(pRspInfo) {

		// 没有数据 pInvestorPosition 会等于 0
		// pInvestorPositionCode := fmt.Sprintf("%v", pInvestorPosition)

		if pInvestorPosition != nil {

			// 获得持仓结构体数据
			mInvestorPosition := p.GetInvestorPositionStruct(pInvestorPosition)

			if mInvestorPosition.Position != 0 {
				p.Strategy.UpdatePosition(mInvestorPosition)
				fmt.Printf("- 合约：%v   \t%v:%v   \t总持仓: %v   \t持仓均价: %v   \t持仓盈亏：%v\n", mInvestorPosition.InstrumentID, mInvestorPosition.PositionDateTitle, mInvestorPosition.PosiDirectionTitle, mInvestorPosition.Position, mInvestorPosition.OpenCost, mInvestorPosition.PositionProfit)
			}
		}

		if bIsLast {

			fmt.Println("-------------------------------------------------------------------------------------------------")

			if !p.IsTraderInitFinish.Load() {
				// 交易程序初始化流程走完了
				p.IsTraderInitFinish.Store(true)
				// 订阅行情Subscribe := []string{"FG209"}
				//p.MdSpi.SubscribeMarketData(p.SubSymbols)
				//p.Strategy.OnStart(p)
			}
		}
	}
}

// OnRtnOrder 报单通知（委托单）
func (p *FtdcTraderSpi) OnRtnOrder(pOrder *goctp.OrderField) {

	// 报单编号
	OrderSysID := pOrder.OrderSysID

	// 报单状态
	OrderStatus := pOrder.OrderStatus

	// 获得报单结构体数据
	mOrder := GetOrderListStruct(pOrder)

	// 报单列表数据 key 键
	mOrder.MapKey = pOrder.InstrumentID + "_" + utils.TrimSpace(pOrder.OrderSysID)

	if OrderSysID == "" {

		// 报单就自动撤单，且没有编号的 都视为报错
		if OrderStatus == thost.THOST_FTDC_OST_Canceled {

			OrderErrorStr := "-------------------------------------------------------------------------------------------------\n" +
				"- 报单出错了\n" +
				"- 报单合约：" + mOrder.InstrumentID + "\t报单引用: " + mOrder.OrderRef + "\n" +
				"- 报单方向：" + mOrder.DirectionTitle + "   \t报单价格: " + utils.Float64ToString(mOrder.LimitPrice) + "\n" +
				"- 报单开平：" + mOrder.CombOffsetFlagTitle + " \t报单数量: " + utils.IntToString(mOrder.Volume) + "\n" +
				"- 错误代码：-1   \t错误消息: " + mOrder.StatusMsg + "\n" +
				"-------------------------------------------------------------------------------------------------"
			fmt.Println(OrderErrorStr)
		}

		return
	}

	// 未成交和撤单的报单（已成交的通知在 OnRtnTrade 函数中处理）
	if OrderStatus == thost.THOST_FTDC_OST_NoTradeQueueing || OrderStatus == thost.THOST_FTDC_OST_Canceled {

		OrderStr := "-------------------------------------------------------------------------------------------------\n" +
			"- 报单通知 " + mOrder.InsertTime + "\n" +
			"- 报单合约：" + mOrder.InstrumentID + " \t报单编号: " + utils.TrimSpace(mOrder.OrderSysID) + "\n" +
			"- 报单方向：" + mOrder.DirectionTitle + "   \t报单价格: " + utils.Float64ToString(mOrder.LimitPrice) + "\n" +
			"- 报单开平：" + mOrder.CombOffsetFlagTitle + " \t报单数量: " + utils.IntToString(mOrder.Volume) + "\n" +
			"- 报单状态：" + mOrder.OrderStatusTitle + " \t状态信息: " + mOrder.StatusMsg + "\n" +
			"-------------------------------------------------------------------------------------------------"
		fmt.Println(OrderStr)
	}
	p.Strategy.OnOrderChange(mOrder)
	// 将报单数据记录下来
	p.MapOrderList.Set(mOrder.MapKey, mOrder)
}

// OnRtnTrade 成交通知（委托单在交易所成交了）
func (p *FtdcTraderSpi) OnRtnTrade(pTrade *goctp.TradeField) {

	// 报单方向
	DirectionTitle := utils.GetDirectionTitle(string(pTrade.Direction))

	// 报单开平
	OffsetFlagTitle := utils.GetOffsetFlagTitle(string(pTrade.OffsetFlag))
	//_, r := utils.Find(main.SubSymbols, pTrade.GetInstrumentID())
	//if !r {
	//	main.MdSpi.SubscribeMarketData([]string{pTrade.GetInstrumentID()}) // 订阅新行情
	//}
	p.Strategy.OnTradeDeal(pTrade)
	OrderStr := "-------------------------------------------------------------------------------------------------\n" +
		"- 成交通知 " + pTrade.TradeTime + "\n" +
		"- 成交合约：" + pTrade.InstrumentID + "\t成交编号: " + utils.TrimSpace(pTrade.TradeID) + " \t报单编号: " + utils.TrimSpace(pTrade.OrderSysID) + "\n" +
		"- 成交方向：" + DirectionTitle + "   \t成交价格: " + utils.Float64ToString(pTrade.Price) + "\n" +
		"- 成交开平：" + OffsetFlagTitle + " \t成交数量: " + utils.IntToString(pTrade.Volume) + "\n" +
		"-------------------------------------------------------------------------------------------------"
	fmt.Println(OrderStr)
}

// OnRspOrderInsert 报单出错响应（综合交易平台交易核心返回的包含错误信息的报单响应）
func (p *FtdcTraderSpi) OnRspOrderInsert(pInputOrder *goctp.InputOrderField, pRspInfo *goctp.RspInfoField, nRequestID int, bIsLast bool) {

	// 报单方向
	DirectionTitle := utils.GetDirectionTitle(string(pInputOrder.Direction))

	// 报单开平
	OffsetFlagTitle := utils.GetOffsetFlagTitle(string(pInputOrder.CombOffsetFlag))

	OrderStr := "-------------------------------------------------------------------------------------------------\n" +
		"- 报单出错了\n" +
		"- 报单合约：" + pInputOrder.InstrumentID + "\t报单引用: " + pInputOrder.OrderRef + "\n" +
		"- 报单方向：" + DirectionTitle + "   \t报单价格: " + utils.Float64ToString(pInputOrder.LimitPrice) + "\n" +
		"- 报单开平：" + OffsetFlagTitle + " \t报单数量: " + utils.IntToString(pInputOrder.VolumeTotalOriginal) + "\n" +
		"- 错误代码：" + string(pRspInfo.ErrorID) + "    \t错误消息: " + pRspInfo.ErrorMsg + "\n" +
		"-------------------------------------------------------------------------------------------------"
	fmt.Println(OrderStr)
}

// 错误应答
func (p *FtdcTraderSpi) OnRspError(pRspInfo *goctp.RspInfoField, nRequestID int, bIsLast bool) {
	p.IsErrorRspInfo(pRspInfo)
}

// 报单操作错误回报
func (p *FtdcTraderSpi) OnErrRtnOrderAction(pOrderAction *goctp.OrderActionField, pRspInfo *goctp.RspInfoField) {
	p.IsErrorRspInfo(pRspInfo)
}

// 报单操作请求响应（撤单失败会触发）
func (p *FtdcTraderSpi) OnRspOrderAction(pInputOrderAction *goctp.InputOrderActionField, pRspInfo *goctp.RspInfoField, nRequestID int, bIsLast bool) {
	p.IsErrorRspInfo(pRspInfo)
}

// 交易系统错误通知
func (p *FtdcTraderSpi) IsErrorRspInfo(pRspInfo *goctp.RspInfoField) bool {

	// 容错处理 pRspInfo ，部分响应函数中，pRspInfo 为 0
	if pRspInfo == nil {
		// log.Printf("RspInfo: %+v\n", pRspInfo)
		return false

	} else {

		// 如果ErrorID != 0, 说明收到了错误的响应
		bResult := (pRspInfo.ErrorID != 0)
		if bResult {
			log.Printf("ErrorID=%v ErrorMsg=%v\n", pRspInfo.ErrorID, pRspInfo.ErrorMsg)
		}

		return bResult
	}
}

// OnHeartBeatWarning 心跳超时警告。当长时间未收到报文时，该方法被调用。
func (p *FtdcTraderSpi) OnHeartBeatWarning(nTimeLapse int) {
	fmt.Println("心跳超时警告 (OnHeartBeatWarning) nTimerLapse=", nTimeLapse)
}

// OrderOpen 开仓
func (p *FtdcTraderSpi) OrderOpen(input InputOrderStruct) int {

	iRequestID := p.GetTraderRequestId()

	mInstrumentInfo, mapRes := p.Strategy.GetInstrumentInfo(input.InstrumentID)
	if !mapRes {
		fmt.Println("开仓失败，合约不存在！")
		return 0
	}

	req := &goctp.InputOrderField{
		// 经纪公司代码
		BrokerID: p.Config.BrokerID,
		// 投资者代码
		InvestorID: p.Config.InvestorID,
		// 合约代码
		InstrumentID: input.InstrumentID,
		// 报单引用
		OrderRef: utils.IntToString(iRequestID),
		// 买卖方向:买(THOST_FTDC_D_Buy),卖(THOST_FTDC_D_Sell)
		Direction: input.Direction,
		// 交易所代码
		ExchangeID: mInstrumentInfo.ExchangeID,
		// 组合开平标志: 开仓
		CombOffsetFlag: string(thost.THOST_FTDC_OF_Open),
		// 组合投机套保标志: 投机
		CombHedgeFlag: string(thost.THOST_FTDC_HF_Speculation),
		// 报单价格条件: 限价
		OrderPriceType: thost.THOST_FTDC_OPT_LimitPrice,
		// 价格
		LimitPrice: input.Price,
		// 数量
		VolumeTotalOriginal: input.Volume,
		// 有效期类型: 当日有效
		TimeCondition: thost.THOST_FTDC_TC_GFD,
		// 成交量类型: 任何数量
		VolumeCondition: thost.THOST_FTDC_VC_AV,
		// 最小成交量
		MinVolume: 1,
		// 触发条件: 立即
		ContingentCondition: thost.THOST_FTDC_CC_Immediately,
		// 强平原因: 非强平
		ForceCloseReason: thost.THOST_FTDC_FCC_NotForceClose,
		// 自动挂起标志: 否
		IsAutoSuspend: 0,
		// 用户强评标志: 否
		UserForceClose: 0,
	}

	iResult := p.TraderApi.ReqOrderInsert(req, iRequestID)

	if iResult != 0 {
		utils.ReqFailMsg("提交报单失败！", iResult)
		return 0
	}

	return iRequestID
}

// 平仓
func (p *FtdcTraderSpi) OrderClose(input InputOrderStruct) int {

	iRequestID := p.GetTraderRequestId()

	mInstrumentInfo, mapRes := p.Strategy.GetInstrumentInfo(input.InstrumentID)
	if !mapRes {
		fmt.Println("平仓失败，合约不存在！")
		return 0
	}

	// 没有设置平仓类型
	if input.CombOffsetFlag == 0 {

		if mInstrumentInfo.ExchangeID == "SHFE" {
			// 上期所（默认平今仓）
			input.CombOffsetFlag = thost.THOST_FTDC_OF_CloseToday
		} else {
			// 非上期所，不用区分今昨仓，直接使用平仓即可
			input.CombOffsetFlag = thost.THOST_FTDC_OF_Close
		}
	}

	req := &goctp.InputOrderField{

		// 经纪公司代码
		BrokerID: p.Config.BrokerID,
		// 投资者代码
		InvestorID: p.Config.InvestorID,
		// 合约代码
		InstrumentID: input.InstrumentID,
		// 报单引用
		OrderRef: utils.IntToString(iRequestID),
		// 买卖方向:买(THOST_FTDC_D_Buy),卖(THOST_FTDC_D_Sell)
		Direction: input.Direction,
		// 交易所代码
		ExchangeID: mInstrumentInfo.ExchangeID,
		// 组合开平标志: 平仓 (针对上期所，区分昨仓、今仓)
		CombOffsetFlag: string(input.CombOffsetFlag),
		// 组合投机套保标志: 投机
		CombHedgeFlag: string(thost.THOST_FTDC_HF_Speculation),
		// 报单价格条件: 限价
		OrderPriceType: thost.THOST_FTDC_OPT_LimitPrice,
		// 价格
		LimitPrice: input.Price,
		// 数量
		VolumeTotalOriginal: input.Volume,
		// 有效期类型: 当日有效
		TimeCondition: thost.THOST_FTDC_TC_GFD,
		// 成交量类型: 任何数量
		VolumeCondition: thost.THOST_FTDC_VC_AV,
		// 最小成交量
		MinVolume: 1,
		// 触发条件: 立即
		ContingentCondition: thost.THOST_FTDC_CC_Immediately,
		// 强平原因: 非强平
		ForceCloseReason: thost.THOST_FTDC_FCC_NotForceClose,
		// 自动挂起标志: 否
		IsAutoSuspend: 0,
		// 用户强评标志: 否
		UserForceClose: 0,
	}

	iResult := p.TraderApi.ReqOrderInsert(req, iRequestID)

	if iResult != 0 {
		utils.ReqFailMsg("提交报单失败！", iResult)
		return 0
	}

	return iRequestID
}

// OrderCancel 撤消报单
func (p *FtdcTraderSpi) OrderCancel(InstrumentID string, OrderSysID string) int {

	iRequestID := p.GetTraderRequestId()

	mapKey := InstrumentID + "_" + OrderSysID

	// 检查报单列表数据是否存在
	mOrderVal, mOrderOk := p.MapOrderList.Get(mapKey)
	if !mOrderOk {
		fmt.Printf("撤消报单失败：合约 %v 报单编号 %v 不存在！\n", InstrumentID, OrderSysID)
		return 0
	}

	mOrder := mOrderVal.(OrderListStruct)

	req := &goctp.InputOrderActionField{
		// 经纪公司代码
		BrokerID: mOrder.BrokerID,
		// 投资者代码
		InvestorID: mOrder.InvestorID,
		// 合约代码
		InstrumentID: InstrumentID,
		// 报单引用
		OrderRef: mOrder.OrderRef,
		// 交易所代码
		ExchangeID: mOrder.ExchangeID,
		// 前置编号
		FrontID: mOrder.FrontID,
		// 会话编号
		SessionID: mOrder.SessionID,
		// 报单编号
		OrderSysID: mOrder.OrderSysID,
		// 操作标志
		ActionFlag: thost.THOST_FTDC_AF_Delete,
	}

	iResult := p.TraderApi.ReqOrderAction(req, iRequestID)

	if iResult != 0 {
		utils.ReqFailMsg("提交报单失败！", iResult)
		return 0
	}

	return iRequestID
}

// GetOrderListStruct 获得报单结构体数据
func GetOrderListStruct(pOrder *goctp.OrderField) *OrderListStruct {

	var mOrder OrderListStruct

	mOrder.BrokerID = pOrder.BrokerID
	mOrder.InvestorID = pOrder.InvestorID
	mOrder.InstrumentID = pOrder.InstrumentID
	mOrder.ExchangeID = pOrder.ExchangeID
	mOrder.FrontID = pOrder.FrontID
	mOrder.OrderRef = pOrder.OrderRef
	mOrder.SessionID = pOrder.SessionID
	mOrder.InsertTime = pOrder.InsertTime
	mOrder.OrderSysID = pOrder.OrderSysID
	mOrder.LimitPrice = pOrder.LimitPrice
	mOrder.Volume = pOrder.VolumeTotalOriginal
	mOrder.Direction = string(pOrder.Direction)
	mOrder.CombOffsetFlag = string(pOrder.CombOffsetFlag)
	mOrder.CombHedgeFlag = string(pOrder.CombHedgeFlag)
	mOrder.OrderStatus = string(pOrder.OrderStatus)
	mOrder.StatusMsg = pOrder.StatusMsg
	mOrder.DirectionTitle = utils.GetDirectionTitle(mOrder.Direction)
	mOrder.OrderStatusTitle = utils.GetOrderStatusTitle(mOrder.OrderStatus)
	mOrder.CombOffsetFlagTitle = utils.GetOffsetFlagTitle(mOrder.CombOffsetFlag)

	return &mOrder
}

// GetInvestorPositionStruct 获得持仓结构体数据
func (p *FtdcTraderSpi) GetInvestorPositionStruct(pInvestorPosition *goctp.InvestorPositionField) InvestorPositionStruct {

	var mInvestorPosition InvestorPositionStruct

	// 检查合约详情是否存在
	mInstrumentInfo, mapRes := p.Strategy.GetInstrumentInfo(pInvestorPosition.InstrumentID)
	if !mapRes {
		fmt.Printf("合约 %v 不存在！\n", pInvestorPosition.InstrumentID)
		return mInvestorPosition
	}

	// 合约乘数
	var VolumeMultiple int = mInstrumentInfo.VolumeMultiple

	// 开仓成本
	var OpenCost float64 = pInvestorPosition.OpenCost / float64(pInvestorPosition.Position*VolumeMultiple)

	mInvestorPosition.BrokerID = pInvestorPosition.BrokerID
	mInvestorPosition.InvestorID = pInvestorPosition.InvestorID
	mInvestorPosition.InstrumentID = pInvestorPosition.InstrumentID
	mInvestorPosition.InstrumentName = mInstrumentInfo.InstrumentName
	mInvestorPosition.PosiDirection = string(pInvestorPosition.PosiDirection)
	mInvestorPosition.PosiDirectionTitle = utils.GetPosiDirectionTitle(mInvestorPosition.PosiDirection)
	mInvestorPosition.HedgeFlag = string(pInvestorPosition.HedgeFlag)
	mInvestorPosition.HedgeFlagTitle = utils.GetHedgeFlagTitle(mInvestorPosition.HedgeFlag)
	mInvestorPosition.PositionDate = string(pInvestorPosition.PositionDate)
	mInvestorPosition.PositionDateTitle = utils.GetPositionDateTitle(mInvestorPosition.PositionDate)
	mInvestorPosition.Position = pInvestorPosition.Position
	mInvestorPosition.YdPosition = pInvestorPosition.YdPosition
	mInvestorPosition.TodayPosition = pInvestorPosition.TodayPosition
	mInvestorPosition.LongFrozen = pInvestorPosition.LongFrozen
	mInvestorPosition.ShortFrozen = pInvestorPosition.ShortFrozen

	// 冻结的持仓量（多空并成一个字段）
	if mInvestorPosition.PosiDirection == string(thost.THOST_FTDC_PD_Long) {
		// 多头冻结的持仓量
		mInvestorPosition.ShortVolume = pInvestorPosition.ShortFrozen
	} else {
		// 空头冻结的持仓量
		mInvestorPosition.ShortVolume = pInvestorPosition.LongFrozen
	}

	mInvestorPosition.OpenVolume = pInvestorPosition.OpenVolume
	mInvestorPosition.CloseVolume = pInvestorPosition.CloseVolume
	mInvestorPosition.PositionCost = utils.Decimal(pInvestorPosition.PositionCost, 2)
	mInvestorPosition.Commission = utils.Decimal(pInvestorPosition.Commission, 2)
	mInvestorPosition.CloseProfit = pInvestorPosition.CloseProfit
	mInvestorPosition.PositionProfit = utils.Decimal(pInvestorPosition.PositionProfit, 2)
	mInvestorPosition.PreSettlementPrice = pInvestorPosition.PreSettlementPrice
	mInvestorPosition.SettlementPrice = utils.Decimal(pInvestorPosition.SettlementPrice, 2)
	mInvestorPosition.SettlementID = pInvestorPosition.SettlementID
	mInvestorPosition.OpenCost = utils.Decimal(OpenCost, 2)
	mInvestorPosition.ExchangeID = pInvestorPosition.ExchangeID

	return mInvestorPosition
}
