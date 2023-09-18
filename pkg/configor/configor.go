package configor

import (
	"fmt"
	"sync"
	"time"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/jinzhu/configor"
)

type (
	ReloadFn func() error
	Loader   interface {
		Reload() error
	}

	Configor interface {
		RegisterReload(name string, reload ReloadFn) error
	}
	simpleConfigor struct {
		cr *configor.Configor

		sync.Mutex
		reloadBroadcasts map[string]ReloadFn
	}

	Option struct {
		Paths []string
	}
)

func New(c interface{}, opt *Option) (Configor, error) {

	cr := &simpleConfigor{
		reloadBroadcasts: make(map[string]ReloadFn, 10),
	}
	cr.cr = configor.New(&configor.Config{
		AutoReload:           true,
		AutoReloadInterval:   time.Second * 5,
		AutoReloadCallback:   cr.reloadCallback,
		ErrorOnUnmatchedKeys: true,
	})

	if err := cr.cr.Load(c, opt.Paths...); nil != err {
		return nil, err
	}
	return cr, nil
}

func (cr *simpleConfigor) RegisterReload(name string, reload ReloadFn) error {
	if nil != cr {
		cr.Lock()
		defer cr.Unlock()
		if _, exist := cr.reloadBroadcasts[name]; exist {
			return fmt.Errorf("配置回调器[%s]已存在，不允许重复注册", name)
		}
		cr.reloadBroadcasts[name] = reload
	}
	return nil
}

func (cr *simpleConfigor) reloadCallback(config interface{}) {
	if nil != cr {
		cr.Lock()
		defer cr.Unlock()

		for name, reload := range cr.reloadBroadcasts {
			go func(name string, reload ReloadFn) error {
				if err := reload(); nil != err {
					hlog.Errorf("配置变更，重置回调器[%s]执行异常，异常原因是%s", name, err.Error())
					return err
				}
				hlog.Infof("配置变更，重置回调器[%s]执行完成", name)
				return nil
			}(name, reload)
		}
	}
}
