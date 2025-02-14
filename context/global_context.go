package context

import (
	"zmq-speed-test/config"
	"zmq-speed-test/container"
)

type GlobalContext struct {
	BestPathBroadcast *container.BroadcastChannel
	BestPath          *container.BestPath
}

func (context *GlobalContext) Init(globalConfig *config.Config) {
	context.BestPath = &container.BestPath{}
	context.BestPath.Init(globalConfig.InitSourceIP, globalConfig.InitTargetIP)
}
