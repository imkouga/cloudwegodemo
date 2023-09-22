package server

import (
	"errors"

	"cloudwegodemo/pkg/contrib/registry"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/common/utils"
)

var (
	ErrHttpServerOptNotProvider = errors.New("HTTPServer option not provided")
)

type (
	GetHTTPServerOption func() (*HTTPOption, error)
)

func NewHTTPServer(optFn GetHTTPServerOption, r registry.Registry) (*server.Hertz, error) {

	if nil == optFn {
		return nil, ErrHttpServerOptNotProvider
	}
	opt, err := optFn()
	if nil != err {
		return nil, err
	}

	serverOpts := make([]config.Option, 0, 10)
	serverOpts = append(serverOpts, server.WithHostPorts(opt.GetAddr()))

	if nil != r {
		info := r.GetInfo()
		if nil == info {
			return nil, errors.New("registry info没有实例化")
		}
		info.Addr = utils.NewNetAddr("tcp", opt.GetAddr())

		rr := r.GetRegistry()
		if nil == rr {
			return nil, errors.New("registry 没有实例化")
		}

		serverOpts = append(serverOpts, server.WithRegistry(r, info))
	}

	h := server.New(serverOpts...)

	return h, nil
}
