package config

import "os"

type Config struct {
	StorageType string
	PostgresDSN string
	ServerPort  string
}

func Load() *Config {
	postgresDSN := getEnv("POSTGRES_DSN", "")
	if postgresDSN == "" {
		postgresDSN = getEnv("POSTGRES_DNS", "")
	}

	return &Config{
		StorageType: getEnv("STORAGE", "memory"),
		PostgresDSN: postgresDSN,
		ServerPort:  getEnv("PORT", "8080"),
	}
}

func getEnv(key, definition string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return definition
}
