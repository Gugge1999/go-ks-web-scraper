package main

type apiStatus struct {
	scrapingIntervalInMinutes uint8 // active | inactive | pending
	memoryUsage               uint16
	uptime                    apiUptime
}

type apiUptime struct {
	years   uint8
	months  uint8
	days    uint8
	hours   uint8
	minutes uint8
	seconds uint8
}
