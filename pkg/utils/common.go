package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/axgle/mahonia"
	"github.com/tidwall/gjson"
)

// GetDirectionTitle 获得报单多空方向
func GetDirectionTitle(Direction string) string {
	var title string

	switch Direction {
	case "0":
		title = "买"

	case "1":
		title = "卖"

	default:
		title = "未知"
	}

	return title
}

// GetPosiDirectionTitle 获得持仓多空方向
func GetPosiDirectionTitle(PosiDirection string) string {

	title := ""

	switch PosiDirection {
	case "1":
		title = "净"

	case "2":
		title = "买"

	case "3":
		title = "卖"

	default:
		title = "未知"
	}

	return title
}

// GetOrderStatusTitle 获得报单状态
func GetOrderStatusTitle(OrderStatus string) string {

	title := ""

	switch OrderStatus {
	case "0":
		title = "已成交"

	case "1":
		title = "部分成交还在队列中"

	case "2":
		title = "部分成交不在队列中"

	case "3":
		title = "未成交"

	case "4":
		title = "未成交不在队列中"

	case "5":
		title = "已撤单"

	case "a":
		title = "未知"

	case "b":
		title = "尚未触发"

	case "c":
		title = "已触发"

	default:
		title = "未知状态"
	}

	return title
}

// GetOffsetFlagTitle 获得开平标志
func GetOffsetFlagTitle(OrderStatus string) string {

	title := ""

	switch OrderStatus {
	case "0":
		title = "开仓"

	case "1":
		title = "平仓"

	case "2":
		title = "强平"

	case "3":
		title = "平今"

	case "4":
		title = "平昨"

	case "5":
		title = "强减"

	case "6":
		title = "本地强平"

	default:
		title = "未知"
	}

	return title
}

// GetHedgeFlagTitle 获得投机套保标志
func GetHedgeFlagTitle(HedgeFlag string) string {

	title := ""

	switch HedgeFlag {

	case "1":
		title = "投机"

	case "2":
		title = "套利"

	case "3":
		title = "套保"

	case "5":
		title = "做市商"

	case "6":
		title = "第一腿投机第二腿套保"

	case "7":
		title = "第一腿套保第二腿投机"

	default:
		title = "未知"
	}

	return title
}

// GetPositionDateTitle 获得持仓日期类型
func GetPositionDateTitle(PositionDate string) string {

	title := ""

	switch PositionDate {

	case "1":
		title = "今仓"

	case "2":
		title = "昨仓"

	default:
		title = "未知"
	}

	return title
}

// IsNullPointer 是否空指针
func IsNullPointer(p interface{}) bool {

	if p == nil {
		return true
	}

	pv := Sprintf("%v", p)
	if pv == "0" {
		return true
	}

	return false
}

// ReqMsg 请求日志
func ReqMsg(Msg string) {
	log.Println(Msg)
}

// ReqFailMsg 请求 api 出现错误
func ReqFailMsg(Msg string, iResult int) {
	fmt.Printf("%v [%d: %s]\n", Msg, iResult, iResultMsg(iResult))
}

// 请求失败的错误码对应消息
func iResultMsg(iResult int) string {

	msg := ""

	switch iResult {
	case 0:
		msg = "成功"
		break

	case -1:
		msg = "请检查账号是否已经登陆"
		break

	case -2:
		msg = "未处理请求超过许可数"
		break

	case -3:
		msg = "每秒发送请求数超过许可数"
		break

	default:
		msg = "未知错误"
		break
	}

	return msg
}

// CheckErr 检查错误，有就抛出
func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

// Decimal float64 保留几位小数点
func Decimal(f float64, n int) float64 {
	value, _ := strconv.ParseFloat(fmt.Sprintf("%."+strconv.Itoa(n)+"f", f), 64)
	return value
}

// IntToString int 转 string
func IntToString(i int) string {
	return strconv.Itoa(i)
}

// Int64ToString int 转 string
func Int64ToString(i int64) string {
	//n, err := strconv.ParseInt(i, 10, 64)
	return strconv.FormatInt(i, 10)
}
func StringToInt64(s string) int64 {
	if i64, err := strconv.ParseInt(s, 10, 64); err == nil {
		return i64
	} else {
		if i64, err := strconv.ParseInt(s, 16, 64); err == nil {
			return i64
		}
	}
	return 0
}

// Float64ToString float64 转 string
func Float64ToString(f float64) string {
	return strconv.FormatFloat(f, 'f', 2, 64)
}

