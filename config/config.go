package config

import (
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

var once sync.Once
var realPath string
var Conf *Config

const (
	SuccessReplyCode   = 0
	FailReplyCode      = 1
	SuccessReplyMsg    = "success"
	RedisBaseValidTime = 86400
	RedisPrefix        = "goTestProject_"
)

type Config struct {
	Common Common
}

func init() {
	Init()
}

func getCurrentDir() string {
	_, fileName, _, _ := runtime.Caller(1)
	aPath := strings.Split(fileName, "/")
	dir := strings.Join(aPath[0:len(aPath)-1], "/")
	return dir
}

func Init() {
	once.Do(func() {
		env := GetMode()
		//realPath, _ := filepath.Abs("./")
		realPath := getCurrentDir()
		configFilePath := realPath + "/" + env + "/"
		viper.SetConfigType("toml")
		viper.AddConfigPath(configFilePath)
		viper.SetConfigName("/common")
		err := viper.MergeInConfig()
		if err != nil {
			panic(err)
		}
		Conf = new(Config)
		viper.Unmarshal(&Conf.Common)
	})
}

func GetMode() string {
	env := os.Getenv("RUN_MODE")
	if env == "" {
		env = "dev"
	}
	return env
}

func GetGinRunMode() string {
	env := GetMode()
	//gin have debug,test,release mode
	if env == "dev" {
		return "debug"
	}
	if env == "test" {
		return "debug"
	}
	if env == "prod" {
		return "release"
	}
	return "release"
}

type CommonRedis struct {
	RedisAddress  string `mapstructure:"redisAddress"`
	RedisPassword string `mapstructure:"redisPassword"`
	Db            int    `mapstructure:"db"`
}

type CommonMysql struct {
	MysqlAddress string `mapstructure:"mysqlAddress"`
}

type Common struct {
	CommonRedis CommonRedis `mapstructure:"common-redis"`
	CommonMysql CommonMysql `mapstructure:"common-mysql"`
}
