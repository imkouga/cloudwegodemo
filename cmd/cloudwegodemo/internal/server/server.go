package server

import (
	"errors"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/google/wire"
)

var (
	ErrHttpServerOptNotProvider = errors.New("HTTPServer option not provided")

	ServerProvider = wire.NewSet(NewHTTPServer)
)

type (
	GetHTTPServerOption func() (*HTTPServer, error)
	GetRPCServerOption  func() (*RPCServer, error)
)

func NewHTTPServer(optFn GetHTTPServerOption) (*server.Hertz, error) {

	if nil == optFn {
		return nil, ErrHttpServerOptNotProvider
	}
	opt, err := optFn()
	if nil != err {
		return nil, err
	}

	h := server.New(server.WithHostPorts(opt.GetAddr()))

	return h, nil
}

func NewRPCServer() {

}

func Run() error {
	return nil
}
