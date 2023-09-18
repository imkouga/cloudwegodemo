package config

import (
	"cloudwegodemo/cmd/cloudwegodemo/internal/server"
	"cloudwegodemo/pkg/configor"
	"cloudwegodemo/pkg/database/mysql"
	redis "cloudwegodemo/pkg/database/redis"

	"github.com/google/wire"
)

var (
	globalC *Config
)
var ConfigProvider = wire.NewSet(NewGlobalConfigor, GetHTTPServerOptionFn, GetMySQLOptionFn, GetRedisOptionFn)

type (
	Config struct {
		Server   *Server   `json:"server,omitempty" yaml:"server,omitempty"`
		Database *DataBase `json:"database,omitempty" yaml:"database,omitempty"`
	}
	Server struct {
		Http *server.HTTPServer `json:"http,omitempty" yaml:"http,omitempty"`
		Rpc  *server.RPCServer  `json:"rpc,omitempty" yaml:"rpc,omitempty"`
	}
	DataBase struct {
		Mysql *mysql.Option `json:"mysql,omitempty" yaml:"mysql,omitempty"`
		Redis *redis.Option `json:"redis,omitempty" yaml:"redis,omitempty"`
	}
)

func NewConfigor(c *Config, opt *configor.Option) (configor.Configor, error) {
	return configor.New(c, opt)
}

func NewGlobalConfigor(opt *configor.Option) (configor.Configor, error) {

	c := &Config{}

	cr, err := NewConfigor(c, opt)
	if nil != err {
		return nil, err
	}
	globalC = c

	return cr, nil
}

func GetHTTPServerOptionFn() server.GetHTTPServerOption {
	return GetHTTPServerOption
}

func GetHTTPServerOption() (*server.HTTPServer, error) {
	return globalC.Server.Http, nil
}

func GetMySQLOptionFn() mysql.GetOptionFn {
	return GetMySQLOption
}

func GetMySQLOption() (*mysql.Option, error) {
	return globalC.Database.Mysql, nil
}

func GetRedisOptionFn() redis.GetOptionFn {
	return GetRedisOption
}

func GetRedisOption() (*redis.Option, error) {
	return globalC.Database.Redis, nil
}
