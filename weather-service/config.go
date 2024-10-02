package main

import "os"

func NewConfig() *Config {
	RedisHost := os.Getenv("REDIS_HOST")
	if RedisHost == "" {
		RedisHost = "localhost"
	}
	RedisPort := os.Getenv("REDIS_PORT")
	if RedisPort == "" {
		RedisPort = "6379"
	}
	AppPort := os.Getenv("APP_PORT")
	if AppPort == "" {
		AppPort = "8083"
	}
	CacheType := os.Getenv("CACHE_TYPE")
	if CacheType == "" {
		CacheType = "memory"
	}
	return &Config{
		RedisHost: RedisHost,
		RedisPort: RedisPort,
		AppPort:   AppPort,
		CacheType: CacheType,
	}
}

type Config struct {
	RedisHost string
	RedisPort string
	AppPort   string
	CacheType string
}
