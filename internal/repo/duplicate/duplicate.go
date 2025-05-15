package duplicate

import (
	"context"
	"errors"
	"github.com/milad-rasouli/jaabz/internal/error_list"
	"github.com/milad-rasouli/jaabz/internal/infra/redis"
	"github.com/redis/rueidis"
	"log/slog"
)

const (
	hllScript = `
		local key = KEYS[1]
		local data = ARGV[1]
		local count_before = redis.call("PFCOUNT", key)
		redis.call("PFADD", key, data)
		local count_after = redis.call("PFCOUNT", key)
		if count_before == count_after then
			return 1
		end
		return 0
	`
)

type Duplicate struct {
	rdis      *redis.Redis
	logger    *slog.Logger
	hllScript *rueidis.Lua
}

func New(logger *slog.Logger, rds *redis.Redis) *Duplicate {
	return &Duplicate{
		rdis:      rds,
		logger:    logger.With("repo", "duplicate"),
		hllScript: rueidis.NewLuaScript(hllScript),
	}
}

func (d *Duplicate) getKey() string {
	return "jaabz"
}

func (d *Duplicate) SaveAndCheckDuplicate(ctx context.Context, data string) error {
	lg := d.logger.With("method", "SaveAndCheckDuplicate")
	key := d.getKey()
	result, err := d.hllScript.Exec(ctx, d.rdis.Client, []string{key}, []string{data}).ToInt64()
	if err != nil {
		lg.Error("failed to check", "error", err)
		return errors.New("failed to execute HyperLogLog script: " + err.Error())
	}

	if result == 1 {
		lg.Error("duplicate found", "data", data)
		return error_list.ErrDuplicate
	}

	lg.Info("not duplicate", "data", data)
	return nil
}
