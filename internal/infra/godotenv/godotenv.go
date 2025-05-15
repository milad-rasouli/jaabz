package godotenv

import (
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	AppName           string
	Environment       string
	JaabzHost         string
	RedisHost         string
	TelegramChannelID string
	TelegramBotToken  string
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
	e.TelegramChannelID = os.Getenv("TELEGRAM_CHANNEL_ID")
	e.TelegramBotToken = os.Getenv("TELEGRAM_BOT_TOKEN")
}
