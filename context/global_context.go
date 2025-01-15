package context

import (
	"zmq-speed-test/config"
	"zmq-speed-test/container"
)

type GlobalContext struct {
	BestPathChangedCh chan struct{}
	BestPath          *container.BestPath
}

func (context *GlobalContext) Init(globalConfig *config.Config) {
	context.BestPathChangedCh = make(chan struct{})

	context.BestPath = &container.BestPath{}
	context.BestPath.Init(globalConfig.InitSourceIP, globalConfig.InitTargetIP)
}
