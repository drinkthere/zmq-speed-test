package message

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
	"syscall"
	"time"
	"unsafe"
	"zmq-speed-test/config"
	"zmq-speed-test/container"
	"zmq-speed-test/context"
	"zmq-speed-test/utils/logger"
)

func StartLocalSharedMemory(
	globalConfig *config.Config,
	globalContext *context.GlobalContext,
	sharedMemoryPath string) {

	r := rand.New(rand.NewSource(2))
	go gatherBnSmTicker(globalConfig, globalContext, r, config.FuturesInstrument, sharedMemoryPath)
	logger.Info("[GatherBtcEthFuturesSm] Start Gather Binance Futures BookTicker %s", sharedMemoryPath)
}

func gatherBnSmTicker(
	globalConfig *config.Config,
	globalContext *context.GlobalContext,
	r *rand.Rand,
	instType config.InstrumentType,
	sharedMemoryPath string) {

	md, data, err := mapSharedMemory(sharedMemoryPath)
	if err != nil {
		fmt.Printf("Error mapping shared memory: %v\n", err)
		return
	}
	defer func() {
		logger.Warn("[GatherSm] %s Gather Exited.", instType)
		syscall.Munmap(data) // 在程序退出时释放内存映射
	}()

	if md == nil {
		logger.Fatal("[GatherSm] Failed to access shared memory data.")

	}

	previousCoinUpdateIdx := make([]int32, container.MaxCoinCount)
	mutex := &sync.RWMutex{}
	ticker := time.NewTicker(time.Duration(globalConfig.CheckUpdateIntervalUs) * time.Microsecond)
	for {
		_ = <-ticker.C

		mutex.RLock()
		for i := 0; i < int(md.CoinCount); i++ {
			idx := md.CoinUpdateIdx[i]
			if idx > previousCoinUpdateIdx[i] {
				bt := md.Data[i][idx%container.CoinDataCount]
				instID := bytesToString(md.CoinName[i][:])
				previousCoinUpdateIdx[i] = idx
				logger.Info("=stat= %s|%f|%f|%f|%f|%d|%d|%d", instID, bt.BuyPrice, bt.BuyNum,
					bt.SellPrice, bt.SellNum, bt.UpdateID, bt.Ets, time.Now().UnixNano())

			}
		}
		mutex.RUnlock()
	}
}

func mapSharedMemory(filePath string) (*container.SharedMemDataMap, []byte, error) {
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open shared memory file: %v", err)
	}

	// 获取文件大小
	fileInfo, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, nil, fmt.Errorf("failed to get file info: %v", err)
	}
	fileSize := fileInfo.Size()

	// 映射共享内存
	data, err := syscall.Mmap(int(file.Fd()), 0, int(fileSize), syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		file.Close()
		return nil, nil, fmt.Errorf("failed to mmap file: %v", err)
	}

	// 关闭文件描述符（映射后可以关闭文件）
	file.Close()

	// 将共享内存解析为 MAP_DATA
	md := (*container.SharedMemDataMap)(unsafe.Pointer(&data[0]))
	return md, data, nil
}

func bytesToString(b []byte) string {
	n := 0
	for ; n < len(b); n++ {
		if b[n] == 0 {
			break
		}
	}
	return string(b[:n])
}
