package main

import (
	"flag"
	"log"
	"time"

	"github.com/pseudocodes/goctp/thost"
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	flag.Parse()
}

func sample1() {
	log.Println("启动交易程序")

	gCtp.MdApi.RegisterSpi(gMdSpi.MdSpiLite)
	gCtp.MdApi.RegisterFront(gCfg.MdFront)
	gCtp.MdApi.Init()

	gCtp.TraderApi.RegisterSpi(gTraderSpi.TraderSpiLite)
	gCtp.TraderApi.RegisterFront(gCfg.TraderFront)

	gCtp.TraderApi.SubscribePublicTopic(thost.THOST_TERT_QUICK)
	gCtp.TraderApi.SubscribePrivateTopic(thost.THOST_TERT_QUICK)
	gCtp.TraderApi.Init()

	for {
		time.Sleep(10 * time.Second)
	}
}

func main() {
	sample1()
}
