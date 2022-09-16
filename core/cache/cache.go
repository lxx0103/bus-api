package cache

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"bus-api/core/config"
	"bus-api/core/log"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
)

var (
	mycache = &cache.Cache{}
	ctx     = context.Background()
)

func ConfigCache() *cache.Cache {
	addr := config.ReadConfig("cache.host")
	password := config.ReadConfig("cache.password")
	db, _ := strconv.Atoi(config.ReadConfig("cache.db"))
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	if err := rdb.Ping(ctx).Err(); err != nil {
		panic("Unable to connect to redis " + err.Error())
	}
	mycache = cache.New(&cache.Options{
		Redis: rdb,
	})

	return mycache
}

//GetKey get key
func GetKey(key string, value string) {
	err := mycache.Get(ctx, key, &value)
	if err != nil {
		log.Debug(err.Error())
	}
	fmt.Println(value)
}

//SetKey set key
func SetKey(key string, value interface{}) error {
	err := mycache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: value,
		TTL:   time.Hour,
	})
	if err != nil {
		log.Debug(err.Error())
		return err
	}
	return nil
}
