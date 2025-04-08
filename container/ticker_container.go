package container

// 定义常量
const (
	MaxCoinCount  = 512
	CoinDataCount = 32
)

type SharedMemBookTicker struct {
	UpdateID  uint64
	Ets       uint64
	BuyPrice  float32
	BuyNum    float32
	SellPrice float32
	SellNum   float32
	Name      [16]byte
}

type SharedMemDataMap struct {
	CoinCount     int32
	CoinName      [MaxCoinCount][16]byte
	CoinUpdateIdx [MaxCoinCount]int32
	Data          [MaxCoinCount][CoinDataCount]SharedMemBookTicker
}
