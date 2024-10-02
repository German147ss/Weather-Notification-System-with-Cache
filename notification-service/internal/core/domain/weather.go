package domain

type CityWeather struct {
	Name       string     `json:"name"`
	State      string     `json:"state"`
	LastUpdate string     `json:"last_update"`
	Forecasts  []Forecast `json:"forecasts"`
}

type Forecast struct {
	Day     string  `json:"day"`
	Weather string  `json:"weather"`
	MaxTemp int     `json:"max_temp"`
	MinTemp int     `json:"min_temp"`
	UvIndex float32 `json:"uv_index"`
}
