package redis

import (
	"cloudwegodemo/pkg/configor"
	"context"
	"errors"

	durationpb "google.golang.org/protobuf/types/known/durationpb"

	"github.com/redis/go-redis/v9"
)

type (
	GetOptionFn func() (*Option, error)
	Option      struct {
		Addr            string               `json:"addr,omitempty" yaml:"addr,omitempty"`
		Password        string               `json:"password,omitempty" yaml:"password,omitempty"`
		DB              int64                `json:"db,omitempty" yaml:"db,omitempty"`
		MaxIdleConns    int64                `json:"max_idle_conns,omitempty" yaml:"max_idle_conns,omitempty"`
		MinIdleConns    int64                `json:"min_idle_conns,omitempty" yaml:"min_idle_conns,omitempty"`
		ConnMaxIdleTime *durationpb.Duration `json:"conn_max_idle_time,omitempty" yaml:"conn_max_idle_time,omitempty"`
		ConnMaxLifeTime *durationpb.Duration `json:"conn_max_life_time,omitempty" yaml:"conn_max_life_time,omitempty"`
	}
	Redis struct {
		optFn GetOptionFn

		rdb *redis.Client
	}
)

func New(opt *Option) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:            opt.Addr,
		Password:        opt.Password,
		DB:              int(opt.DB),
		MaxIdleConns:    int(opt.MaxIdleConns),
		MinIdleConns:    int(opt.MinIdleConns),
		ConnMaxIdleTime: opt.ConnMaxIdleTime.AsDuration(),
		ConnMaxLifetime: opt.ConnMaxLifeTime.AsDuration(),
	})
	return rdb, nil
}

func NewRedisPool(cr configor.Configor, optFn GetOptionFn) (*Redis, error) {

	if nil == cr || nil == optFn {
		return nil, errors.New("实例化redis连接池异常，没有传入配置")
	}

	opt, err := optFn()
	if nil != err {
		return nil, err
	}

	r := &Redis{
		optFn: optFn,
	}
	rdb, err := New(opt)
	if nil != err {
		return nil, err
	}
	r.rdb = rdb

	if err := cr.RegisterReload("redis", r.reload); nil != err {
		return nil, err
	}

	return r, nil
}

func (r *Redis) reload() error {

	if nil == r.optFn {
		return errors.New("系统异常，重置Redis连接池失败，没有提供配置加载器")
	}
	opt, err := r.optFn()
	if nil != err {
		return err
	}

	rdb, err := New(opt)
	if nil != err {
		return err
	}
	if err := initRDB(rdb, opt); nil != err {
		return err
	}
	r.rdb = rdb

	return nil
}

func (r *Redis) Init() error {

	if nil == r || nil == r.rdb || nil == r.optFn {
		return errors.New("系统异常，Redis连接池没有实例化")
	}

	opt, err := r.optFn()
	if nil != err {
		return err
	}

	return initRDB(r.rdb, opt)
}

func (r *Redis) RDB() *redis.Client {
	return r.rdb
}

func initRDB(rdb *redis.Client, opt *Option) error {

	if nil == rdb || nil == opt {
		return nil
	}
	return rdb.Ping(context.Background()).Err()
}
