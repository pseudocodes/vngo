package strategy

import (
	"strings"

	"github.com/asmcos/requests"
	"github.com/tidwall/gjson"

	"github.com/pseudocodes/vngo/pkg/utils"
)

// SYMBOLS "KQ.m@CZCE.SM" 主链格式
var SYMBOLS = map[string][]string{
	"CZCE":  {"SM", "UR", "LR", "AP", "RS", "CY", "RI", "TA", "SA", "RM", "PF", "SR", "SF", "CJ", "MA", "ZC", "FG", "CF", "PK", "PM", "OI", "WH", "JR"},
	"DCE":   {"m", "j", "eg", "lh", "jm", "a", "l", "p", "b", "pb", "i", "jd", "c", "pp", "rr", "bb", "eb", "pg", "y", "v", "fb", "cs"},
	"CFFEX": {"IF", "IC", "TS", "IH", "TF", "T"},
	"INE":   {"sc", "lu", "nr", "bc"},
	"SHFE":  {"zn", "ni", "sn", "wr", "rb", "ru", "fu", "ag", "bu", "hc", "pb", "al", "ss", "cu", "au", "sp"},
}

type Ticker struct {
	InstrumentID string  `json:"instrument_id"` //symbol
	UpdateTime   int64   `json:"update_time"`   //时间
	Price        float64 `json:"price"`
	DayVolume    int     `json:"day_volume"`
	Volume       int     `json:"volume"`        //现手
	OpenInterest float64 `json:"open_interest"` //持仓量变化
	Interest     float64 `json:"interest"`      //持仓量变化
	Average      float64 `json:"average"`       //均价
	OpenPrice    float64 `json:"open_price"`    //当日开盘
	HighestPrice float64 `json:"highest_price"`
	LowestPrice  float64 `json:"lowest_price"`
	BidPrice1    float64 `json:"bid_price_1"`
	BidVolume1   int     `json:"bid_volume_1"`
	AskPrice1    float64 `json:"ask_price_1"`
	AskVolume1   int     `json:"ask_volume_1"`
}

type Kline struct {
	Datetime float64 `json:"datetime"`
	Open     float64 `json:"open"`
	High     float64 `json:"high"`
	Low      float64 `json:"low"`
	Close    float64 `json:"close"`
	OpenOi   float64 `json:"open_oi"`
	CloseOi  float64 `json:"close_oi"`
	Volume   float64 `json:"volume"`
}

// GetSymZL 获取主链代码
func GetSymZL(sym string) string {
	for ex, symbols := range SYMBOLS {
		fm := make(map[string]int)
		for i, v := range symbols {
			fm[v] = i
		}
		if _, ok := fm[sym]; ok {
			return "KQ.m@" + ex + "." + sym
		}
	}
	return ""
}

// GetMainContact 获取主力合约代码，不包含交易所代码： FG209
func GetMainContact(sym string) string {
	req := requests.Requests()
	req.Header.Set("Content-Type", "application/json")
	resp, _ := req.Get("http://localhost:3000/mct/" + sym) //
	var json []string
	err := resp.Json(&json)
	if err != nil {
		return ""
	}
	if len(json) > 0 {
		symbol := strings.Split(json[0], ".")
		return symbol[1]
	}
	return ""
}

//定义全局变量 gKlines, gKlines["FG209"] = { {20220202, 1,2,3,4,5,100},}

// InitKline 初始化某一个周期某一个品种的K线
func InitKline(Symbol string, period int64, maxlen int64) (string, []Kline) {
	ks := make([]Kline, 0)
	req := requests.Requests()
	req.Header.Set("Content-Type", "application/json")
	//http://localhost:3000/kline/CZCE.FG209?period=300&len=1500
	p := requests.Params{
		"len":    utils.Int64ToString(maxlen),
		"period": utils.Int64ToString(period),
	}
	resp, _ := req.Get("http://127.0.0.1:3000/kline/"+Symbol, p)
	if resp.R.StatusCode != 200 {
		return "", nil
	}

	oJKline := gjson.Get("{\"ROOT\": "+resp.Text()+"}", "ROOT")

	oJKline.ForEach(func(key, value gjson.Result) bool {
		_k := Kline{
			Datetime: value.Get("datetime").Float() / 1000000000,
			Open:     value.Get("open").Float(),
			High:     value.Get("high").Float(),
			Low:      value.Get("low").Float(),
			Close:    value.Get("close").Float(),
			OpenOi:   value.Get("open_oi").Float(),
			CloseOi:  value.Get("close_oi").Float(),
			Volume:   value.Get("volume").Float(),
		}
		ks = append(ks, _k)
		return true
	})
	if ks[0].Datetime < 946656000 { // 再次获取
		return InitKline(Symbol, period, maxlen)
	}
	return Symbol, ks
}

// MakeKLine 根据ticker序列，生成初始K线
func MakeKLine(tickers Ticker) Kline {
	nk := Kline{
		Datetime: float64(tickers.UpdateTime),
		Open:     tickers.Price,
		Close:    tickers.Price,
		High:     tickers.Price,
		Low:      tickers.Price,
		Volume:   float64(tickers.Volume),
		OpenOi:   tickers.OpenInterest - tickers.Interest,
		CloseOi:  tickers.Interest,
	}
	return nk
}

// UpdateKLine 根据ticker序列，生成初始K线
func UpdateKLine(thisKline Kline, tickers Ticker) Kline {
	if thisKline.Datetime <= 0 {
		return thisKline
	}
	thisKline.Volume += float64(tickers.Volume)
	thisKline.CloseOi += tickers.Interest
	thisKline.Close = tickers.Price
	if tickers.Price < thisKline.Low {
		thisKline.Low = tickers.Price
	}
	if tickers.Price > thisKline.High {
		thisKline.High = tickers.Price
	}
	return thisKline
}
