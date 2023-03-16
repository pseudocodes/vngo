package strategy

/** SuperTrend 超级趋势策略
  SuperTrend 指标 + ema 144均线
  SuperTrend：https://www.elearnmarkets.com/blog/supertrend-indicator-strategy-trading/
	Up = (high + low） / 2 + multiplier  x  ATR
	Down = (high + low) / 2 – multiplier x ATR
*/

// type SuperTrend struct {
// 	Strategy
// }

// func init() { //所有的策略编写完毕后，需要在这地方进行注册登记
// 	RegisterStrategy["SuperTrend"] = &SuperTrend{}
// }

// func (s *SuperTrend) OnStart(mdspi *FtdcMdSpi) {
// 	s.Strategy.OnStart(mdspi)
// 	println("SuperTrend Strategy Start")
// }

// func (s *SuperTrend) OnBar(sym string, kdata Kline) {
// 	//println("SuperTrend OnBar")
// 	s.Strategy.OnBar(sym, kdata)
// 	st := GetSuperTrend(s.AM[sym], 14, 3)
// 	//fmt.Println(s.AM[sym][1480:])
// 	//fmt.Println(st[1480:])
// 	//IF PREV.ST > PREV.CLOSE AND CUR.ST < CUR.CLOSE ==> BUY SIGNAL
// 	//IF PREV.ST < PREV.CLOSE AND CUR.ST > CUR.CLOSE ==> SELL SIGNAL
// 	kLen := len(s.AM[sym])
// 	if (st[kLen-2] > s.AM[sym][kLen-2].Close) && (st[kLen-1] < s.AM[sym][kLen-1].Close) {
// 		s.OpenOrder(sym, 0, s.AM[sym][kLen-1].Close, 10)
// 		fmt.Println(utils.TimeToStr(int64(s.AM[sym][kLen-1].Datetime), ""), sym, "开多单")
// 	}
// 	if (st[kLen-2] < s.AM[sym][kLen-2].Close) && (st[kLen-1] > s.AM[sym][kLen-1].Close) {
// 		s.OpenOrder(sym, 1, s.AM[sym][kLen-1].Close, 10)
// 		fmt.Println(utils.TimeToStr(int64(s.AM[sym][kLen-1].Datetime), ""), sym, "开空单")
// 	}
// }

// func GetSuperTrend(klines []Kline, period int, multiplier float64) []float64 {
// 	klen := len(klines)
// 	closes := make([]float64, 0)
// 	highs := make([]float64, 0)
// 	lows := make([]float64, 0)
// 	hlAvg := make([]float64, 0)
// 	for _, k := range klines {
// 		closes = append(closes, k.Close)
// 		highs = append(highs, k.High)
// 		lows = append(lows, k.Low)
// 		hlAvg = append(hlAvg, (k.High+k.Low)/2)
// 	}

// 	arts := utils.Atr(highs, lows, closes, period)
// 	upperBand := make([]float64, 0)
// 	lowerBand := make([]float64, 0)
// 	for i, atr := range arts {
// 		upperBand = append(upperBand, hlAvg[i]+multiplier*atr)
// 		lowerBand = append(lowerBand, hlAvg[i]-multiplier*atr)
// 	}
// 	// FINAL UPPER BAND
// 	finalUpBands := make([]float64, 0)
// 	for i, up := range upperBand {
// 		if i == 0 {
// 			finalUpBands = append(finalUpBands, 0)
// 		} else {
// 			if (up < finalUpBands[i-1]) || (closes[i-1] > finalUpBands[i-1]) {
// 				finalUpBands = append(finalUpBands, up)
// 			} else {
// 				finalUpBands = append(finalUpBands, finalUpBands[i-1])
// 			}
// 		}
// 	}
// 	// FINAL LOWER BAND
// 	finalLowBands := make([]float64, 0)
// 	for i, low := range lowerBand {
// 		if i == 0 {
// 			finalLowBands = append(finalLowBands, 0)
// 		} else {
// 			if (low > finalLowBands[i-1]) || (closes[i-1] < finalLowBands[i-1]) {
// 				finalLowBands = append(finalLowBands, low)
// 			} else {
// 				finalLowBands = append(finalLowBands, finalLowBands[i-1])
// 			}
// 		}
// 	}
// 	//SUPER-TREND
// 	superTrent := make([]float64, klen)
// 	for i := period; i < klen; i++ {
// 		if (superTrent[i-1] == finalUpBands[i-1]) && (closes[i] <= finalUpBands[i]) {
// 			superTrent[i] = finalUpBands[i]
// 		} else if (superTrent[i-1] == finalUpBands[i-1]) && (closes[i] > finalUpBands[i]) {
// 			superTrent[i] = finalLowBands[i]
// 		} else if (superTrent[i-1] == finalLowBands[i-1]) && (closes[i] >= finalLowBands[i]) {
// 			superTrent[i] = finalLowBands[i]
// 		} else if (superTrent[i-1] == finalLowBands[i-1]) && (closes[i] < finalLowBands[i]) {
// 			superTrent[i] = finalUpBands[i]
// 		} else {
// 			superTrent[i] = finalLowBands[i]
// 		}
// 	}

// 	return superTrent
// }
