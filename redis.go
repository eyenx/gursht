package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis"
)

// create redis connection
func redisConn() *redis.Client {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	if redisHost == "" {
		// default host
		redisHost = "localhost"
	}
	if redisPort == "" {
		// default port
		redisPort = "6379"
	}
	client := redis.NewClient(&redis.Options{
		Addr:     redisHost + ":" + redisPort,
		Password: "",
		DB:       0,
	})
	_, err := client.Ping().Result()
	if err != nil {
		log.Fatal(err)
	}
	return client
}

// write to redis
func redisWrite(u Url) Url {
	c := redisConn()
	err := c.Set(u.ShortUrl, u.LongUrl, 0).Err()
	c.Close()
	if err != nil {
		log.Fatal(err)
	}
	return u
}

// read from redis
func redisRead(s string) string {
	c := redisConn()
	val, err := c.Get(s).Result()
	c.Close()
	if err == redis.Nil {
		val = ""
	} else if err != nil {
		fmt.Println(err)
	}
	return val
}
