package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"os"
	"time"
)

func newPool(addr string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", addr) },
	}
}

/*
keyがない場合は0を返す
*/
func get(key string, c redis.Conn) int {
	i, err := redis.Int(c.Do("GET", key))

	if i == 0 {
		return 0
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return i
}

func increment(key string, c redis.Conn) int {
	i, err := redis.Int(c.Do("INCR", key))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return i
}

func reset(key string, c redis.Conn) error {
	_, err := c.Do("DEL", key)
	return err
}
