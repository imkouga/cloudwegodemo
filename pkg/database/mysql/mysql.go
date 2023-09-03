package mysql

import (
	"context"
	"errors"

	"cloudwegodemo/pkg/configor"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

type GetOptionFn func() (*Option, error)

type contextTxKey struct{}

type MySQL struct {
	db *gorm.DB

	opt   *Option
	optFn GetOptionFn
}

func NewMySQLPool(cr *configor.Configor, optFn GetOptionFn) (*MySQL, error) {

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
		cr.RegisterReload(mysqlDB.reload)
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
	db.SetConnMaxIdleTime(d.opt.GetConnMaxIdleTime().AsDuration())
	db.SetConnMaxLifetime(d.opt.GetConnMaxLifeTime().AsDuration())
	db.SetMaxIdleConns(int(d.opt.GetMaxIdleConns()))
	db.SetMaxOpenConns(int(d.opt.GetMaxOpenConns()))
	return nil
}

func (d *MySQL) reload() {

	if nil == d || nil == d.optFn {
		return
	}

	opt, err := d.optFn()
	if nil != err {
		return
	}
	db, err := NewGorm(opt)
	if nil != err {
		return
	}
	d.db, d.opt = db, opt

	_ = d.init()
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
