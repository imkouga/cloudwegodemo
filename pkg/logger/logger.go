package logger

import (
	"github.com/cloudwego/hertz/pkg/common/hlog"
	hertzzap "github.com/hertz-contrib/logger/zap"
	"go.uber.org/zap"
)

type Option interface {
	apply(cfg *logConfig)
}

type option func(cfg *logConfig)

type logConfig struct {
	logger hlog.FullLogger
	lv     hlog.Level

	serviceName, serviceVersion, serviceHostname string
}

func (fn option) apply(cfg *logConfig) {
	fn(cfg)
}

func Init(opts ...Option) error {

	cfg := &logConfig{
		logger: nil,
		lv:     hlog.LevelInfo,
	}
	for _, opt := range opts {
		opt.apply(cfg)
	}

	if cfg.logger == nil {
		cfg.logger = hertzzap.NewLogger(
			hertzzap.WithZapOptions(zap.Fields(
				zap.String("service.name", cfg.serviceName),
				zap.String("service.version", cfg.serviceVersion),
				zap.String("service.hostname", cfg.serviceHostname),
			)),
		)
	}

	hlog.SetLogger(cfg.logger)
	hlog.SetLevel(cfg.lv)

	return nil
}

func WithLogger(logger hlog.FullLogger) Option {
	return option(func(cfg *logConfig) {
		cfg.logger = logger
	})
}

func WithLevel(lv hlog.Level) Option {
	return option(func(cfg *logConfig) {
		cfg.lv = lv
	})
}

func WithServiceAbstract(serviceName, serviceVersion, serviceHostname string) Option {
	return option(func(cfg *logConfig) {
		cfg.serviceName = serviceName
		cfg.serviceVersion = serviceVersion
		cfg.serviceHostname = serviceHostname
	})
}
