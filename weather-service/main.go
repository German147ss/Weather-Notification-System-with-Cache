// TODO 1 : usar inyeccion de dependicias para no ser necesario usar redis y el cptc , si no mocks.
// TODO 2 : usar el logger para mostrar los errores y respuestas
// TODO 3 : implementar el cron job para cachear los climas

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"golang.org/x/net/html/charset"
)

// define ctx for redis
var ctx = context.Background()

// Redis client
var rdb *redis.Client

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

// CharsetReader function for reading ISO-8859-1 encoded documents
func CharsetReader(char string, input io.Reader) (io.Reader, error) {
	if strings.EqualFold(char, "ISO-8859-1") {
		return charset.NewReaderLabel(char, input)
	}
	return nil, fmt.Errorf("unsupported charset: %s", char)
}

// Función para obtener el clima de una ciudad desde el CPTEC
func getWeatherFromCptc(city string) (*CptcResponse, error) {
	url := fmt.Sprintf("http://servicos.cptec.inpe.br/XML/cidade/%s/previsao.xml", city)
	resp, err := http.Get(url)
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

	// Esta función ajusta automáticamente la codificación
	decoder.CharsetReader = charset.NewReaderLabel

	var weather CptcResponse
	err = decoder.Decode(&weather)
	if err != nil {
		return nil, err
	}

	return &weather, nil
}

// Función para obtener el clima con Redis como caché
func getWeatherWithCache(city string) (*CptcResponse, error) {
	// Generar la clave para Redis
	cacheKey := "weather:" + city

	// Intentar obtener el clima desde el caché de Redis
	val, err := rdb.Get(ctx, cacheKey).Result()

	if err == redis.Nil {
		// Clima no encontrado en caché, hacer la solicitud a la API del CPTEC
		fmt.Println("Cache miss. Obteniendo datos del CPTEC...")
		weather, err := getWeatherFromCptc(city)
		if err != nil {
			return nil, err
		}

		// Convertir la estructura a JSON para almacenarla en Redis
		jsonData, err := json.Marshal(weather)
		if err != nil {
			return nil, err
		}

		// Almacenar el resultado en Redis con una expiración de 1 hora
		err = rdb.Set(ctx, cacheKey, jsonData, time.Hour).Err()
		if err != nil {
			return nil, err
		}

		fmt.Println("Datos almacenados en caché para la ciudad:", city)

		return weather, nil
	} else if err != nil {
		return nil, err
	}

	// Clima encontrado en el caché, deserializarlo
	fmt.Println("Cache hit. Obteniendo datos de Redis...")
	var weather CptcResponse
	err = json.Unmarshal([]byte(val), &weather)
	if err != nil {
		return nil, err
	}

	return &weather, nil
}

/* // Función para cachear el clima de una ciudad
func cacheWeather(city string) error {
	weather, err := getWeatherFromCptc(city)
	if err != nil {
		return err
	}

	// Convertir la estructura a JSON para almacenarla en Redis
	jsonData, err := json.Marshal(weather)
	if err != nil {
		return err
	}

	// Almacenar el resultado en Redis con una expiración de 1 hora
	cacheKey := "weather:" + city
	err = rdb.Set(ctx, cacheKey, jsonData, time.Hour).Err()
	if err != nil {
		return err
	}

	fmt.Println("Datos almacenados en caché para la ciudad:", city)
	return nil
}

// Job que se ejecutará cada hora para cachear el clima de varias ciudades
func cacheWeatherJob() {
	cities := []string{"3477", "244", "3956"} // IDs de ejemplo de ciudades: São Paulo, Rio de Janeiro, Belo Horizonte
	for _, city := range cities {
		fmt.Println("Consultando y cacheando clima para la ciudad:", city)
		err := cacheWeather(city)
		if err != nil {
			fmt.Println("Error cacheando el clima para la ciudad:", city, err)
		}
	}
} */

func main() {
	// Inicializar el cliente Redis
	port := os.Getenv("REDIS_PORT")
	host := os.Getenv("REDIS_HOST")

	address := host + ":" + port
	fmt.Println("Conectando a Redis en:", address)
	rdb = redis.NewClient(&redis.Options{Addr: address, DB: 0})

	//pin redis
	status := rdb.Ping(ctx)
	if status.Err() != nil {
		fmt.Println("Error al conectar a Redis:", status.Err())
		panic(status.Err())
	}

	/* // Iniciar el cron job
	c := cron.New()

	// Ejecutar cada hora
	c.AddFunc("@hourly", func() {
		fmt.Println("Ejecutando el job para cachear clima...")
		cacheWeatherJob()
	})

	c.Start() */

	http.HandleFunc("/weather/", func(w http.ResponseWriter, r *http.Request) {
		// Obtener la ciudad desde la URL
		city := r.URL.Path[len("/weather/"):]

		// Obtener el clima con cache
		weather, err := getWeatherWithCache(city)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Convertir la estructura a JSON
		jsonData, err := json.Marshal(weather)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Escribir la respuesta JSON
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	})

	fmt.Println("Servidor iniciado en el puerto 8083")
	http.ListenAndServe(":8083", nil)
}
