package redis

import (
	"context"
	"errors"
	"github.com/milad-rasouli/jaabz/internal/infra/godotenv"

	"github.com/redis/rueidis"
)

type Redis struct {
	Client rueidis.Client
	Env    *godotenv.Env
}

func NewRedis(env *godotenv.Env) *Redis {
	return &Redis{
		Env: env,
	}
}

func (r *Redis) Setup(ctx context.Context) error {
	if r.Client != nil {
		r.Client.Close()
	}

	client, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress: []string{r.Env.RedisHost},
	})
	if err != nil {
		return err
	}

	if err := client.Do(ctx, client.B().Ping().Build()).Error(); err != nil {
		return errors.New("failed to ping Redis: " + err.Error())
	}

	r.Client = client
	return nil
}

func (r *Redis) HealthCheck(ctx context.Context) error {
	if r.Client == nil {
		return errors.New("Redis client is not initialized")
	}

	if err := r.Client.Do(ctx, r.Client.B().Ping().Build()).Error(); err != nil {
		return errors.New("Redis health check failed: " + err.Error())
	}

	return nil
}

func (r *Redis) Close() error {
	if r.Client != nil {
		r.Client.Close()
	}
	return nil
}
