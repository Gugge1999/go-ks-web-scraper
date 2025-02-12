package types

type ApiStatus struct {
	Status                    string `json:"status"` // active | inactive | pending
	ScrapingIntervalInMinutes uint   `json:"scrapingIntervalInMinutes"`
	MemoryUsage               uint64 `json:"memoryUsage"`
	NumberOfCpus              int    `json:"numberOfCpus"`
	Uptime                    Uptime `json:"uptime"`
}

type Uptime struct {
	Years   uint `json:"years"`
	Months  uint `json:"months"`
	Days    uint `json:"days"`
	Hours   uint `json:"hours"`
	Minutes uint `json:"minutes"`
	Seconds uint `json:"seconds"`
}
