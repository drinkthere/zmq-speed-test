package message

import (
	zmq "github.com/pebbe/zmq4"
	"google.golang.org/protobuf/proto"
	"zmq-speed-test/context"
	"zmq-speed-test/protocol/pb"

	"time"
	//"strconv"
	"zmq-speed-test/config"
	"zmq-speed-test/utils/logger"
)

func StartOkxLocalTickerForward(cfg *config.Config, globalContext *context.GlobalContext) {
	forwardSvc := LocalTickerForwardService{}
	forwardSvc.Init()
	forwardSvc.StartSubService(cfg, globalContext)
	if cfg.UseBestPath {
		forwardSvc.StartSubBestPathChange(cfg, globalContext)
	}

}

type OkxLocalTickerForwardService struct {
	msgChan           chan *string
	isSubStopped      bool
	isPubStopped      bool
	isBestPathStopped bool
}

func (s *OkxLocalTickerForwardService) Init() {
	s.msgChan = make(chan *string)
	s.isSubStopped = true
	s.isPubStopped = true
	s.isBestPathStopped = true
}

func (s *OkxLocalTickerForwardService) StartSubBestPathChange(cfg *config.Config, globalContext *context.GlobalContext) {
	go func() {
		defer func() {
			logger.Warn("[SubBestPathChange] Sub Service Listening Exited.")
		}()

		logger.Info("[SubBestPathChange] Start Best Path Service.")
		var ctx *zmq.Context
		var sub *zmq.Socket
		for {
			if s.isBestPathStopped {
				ctx, _ = zmq.NewContext()
				sub, _ = ctx.NewSocket(zmq.SUB)
				err := sub.Connect(cfg.BestPathChangedIPC)

				if err != nil {
					logger.Error("[SubBestPathChange] Connect to Best Path Changed ZMQ %s Error: %s", cfg.BestPathChangedIPC, err.Error())
					ctx.Term()
					time.Sleep(time.Second * 1)
					continue
				}
				err = sub.SetSubscribe("")
				if err != nil {
					logger.Error("[SubBestPathChange] Subscribe Best Path Changed ZMQ Subscription Error: %s", err.Error())
					sub.Close()
					ctx.Term()
					time.Sleep(time.Second * 1)
					continue
				}
				s.isBestPathStopped = false
			}

			logger.Info("[SubBestPathChange] Start Receiving Data.")
			msg, err := sub.Recv(0)
			if err != nil {
				logger.Error("[SubBestPathChange] Receive Best Path Changed ZMQ Msg Error: %s", err.Error())
				sub.Close()
				ctx.Term()
				s.isBestPathStopped = true
				time.Sleep(time.Second * 1)
				continue
			}

			var bestPath pb.BestPath
			// parse message
			err = proto.Unmarshal([]byte(msg), &bestPath)
			if err != nil {
				logger.Error("[SubBestPathChange] Parse Best Path Changed ZMQ Msg %s Error: %s", msg, err.Error())
				sub.Close()
				ctx.Term()
				s.isBestPathStopped = true
				time.Sleep(time.Second * 1)
				continue
			}
			// best path的sourceIP是aliyun的ip，targetIP是aws的ip，这里okx的localForward部署在AWS上，所以，sourceIP和targetIP要反过来
			logger.Info("[SubBestPathChange] Best Path Has Been Changed from %s->%s to %s->%s",
				globalContext.BestPath.SourceIP, globalContext.BestPath.TargetIP,
				bestPath.TargetIP, bestPath.SourceIP)
			globalContext.BestPath.SourceIP = bestPath.TargetIP
			globalContext.BestPath.TargetIP = bestPath.SourceIP
			globalContext.BestPathChangedCh <- struct{}{}
		}
	}()
}

func (s *OkxLocalTickerForwardService) StartSubService(cfg *config.Config, globalContext *context.GlobalContext) {
	go func() {
		defer func() {
			logger.Warn("[LocalTickerForward] Sub Service Listening Exited.")
		}()

		logger.Info("[LocalTickerForward] Start Local Sub Service.")
		var ctx *zmq.Context
		var sub *zmq.Socket
		for {
			select {
			case <-globalContext.BestPathChangedCh:
				logger.Info("[LocalTickerForward] Best path changed, closing current connection and restarting.")
				sub.Close()
				ctx.Term()
				s.isSubStopped = true
				time.Sleep(time.Second * 1)
				continue
			default:
				if s.isSubStopped {
					ctx, _ = zmq.NewContext()
					sub, _ = ctx.NewSocket(zmq.SUB)
					target := globalContext.BestPath.TargetIP + ":" + cfg.TargetPort
					err := sub.Connect("tcp://" + globalContext.BestPath.SourceIP + ":0;" + target)

					if err != nil {
						logger.Error("[LocalTickerForward] Connect to Remote ZMQ %s Error: %s", target, err.Error())
						ctx.Term()
						time.Sleep(time.Second * 1)
						continue
					}
					err = sub.SetSubscribe("")
					if err != nil {
						logger.Error("[LocalTickerForward] Subscribe Remote ZMQ Subscription Error: %s", err.Error())
						sub.Close()
						ctx.Term()
						time.Sleep(time.Second * 1)
						continue
					}
					s.isSubStopped = false
				}
			}
			msg, err := sub.Recv(0)
			if err != nil {
				logger.Error("[LocalTickerForward] Receive Remote ZMQ Msg Error: %s", err.Error())
				sub.Close()
				ctx.Term()
				s.isSubStopped = true
				time.Sleep(time.Second * 1)
				continue
			}
			s.msgChan <- &msg
		}
	}()
}
func (s *OkxLocalTickerForwardService) StartPubService(cfg *config.Config) {
	go func() {
		defer func() {
			logger.Warn("[LocalTickerForward] Pub Service Listening Exited.")
		}()

		logger.Info("[LocalTickerForward] Start Local Pub Service.")
		var ctx *zmq.Context
		var pub *zmq.Socket
		for {
			if s.isPubStopped {
				ctx, _ = zmq.NewContext()
				pub, _ = ctx.NewSocket(zmq.PUB)
				err := pub.Bind(cfg.LocalForwardIPC)
				if err != nil {
					logger.Error("[LocalTickerForward] Bind to Local ZMQ %s Error: %s", cfg.LocalForwardIPC, err.Error())
					ctx.Term()
					time.Sleep(time.Second * 1)
					continue
				}
				s.isPubStopped = false
			}
			msg := <-s.msgChan
			_, err := pub.Send(*msg, 0)
			if err != nil {
				logger.Warn("[LocalTickerForward] Error sending MarketData: %v", err)
				s.isPubStopped = true
				pub.Close()
				ctx.Term()
				time.Sleep(time.Second * 1)
				continue
			}
		}
	}()
}
