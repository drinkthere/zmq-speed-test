package config

import (
	"encoding/json"
	"go.uber.org/zap/zapcore"
	"os"
)

type Config struct {
	// 日志配置
	LogLevel zapcore.Level
	LogPath  string

	IsOkxLocalForward bool     // 是否是OKX的本地转发，接受okx（aliyun）的数据并转发给币安（aws)本地服务
	LocalForwardIPC   string   // 将消息转发到本地ZMQ服务的IPC
	InitSourceIP      string   // 初始化时的本地最优IP
	InitTargetIP      string   // 初始化时的aws最优IP
	TargetPorts       []string // 不同数据源对应的Port

	UseBestPath        bool     // 是否使用BestPath
	BestPathChangedIPC string   // ZMQ地址，监听并获取最优路径变换的消息
	InstIDs            []string // 要套利的交易对

	PprofListenAddress             string   // profiling 监听的地址
	BinanceFuturesSharedMemoryPath []string // 共享内存文件路径
	CheckUpdateIntervalUs          int64    // 检查共享内存是否有更新的频率，单位Us
}

func LoadConfig(filename string) *Config {
	config := new(Config)
	reader, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer reader.Close()

	// 加载配置
	decoder := json.NewDecoder(reader)
	err = decoder.Decode(&config)
	if err != nil {
		panic(err)
	}

	return config
}
