// internal/cache/cache.go
package cache

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"crypto/tls"
	"time"
	// "fmt"

	"coding-profile-service/pkg/model"
	"github.com/redis/go-redis/v9"
)

var (
	client *redis.Client
	ctx    = context.Background()
)

func Init() {
    addr := os.Getenv("REDIS_ADDR")
    if addr == "" {
        addr = "localhost:6379"
    }

    password := os.Getenv("REDIS_PASSWORD")

    // check if its local or upstash
    var tlsConfig *tls.Config
    if addr != "localhost:6379" {
        tlsConfig = &tls.Config{} // ← Upstash needs TLS
    }

    client = redis.NewClient(&redis.Options{
        Addr:      addr,
        Password:  password,
        DB:        0,
        TLSConfig: tlsConfig, // ← add this
    })

    if err := client.Ping(ctx).Err(); err != nil {
        log.Printf("⚠️  Redis not connected: %v — requests will hit scrapers directly", err)
        client = nil
    } else {
        log.Println("✅ Redis connected at", addr)
    }
}

func SetCache(key string, resp model.StatsResponse, ttl time.Duration) {
	if client == nil {
		return
	}
	data, err := json.Marshal(resp)
	if err != nil {
		return
	}
	client.Set(ctx, key, data, ttl)
}

func GetCache(key string) (model.StatsResponse, bool) {
	if client == nil {
		return model.StatsResponse{}, false
	}
	val, err := client.Get(ctx, key).Result()
	if err != nil {
		return model.StatsResponse{}, false // cache miss
	}
	var resp model.StatsResponse
	if err := json.Unmarshal([]byte(val), &resp); err != nil {
		return model.StatsResponse{}, false
	}
	return resp, true
}