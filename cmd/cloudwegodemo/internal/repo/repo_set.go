package repo

import (
	"cloudwegodemo/pkg"
	"cloudwegodemo/pkg/database/mysql"
	"cloudwegodemo/pkg/database/redis"

	"github.com/google/wire"
)

var (
	RepoProvider = wire.NewSet(NewRepoSet)
)

type (
	RepoSet struct {
		pkg.Base

		mysqlDB *mysql.MySQL
		rdb     *redis.Redis
	}
)

func NewRepoSet(mysqlDB *mysql.MySQL, rdb *redis.Redis) (*RepoSet, error) {
	return &RepoSet{mysqlDB: mysqlDB, rdb: rdb}, nil
}

func (r *RepoSet) Init() error {

	if err := r.mysqlDB.Init(); nil != err {
		return err
	}
	if err := r.rdb.Init(); nil != err {
		return err
	}

	return nil
}
