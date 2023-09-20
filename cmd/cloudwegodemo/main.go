package main

import (
	"flag"
	"os"

	"cloudwegodemo/cmd/cloudwegodemo/internal/biz"
	"cloudwegodemo/cmd/cloudwegodemo/internal/router"
	"cloudwegodemo/pkg/configor"
	"cloudwegodemo/pkg/logger"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

var (
	Version     = "v1.0.0"
	ServerName  = "cloudwegodemo"
	LogLevel    int
	HostName, _ = os.Hostname()
	cfgPath     string
)

type APP struct {
	httpServer *server.Hertz
	bizSet     *biz.BaseBiz
}

func init() {
	flag.StringVar(&cfgPath, "conf", "../../etc", "config path, eg: -conf config.yaml")
	flag.IntVar(&LogLevel, "log", 2, "log level, eg: -log 2; use 0 1 2 3 4 5 6")
	flag.Parse()
}

func initLogger() error {
	return logger.Init(
		logger.WithLevel(hlog.Level(LogLevel)),
		logger.WithServiceAbstract(ServerName, Version, HostName),
	)
}

// initInspect 初始化巡检模块 不返回异常，永远返回nil
func initInspect() error {
	return nil
}

func initAPP(app *APP) error {

	if err := app.bizSet.Init(); nil != err {
		return err
	}
	return nil
}

func runAPP(app *APP) error {
	router.GeneratedRegister(app.httpServer)
	app.httpServer.Spin()
	return nil
}

func newApp(c configor.Configor, bizSet *biz.BaseBiz, h *server.Hertz) (*APP, func(), error) {
	return &APP{
		httpServer: h,
		bizSet:     bizSet,
	}, func() {}, nil
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

	app, cancel, err := wireApp(&configor.Option{
		Paths: []string{cfgPath},
	})
	if nil != err {
		hlog.Fatal(err)
	}
	defer cancel()

	if err := initAPP(app); nil != err {
		hlog.Fatal(err)
	}
	if err := runAPP(app); nil != err {
		hlog.Fatal(err)
	}
}
