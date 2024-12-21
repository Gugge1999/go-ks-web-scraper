package types

type ApiStatus struct {
	Active                    bool   `json:"active"` // active | inactive | pending
	ScrapingIntervalInMinutes uint   `json:"scrapingIntervalInMinutes"`
	MemoryUsage               uint16 `json:"memoryUsage"`
	Uptime                    Uptime `json:"uptime"`
}

type Uptime struct {
	Years   uint8 `json:"years"`
	Months  uint8 `json:"months"`
	Days    uint8 `json:"days"`
	Hours   uint8 `json:"hours"`
	Minutes uint8 `json:"minutes"`
	Seconds uint8 `json:"seconds"`
}
