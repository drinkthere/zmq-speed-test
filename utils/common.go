package utils

import (
	"math"
	"strconv"
	"sync/atomic"
	"time"
)

func MaxFloat64(list []float64) (max float64) {
	max = list[0]
	for _, v := range list {
		if v > max {
			max = v
		}
	}
	return
}

func MinFloat64(list []float64) (min float64) {
	min = list[0]
	for _, v := range list {
		if v < min {
			min = v
		}
	}
	return
}

func Round(value float64, decimals int) float64 {
	shift := math.Pow(10, float64(decimals))
	rounded := math.Round(value*shift) / shift
	return rounded
}

func InArray(target string, strArray []string) bool {
	for _, element := range strArray {
		if target == element {
			return true
		}
	}
	return false
}

func GetTimestampInMS() int64 {
	return time.Now().UnixNano() / 1e6
}

var gClientOrderID = GetTimestampInMS()

func GetClientOrderID() string {
	atomic.AddInt64(&gClientOrderID, 1)
	return strconv.FormatInt(atomic.LoadInt64(&gClientOrderID), 10)
}

// IsSettlement 判断是否是结算时间
// @param timeStamp: 当前时间戳，单位s
func IsSettlement(timeStamp int64) bool {
	// 不用判断时区，因为每8小时结算一次，所以东八区结算时间一样
	// 加60是为了方便后面比较
	tmp := (timeStamp + 60) % (60 * 60 * 24)
	if (tmp >= (23*60*60+59*60+50) || tmp <= 10) ||
		(tmp >= (7*60*60+59*60+50) && tmp <= 8*60*60+10) ||
		(tmp >= (15*60*60+59*60+50) && tmp <= 16*60*60+10) {
		return true
	}
	return false
}
