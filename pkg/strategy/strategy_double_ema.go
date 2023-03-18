package strategy

import (
	"fmt"

	"github.com/pseudocodes/vngo/pkg/utils"
)

type SuperDoubleEMA struct {
	Strategy
}

func init() { //所有的策略编写完毕后，需要在这地方进行注册登记
	RegisterStrategy["SuperDoubleEMA"] = &SuperDoubleEMA{}
}

func (s *SuperDoubleEMA) OnStart(mdspi *FtdcMdSpi) {
	s.Strategy.OnStart(mdspi)

}

func (s *SuperDoubleEMA) OnBar(sym string, kdata Kline) {
	s.Strategy.OnBar(sym, kdata)
	closes := make([]float64, 0)
	for _, k := range s.AM[sym] {
		closes = append(closes, k.Close)
	}
	p1 := s.Params[0] //短周期
	p2 := s.Params[1] //长周期
	ema12 := utils.Ema(closes, int(p1))
	ema144 := utils.Ema(closes, int(p2))
	kLen := len(s.AM[sym])
	len12 := len(ema12)
	len114 := len(ema144)
	// 上穿
	if (ema12[len12-1] > ema144[len114-1]) && (ema12[len12-2] <= ema144[len114-2]) {
		// s.OpenOrder(sym, 0, s.AM[sym][kLen-1].Close, 10)
		fmt.Println(utils.TimeToStr(int64(s.AM[sym][kLen-1].Datetime), ""), sym, "开多单")
		// todo 平空单
	}
	//下穿
	if (ema12[len12-1] < ema144[len114-1]) && (ema12[len12-2] >= ema144[len114-2]) {
		// s.OpenOrder(sym, 1, s.AM[sym][kLen-1].Close, 10)
		fmt.Println(utils.TimeToStr(int64(s.AM[sym][kLen-1].Datetime), ""), sym, "开空单")
		// todo 平空单
	}

}
