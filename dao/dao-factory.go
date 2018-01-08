package dao

import (
	"errors"
	"fmt"
	"gopkg.in/redis.v5"
	"github.com/BurntSushi/toml"
)

type DBType int

type DBConfig struct {
	Url      string
	Port     string
	Password string
	Database int
}

const (
	// RedisDAO is used for Redis implementation of TaskDAO
	RedisDAO DBType = iota
	// MockDAO is used for mocked implementation of TaskDAO
	MockDAO
)

var redisLocalConfig = DBConfig{
	Url:      "localhost",
	Password: "",
	Database: 0,
	Port:     "6379",
}

// GetDAO returns a TaskDAO according to type and params
func GetDAO(daoType DBType, dbConfigFile string) (TaskDAO, error) {
	switch daoType {
	case RedisDAO:
		config := getConfig(dbConfigFile)
		redisCli := initRedis(config)
		return NewTaskDAORedis(redisCli), nil
	case MockDAO:
		return NewTaskDAOMock(), nil
	default:
		return nil, errors.New("unknown DAO type")
	}
}

// Initialize Redis database
func initRedis(dbConfig DBConfig) *redis.Client {
	fmt.Println("redis connexion " + dbConfig.Url)

	// Connection to the Redis database
	redisCli := redis.NewClient(&redis.Options{
		Addr:     dbConfig.Url + ":" + dbConfig.Port,
		Password: dbConfig.Password,
		DB:       int(RedisDAO),
	})

	// Verification of connection
	ok, err := redisCli.Ping().Result()
	if err != nil {
		fmt.Println("redis connexion KO :", err)
		panic(err)
	} else {
		fmt.Println("redis connexion OK :", ok)
	}

	return redisCli
}

func getConfig(dbConfigFile string) DBConfig {
	var config DBConfig
	if dbConfigFile == "" {
		config = redisLocalConfig
	} else {
		if _, err := toml.DecodeFile(dbConfigFile, &config); err != nil {
			fmt.Println("configuration file error :", err)
			panic(err)
		}
	}
	return config
}