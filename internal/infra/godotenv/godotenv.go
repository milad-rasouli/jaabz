package godotenv

import (
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	AppName     string
	Environment string
	JaabzHost   string
	RedisHost   string
}

func NewEnv() *Env {
	e := &Env{}
	e.Load()
	return e
}

func (e *Env) Load() {
	godotenv.Load(".env")
	e.AppName = os.Getenv("APP_NAME")
	e.Environment = os.Getenv("ENVIRONMENT")
	e.JaabzHost = os.Getenv("JAABZ_HOST")
	e.RedisHost = os.Getenv("REDIS_HOST")
}
