package NetSurf

const (
	Version = "0.0.1"

	// TickCommonOperation common info update period.
	TickCommonOperation = 10

	// TickScanOperation Wi-Fi scan period.
	TickScanOperation = 5

	// CtxKeyWifiController context key for controller.Controller object.
	CtxKeyWifiController = "wifi_controller"

	// CtxKeyLoggerChannel context key for logger (ui.writer).
	CtxKeyLoggerChannel = "logger_channel"

	// CtxKeyHotKeys context key for ui.HotKeys object.
	CtxKeyHotKeys = "hot_keys"
)
