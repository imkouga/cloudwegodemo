package config

import (
	"cloudwegodemo/cmd/cloudwegodemo/internal/server"
	"cloudwegodemo/pkg/configor"
	"cloudwegodemo/pkg/database/mysql"

	"github.com/google/wire"
)

var (
	globalC *Config
)

var ConfigProvider = wire.NewSet(NewGlobalConfigor, GetHTTPServerOptionFn, GetMySQLOptionFn)

func NewConfigor(c *Config, opt *configor.Option) (*configor.Configor, error) {
	return configor.New(globalC, opt)
}

func NewGlobalConfigor(opt *configor.Option) *configor.Configor {
	globalC = &Config{}
	k, _ := NewConfigor(globalC, opt)
	return k
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
	return nil, nil
}
