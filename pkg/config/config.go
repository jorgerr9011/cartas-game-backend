package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Db_redis string
}

func LoadConfig() (*Config, error) {
	// Carga condicional del .env (solo en desarrollo)
	if os.Getenv("ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Println("No se pudo cargar el archivo .env, usando variables del entorno")
		}
	}

	// Leer las variables de entorno
	return &Config{
		Db_redis: getEnv("REDIS_ADDR", "redis"),
	}, nil
}

// getEnv obtiene una variable de entorno o devuelve un valor predeterminado
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
