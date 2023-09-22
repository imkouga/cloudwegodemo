package main

import (
	"flag"
	"os"

	"cloudwegodemo/cmd/cloudwegodemo/internal/biz"
	"cloudwegodemo/cmd/cloudwegodemo/internal/router"
	"cloudwegodemo/pkg"

	"cloudwegodemo/pkg/bootstrap"
	"cloudwegodemo/pkg/configor"
	"cloudwegodemo/pkg/contrib/registry"
	"cloudwegodemo/pkg/logger"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

var (
	Version     = "v1.0.0"
	ServiceName = "cloudwegodemo"
	LogLevel    int
	HostName, _ = os.Hostname()
	cfgPath     string
)

func init() {
	flag.StringVar(&cfgPath, "conf", "../../etc", "config path, eg: -conf config.yaml")
	flag.IntVar(&LogLevel, "log", 2, "log level, eg: -log 2; use 0 1 2 3 4 5 6")
	flag.Parse()
}

func initLogger() error {
	return logger.Init(
		logger.WithLevel(hlog.Level(LogLevel)),
		logger.WithServiceAbstract(ServiceName, Version, HostName),
	)
}

// initInspect 初始化巡检模块 不返回异常，永远返回nil
func initInspect() error {
	return nil
}

func newApp(c configor.Configor, h *server.Hertz, r registry.Registry, bizSet *biz.BaseBiz) (*bootstrap.APP, func(), error) {

	app := bootstrap.New(
		bootstrap.WithName(ServiceName),
		bootstrap.WithVersion(Version),

		bootstrap.WithConfigor(c),

		bootstrap.WithRouteRegister(router.GeneratedRegister),
		bootstrap.WithHertz(h),

		bootstrap.WithRegistry(r),

		bootstrap.WithBase(bizSet),
	)

	return app, func() {}, nil
}

func main() {

	// 初始化顺序强约束
	inits := []func() error{
		initLogger,
		initInspect,
	}
	for _, init := range inits {
		if err := init(); nil != err {
			hlog.Fatal(err)
		}
	}

	app, cancel, err := wireApp(pkg.ServiceName(ServiceName), &configor.Option{
		Paths: []string{cfgPath},
	})
	if nil != err {
		hlog.Fatal(err)
	}
	defer cancel()

	if err := app.Init(); nil != err {
		hlog.Fatal(err)
	}
	if err := app.Run(); nil != err {
		hlog.Fatal(err)
	}
}
