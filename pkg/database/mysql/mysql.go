package mysql

import (
	"context"
	"database/sql"
	"errors"

	"cloudwegodemo/pkg/configor"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

type (
	GetOptionFn  func() (*Option, error)
	contextTxKey struct{}
	MySQL        struct {
		db *gorm.DB

		opt   *Option
		optFn GetOptionFn
	}
	Option struct {
		Dsn             string               `json:"dsn,omitempty" yaml:"dsn,omitempty"`
		MaxOpenConns    int64                `json:"max_open_conns,omitempty" yaml:"max_open_conns,omitempty"`
		MaxIdleConns    int64                `json:"max_idle_conns,omitempty" yaml:"max_idle_conns,omitempty"`
		ConnMaxIdleTime *durationpb.Duration `json:"conn_max_idle_time,omitempty" yaml:"conn_max_idle_time,omitempty"`
		ConnMaxLifeTime *durationpb.Duration `json:"conn_max_life_time,omitempty" yaml:"conn_max_life_time,omitempty"`
	}
)

func NewMySQLPool(cr configor.Configor, optFn GetOptionFn) (*MySQL, error) {

	if nil == optFn {
		return nil, errors.New("获取MySQL配置信息失败")
	}

	opt, err := optFn()
	if nil != err {
		return nil, err
	}
	db, err := NewGorm(opt)
	if nil != err {
		return nil, err
	}
	mysqlDB := &MySQL{db: db, opt: opt, optFn: optFn}

	if nil != cr {
		if err := cr.RegisterReload("MySQL", mysqlDB.reload); nil != err {
			return nil, err
		}
	}

	return mysqlDB, nil
}

func NewGorm(opt *Option) (*gorm.DB, error) {

	var (
		driver gorm.Dialector
		opts   = &gorm.Config{
			Logger:          gormLogger.Discard,
			CreateBatchSize: 1000,
		}
	)

	// only support MySQL.
	driver = mysql.New(mysql.Config{
		DSN:                      opt.Dsn,
		DisableDatetimePrecision: true,
		DontSupportRenameIndex:   true,
	})
	hlog.Infof("mysql dsn is %s", opt.Dsn)

	conn, err := gorm.Open(driver, opts)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (d *MySQL) Init() error {

	if nil == d || nil == d.db || nil == d.opt {
		return errors.New("MySQL连接池没有实例化")
	}
	return d.init()
}

func (d *MySQL) init() error {

	db, err := d.db.DB()
	if nil != err {
		return err
	}
	return initDB(db, d.opt)
}

func (d *MySQL) reload() error {

	if nil == d || nil == d.optFn {
		return nil
	}

	opt, err := d.optFn()
	if nil != err {
		return err
	}
	gdb, err := NewGorm(opt)
	if nil != err {
		hlog.Error(err)
		return err
	}
	db, err := gdb.DB()
	if nil != err {
		hlog.Error(err)
		return err
	}
	if err := initDB(db, opt); nil != err {
		hlog.Error(err)
		return err
	}

	d.db, d.opt = gdb, opt
	return nil
}

// ExecTx gorm Transaction
func (d *MySQL) ExecTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return d.DB(ctx).WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, contextTxKey{}, tx)
		return fn(ctx)
	})
}

// DB 根据此方法来判断当前的 db 是不是使用 事务的 DB
func (d *MySQL) DB(ctx context.Context) *gorm.DB {

	tx, ok := ctx.Value(contextTxKey{}).(*gorm.DB)
	if !ok {
		return d.db
	}
	return tx
}

func initDB(db *sql.DB, opt *Option) error {

	if nil == db || nil == opt {
		return nil
	}

	db.SetConnMaxIdleTime(opt.ConnMaxIdleTime.AsDuration())
	db.SetConnMaxLifetime(opt.ConnMaxLifeTime.AsDuration())
	db.SetMaxIdleConns(int(opt.MaxIdleConns))
	db.SetMaxOpenConns(int(opt.MaxOpenConns))
	return nil
}
