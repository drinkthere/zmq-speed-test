package main

import (
	"zmq-speed-test/message"
)

func startTickerMessage() {
	// 定远远程数据，转发到本地
	if globalConfig.IsOkxLocalForward {
		message.StartOkxLocalTickerForward(&globalConfig, &globalContext)
	} else {
		message.StartLocalTickerForward(&globalConfig, &globalContext)
	}
}
