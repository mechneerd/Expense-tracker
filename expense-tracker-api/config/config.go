package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"expense-tracker-api/pkg/jwt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	DB           *pgxpool.Pool
	JWTSecretKey string
	TokenManager *jwt.TokenManager
	RedisAddr    string
	RedisPassword string
	OTLifetime   int
	Rdb          *redis.Client
}

var (
	configInstance *Config
	configOnce     sync.Once
)

func getEnv(key, defaultValue string) string {
	if val, exists := os.LookupEnv(key); exists {
		return val
	}
	return defaultValue
}

func containsParam(url, param string) bool {
	for i := 0; i < len(url)-len(param); i++ {
		if url[i] == '?' || url[i] == '&' {
			if url[i+1:i+len(param)+1] == param+"=" {
				return true
			}
		}
	}
	return false
}

func containsChar(s string, c byte) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			return true
		}
	}
	return false
}

func getEnvAsInt(key string, defaultValue int) int {
	if val, exists := os.LookupEnv(key); exists {
		var v int
		fmt.Sscanf(val, "%d", &v)
		return v
	}
	return defaultValue
}

func LoadConfig() *Config {
	configOnce.Do(func() {
		cfg := &Config{
			JWTSecretKey: getEnv("JWT_SECRET_KEY", "your-default-secret-key"),
			RedisAddr:    getEnv("REDIS_ADDR", "localhost:6379"),
			RedisPassword: getEnv("REDIS_PASSWORD", ""),
			OTLifetime:   getEnvAsInt("OT_LIFETIME_MINUTES", 5),
		}

		cfg.TokenManager = jwt.NewTokenManager(
			cfg.JWTSecretKey,
			15*time.Minute,  // Access token: 15 minutes
			30*24*time.Hour, // Refresh token: 30 days
		)

		cfg.Rdb = redis.NewClient(&redis.Options{
			Addr:     cfg.RedisAddr,
			Password: cfg.RedisPassword,
			DB:       0,
		})

		// Test Redis connection
		ctx := context.Background()
		if err := cfg.Rdb.Ping(ctx).Err(); err != nil {
			log.Printf("Warning: Could not connect to Redis: %v", err)
			cfg.Rdb = nil
		}

		dbURL := getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/expense_tracker?sslmode=disable")

		pgSSLMode := getEnv("PGSSLMODE", "")
		if pgSSLMode != "" && !containsParam(dbURL, "sslmode") {
			separator := "&"
			if !containsChar(dbURL, '?') {
				separator = "?"
			}
			dbURL = dbURL + separator + "sslmode=" + pgSSLMode
		}

		pool, err := pgxpool.New(context.Background(), dbURL)
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}

		cfg.DB = pool
		configInstance = cfg
	})

	return configInstance
}
