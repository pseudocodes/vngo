package main

import (
	"fmt"
	"os"
	"sync/atomic"

	"github.com/pseudocodes/vngo/pkg/strategy"
	"github.com/pseudocodes/vngo/pkg/utils"

	"github.com/pseudocodes/goctp"
	"github.com/tidwall/gjson"
)

var (
	gCtp       strategy.CtpClient                           // ctp 句柄及配置项
	gMdSpi     *strategy.FtdcMdSpi                          // 行情模块函数 句柄
	gTraderSpi *strategy.FtdcTraderSpi                      // 交易模块函数 句柄
	gCfg       strategy.Config                              // 加载配置文件项目
	StreamFile = utils.GetCurrentExePath() + "/StreamFile/" // StreamFile ctp 流文件，绝对路径
)

// LoadCfg 设置交易账号等相关的参数
func LoadCfg(RunMode string) {
	cfg := utils.LoadJson("cfg.json")
	if RunMode == "" {
		RunMode = cfg.Get("USERS.default").String()
	}
	ukey := "USERS." + RunMode
	if !cfg.Get(ukey).Exists() {
		_, err := utils.Println("该模式未设置交易账号信息")
		if err != nil {
		}
		os.Exit(1)
	}
	class := cfg.Get("STRATEGYS.default").String() //策略配置默认参数
	sKey := "STRATEGYS." + class

	gCfg = strategy.Config{
		MdFront:     cfg.Get(ukey + ".m_host").String(),
		TraderFront: cfg.Get(ukey + ".t_host").String(),
		BrokerID:    cfg.Get(ukey + ".bid").String(),
		InvestorID:  cfg.Get(ukey + ".uid").String(),
		Password:    cfg.Get(ukey + ".pwd").String(),
		AppID:       cfg.Get(ukey + ".app_id").String(),
		AuthCode:    cfg.Get(ukey + ".auth_code").String(),

		Class:   class, //策略struct名称
		MaxKlen: cfg.Get(sKey + ".max_klen").Int(),
		Period:  cfg.Get(sKey + ".period").Int(), //策略struct名称
	}

	// 加载策略获取行情订阅
	cfg.Get(sKey + ".symbol").ForEach(func(key, value gjson.Result) bool {
		gCfg.Symbol = append(gCfg.Symbol, value.String())
		return true
	})
	// 加载策略参数
	cfg.Get(sKey + ".params").ForEach(func(key, value gjson.Result) bool {
		gCfg.Params = append(gCfg.Params, value.Int())
		return true
	})
}

func init() {
	// 全局 行情、交易 函数句柄
	RunMode := "dev" // 运行模式【运行程序时带上参数可设置】,需要在cfg.json中配置参数
	if len(os.Args) > 1 {
		RunMode = os.Args[1]
	}
	LoadCfg(RunMode) // 设置交易相关参数，账号

	// 检查流文件目录是否存在
	fileExists := utils.IsDirExist(StreamFile)
	if !fileExists {
		err := os.Mkdir(StreamFile, os.ModePerm)
		if err != nil {
			fmt.Println("创建目录失败，请检查是否有操作权限")
			os.Exit(2)
		}
	}

	gCtp = strategy.CtpClient{
		MdApi:     goctp.CreateMdApiLite(goctp.MdFlowPath(StreamFile)),
		TraderApi: goctp.CreateTraderApiLite(goctp.TraderFlowPath(StreamFile)),
		Config:    gCfg,

		MdRequestId:        atomic.Int32{},
		TraderRequestId:    atomic.Int32{},
		IsTraderInit:       atomic.Bool{},
		IsTraderInitFinish: atomic.Bool{},
		IsMdLogin:          atomic.Bool{},
		IsTraderLogin:      atomic.Bool{},
	}
	fmt.Printf("%+v\n", gCfg)
	if gCtp.Register(&gCfg, gTraderSpi) { //注册策略
		gMdSpi = strategy.CreateFtdcMdSpi(&gCtp)
		gTraderSpi = strategy.CreateFtdcTraderSpi(&gCtp)
		// gMdSpi = strategy.FtdcMdSpi{CtpClient: &gCtp}
		// gTraderSpi = strategy.FtdcTraderSpi{CtpClient: &gCtp}
	} else {
		fmt.Printf("注册策略： %v 失败！\n", gCfg.Class)
		os.Exit(3)
	}
}
