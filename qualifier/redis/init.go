import (
	"database/sql"
	"fmt"
	"redigo"
)

db := *sql.DB
c := *redis.Pool

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
	i, err := redis.String(c.Do("DEL", key))
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    return i
}

func init() {


	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local",
		getEnv("ISU4_DB_USER", "root"),
		getEnv("ISU4_DB_PASSWORD", ""),
		getEnv("ISU4_DB_HOST", "localhost"),
		getEnv("ISU4_DB_PORT", "3306"),
		getEnv("ISU4_DB_NAME", "isu4_qualifier"),
	)

	var err error

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	c = newPool("redis:6379")

}

func initLoginLog() {

	query:= "SELECT ip, user_id, succeeded FROM login_log ORDER BY created_at"
	rows, err := db.Query(query)
	if err != nil{
		panic(err)
	}
	for rows.Next() {
		var ip string
		var id int
		var scceeded int
		if err := rows.Scan(&ip, &user_id, &succeeded); err != nil{
			log.Fatal(err)
		}
		if succeeded == 1 {
			del(ip, c)
			del(id, c)
		}else if succeeded == 0 {
			increment(ip, c)
			increment(id, c)
		}
	}
	defer c.Close()
}

func main(){
	init()
	initLoginLog()
}