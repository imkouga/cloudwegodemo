package bootstrap

import (
	"errors"

	"cloudwegodemo/pkg"
	"cloudwegodemo/pkg/configor"
	"cloudwegodemo/pkg/contrib/registry"

	"github.com/cloudwego/hertz/pkg/app/server"
)

type (
	RouterRegister func(*server.Hertz)

	Option interface {
		apply(app *APP)
	}

	option func(app *APP)

	APP struct {
		name    string
		version string

		cr configor.Configor

		rr          RouterRegister
		hertzServer *server.Hertz
		// kitexServer server.Kitex

		registry registry.Registry

		bases []pkg.Base
	}
)

func (app *APP) Init() error {

	if nil == app {
		return errors.New("app没有实例化")
	}

	if nil != app.registry {
		app.registry.SetServiceName(app.name)
	}

	for _, base := range app.bases {
		if err := base.Init(); nil != err {
			return err
		}
	}

	return nil
}

func (app *APP) Run() error {

	if nil == app {
		return errors.New("app没有实例化")
	}

	if nil != app.hertzServer {
		app.hertzServer.Spin()
	}

	return nil
}

func (fn option) apply(app *APP) {
	fn(app)
}

func New(opts ...Option) *APP {

	app := &APP{}
	for _, opt := range opts {
		opt.apply(app)
	}

	return app
}

func WithName(name string) Option {
	return option(func(app *APP) {
		app.name = name
	})
}

func WithVersion(version string) Option {
	return option(func(app *APP) {
		app.version = version
	})
}

func WithConfigor(cr configor.Configor) Option {
	return option(func(app *APP) {
		app.cr = cr
	})
}

func WithRegistry(r registry.Registry) Option {
	return option(func(app *APP) {
		app.registry = r
	})
}

func WithRouteRegister(rr RouterRegister) Option {
	return option(func(app *APP) {
		app.rr = rr
	})
}

func WithHertz(h *server.Hertz) Option {
	return option(func(app *APP) {
		app.hertzServer = h
	})
}

func WithBase(bases ...pkg.Base) Option {
	return option(func(app *APP) {
		app.bases = append(app.bases, bases...)
	})
}
