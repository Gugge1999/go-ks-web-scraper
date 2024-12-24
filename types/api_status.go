package types

type ApiStatus struct {
	Status                    string `json:"status"` // active | inactive | pending
	ScrapingIntervalInMinutes uint   `json:"scrapingIntervalInMinutes"`
	MemoryUsage               uint64 `json:"memoryUsage"`
	NumberOfCpus              int    `json:"numberOfCpus"`
	Uptime                    Uptime `json:"uptime"`
}

type Uptime struct {
	Years   uint16 `json:"years"`
	Months  uint8  `json:"months"`
	Days    uint16 `json:"days"`
	Hours   uint8  `json:"hours"`
	Minutes uint8  `json:"minutes"`
	Seconds uint8  `json:"seconds"`
}
