package libredis

import "time"

// Client interface contract
type Client interface {
	Get(key string) (string, error)
	HGet(key, field string) (string, error)
	Del(key string) (int64, error)
	Set(key string, value string, duration time.Duration) (string, error)
	SetOnce(key string, value string, duration time.Duration) (bool, error)
	HSet(key, field, value string) (bool, error)
	Ping() (string, error)
	Expire(key string, exp time.Duration) (bool, error)
	Keys(key string) ([]string, error)
	Incr(key string) (int64, error)
	ScanIterator(key string) Iterator
}
