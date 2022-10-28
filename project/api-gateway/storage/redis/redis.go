package redis

import (
	"github.com/gomodule/redigo/redis"
	"github.com/project/api-gateway/storage/repo"
)

type redisRepo struct {
	rds *redis.Pool
}

func NewRedisRepo(rds *redis.Pool) repo.InMemoryStorageI {
	return &redisRepo{
		rds: rds,
	}
}

func (r *redisRepo) SetWithTTL(key, value string, seconds int64) (err error) {
	conn := r.rds.Get()
	defer conn.Close()

	_, err = conn.Do("SETEX", key, seconds, value)
	return
}

func (r *redisRepo) Get(key string) (interface{}, error) {
	conn := r.rds.Get()
	defer conn.Close()
	
	return conn.Do("GET", key)
}
