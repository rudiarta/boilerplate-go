package libredis

import (
	"crypto/tls"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

// Conn base struct
type Conn struct {
	Client *redis.Client
}

type Iterator interface {
	Err() error
	Next() bool
	Val() string
}

// ConnectRedis init redis
func ConnectRedis(redisHost, redisPort, redisPassword, redisDB, redisTLS string) (Client, error) {
	cl, err := GetRedis(
		redisHost,
		redisPort,
		redisPassword,
		redisDB,
		redisTLS,
	)
	if err != nil {
		return nil, err
	}
	return Conn{
		Client: cl,
	}, nil
}

// Del key
func (r Conn) Del(key string) (int64, error) {
	return r.Client.Del(key).Result()
}

// HGet Key
func (r Conn) HGet(key, field string) (string, error) {
	return r.Client.HGet(key, field).Result()
}

// Get key
func (r Conn) Get(key string) (string, error) {
	return r.Client.Get(key).Result()
}

func (r Conn) HSet(key, field, value string) (bool, error) {
	return r.Client.HSet(key, field, value).Result()
}

// Set redis
func (r Conn) Set(key string, value string, exp time.Duration) (string, error) {
	return r.Client.Set(key, value, exp).Result()
}

// Set redis once
func (r Conn) SetOnce(key string, value string, exp time.Duration) (bool, error) {
	return r.Client.SetNX(key, value, exp).Result()
}

// Ping result
func (r Conn) Ping() (string, error) {
	return r.Client.Ping().Result()
}

// Expire result
func (r Conn) Expire(key string, exp time.Duration) (bool, error) {
	return r.Client.Expire(key, exp).Result()
}

// Scan Iterator
func (r Conn) ScanIterator(key string) Iterator {
	return r.Client.Scan(0, key, 0).Iterator()
}

// Keys get multi key
func (r Conn) Keys(key string) ([]string, error) {
	return r.Client.Keys(key).Result()
}

// Keys get multi key
func (r Conn) Incr(key string) (int64, error) {
	return r.Client.Incr(key).Result()
}

// GetRedis function
func GetRedis(redisHost, redisPort, redisPassword, redisDB, redisTLS string) (*redis.Client, error) {
	tlsSecured, err := strconv.ParseBool(redisTLS)
	if err != nil {
		return nil, err
	}

	var conf *tls.Config

	// force checking for unsecured aws redis
	if tlsSecured {
		conf = &tls.Config{
			InsecureSkipVerify: tlsSecured,
		}
	} else {
		conf = nil
	}

	useDB, _ := strconv.Atoi(redisDB)
	cl := redis.NewClient(&redis.Options{
		Addr:      fmt.Sprintf("%v:%v", redisHost, redisPort),
		Password:  redisPassword,
		DB:        useDB, // use default DB
		TLSConfig: conf,
	})

	return cl, nil
}
