package types

type apiStatus struct {
	ScrapingIntervalInMinutes uint8 // active | inactive | pending
	MemoryUsage               uint16
	Uptime                    apiUptime
}

type apiUptime struct {
	years   uint8
	months  uint8
	days    uint8
	hours   uint8
	minutes uint8
	seconds uint8
}
