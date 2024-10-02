package entities

type User struct {
	ID                   int    `json:"id"`
	LocationCode         string `json:"location_code"`
	NotificationSchedule int    `json:"notification_schedule"`
	IsEnabled            bool   `json:"is_enabled"`
	State                string `json:"state"`
}

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

type CityWaves struct {
	Name       string         `json:"name"`
	State      string         `json:"state"`
	LastUpdate string         `json:"last_update"`
	Morning    WavePrediction `json:"morning"`
	Afternoon  WavePrediction `json:"afternoon"`
	Night      WavePrediction `json:"night"`
}

type WavePrediction struct {
	Day        string  `json:"day"`
	SeaStatus  string  `json:"sea_status"`
	WaveHeight float64 `json:"wave_height"`
	WaveDir    string  `json:"wave_direction"`
	WindSpeed  float64 `json:"wind_speed"`
	WindDir    string  `json:"wind_direction"`
}

type WeatherAndWaves struct {
	Weather *CityWeather `json:"weather"`
	Waves   *CityWaves   `json:"waves"`
}

type WeatherResponse struct {
	Nome        string     `json:"nome"`
	UF          string     `json:"uf"`
	Atualizacao string     `json:"atualizacao"`
	Previsoes   []Previsao `json:"previsao"`
}

type Previsao struct {
	Dia    string  `json:"dia"`
	Tempo  string  `json:"tempo"`
	Maxima int     `json:"maxima"`
	Minima int     `json:"minima"`
	IUV    float32 `json:"iuv"`
}
