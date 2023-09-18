package database

import (
	"cloudwegodemo/pkg/database/mysql"
	"cloudwegodemo/pkg/database/redis"

	"github.com/google/wire"
)

var (
	DatabaseProvider = wire.NewSet(mysql.NewMySQLPool, redis.NewRedisPool)
)
