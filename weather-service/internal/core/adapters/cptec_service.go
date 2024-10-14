package adapters

import (
	"app/internal/core/domain"
	"app/internal/core/ports"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/net/html/charset"
)

type CptcResponse struct {
	XMLName     xml.Name   `xml:"cidade"`
	Nome        string     `xml:"nome"`
	UF          string     `xml:"uf"`
	Atualizacao string     `xml:"atualizacao"`
	Previsoes   []Previsao `xml:"previsao"`
}

type Previsao struct {
	Dia    string  `xml:"dia"`
	Tempo  string  `xml:"tempo"`
	Maxima int     `xml:"maxima"`
	Minima int     `xml:"minima"`
	IUV    float32 `xml:"iuv"`
}

func (c *CptcResponse) GetWeather() domain.CityWeather {
	if c == nil {
		return domain.CityWeather{}
	}

	if len(c.Previsoes) == 0 {
		return domain.CityWeather{}
	}
	var forecasts []domain.Forecast
	for _, p := range c.Previsoes {
		forecasts = append(forecasts, domain.Forecast{
			Day:     p.Dia,
			Weather: p.Tempo,
			MaxTemp: p.Maxima,
			MinTemp: p.Minima,
			UvIndex: p.IUV,
		})
	}

	return domain.CityWeather{
		Name:       c.Nome,
		State:      c.UF,
		LastUpdate: c.Atualizacao,
		Forecasts:  forecasts,
	}
}

type CptecCiudadesReponse struct {
	XMLName xml.Name `xml:"cidades"`
	Cidade  Cidade   `xml:"cidade"`
}

type Cidade struct {
	Nome string `xml:"nome"`
	UF   string `xml:"uf"`
	ID   string `xml:"id"`
}

func (c *CptecCiudadesReponse) GetID() string {
	return c.Cidade.ID
}

type CptecWavesResponse struct {
	XMLName     xml.Name        `xml:"cidade"`
	Nome        string          `xml:"nome"`
	UF          string          `xml:"uf"`
	Atualizacao string          `xml:"atualizacao"`
	Manha       WavesPrediction `xml:"manha"`
	Tarde       WavesPrediction `xml:"tarde"`
	Noite       WavesPrediction `xml:"noite"`
}

type WavesPrediction struct {
	Dia      string  `xml:"dia"`
	Agitacao string  `xml:"agitacao"`
	Altura   float64 `xml:"altura"`
	Direcao  string  `xml:"direcao"`
	Vento    float64 `xml:"vento"`
	VentoDir string  `xml:"vento_dir"`
}

func (w *CptecWavesResponse) GetWaves() domain.CityWaves {
	return domain.CityWaves{
		Name:       w.Nome,
		State:      w.UF,
		LastUpdate: w.Atualizacao,
		Morning: domain.WavePrediction{
			Day:        w.Manha.Dia,
			SeaStatus:  w.Manha.Agitacao,
			WaveHeight: w.Manha.Altura,
			WaveDir:    w.Manha.Direcao,
			WindSpeed:  w.Manha.Vento,
			WindDir:    w.Manha.VentoDir,
		},
		Afternoon: domain.WavePrediction{
			Day:        w.Tarde.Dia,
			SeaStatus:  w.Tarde.Agitacao,
			WaveHeight: w.Tarde.Altura,
			WaveDir:    w.Tarde.Direcao,
			WindSpeed:  w.Tarde.Vento,
			WindDir:    w.Tarde.VentoDir,
		},
		Night: domain.WavePrediction{
			Day:        w.Noite.Dia,
			SeaStatus:  w.Noite.Agitacao,
			WaveHeight: w.Noite.Altura,
			WaveDir:    w.Noite.Direcao,
			WindSpeed:  w.Noite.Vento,
			WindDir:    w.Noite.VentoDir,
		},
	}
}

// WEATHER AND WAVES
type WeatherAndWaves struct {
	Weather *CptcResponse
	Waves   *CptecWavesResponse
}

const maxRetries = 3
const backoff = 1 * time.Second

type CPTECWeatherService struct {
	client *http.Client
}

func NewCPTECWeatherService() ports.WeatherService {
	client := &http.Client{
		Timeout: time.Second * 30, // Incrementar tiempo de espera
	}
	return &CPTECWeatherService{client: client}
}

func (c *CPTECWeatherService) SearchIdByName(cityName string) (string, error) {
	url := fmt.Sprintf("http://servicos.cptec.inpe.br/XML/listaCidades?city=%s", cityName)
	resp, err := c.client.Get(url)
	if err != nil {
		fmt.Println("Error al buscar el ID de la ciudad:", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	reader := bytes.NewReader(body)
	decoder := xml.NewDecoder(reader)

	decoder.CharsetReader = charset.NewReaderLabel

	var city CptecCiudadesReponse
	err = decoder.Decode(&city)
	if err != nil {
		fmt.Println("SearchIdByName - Error al decodificar el XML:", err)
		return "", err
	}

	return city.GetID(), nil
}

func (c *CPTECWeatherService) GetWeather(city string) (*domain.CityWeather, error) {
	url := fmt.Sprintf("http://servicos.cptec.inpe.br/XML/cidade/%s/previsao.xml", city)
	resp, err := c.retryRequest(url, maxRetries, backoff)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(body)
	decoder := xml.NewDecoder(reader)

	// Adjust the encoding automatically
	decoder.CharsetReader = charset.NewReaderLabel

	var weather CptcResponse
	err = decoder.Decode(&weather)
	if err != nil {
		fmt.Println("GetWeather - Error al decodificar el XML:", err)
		return nil, err
	}

	result := weather.GetWeather()

	return &result, nil
}

// Get waves from CPTEC
func (c *CPTECWeatherService) GetWaves(city string) (*domain.CityWaves, error) {
	url := fmt.Sprintf("http://servicos.cptec.inpe.br/XML/cidade/%s/dia/0/ondas.xml", city)
	resp, err := c.retryRequest(url, maxRetries, backoff)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(body)
	decoder := xml.NewDecoder(reader)

	// Adjust the encoding automatically
	decoder.CharsetReader = charset.NewReaderLabel

	var waves CptecWavesResponse
	err = decoder.Decode(&waves)
	if err != nil {
		fmt.Printf("GetWaves - error al decodificar el XML: %v, se procede a ignorar", err)
		return nil, nil
	}
	result := waves.GetWaves()
	return &result, nil
}

func (c *CPTECWeatherService) retryRequest(url string, maxRetries int, backoff time.Duration) (*http.Response, error) {
	var resp *http.Response
	var err error

	for i := 0; i < maxRetries; i++ {
		resp, err = c.client.Get(url)
		if err == nil {
			return resp, nil
		}

		fmt.Printf("Attempt %d: Error making request to %s: %v\n", i+1, url, err)
		time.Sleep(backoff)
	}

	return nil, fmt.Errorf("failed to make request to %s after %d attempts: %v", url, maxRetries, err)
}
