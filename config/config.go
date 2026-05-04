package config

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

var (
	pgPool      *pgxpool.Pool
	pgPoolOnce  sync.Once

	redisClient     *redis.Client
	redisClientOnce sync.Once

	limiterClient     *redis.Client
	limiterClientOnce sync.Once
)

func GetPgPool() *pgxpool.Pool {
	pgPoolOnce.Do(func() {
		v_dsn := fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s pool_min_conns=%s pool_max_conns=%s connect_timeout=5",
			os.Getenv("POSTGRES_DB_HOST"),
			os.Getenv("POSTGRES_DB_PORT"),
			os.Getenv("POSTGRES_DB_USER"),
			os.Getenv("POSTGRES_DB_PASS"),
			os.Getenv("POSTGRES_DB_DATA"),
			os.Getenv("POSTGRES_POOL_MIN"),
			os.Getenv("POSTGRES_POOL_MAX"),
		)
		v_pool, v_err := pgxpool.New(context.Background(), v_dsn)
		if v_err != nil {
			panic(fmt.Sprintf("Failed to connect to PostgreSQL: %v", v_err))
		}
		pgPool = v_pool
	})
	return pgPool
}

func GetRedisClient() *redis.Client {
	redisClientOnce.Do(func() {
		v_opt, v_err := redis.ParseURL(os.Getenv("REDIS_URL"))
		if v_err != nil {
			panic(fmt.Sprintf("Failed to parse Redis URL: %v", v_err))
		}
		v_opt.PoolSize         = 20
		v_opt.DialTimeout      = 3
		v_opt.ReadTimeout      = 3
		v_opt.WriteTimeout     = 3
		redisClient            = redis.NewClient(v_opt)
	})
	return redisClient
}

func GetLimiterClient() *redis.Client {
	limiterClientOnce.Do(func() {
		v_opt, v_err := redis.ParseURL(os.Getenv("REDIS_URL"))
		if v_err != nil {
			panic(fmt.Sprintf("Failed to parse Redis URL: %v", v_err))
		}
		v_opt.PoolSize     = 10
		v_opt.DialTimeout  = 3
		v_opt.ReadTimeout  = 3
		v_opt.WriteTimeout = 3
		limiterClient      = redis.NewClient(v_opt)
	})
	return limiterClient
}

func GetPort() string {
	v_port := os.Getenv("PORT")
	if v_port == "" {
		return "8080"
	}
	return v_port
}
