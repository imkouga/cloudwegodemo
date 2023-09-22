package registry

import (
	"errors"

	"github.com/cloudwego/hertz/pkg/app/server/registry"
	"github.com/cloudwego/hertz/pkg/common/hlog"

	"cloudwegodemo/pkg"
	"cloudwegodemo/pkg/configor"

	"github.com/hertz-contrib/registry/etcd"
)

type (
	Registry interface {
		registry.Registry

		SetServiceName(name string)
		GetRegistry() registry.Registry
		GetInfo() *registry.Info
	}
	simpleRegistry struct {
		r     registry.Registry
		info  *registry.Info
		optFn GetRegistryOption
		cr    configor.Configor
	}
)

func NewRegistry(serviceName pkg.ServiceName, cr configor.Configor, optFn GetRegistryOption) (Registry, error) {

	if nil == optFn {
		return nil, errors.New("实例化Registry异常，没有传入配置项")
	}

	opt, err := optFn()
	if nil != err {
		return nil, err
	}
	if !opt.Enabled {
		hlog.Info("服务注册与发现没有启用")
		return nil, nil
	}

	r, err := etcd.NewEtcdRegistry(opt.Endpoints)
	if nil != err {
		return nil, err
	}
	info := &registry.Info{
		ServiceName: string(serviceName),
		Addr:        nil,
		Weight:      10,
		Tags:        nil,
	}

	return &simpleRegistry{
		r:     r,
		info:  info,
		optFn: optFn,
		cr:    cr,
	}, nil
}

func (sr *simpleRegistry) SetServiceName(name string) {
	if nil != sr && nil != sr.info {
		sr.info.ServiceName = name
	}
}

func (sr *simpleRegistry) GetRegistry() registry.Registry {
	if nil != sr {
		return sr.r
	}
	return nil
}

func (sr *simpleRegistry) GetInfo() *registry.Info {
	if nil != sr {
		return sr.info
	}
	return nil
}

func (sr *simpleRegistry) Register(info *registry.Info) error {
	if nil != sr && nil != sr.r {
		return sr.r.Register(info)
	}
	return nil
}

func (sr *simpleRegistry) Deregister(info *registry.Info) error {
	if nil != sr && nil != sr.r {
		return sr.r.Deregister(info)
	}
	return nil
}
