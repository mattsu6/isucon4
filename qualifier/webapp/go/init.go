package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
)

var p *redis.Pool

func newPool(addr string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", addr) },
	}
}

func increment(key string, c redis.Conn) int {
	i, err := redis.Int(c.Do("INCR", key))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return i
}

func setZero(key string, c redis.Conn) int {
	i, err := redis.Int(c.Do("SET", key, 0))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return i
}

func del(key string, c redis.Conn) string {
	i, err := redis.String(c.Do("DEL", key))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return i
}

func initLoginLog() {
	p = newPool("redis:6379")
	query := "SELECT ip, user_id, succeeded FROM login_log ORDER BY created_at"
	rows, err := db.Query(query)
	c := p.Get()
	defer c.Close()
	if err != nil {
		panic(err)
	}
	defer c.Close()
	for rows.Next() {
		var ip string
		var id int
		var succeeded int
		if err := rows.Scan(&ip, &id, &succeeded); err != nil {
			log.Fatal(err)
		}
		if succeeded == 1 {
			del(ip, c)
			del(strconv.Itoa(id), c)
		} else if succeeded == 0 {
			increment(ip, c)
			increment(strconv.Itoa(id), c)
		}
	}

}

func initredis() {
	initLoginLog()
}
