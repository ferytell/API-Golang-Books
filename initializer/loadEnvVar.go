package initializer

import "github.com/joho/godotenv"

func LoadEnvVar() {
	_ = godotenv.Load()
}
