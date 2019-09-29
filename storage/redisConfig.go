package storage

import (
	"github.com/go-redis/redis"
	"time"
)

const (
	//全局自增ID
	URLIDKEY="next.url.id"
	//短地址和原地址映射
	ShortLinkKey="shortlink:%s:url"
	//地址hash
	UrlHashKey="urlhash:%s:url"
	//短地址详情
	ShortlinkDetailKey="shortlink:%s:detail"
)


type RedisClient struct {
	Cli *redis.Client
}

type UrlDetail struct {
	URL string `json:"url"`
	CreatedAt string `json:"create_at"`
	ExpiredInMinutes time.Duration `json:"exp"`
}


func NewRedisClient()*RedisClient{
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		DB:       0,
		Password: "",
	})
	if _,err := client.Ping().Result();err!=nil{
		panic(err)
	}
	return &RedisClient{
		Cli:client,
	}
}