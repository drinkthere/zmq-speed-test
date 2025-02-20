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

func StartLocalTickerForward(cfg *config.Config, globalContext *context.GlobalContext) {
	forwardSvc := &LocalTickerForwardService{}

	for _, targetPort := range cfg.TargetPorts {
		forwardSvc.StartSubService(cfg.UseBestPath, targetPort, globalContext)
	}

	if cfg.UseBestPath {
		forwardSvc.StartSubBestPathChange(cfg, globalContext)
	}
}

type LocalTickerForwardService struct{}

func (s *LocalTickerForwardService) StartSubBestPathChange(cfg *config.Config, globalContext *context.GlobalContext) {
	go func() {
		defer func() {
			logger.Warn("[SubBestPathChange] Sub Service Listening Exited.")
		}()

		logger.Warn("[SubBestPathChange] Start Best Path Service.")
		var ctx *zmq.Context
		var sub *zmq.Socket
		isBestPathStopped := true
		for {
			if isBestPathStopped {
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
				isBestPathStopped = false
			}

			logger.Warn("[SubBestPathChange] Start Receiving Data.")
			msg, err := sub.Recv(0)
			if err != nil {
				logger.Error("[SubBestPathChange] Receive Best Path Changed ZMQ Msg Error: %s", err.Error())
				sub.Close()
				ctx.Term()
				isBestPathStopped = true
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
				isBestPathStopped = true
				time.Sleep(time.Second * 1)
				continue
			}
			logger.Warn("[SubBestPathChange] Best Path Has Been Changed from %s->%s to %s->%s",
				globalContext.BestPath.SourceIP, globalContext.BestPath.TargetIP,
				bestPath.SourceIP, bestPath.TargetIP)
			globalContext.BestPath.SourceIP = bestPath.SourceIP
			globalContext.BestPath.TargetIP = bestPath.TargetIP
			globalContext.BestPathBroadcast.Broadcast()
		}
	}()
}

func (s *LocalTickerForwardService) StartSubService(useBestPath bool, targetPort string, globalContext *context.GlobalContext) {
	go func() {
		defer func() {
			logger.Warn("[LocalTickerForward] Sub Service Listening Exited.")
		}()

		var bestPathCh chan struct{}
		if useBestPath {
			bestPathCh := globalContext.BestPathBroadcast.Subscribe()
			defer func() {
				if bestPathCh != nil {
					globalContext.BestPathBroadcast.Unsubscribe(bestPathCh)
					logger.Info("[LocalTickerForward] Unsubscribed from best path changes")
				}
			}()
		}

		logger.Warn("[LocalTickerForward] Start Local Sub Service.")
		var ctx *zmq.Context
		var sub *zmq.Socket
		isSubStopped := true
		for {
			select {
			case <-bestPathCh:
				logger.Warn("[LocalTickerForward] Best path changed, closing current connection and restarting.")
				sub.Close()
				ctx.Term()
				isSubStopped = true
				time.Sleep(time.Second * 1)
				continue
			default:
				if isSubStopped {
					ctx, _ = zmq.NewContext()
					sub, _ = ctx.NewSocket(zmq.SUB)
					target := globalContext.BestPath.TargetIP + ":" + targetPort
					logger.Info("%s", target)
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
					isSubStopped = false
				}
			}
			msg, err := sub.Recv(0)
			if err != nil {
				logger.Error("[LocalTickerForward] Receive Remote ZMQ Msg Error: %s", err.Error())
				sub.Close()
				ctx.Term()
				isSubStopped = true
				time.Sleep(time.Second * 1)
				continue
			}
			var ticker pb.TickerInfo
			// parse message
			err = proto.Unmarshal([]byte(msg), &ticker)
			if err != nil {
				logger.Error("[BSTickerForward] Parse ZMQ Msg Error: %s", err.Error())
				sub.Close()
				ctx.Term()
				isSubStopped = true
				time.Sleep(time.Second * 1)
				continue
			}
			logger.Info("=stat= %s|%f|%f|%f|%f|%d|%d|%d", ticker.InstID, ticker.BestBid, ticker.BidSz,
				ticker.BestAsk, ticker.AskSz, ticker.UpdateID, ticker.EventTs, time.Now().UnixNano())
		}
	}()
}
