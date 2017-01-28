package model

import (
  "fmt"
  "github.com/garyburd/redigo/redis"
)

type RedisMq struct {
  RedisClient *redis.Pool
  redisHost string
  redisDB int
  C       redis.Conn
}

func InitRedisMq(redisHost string, redisDB int) (*RedisMq, error) {
  rmq := &RedisMq{
    redisHost : redisHost,
    redisDB : redisDB,
  }
  // rmq.RedisClient = &redis.Pool{
  //   MaxIdle:     1,
	// 	MaxActive:   10,
	// 	IdleTimeout: 180 * time.Second,
	// 	Dial: func() (redis.Conn, error) {
	// 		c, err := redis.Dial("tcp", rmq.redisHost)
	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 		// 选择db
	// 		c.Do("SELECT", rmq.redisDB)
	// 		return c, nil
	// 	},
  // }
  c, err := redis.Dial("tcp", redisHost)
  if err != nil {
      // handle error
      fmt.Println("redis dial error: " + err.Error())
      return nil, err
  }
  c.Do("SELECT", rmq.redisDB)
  rmq.C = c
  return rmq, nil
}

func (rmq *RedisMq) GetUrls() (urls []string){
  c := rmq.C
  // defer rc.Close()
  //values, _ := redis.Values(rc.Do("lrange", "redlist", "0", "100")))
  n, _ := redis.Int(c.Do("LLEN", "url"))

  fmt.Printf("urls length: %#v/n", n)
  // if len(urls) < 100 then load data from mongodb
  if n < 100 {
    go loadDataFromMongod()
    urls, _ = redis.Strings(c.Do("lrange", "url", "0", "-1"))
    for _, url := range urls {
      fmt.Printf("get urls from redis: " + url)
    }
  } else {
    urls, _ = redis.Strings(c.Do("lrange", "url", "0", "100"))
  }
  return urls
}

func (rmq *RedisMq) PushUrls(urls []string) {
  // rc := rmq.RedisClient.Get()
  rc := rmq.C
  // defer rc.Close()
  //values, _ := redis.Values(rc.Do("lrange", "redlist", "0", "100")))
  // for url := l.Front; url != nil; url = url.Next() {
  rc.Do("lpush", "url", urls)
  // }
}

func loadDataFromMongod() {
  //1) queru 1000 urls from mongodb
  //2) push urls to redismq
}
