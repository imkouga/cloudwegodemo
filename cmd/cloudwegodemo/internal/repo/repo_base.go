package repo

import (
	"cloudwegodemo/pkg"
	"cloudwegodemo/pkg/database/mysql"
	"cloudwegodemo/pkg/database/redis"

	"github.com/google/wire"
)

var (
	RepoProvider = wire.NewSet(NewBaseRepo)
)

type (
	BaseRepo interface {
		pkg.Base
	}
	baseRepo struct {
		mysqlDB *mysql.MySQL
		rdb     *redis.Redis
	}
)

func NewBaseRepo(mysqlDB *mysql.MySQL, rdb *redis.Redis) (BaseRepo, error) {

	r := &baseRepo{mysqlDB: mysqlDB, rdb: rdb}
	if err := r.Init(); nil != err {
		return nil, err
	}
	return r, nil
}

func (r *baseRepo) Init() error {

	if err := r.mysqlDB.Init(); nil != err {
		return err
	}
	if err := r.rdb.Init(); nil != err {
		return err
	}

	return nil
}
