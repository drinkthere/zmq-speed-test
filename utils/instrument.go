package utils

func ConvertBinanceFuturesInstIDToSpotInstID(binanceFuturesInstID string) string {
	// BTCUSDT => BTCUSDT
	// @TODO 兼容1000LUNC => LUNC
	return binanceFuturesInstID
}

func ConvertBybitLinearInstIDToSpotInstID(bybitLinearInstID string) string {
	// BTCUSDT => BTCUSDT
	// @TODO 兼容1000LUNC => LUNC
	return bybitLinearInstID
}
