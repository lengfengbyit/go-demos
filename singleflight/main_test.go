package main

import (
	"context"
	"database/sql"
	"log"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/singleflight"
)

var group singleflight.Group
var rdb *redis.Client
var db *sql.DB

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:            ":6379",
		Password:        "",
		DB:              0,
		PoolSize:        20,
		PoolTimeout:     time.Second * 10,
		MaxIdleConns:    5,
		ConnMaxIdleTime: time.Second * 10,
	})
	rdb.Ping(context.Background())
	log.Println("redis connect success")

	var err error
	db, err = sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/test")
	failOnError(err)

	log.Println(sql.Drivers())
}

func failOnError(err error) {
	if err != nil && err != redis.Nil {
		panic(err)
	}
}

// go test -v -run TestRedis
// 并发查询一个不存在的 key
func TestRedis(t *testing.T) {
	ctx := context.Background()
	_, err := rdb.Get(ctx, "key").Result()
	failOnError(err)
}

func TestMysql(t *testing.T) {
	id := queryMysql()
	log.Printf("id: %d", id)
}

func queryMysql() (id int64) {
	err := db.QueryRow(`select id  from users where name = ?`, "Mike").Scan(&id)
	//failOnError(err)
	prepare, err := db.Prepare("select id from users where name = ?")
	failOnError(err)

	err = prepare.QueryRow("Mike").Scan(&id)
	failOnError(err)
	return
}

func TestCache(t *testing.T) {
	name := "Mike"
	ctx := context.Background()
	result, err := rdb.Get(ctx, name).Result()
	failOnError(err)
	if err == redis.Nil {
		val, err, shared := group.Do(name, func() (interface{}, error) {
			return queryMysql(), nil
		})
		failOnError(err)
		log.Printf("shared: %v, val: %v", shared, val)
		//result = val.(int64)
	} else {
		log.Printf("result: %s", result)
	}

}
