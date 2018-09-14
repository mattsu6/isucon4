package main

import (
	"database/sql"
	"fmt"
	"log"
	_ "net/http/pprof"
	"os"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
)

var p *redis.Pool
var db *sql.DB

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

func del(key string, c redis.Conn) int {
	i, err := redis.Int(c.Do("DEL", key))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return i
}

func flush(c redis.Conn) string {
	i, err := redis.String(c.Do("FLUSHALL"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return i
}

func initLoginLog() {

	dsn := "root:@tcp(mysql:3306)/isu4_qualifier?parseTime=true&loc=Local"

	var err error

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	p = newPool("redis:6379")
	query := "SELECT ip, user_id, succeeded FROM login_log ORDER BY created_at"
	rows, err := db.Query(query)
	c := p.Get()
	flush(c)
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
		fmt.Println(ip + " : " + strconv.Itoa(id))
		if succeeded == 1 {
			del(ip, c)
			del(strconv.Itoa(id), c)
		} else if succeeded == 0 {
			increment(ip, c)
			increment(strconv.Itoa(id), c)
		}
	}

}

func getEnv(key string, def string) string {
	v := os.Getenv(key)
	if len(v) == 0 {
		return def
	}
	return v
}

func main() {
	initLoginLog()
}
