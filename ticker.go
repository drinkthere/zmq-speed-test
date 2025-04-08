package main

import (
	"time"
	"zmq-speed-test/message"
)

func startTickerMessage() {
	if len(globalConfig.TargetPorts) > 0 {
		message.StartLocalTickerForward(&globalConfig, &globalContext)
	} else if len(globalConfig.BinanceFuturesSharedMemoryPath) > 0 {
		for _, sharedMemoryPath := range globalConfig.BinanceFuturesSharedMemoryPath {
			message.StartLocalSharedMemory(&globalConfig, &globalContext, sharedMemoryPath)
			time.Sleep(100 * time.Millisecond)
		}

	}

}
