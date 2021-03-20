package wcache

import (
	"log"
	"time"

	cache "github.com/victorspringer/http-cache"
	"github.com/victorspringer/http-cache/adapter/redis"
)

// TODO: Make this more configurable
func CacheClient() *cache.Client {
	ringOpt := &redis.RingOptions{
		Addrs: map[string]string{
			"server": "127.0.0.1:6378",
		},
	}

	cacheClient, err := cache.NewClient(
		cache.ClientWithAdapter(redis.NewAdapter(ringOpt)),
		cache.ClientWithTTL(10*time.Minute),
		cache.ClientWithRefreshKey("opn"),
	)

	if err != nil {
		log.Fatal("Error with cacheClient")
	}

	return cacheClient
}
