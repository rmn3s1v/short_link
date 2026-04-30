package config

import "os"

type Config struct {
	StorageType string
	PostgresDNS string
	ServerPort  string
}

func Load() *Config {
	return &Config{
		StorageType: getEnv("STORAGE", "memory"),
		PostgresDNS: getEnv("POSTGRES_DNS", ""),
		ServerPort:  getEnv("PORT", ""),
	}
}

func getEnv(key, definition string) string{
	if v := os.Getenv(key); v != ""{
		return  v
	}

	return definition
}
