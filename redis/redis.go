package redis

import (
	"fmt"
	"sync"

	redis0 "github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

var db *redis0.Client
var once sync.Once

func GetRedis() *redis0.Client {
	if db == nil {
		once.Do(func() {
			db = redis0.NewClient(&redis0.Options{
				Addr:     fmt.Sprintf("%s:%s", viper.GetString("redis.addr"), viper.GetString("redis.port")),
				Password: viper.GetString("redis.password"),
				DB:       0,
			})
		})
	}
	return db
}
