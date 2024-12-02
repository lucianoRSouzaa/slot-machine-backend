package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	envPath := filepath.Join(".env")
	err := godotenv.Load(envPath)
	if err != nil {
		log.Println(err)
		log.Println("Erro ao carregar o arquivo .env, usando variáveis de ambiente do sistema")
	}
	return err
}

func GetEnv(key string) string {
	return os.Getenv(key)
}

func GetRequiredEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("A variável de ambiente %s não está definida", key)
	}
	return value
}
