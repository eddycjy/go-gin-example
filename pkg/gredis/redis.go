package gredis

import (
	"encoding/json"
	"time"

	"github.com/gomodule/redigo/redis"

	"github.com/EDDYCJY/go-gin-example/pkg/setting"
)

var RedisConn *redis.Pool

// Setup Initialize the Redis instance
func Setup() error {
	RedisConn = &redis.Pool{
		MaxIdle:     setting.RedisSetting.MaxIdle,
		MaxActive:   setting.RedisSetting.MaxActive,
		IdleTimeout: setting.RedisSetting.IdleTimeout,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", setting.RedisSetting.Host)
			if err != nil {
				return nil, err
			}
			if setting.RedisSetting.Password != "" {
				if _, err := c.Do("AUTH", setting.RedisSetting.Password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	return nil
}

// Set a key/value
func Set(key string, data interface{}, time int) error {
	conn := RedisConn.Get()
	defer conn.Close()

	value, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = conn.Do("SET", key, value)
	if err != nil {
		return err
	}

	_, err = conn.Do("EXPIRE", key, time)
	if err != nil {
		return err
	}

	return nil
}

// Exists check a key
func Exists(key string) bool {
	conn := RedisConn.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false
	}

	return exists
}

// Get get a key
func Get(key string) ([]byte, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	reply, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}

	return reply, nil
}

// Delete delete a kye
func Delete(key string) (bool, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	return redis.Bool(conn.Do("DEL", key))
}

// LikeDeletes batch delete
func LikeDeletes(key string) error {
	conn := RedisConn.Get()
	defer conn.Close()

	keys, err := redis.Strings(conn.Do("KEYS", "*"+key+"*"))
	if err != nil {
		return err
	}

	for _, key := range keys {
		_, err = Delete(key)
		if err != nil {
			return err
		}
	}

	return nil
}


//原子增加数值increment
func IncrBy(key string, increment int) (int, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	value, err := redis.Int(conn.Do("INCR", key))
	if err != nil {
		return 0, err
	}

	return value, nil

}

//原子增加数值1
func Incr(key string) (int, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	value, err := redis.Int(conn.Do("INCR", key))
	if err != nil {
		return 0, err
	}

	return value, nil
}

//设置过期时间戳
func ExpireAt(key string, timestamp int64) error {
	conn := RedisConn.Get()
	defer conn.Close()

	_, err := conn.Do("EXPIREAT", key, timestamp)
	if err != nil {
		return err
	}

	return nil
}

//原子增加数值1,设置过期时间戳
func IncrExpireAt(key string, timestamp int64) (int, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	value, err := redis.Int(conn.Do("INCR", key))
	if err != nil {
		return 0, err
	}

	_, err = conn.Do("EXPIREAT", key, timestamp)
	if err != nil {
		return 0, err
	}

	return value, nil
}

//LPush
func LPush(key string, value interface{}) error {
	conn := RedisConn.Get()
	defer conn.Close()

	_, err := conn.Do("LPUSH", key, value)
	if err != nil {
		return err
	}

	return nil

}

//LPush
func LRem(key string, count int, value interface{}) error {
	conn := RedisConn.Get()
	defer conn.Close()

	_, err := conn.Do("LREM", key, count, value)
	if err != nil {
		return err
	}

	return nil

}

//BRPop
func BRPop(key string, timeout int) (interface{}, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	var value1 string
	var value2 string

	reply, err := redis.Values(conn.Do("BRPOP", key, timeout))
	if err != nil {
		return value2, err
	}

	if _, err := redis.Scan(reply, &value1, &value2); err != nil {
		return value2, err
	}

	return value2, nil
}

//ZRangeByScore
func ZRangeByScore(key string, startScore, endScore float64) ([]string, error) {
	conn := RedisConn.Get()
	defer conn.Close()
	var resultList []string

	values, err := redis.Values(conn.Do("ZRANGEBYSCORE", key, startScore, endScore))
	if err != nil {
		return resultList, err
	}

	if err := redis.ScanSlice(values, &resultList); err != nil {
		return resultList, err
	}

	return resultList, nil
}