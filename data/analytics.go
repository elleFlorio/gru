package data

type GruAnalytics struct {
	Service map[string]AnalyticData `json:"service"`
	System  AnalyticData            `json:"system"`
}

type AnalyticData struct {
	BaseAnalytics map[string]float64
	UserAnalytics map[string]float64
}
