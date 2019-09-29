package storage

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

type Storage interface {
	//长地址转换为短地址
	Shorten(url string,exp int64)(string,error)
	//短地址信息
	ShortInfo(eid string)(interface{},error)
	//短地址转换为长地址
	UnShorten(url string)(string,error)
}



func (r *RedisClient)Shorten(url string,exp int64)(string,error){
	h:=Sha1String(url)
	if v,err := r.Cli.Get(fmt.Sprintf(UrlHashKey,h)).Result();err!=nil{
		if err==redis.Nil{
			panic(err)
		}else{
			return "",err
		}
	}else{
		if v =="{}"{
			//过期
		}else{
			return v,nil
		}
	}

	err := r.Cli.Incr(URLIDKEY).Err()
	if err!=nil{
		return "",err
	}

	id, err := r.Cli.Get(URLIDKEY).Int64()
	if err!=nil{
		return "",err
	}
	eid:=base62encode(id)

	err = r.Cli.Set(fmt.Sprintf(ShortLinkKey, eid), url, time.Minute*time.Duration(exp)).Err()
	if err!=nil{
		return "",err
	}

	err= r.Cli.Set(fmt.Sprintf(UrlHashKey, h), eid, time.Minute*time.Duration(exp)).Err()
	if err!=nil{
		return "",err
	}

	detail, err := json.Marshal(&UrlDetail{
		URL:              url,
		CreatedAt:        time.Now().String(),
		ExpiredInMinutes: time.Minute * time.Duration(exp),
	})

	if err!=nil{
		return "",err
	}

	err = r.Cli.Set(fmt.Sprintf(ShortlinkDetailKey, eid), detail, time.Minute*time.Duration(exp)).Err()
	if err!=nil{
		return "",err
	}
	return eid,nil
}




func (r *RedisClient)ShortInfo(eid string)(interface{},error){
	info, e := r.Cli.Get(fmt.Sprintf(ShortlinkDetailKey, eid)).Result()
	if e!=nil{
		return nil,e
	}
	return info,nil
}


func (r *RedisClient)UnShorten(url string)(string,error){
	v, e := r.Cli.Get(fmt.Sprintf(ShortLinkKey, url)).Result()
	if e!=nil{
		return "",e
	}
	return v,nil
}