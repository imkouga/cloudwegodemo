package repo

import (
	"cloudwegodemo/pkg"
	"cloudwegodemo/pkg/database/mysql"

	"github.com/google/wire"
)

var (
	RepoProvider = wire.NewSet(NewRepoSet)
)

type (
	RepoSet struct {
		pkg.Base
		mysqlDB *mysql.MySQL
	}
)

func NewRepoSet(mysqlDB *mysql.MySQL) (*RepoSet, error) {
	return &RepoSet{mysqlDB: mysqlDB}, nil
}

func (r *RepoSet) Init() error {

	if err := r.mysqlDB.Init(); nil != err {
		return err
	}

	return nil
}
