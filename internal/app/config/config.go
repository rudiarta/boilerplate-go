package config

import (
	"os"

	"github.com/rudiarta/boilerplate-go/pkg/libelastic"
	"github.com/rudiarta/boilerplate-go/pkg/libgorm"
	"github.com/rudiarta/boilerplate-go/pkg/libredis"
	"gorm.io/gorm"
)

type environmentValue struct {
	Port        string
	ServiceName string
}

type configCtx struct {
}

type ConfigCtx interface {
	GormDbConnection() (*gorm.DB, error)
	ElasticClientConnection() (libelastic.ElasticClient, error)
	RedisClientConnection() (libredis.Client, error)
	GetEnvironmentValue() (*environmentValue, error)
}

func NewConfigCtx() ConfigCtx {
	return &configCtx{}
}

func (c *configCtx) GormDbConnection() (*gorm.DB, error) {
	var username string = os.Getenv("DB_USERNAME")
	var password string = os.Getenv("DB_PASSWORD")
	var host string = os.Getenv("DB_HOST")
	var port string = os.Getenv("DB_PORT")
	var dbname string = os.Getenv("DB_NAME")
	return libgorm.InitGorm(username, password, host, port, dbname)
}

func (c *configCtx) GetEnvironmentValue() (*environmentValue, error) {
	return &environmentValue{
		Port:        os.Getenv("APP_PORT"),
		ServiceName: os.Getenv("SERVICE_NAME"),
	}, nil
}

func (c *configCtx) ElasticClientConnection() (libelastic.ElasticClient, error) {
	var elasticUrl string = os.Getenv("ELASTIC_URL")
	var elasticUsername string = os.Getenv("ELASTIC_USERNAME")
	var elasticPassword string = os.Getenv("ELASTIC_PASSWORD")
	return libelastic.NewClient(elasticUrl, elasticUsername, elasticPassword)
}

func (c *configCtx) RedisClientConnection() (libredis.Client, error) {
	var redisHost string = os.Getenv("REDIS_HOST")
	var redisPort string = os.Getenv("REDIS_PORT")
	var redisPassword string = os.Getenv("REDIS_PASSWORD")
	var redisDB string = os.Getenv("REDIS_DB")
	var redisTLS string = os.Getenv("REDIS_TLS")
	return libredis.ConnectRedis(redisHost, redisPort, redisPassword, redisDB, redisTLS)
}
