package config

type (
	RiskType       int
	Exchange       string
	InstrumentType string
	OrderSide      string
	OrderStatus    string
	TickerSource   string
)

const (
	// NoRisk 可以挂单，其他RiskType暂停挂单。1表示出错，2表示处于结算时间，3系统暂停等待价格更新
	NoRisk                = RiskType(iota)
	FatalErrorRisk        = NoRisk + 1
	SettlingRisk          = FatalErrorRisk + 1
	ExitRisk              = SettlingRisk + 1
	Price10sNotUpdateRisk = ExitRisk + 1

	BinanceExchange = Exchange("Binance")
	OkxExchange     = Exchange("Okx")
	BybitExchange   = Exchange("Bybit")

	UnknownInstrument = InstrumentType("UNKNOWN")
	SpotInstrument    = InstrumentType("SPOT")
	FuturesInstrument = InstrumentType("FUTURES")
	LinearInstrument  = InstrumentType("LINEAR")

	BinanceSpotTickerSource    = TickerSource("BinanceSpot")
	BinanceFuturesTickerSource = TickerSource("BinanceFutures")

	BybitSpotTickerSource   = TickerSource("BybitSpot")
	BybitLinearTickerSource = TickerSource("BybitLinear")
)
