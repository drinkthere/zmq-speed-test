package main

import (
	"zmq-speed-test/message"
)

func startTickerMessage() {
	message.StartLocalTickerForward(&globalConfig, &globalContext)
}
