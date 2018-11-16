package models

import (
    "github.com/go-redis/redis"
    "log"
)

var RedisClient *redis.Client

func InitRedis(server string, password string, database int) {

    RedisClient = redis.NewClient(&redis.Options{
	    Addr:     server,
	    Password: password,  // no password set
	    DB:       database,  // use default DB
	})

    pong, err := RedisClient.Ping().Result()
    if err != nil {
        log.Panic(err)
    } else {
        log.Print(pong)
    }
}