// StringToFloat64 string 转 float64
func StringToFloat64(str string) float64 {
	f64, _ := strconv.ParseFloat(str, 64)
	return f64
}

// string 转 int
func StringToInt(str string) int {
	num, _ := strconv.Atoi(str)
	return num
}

// TrimSpace 去掉左右两边空格
func TrimSpace(str string) string {
	return strings.TrimSpace(str)
}

// fmt.Println
func Println(a ...interface{}) (n int, err error) {
	return fmt.Println(a...)
}

// fmt.Sprintf
func Sprintf(format string, a ...interface{}) string {
	return fmt.Sprintf(format, a...)
}

// fmt.Printf
func Printf(format string, a ...interface{}) (n int, err error) {
	return fmt.Printf(format, a...)
}

func StrInArray(str string, arr []string) bool {
	for _, v := range arr {
		if str == v {
			return true
		}
	}
	return false
}

// IsFileExist 文件是否存在
func IsFileExist(path string) bool {
	fi, err := os.Stat(path)
	if err == nil || os.IsExist(err) {
		mode := fi.Mode()
		return mode.IsRegular()
	}
	return false
}

// IsDirExist 目录是否存在
func IsDirExist(path string) bool {
	fi, err := os.Stat(path)
	if err == nil || os.IsExist(err) {
		mode := fi.Mode()
		return mode.IsDir()
	}
	return false
}

func GetCurrentExePath() string {
	dir := getCurrentAbPathByExecutable()
	tmpDir, _ := filepath.EvalSymlinks(os.TempDir())
	if strings.Contains(dir, tmpDir) { //控制台调试方式调用，获取当前utils目录作为路径，exe路径向上一级
		return path.Join(getCurrentAbPathByCaller(), "..")
	}
	return dir
}

// 获取当前执行文件绝对路径
func getCurrentAbPathByExecutable() string {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	res, _ := filepath.EvalSymlinks(filepath.Dir(exePath))
	return res
}

// 获取当前执行文件绝对路径（go run）
func getCurrentAbPathByCaller() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(filename)
	}
	return abPath
}

/**
 *  ConvertToString 编码转换
 *
 * 例： result = ConvertToString(text, "gbk", "utf-8")
 */
func ConvertToString(text string, srcCode string, tagCode string) string {

	srcCoder := mahonia.NewDecoder(srcCode)
	tagCoder := mahonia.NewDecoder(tagCode)
	srcResult := srcCoder.ConvertString(text)

	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)

	result := string(cdata)

	return result
}

func LoadFile(f string) string {
	// 打开json文件
	jsonFile, err := os.Open(f)

	// 最好要处理以下错误
	if err != nil {
		fmt.Println(err)
	}

	// 要记得关闭
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	//fmt.Println(string(byteValue))
	return string(byteValue)
}

func LoadJson(f string) gjson.Result {
	var _f = path.Join(GetCurrentExePath(), f)
	_r := gjson.Result{}
	if IsFileExist(_f) {
		var jsonStr = LoadFile(_f)
		j := `{"root": ` + jsonStr + `}`
		_v := gjson.Get(j, "root")
		// _v, _, _, _e := jsonparser.Get(jsonStr)
		// fmt.Println(_dt, _o)
		if _v.Exists() {
			return _v
		}
	}
	return _r
}

// Find Find获取一个切片并在其中查找元素。如果找到它，它将返回它的密钥，否则它将返回-1和一个错误的bool。
func Find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func IsNil(v interface{}) bool {
	valueOf := reflect.ValueOf(v)

	k := valueOf.Kind()

	switch k {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		return valueOf.IsNil()
	default:
		return v == nil
	}
}

// StrToTime 转换 17:59:24 到具体的秒
func StrToTime(t string) int64 {
	TimeLayout := "20060102 15:04:05"

	now := t
	if len(t) < 10 {
		now = time.Now().Format("20060102 ") + t
	}
	l, _ := time.LoadLocation("Asia/Shanghai")
	lt, _ := time.ParseInLocation(TimeLayout, now, l)
	//ts, _ := time.Parse(TimeLayout, now)
	//ts = ts.Local()
	return lt.Unix()
}

func TimeToStr(t int64, format string) string {
	tm := time.Unix(t, 0)
	tm = tm.Local()
	if format == "" {
		format = "15:04:05"
	}
	return tm.Format(format)
}
