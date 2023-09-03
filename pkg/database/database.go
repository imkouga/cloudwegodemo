package database

import (
	"cloudwegodemo/pkg/database/mysql"

	"github.com/google/wire"
)

var (
	DatabaseProvider = wire.NewSet(mysql.NewMySQLPool)
)
