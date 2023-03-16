package strategy

// GetTraderRequestId 获得交易请求编号
func (p *FtdcTraderSpi) GetTraderRequestId() int {
	return int(p.TraderRequestId.Add(1))
}

// OnFrontDisconnected
/**
  当客户端与交易后台通信连接断开时，该方法被调用。
  当发生这个情况后，API会自动重新连接，客户端可不做处理。
*/
/*
func (p *FtdcTraderSpi) OnFrontDisconnected(nReason int) {

	p.IsTraderLogin = false
	p.IsTraderInit = false
	p.IsTraderInitFinish = false

	log.Println("交易服务器已断线，尝试重新连接中...")
}

// ReqMsg 发送请求日志（仅查询类的函数需要调用）
func (p *FtdcTraderSpi) ReqMsg(Msg string) {

	// 交易程序未初始化完成时，执行查询类的函数需要有1.5秒间隔
	if !p.IsTraderInitFinish {
		time.Sleep(time.Millisecond * 1500)
	}

	fmt.Println("")
	log.Println(Msg)
}

// OnFrontConnected 当客户端与交易后台建立起通信连接时（还未登录前），该方法被调用。
func (p *FtdcTraderSpi) OnFrontConnected() {

	TraderStr := "=================================================================================================\n" +
		"= 交易模块初始化成功，API 版本：" + lib.CThostFtdcTraderApiGetApiVersion() + "\n" +
		"================================================================================================="
	fmt.Println(TraderStr)

	p.IsTraderInit = true

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
	req := lib.NewCThostFtdcReqAuthenticateField()
	req.SetBrokerID(p.Config.BrokerID)
	req.SetUserID(p.Config.InvestorID)
	req.SetAppID(p.Config.AppID)
	req.SetAuthCode(p.Config.AuthCode)
	iResult := p.TraderApi.ReqAuthenticate(req, p.GetTraderRequestId())

	if iResult != 0 {
		utils.ReqFailMsg("发送客户端认证请求失败！", iResult)
	}
}

// OnRspAuthenticate 客户端认证响应
func (p *FtdcTraderSpi) OnRspAuthenticate(pRspAuthenticateField lib.CThostFtdcRspAuthenticateField, pRspInfo lib.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {

	if bIsLast && !p.IsErrorRspInfo(pRspInfo) {

		log.Println("客户端认证成功！")

		//p.MdSpi.ReqUserLogin()

		p.ReqUserLogin()
	}
}

// ReqUserLogin 用户登录请求
func (p *FtdcTraderSpi) ReqUserLogin() {

	time.Sleep(time.Second * 1)

	log.Println("交易系统账号登陆中...")

	req := lib.NewCThostFtdcReqUserLoginField()
	req.SetBrokerID(p.Config.BrokerID)
	req.SetUserID(p.Config.InvestorID)
	req.SetPassword(p.Config.Password)

	iResult := p.TraderApi.ReqUserLogin(req, p.GetTraderRequestId())

	if iResult != 0 {
		utils.ReqFailMsg("发送用户登录请求失败！", iResult)
	}
}

func (p *FtdcTraderSpi) OnRspUserLogin(pRspUserLogin lib.CThostFtdcRspUserLoginField, pRspInfo lib.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {

	if bIsLast && !p.IsErrorRspInfo(pRspInfo) {

		p.IsTraderLogin = true

		log.Printf("交易账号已登录，当前交易日：%v\n", p.TraderApi.GetTradingDay())

		p.ReqSettlementInfoConfirm()
	}
}

// ReqSettlementInfoConfirm 投资者结算单确认
func (p *FtdcTraderSpi) ReqSettlementInfoConfirm() int {

	p.ReqMsg("投资者结算单确认中...")

	req := lib.NewCThostFtdcSettlementInfoConfirmField()
	req.SetBrokerID(p.Config.BrokerID)
	req.SetInvestorID(p.Config.InvestorID)

	iResult := p.TraderApi.ReqSettlementInfoConfirm(req, p.GetTraderRequestId())

	if iResult != 0 {
		utils.ReqFailMsg("确认投资者结算单失败！", iResult)
	}

	return iResult
}

// OnRspSettlementInfoConfirm 发送投资者结算单确认响应
func (p *FtdcTraderSpi) OnRspSettlementInfoConfirm(pSettlementInfoConfirm lib.CThostFtdcSettlementInfoConfirmField, pRspInfo lib.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {

	if bIsLast && !p.IsErrorRspInfo(pRspInfo) {
		log.Println("投资者结算单确认成功")

		p.ReqQryInstrument()
	}
}

// ReqQryInstrument 请求查询合约
func (p *FtdcTraderSpi) ReqQryInstrument() int {

	p.ReqMsg("查询合约中...")

	req := lib.NewCThostFtdcQryInstrumentField()
	req.SetInstrumentID("")

	iResult := p.TraderApi.ReqQryInstrument(req, p.GetTraderRequestId())

	if iResult != 0 {
		utils.ReqFailMsg("查询合约失败！", iResult)
	}

	return iResult
}

// OnRspQryInstrument 请求查询合约响应
func (p *FtdcTraderSpi) OnRspQryInstrument(pInstrument lib.CThostFtdcInstrumentField, pRspInfo lib.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {

	if !p.IsErrorRspInfo(pRspInfo) {

		var mInstrumentInfo InstrumentInfoStruct

		var mapKey string = pInstrument.GetInstrumentID()

		mInstrumentInfo.InstrumentID = pInstrument.GetInstrumentID()
		mInstrumentInfo.ExchangeID = pInstrument.GetExchangeID()
		mInstrumentInfo.InstrumentName = utils.ConvertToString(pInstrument.GetInstrumentName(), "gbk", "utf-8")
		mInstrumentInfo.ExchangeInstID = pInstrument.GetExchangeInstID()
		mInstrumentInfo.ProductID = pInstrument.GetProductID()
		mInstrumentInfo.ProductClass = string(pInstrument.GetProductClass())
		mInstrumentInfo.DeliveryYear = pInstrument.GetDeliveryYear()
		mInstrumentInfo.DeliveryMonth = pInstrument.GetDeliveryMonth()
		mInstrumentInfo.MaxMarketOrderVolume = pInstrument.GetMaxMarketOrderVolume()
		mInstrumentInfo.MinMarketOrderVolume = pInstrument.GetMinMarketOrderVolume()
		mInstrumentInfo.MaxLimitOrderVolume = pInstrument.GetMaxLimitOrderVolume()
		mInstrumentInfo.MinLimitOrderVolume = pInstrument.GetMinLimitOrderVolume()
		mInstrumentInfo.VolumeMultiple = pInstrument.GetVolumeMultiple()
		mInstrumentInfo.PriceTick = pInstrument.GetPriceTick()
		mInstrumentInfo.CreateDate = pInstrument.GetCreateDate()
		mInstrumentInfo.OpenDate = pInstrument.GetOpenDate()
		mInstrumentInfo.ExpireDate = pInstrument.GetExpireDate()
		mInstrumentInfo.StartDelivDate = pInstrument.GetStartDelivDate()
		mInstrumentInfo.EndDelivDate = pInstrument.GetEndDelivDate()
		mInstrumentInfo.InstLifePhase = string(pInstrument.GetInstLifePhase())
		mInstrumentInfo.IsTrading = pInstrument.GetIsTrading()
		mInstrumentInfo.PositionType = string(pInstrument.GetPositionType())
		mInstrumentInfo.PositionDateType = string(pInstrument.GetPositionDateType())
		mInstrumentInfo.LongMarginRatio = pInstrument.GetLongMarginRatio()
		mInstrumentInfo.ShortMarginRatio = pInstrument.GetShortMarginRatio()
		mInstrumentInfo.MaxMarginSideAlgorithm = string(pInstrument.GetMaxMarginSideAlgorithm())
		mInstrumentInfo.UnderlyingInstrID = pInstrument.GetUnderlyingInstrID()
		mInstrumentInfo.StrikePrice = pInstrument.GetStrikePrice()
		mInstrumentInfo.OptionsType = string(pInstrument.GetOptionsType())
		mInstrumentInfo.UnderlyingMultiple = pInstrument.GetUnderlyingMultiple()
		mInstrumentInfo.CombinationType = string(pInstrument.GetCombinationType())

		p.Strategy.SetInstruments(mapKey, mInstrumentInfo)

		if bIsLast {

			log.Printf("合约记录初始化完毕！")

			if !p.IsTraderInitFinish {
				// 请求查询资金账户
				p.ReqQryTradingAccount()
			}
		}
	}
}

// ReqQryTradingAccount 请求查询资金账户
func (p *FtdcTraderSpi) ReqQryTradingAccount() int {

	p.ReqMsg("查询资金账户中...")

	req := lib.NewCThostFtdcQryTradingAccountField()
	req.SetBrokerID(p.Config.BrokerID)
	req.SetInvestorID(p.Config.InvestorID)

	iResult := p.TraderApi.ReqQryTradingAccount(req, p.GetTraderRequestId())

	if iResult != 0 {
		utils.ReqFailMsg("查询资金账户失败！", iResult)
	}

	return iResult
}

// OnRspQryTradingAccount 请求查询资金账户响应
func (p *FtdcTraderSpi) OnRspQryTradingAccount(pTradingAccount lib.CThostFtdcTradingAccountField, pRspInfo lib.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {

	if bIsLast && !p.IsErrorRspInfo(pRspInfo) {

		var mAccountInfo AccountInfoStruct

		mAccountInfo.MapKey = pTradingAccount.GetBrokerID() + "_" + pTradingAccount.GetAccountID()

		mAccountInfo.BrokerID = pTradingAccount.GetBrokerID()
		mAccountInfo.AccountID = pTradingAccount.GetAccountID()
		mAccountInfo.PreMortgage = utils.Decimal(pTradingAccount.GetPreMortgage(), 2)
		mAccountInfo.PreCredit = utils.Decimal(pTradingAccount.GetPreCredit(), 2)
		mAccountInfo.PreDeposit = utils.Decimal(pTradingAccount.GetPreDeposit(), 2)
		mAccountInfo.PreBalance = utils.Decimal(pTradingAccount.GetPreBalance(), 2)
		mAccountInfo.PreMargin = utils.Decimal(pTradingAccount.GetPreMargin(), 2)
		mAccountInfo.InterestBase = utils.Decimal(pTradingAccount.GetInterestBase(), 2)
		mAccountInfo.Interest = utils.Decimal(pTradingAccount.GetInterest(), 2)
		mAccountInfo.Deposit = utils.Decimal(pTradingAccount.GetDeposit(), 2)
		mAccountInfo.Withdraw = utils.Decimal(pTradingAccount.GetWithdraw(), 2)
		mAccountInfo.FrozenMargin = utils.Decimal(pTradingAccount.GetFrozenMargin(), 2)
		mAccountInfo.FrozenCash = utils.Decimal(pTradingAccount.GetFrozenCash(), 2)
		mAccountInfo.FrozenCommission = utils.Decimal(pTradingAccount.GetFrozenCommission(), 2)
		mAccountInfo.CurrMargin = utils.Decimal(pTradingAccount.GetCurrMargin(), 2)
		mAccountInfo.CashIn = utils.Decimal(pTradingAccount.GetCashIn(), 2)
		mAccountInfo.Commission = utils.Decimal(pTradingAccount.GetCommission(), 2)
		mAccountInfo.CloseProfit = utils.Decimal(pTradingAccount.GetCloseProfit(), 2)
		mAccountInfo.PositionProfit = utils.Decimal(pTradingAccount.GetPositionProfit(), 2)
		mAccountInfo.Balance = utils.Decimal(pTradingAccount.GetBalance(), 2)
		mAccountInfo.Available = utils.Decimal(pTradingAccount.GetAvailable(), 2)
		mAccountInfo.WithdrawQuota = utils.Decimal(pTradingAccount.GetWithdrawQuota(), 2)
		mAccountInfo.Reserve = utils.Decimal(pTradingAccount.GetReserve(), 2)
		mAccountInfo.TradingDay = pTradingAccount.GetTradingDay()
		mAccountInfo.SettlementID = pTradingAccount.GetSettlementID()
		mAccountInfo.Credit = utils.Decimal(pTradingAccount.GetCredit(), 2)
		mAccountInfo.Mortgage = utils.Decimal(pTradingAccount.GetMortgage(), 2)
		mAccountInfo.ExchangeMargin = utils.Decimal(pTradingAccount.GetExchangeMargin(), 2)
		mAccountInfo.DeliveryMargin = utils.Decimal(pTradingAccount.GetDeliveryMargin(), 2)
		mAccountInfo.ExchangeDeliveryMargin = utils.Decimal(pTradingAccount.GetExchangeDeliveryMargin(), 2)
		mAccountInfo.ReserveBalance = utils.Decimal(pTradingAccount.GetReserveBalance(), 2)
		mAccountInfo.CurrencyID = pTradingAccount.GetCurrencyID()
		mAccountInfo.PreFundMortgageIn = utils.Decimal(pTradingAccount.GetPreFundMortgageIn(), 2)
		mAccountInfo.PreFundMortgageOut = utils.Decimal(pTradingAccount.GetPreFundMortgageOut(), 2)
		mAccountInfo.FundMortgageIn = utils.Decimal(pTradingAccount.GetFundMortgageIn(), 2)
		mAccountInfo.FundMortgageOut = utils.Decimal(pTradingAccount.GetFundMortgageOut(), 2)
		mAccountInfo.FundMortgageAvailable = utils.Decimal(pTradingAccount.GetFundMortgageAvailable(), 2)
		mAccountInfo.MortgageableFund = utils.Decimal(pTradingAccount.GetMortgageableFund(), 2)
		mAccountInfo.SpecProductMargin = utils.Decimal(pTradingAccount.GetSpecProductMargin(), 2)
		mAccountInfo.SpecProductFrozenMargin = utils.Decimal(pTradingAccount.GetSpecProductFrozenMargin(), 2)
		mAccountInfo.SpecProductCommission = utils.Decimal(pTradingAccount.GetSpecProductCommission(), 2)
		mAccountInfo.SpecProductFrozenCommission = utils.Decimal(pTradingAccount.GetSpecProductFrozenCommission(), 2)
		mAccountInfo.SpecProductPositionProfit = utils.Decimal(pTradingAccount.GetSpecProductPositionProfit(), 2)
		mAccountInfo.SpecProductCloseProfit = utils.Decimal(pTradingAccount.GetSpecProductCloseProfit(), 2)
		mAccountInfo.SpecProductPositionProfitByAlg = utils.Decimal(pTradingAccount.GetSpecProductPositionProfitByAlg(), 2)
		mAccountInfo.SpecProductExchangeMargin = utils.Decimal(pTradingAccount.GetSpecProductExchangeMargin(), 2)
		mAccountInfo.BizType = string(pTradingAccount.GetBizType())
		mAccountInfo.FrozenSwap = utils.Decimal(pTradingAccount.GetFrozenSwap(), 2)
		mAccountInfo.RemainSwap = utils.Decimal(pTradingAccount.GetRemainSwap(), 2)

		AccountInfoStr := "-------------------------------------------------------------------------------------------------\n" +
			"- 公司代码：" + pTradingAccount.GetBrokerID() + "\n" +
			"- 资金账号：" + pTradingAccount.GetAccountID() + "\n" +
			"- 期初资金：" + utils.Float64ToString(mAccountInfo.PreBalance) + "\n" +
			"- 动态权益：" + utils.Float64ToString(mAccountInfo.Balance) + "\n" +
			"- 可用资金：" + utils.Float64ToString(mAccountInfo.Available) + "\n" +
			"- 持仓盈亏：" + utils.Float64ToString(mAccountInfo.PositionProfit) + "\n" +
			"- 平仓盈亏：" + utils.Float64ToString(mAccountInfo.CloseProfit) + "\n" +
			"- 手续费  ：" + utils.Float64ToString(mAccountInfo.Commission) + "\n" +
			"-------------------------------------------------------------------------------------------------"
		fmt.Println(AccountInfoStr)

		if !p.IsTraderInitFinish {
			// 请求查询投资者报单（委托单）
			p.ReqQryOrder()
		}
	}
}

// ReqQryOrder 请求查询投资者报单（委托单）
func (p *FtdcTraderSpi) ReqQryOrder() int {

	p.ReqMsg("查询投资者报单中...")

	req := lib.NewCThostFtdcQryOrderField()
	req.SetBrokerID(p.Config.BrokerID)
	req.SetInvestorID(p.Config.InvestorID)

	iResult := p.TraderApi.ReqQryOrder(req, p.GetTraderRequestId())

	if iResult != 0 {
		utils.ReqFailMsg("查询投资者报单失败！", iResult)
	}

	return iResult
}

// OnRspQryOrder 请求查询投资者报单响应
func (p *FtdcTraderSpi) OnRspQryOrder(pOrder lib.CThostFtdcOrderField, pRspInfo lib.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {

	if !p.IsErrorRspInfo(pRspInfo) {

		// 如果 没有数据 pOrder 等于0
		pOrderCode := fmt.Sprintf("%v", pOrder)

		// 只记录有报单编号的报单数据
		if pOrderCode != "0" && pOrder.GetOrderSysID() != "" {
			// 获得报单结构体数据
			mOrder := GetOrderListStruct(pOrder)

			// 报单列表数据 key 键
			mOrder.MapKey = pOrder.GetInstrumentID() + "_" + utils.TrimSpace(pOrder.GetOrderSysID())

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
				if val.OrderStatus == string(lib.THOST_FTDC_OST_NoTradeQueueing) || val.OrderStatus == string(lib.THOST_FTDC_OST_PartTradedQueueing) {
					MapOrderNoTradeSize += 1
					fmt.Printf("- 合约：%v   \t%v:%v   \t数量：%v   \t价格：%v   \t报单编号：%v (%v)\n", val.InstrumentID, val.DirectionTitle, val.CombOffsetFlagTitle, val.Volume, val.LimitPrice, utils.TrimSpace(val.OrderSysID), val.OrderStatusTitle)
				}
			}

			fmt.Printf("- 共有报单记录 %v 条，未成交 %v 条（不含错单）\n", p.MapOrderList.Size(), MapOrderNoTradeSize)
			fmt.Println("-------------------------------------------------------------------------------------------------")

			if !p.IsTraderInitFinish {
				// 请求查询投资者持仓（汇总）
				p.ReqQryInvestorPosition()
			}
		}
	}
}

// ReqQryInvestorPosition 请求查询投资者持仓（汇总）
func (p *FtdcTraderSpi) ReqQryInvestorPosition() int {

	p.ReqMsg("查询投资者持仓中...")

	req := lib.NewCThostFtdcQryInvestorPositionField()
	req.SetBrokerID(p.Config.BrokerID)
	req.SetInvestorID(p.Config.InvestorID)

	iResult := p.TraderApi.ReqQryInvestorPosition(req, p.GetTraderRequestId())

	if iResult != 0 {
		utils.ReqFailMsg("查询投资者持仓失败！", iResult)
	}

	fmt.Println("-------------------------------------------------------------------------------------------------")
	p.Strategy.EmptyPosition()
	return iResult
}

// OnRspQryInvestorPosition 请求查询投资者持仓（汇总）响应
func (p *FtdcTraderSpi) OnRspQryInvestorPosition(pInvestorPosition lib.CThostFtdcInvestorPositionField, pRspInfo lib.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {

	if !p.IsErrorRspInfo(pRspInfo) {

		// 没有数据 pInvestorPosition 会等于 0
		pInvestorPositionCode := fmt.Sprintf("%v", pInvestorPosition)

		if pInvestorPositionCode != "0" {

			// 获得持仓结构体数据
			mInvestorPosition := p.GetInvestorPositionStruct(pInvestorPosition)

			if mInvestorPosition.Position != 0 {
				p.Strategy.UpdatePosition(mInvestorPosition)
				fmt.Printf("- 合约：%v   \t%v:%v   \t总持仓：%v   \t持仓均价：%v   \t持仓盈亏：%v\n", mInvestorPosition.InstrumentID, mInvestorPosition.PositionDateTitle, mInvestorPosition.PosiDirectionTitle, mInvestorPosition.Position, mInvestorPosition.OpenCost, mInvestorPosition.PositionProfit)
			}
		}

		if bIsLast {

			fmt.Println("-------------------------------------------------------------------------------------------------")

			if !p.IsTraderInitFinish {
				// 交易程序初始化流程走完了
				p.IsTraderInitFinish = true
				// 订阅行情Subscribe := []string{"FG209"}
				//p.MdSpi.SubscribeMarketData(p.SubSymbols)
				//p.Strategy.OnStart(p)
			}
		}
	}
}

// OnRtnOrder 报单通知（委托单）
func (p *FtdcTraderSpi) OnRtnOrder(pOrder lib.CThostFtdcOrderField) {

	// 报单编号
	OrderSysID := pOrder.GetOrderSysID()

	// 报单状态
	OrderStatus := pOrder.GetOrderStatus()

	// 获得报单结构体数据
	mOrder := GetOrderListStruct(pOrder)

	// 报单列表数据 key 键
	mOrder.MapKey = pOrder.GetInstrumentID() + "_" + utils.TrimSpace(pOrder.GetOrderSysID())

	if OrderSysID == "" {

		// 报单就自动撤单，且没有编号的 都视为报错
		if OrderStatus == lib.THOST_FTDC_OST_Canceled {

			OrderErrorStr := "-------------------------------------------------------------------------------------------------\n" +
				"- 报单出错了\n" +
				"- 报单合约：" + mOrder.InstrumentID + "\t报单引用：" + mOrder.OrderRef + "\n" +
				"- 报单方向：" + mOrder.DirectionTitle + "   \t报单价格：" + utils.Float64ToString(mOrder.LimitPrice) + "\n" +
				"- 报单开平：" + mOrder.CombOffsetFlagTitle + " \t报单数量：" + utils.IntToString(mOrder.Volume) + "\n" +
				"- 错误代码：-1   \t错误消息：" + mOrder.StatusMsg + "\n" +
				"-------------------------------------------------------------------------------------------------"
			fmt.Println(OrderErrorStr)
		}

		return
	}

	// 未成交和撤单的报单（已成交的通知在 OnRtnTrade 函数中处理）
	if OrderStatus == lib.THOST_FTDC_OST_NoTradeQueueing || OrderStatus == lib.THOST_FTDC_OST_Canceled {

		OrderStr := "-------------------------------------------------------------------------------------------------\n" +
			"- 报单通知 " + mOrder.InsertTime + "\n" +
			"- 报单合约：" + mOrder.InstrumentID + " \t报单编号：" + utils.TrimSpace(mOrder.OrderSysID) + "\n" +
			"- 报单方向：" + mOrder.DirectionTitle + "   \t报单价格：" + utils.Float64ToString(mOrder.LimitPrice) + "\n" +
			"- 报单开平：" + mOrder.CombOffsetFlagTitle + " \t报单数量：" + utils.IntToString(mOrder.Volume) + "\n" +
			"- 报单状态：" + mOrder.OrderStatusTitle + " \t状态信息：" + mOrder.StatusMsg + "\n" +
			"-------------------------------------------------------------------------------------------------"
		fmt.Println(OrderStr)
	}
	p.Strategy.OnOrderChange(mOrder)
	// 将报单数据记录下来
	p.MapOrderList.Set(mOrder.MapKey, mOrder)
}

// OnRtnTrade 成交通知（委托单在交易所成交了）
func (p *FtdcTraderSpi) OnRtnTrade(pTrade lib.CThostFtdcTradeField) {

	// 报单方向
	DirectionTitle := utils.GetDirectionTitle(string(pTrade.GetDirection()))

	// 报单开平
	OffsetFlagTitle := utils.GetOffsetFlagTitle(string(pTrade.GetOffsetFlag()))
	//_, r := utils.Find(main.SubSymbols, pTrade.GetInstrumentID())
	//if !r {
	//	main.MdSpi.SubscribeMarketData([]string{pTrade.GetInstrumentID()}) // 订阅新行情
	//}
	p.Strategy.OnTradeDeal(pTrade)
	OrderStr := "-------------------------------------------------------------------------------------------------\n" +
		"- 成交通知 " + pTrade.GetTradeTime() + "\n" +
		"- 成交合约：" + pTrade.GetInstrumentID() + "\t成交编号：" + utils.TrimSpace(pTrade.GetTradeID()) + " \t报单编号：" + utils.TrimSpace(pTrade.GetOrderSysID()) + "\n" +
		"- 成交方向：" + DirectionTitle + "   \t成交价格：" + utils.Float64ToString(pTrade.GetPrice()) + "\n" +
		"- 成交开平：" + OffsetFlagTitle + " \t成交数量：" + utils.IntToString(pTrade.GetVolume()) + "\n" +
		"-------------------------------------------------------------------------------------------------"
	fmt.Println(OrderStr)
}

// OnRspOrderInsert 报单出错响应（综合交易平台交易核心返回的包含错误信息的报单响应）
func (p *FtdcTraderSpi) OnRspOrderInsert(pInputOrder lib.CThostFtdcInputOrderField, pRspInfo lib.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {

	// 报单方向
	DirectionTitle := utils.GetDirectionTitle(string(pInputOrder.GetDirection()))

	// 报单开平
	OffsetFlagTitle := utils.GetOffsetFlagTitle(string(pInputOrder.GetCombOffsetFlag()))

	OrderStr := "-------------------------------------------------------------------------------------------------\n" +
		"- 报单出错了\n" +
		"- 报单合约：" + pInputOrder.GetInstrumentID() + "\t报单引用：" + pInputOrder.GetOrderRef() + "\n" +
		"- 报单方向：" + DirectionTitle + "   \t报单价格：" + utils.Float64ToString(pInputOrder.GetLimitPrice()) + "\n" +
		"- 报单开平：" + OffsetFlagTitle + " \t报单数量：" + utils.IntToString(pInputOrder.GetVolumeTotalOriginal()) + "\n" +
		"- 错误代码：" + string(pRspInfo.GetErrorID()) + "    \t错误消息：" + utils.ConvertToString(pRspInfo.GetErrorMsg(), "gbk", "utf-8") + "\n" +
		"-------------------------------------------------------------------------------------------------"
	fmt.Println(OrderStr)
}

// 错误应答
func (p *FtdcTraderSpi) OnRspError(pRspInfo lib.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	p.IsErrorRspInfo(pRspInfo)
}

// 报单操作错误回报
func (p *FtdcTraderSpi) OnErrRtnOrderAction(pOrderAction lib.CThostFtdcOrderActionField, pRspInfo lib.CThostFtdcRspInfoField) {
	p.IsErrorRspInfo(pRspInfo)
}

// 报单操作请求响应（撤单失败会触发）
func (p *FtdcTraderSpi) OnRspOrderAction(pInputOrderAction lib.CThostFtdcInputOrderActionField, pRspInfo lib.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	p.IsErrorRspInfo(pRspInfo)
}

// 交易系统错误通知
func (p *FtdcTraderSpi) IsErrorRspInfo(pRspInfo lib.CThostFtdcRspInfoField) bool {

	rspInfo := fmt.Sprintf("%v", pRspInfo)

	// 容错处理 pRspInfo ，部分响应函数中，pRspInfo 为 0
	if rspInfo == "0" {
		return false

	} else {

		// 如果ErrorID != 0, 说明收到了错误的响应
		bResult := (pRspInfo.GetErrorID() != 0)
		if bResult {
			log.Printf("ErrorID=%v ErrorMsg=%v\n", pRspInfo.GetErrorID(), utils.ConvertToString(pRspInfo.GetErrorMsg(), "gbk", "utf-8"))
		}

		return bResult
	}
}

// OnHeartBeatWarning 心跳超时警告。当长时间未收到报文时，该方法被调用。
func (p *FtdcTraderSpi) OnHeartBeatWarning(nTimeLapse int) {
	fmt.Println("心跳超时警告（OnHeartBeatWarning） nTimerLapse=", nTimeLapse)
}

// OrderOpen 开仓
func (p *FtdcTraderSpi) OrderOpen(Input InputOrderStruct) int {

	iRequestID := p.GetTraderRequestId()

	mInstrumentInfo, mapRes := p.Strategy.GetInstrumentInfo(Input.InstrumentID)
	if !mapRes {
		fmt.Println("开仓失败，合约不存在！")
		return 0
	}

	req := lib.NewCThostFtdcInputOrderField()

	// 经纪公司代码
	req.SetBrokerID(p.Config.BrokerID)
	// 投资者代码
	req.SetInvestorID(p.Config.InvestorID)
	// 合约代码
	req.SetInstrumentID(Input.InstrumentID)
	// 报单引用
	req.SetOrderRef(utils.IntToString(iRequestID))
	// 买卖方向:买(THOST_FTDC_D_Buy),卖(THOST_FTDC_D_Sell)
	req.SetDirection(Input.Direction)
	// 交易所代码
	req.SetExchangeID(mInstrumentInfo.ExchangeID)
	// 组合开平标志: 开仓
	req.SetCombOffsetFlag(string(lib.THOST_FTDC_OF_Open))
	// 组合投机套保标志: 投机
	req.SetCombHedgeFlag(string(lib.THOST_FTDC_HF_Speculation))
	// 报单价格条件: 限价
	req.SetOrderPriceType(lib.THOST_FTDC_OPT_LimitPrice)
	// 价格
	req.SetLimitPrice(Input.Price)
	// 数量
	req.SetVolumeTotalOriginal(Input.Volume)
	// 有效期类型: 当日有效
	req.SetTimeCondition(lib.THOST_FTDC_TC_GFD)
	// 成交量类型: 任何数量
	req.SetVolumeCondition(lib.THOST_FTDC_VC_AV)
	// 最小成交量
	req.SetMinVolume(1)
	// 触发条件: 立即
	req.SetContingentCondition(lib.THOST_FTDC_CC_Immediately)
	// 强平原因: 非强平
	req.SetForceCloseReason(lib.THOST_FTDC_FCC_NotForceClose)
	// 自动挂起标志: 否
	req.SetIsAutoSuspend(0)
	// 用户强评标志: 否
	req.SetUserForceClose(0)

	iResult := p.TraderApi.ReqOrderInsert(req, iRequestID)

	if iResult != 0 {
		utils.ReqFailMsg("提交报单失败！", iResult)
		return 0
	}

	return iRequestID
}

// 平仓
func (p *FtdcTraderSpi) OrderClose(Input InputOrderStruct) int {

	iRequestID := p.GetTraderRequestId()

	mInstrumentInfo, mapRes := p.Strategy.GetInstrumentInfo(Input.InstrumentID)
	if !mapRes {
		fmt.Println("平仓失败，合约不存在！")
		return 0
	}

	// 没有设置平仓类型
	if Input.CombOffsetFlag == 0 {

		if mInstrumentInfo.ExchangeID == "SHFE" {
			// 上期所（默认平今仓）
			Input.CombOffsetFlag = lib.THOST_FTDC_OF_CloseToday
		} else {
			// 非上期所，不用区分今昨仓，直接使用平仓即可
			Input.CombOffsetFlag = lib.THOST_FTDC_OF_Close
		}
	}

	req := lib.NewCThostFtdcInputOrderField()

	// 经纪公司代码
	req.SetBrokerID(p.Config.BrokerID)
	// 投资者代码
	req.SetInvestorID(p.Config.InvestorID)
	// 合约代码
	req.SetInstrumentID(Input.InstrumentID)
	// 报单引用
	req.SetOrderRef(utils.IntToString(iRequestID))
	// 买卖方向:买(THOST_FTDC_D_Buy),卖(THOST_FTDC_D_Sell)
	req.SetDirection(Input.Direction)
	// 交易所代码
	req.SetExchangeID(mInstrumentInfo.ExchangeID)
	// 组合开平标志: 平仓 (针对上期所，区分昨仓、今仓)
	req.SetCombOffsetFlag(string(Input.CombOffsetFlag))
	// 组合投机套保标志: 投机
	req.SetCombHedgeFlag(string(lib.THOST_FTDC_HF_Speculation))
	// 报单价格条件: 限价
	req.SetOrderPriceType(lib.THOST_FTDC_OPT_LimitPrice)
	// 价格
	req.SetLimitPrice(Input.Price)
	// 数量
	req.SetVolumeTotalOriginal(Input.Volume)
	// 有效期类型: 当日有效
	req.SetTimeCondition(lib.THOST_FTDC_TC_GFD)
	// 成交量类型: 任何数量
	req.SetVolumeCondition(lib.THOST_FTDC_VC_AV)
	// 最小成交量
	req.SetMinVolume(1)
	// 触发条件: 立即
	req.SetContingentCondition(lib.THOST_FTDC_CC_Immediately)
	// 强平原因: 非强平
	req.SetForceCloseReason(lib.THOST_FTDC_FCC_NotForceClose)
	// 自动挂起标志: 否
	req.SetIsAutoSuspend(0)
	// 用户强评标志: 否
	req.SetUserForceClose(0)

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

	req := lib.NewCThostFtdcInputOrderActionField()

	// 经纪公司代码
	req.SetBrokerID(mOrder.BrokerID)
	// 投资者代码
	req.SetInvestorID(mOrder.InvestorID)
	// 合约代码
	req.SetInstrumentID(InstrumentID)
	// 报单引用
	req.SetOrderRef(mOrder.OrderRef)
	// 交易所代码
	req.SetExchangeID(mOrder.ExchangeID)
	// 前置编号
	req.SetFrontID(mOrder.FrontID)
	// 会话编号
	req.SetSessionID(mOrder.SessionID)
	// 报单编号
	req.SetOrderSysID(mOrder.OrderSysID)
	// 操作标志
	req.SetActionFlag(lib.THOST_FTDC_AF_Delete)

	iResult := p.TraderApi.ReqOrderAction(req, iRequestID)

	if iResult != 0 {
		utils.ReqFailMsg("提交报单失败！", iResult)
		return 0
	}

	return iRequestID
}

// GetOrderListStruct 获得报单结构体数据
func GetOrderListStruct(pOrder lib.CThostFtdcOrderField) OrderListStruct {

	var mOrder OrderListStruct

	mOrder.BrokerID = pOrder.GetBrokerID()
	mOrder.InvestorID = pOrder.GetInvestorID()
	mOrder.InstrumentID = pOrder.GetInstrumentID()
	mOrder.ExchangeID = pOrder.GetExchangeID()
	mOrder.FrontID = pOrder.GetFrontID()
	mOrder.OrderRef = pOrder.GetOrderRef()
	mOrder.SessionID = pOrder.GetSessionID()
	mOrder.InsertTime = pOrder.GetInsertTime()
	mOrder.OrderSysID = pOrder.GetOrderSysID()
	mOrder.LimitPrice = pOrder.GetLimitPrice()
	mOrder.Volume = pOrder.GetVolumeTotalOriginal()
	mOrder.Direction = string(pOrder.GetDirection())
	mOrder.CombOffsetFlag = string(pOrder.GetCombOffsetFlag())
	mOrder.CombHedgeFlag = string(pOrder.GetCombHedgeFlag())
	mOrder.OrderStatus = string(pOrder.GetOrderStatus())
	mOrder.StatusMsg = utils.ConvertToString(pOrder.GetStatusMsg(), "gbk", "utf-8")
	mOrder.DirectionTitle = utils.GetDirectionTitle(mOrder.Direction)
	mOrder.OrderStatusTitle = utils.GetOrderStatusTitle(mOrder.OrderStatus)
	mOrder.CombOffsetFlagTitle = utils.GetOffsetFlagTitle(mOrder.CombOffsetFlag)

	return mOrder
}

// GetInvestorPositionStruct 获得持仓结构体数据
func (p *FtdcTraderSpi) GetInvestorPositionStruct(pInvestorPosition lib.CThostFtdcInvestorPositionField) InvestorPositionStruct {

	var mInvestorPosition InvestorPositionStruct

	// 检查合约详情是否存在
	mInstrumentInfo, mapRes := p.Strategy.GetInstrumentInfo(pInvestorPosition.GetInstrumentID())
	if !mapRes {
		fmt.Printf("合约 %v 不存在！\n", pInvestorPosition.GetInstrumentID())
		return mInvestorPosition
	}

	// 合约乘数
	var VolumeMultiple int = mInstrumentInfo.VolumeMultiple

	// 开仓成本
	var OpenCost float64 = pInvestorPosition.GetOpenCost() / float64(pInvestorPosition.GetPosition()*VolumeMultiple)

	mInvestorPosition.BrokerID = pInvestorPosition.GetBrokerID()
	mInvestorPosition.InvestorID = pInvestorPosition.GetInvestorID()
	mInvestorPosition.InstrumentID = pInvestorPosition.GetInstrumentID()
	mInvestorPosition.InstrumentName = mInstrumentInfo.InstrumentName
	mInvestorPosition.PosiDirection = string(pInvestorPosition.GetPosiDirection())
	mInvestorPosition.PosiDirectionTitle = utils.GetPosiDirectionTitle(mInvestorPosition.PosiDirection)
	mInvestorPosition.HedgeFlag = string(pInvestorPosition.GetHedgeFlag())
	mInvestorPosition.HedgeFlagTitle = utils.GetHedgeFlagTitle(mInvestorPosition.HedgeFlag)
	mInvestorPosition.PositionDate = string(pInvestorPosition.GetPositionDate())
	mInvestorPosition.PositionDateTitle = utils.GetPositionDateTitle(mInvestorPosition.PositionDate)
	mInvestorPosition.Position = pInvestorPosition.GetPosition()
	mInvestorPosition.YdPosition = pInvestorPosition.GetYdPosition()
	mInvestorPosition.TodayPosition = pInvestorPosition.GetTodayPosition()
	mInvestorPosition.LongFrozen = pInvestorPosition.GetLongFrozen()
	mInvestorPosition.ShortFrozen = pInvestorPosition.GetShortFrozen()

	// 冻结的持仓量（多空并成一个字段）
	if mInvestorPosition.PosiDirection == string(lib.THOST_FTDC_PD_Long) {
		// 多头冻结的持仓量
		mInvestorPosition.ShortVolume = pInvestorPosition.GetShortFrozen()
	} else {
		// 空头冻结的持仓量
		mInvestorPosition.ShortVolume = pInvestorPosition.GetLongFrozen()
	}

	mInvestorPosition.OpenVolume = pInvestorPosition.GetOpenVolume()
	mInvestorPosition.CloseVolume = pInvestorPosition.GetCloseVolume()
	mInvestorPosition.PositionCost = utils.Decimal(pInvestorPosition.GetPositionCost(), 2)
	mInvestorPosition.Commission = utils.Decimal(pInvestorPosition.GetCommission(), 2)
	mInvestorPosition.CloseProfit = pInvestorPosition.GetCloseProfit()
	mInvestorPosition.PositionProfit = utils.Decimal(pInvestorPosition.GetPositionProfit(), 2)
	mInvestorPosition.PreSettlementPrice = pInvestorPosition.GetPreSettlementPrice()
	mInvestorPosition.SettlementPrice = utils.Decimal(pInvestorPosition.GetSettlementPrice(), 2)
	mInvestorPosition.SettlementID = pInvestorPosition.GetSettlementID()
	mInvestorPosition.OpenCost = utils.Decimal(OpenCost, 2)
	mInvestorPosition.ExchangeID = pInvestorPosition.GetExchangeID()

	return mInvestorPosition
}

*/
