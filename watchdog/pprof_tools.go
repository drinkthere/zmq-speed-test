package watchdog

import (
	"net/http"
	_ "net/http/pprof"
	"zmq-speed-test/config"
	"zmq-speed-test/utils/logger"
)

func StartPprofNet(cfg *config.Config) {
	if cfg.PprofListenAddress == "" {
		logger.Warn("[Watchdog] No Need Start Pprof Net")
		return
	}
	go func() {
		http.ListenAndServe(cfg.PprofListenAddress, nil)
	}()
	logger.Warn("[Watchdog] Start Pprof Net")
}
