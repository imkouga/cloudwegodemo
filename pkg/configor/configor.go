package configor

import (
	"time"

	"github.com/jinzhu/configor"
)

type (
	ReloadFn func()
	Loader   interface {
		Reload() error
	}
	Configor struct {
		cr               *configor.Configor
		reloadBroadcasts []ReloadFn
	}
	Option struct {
		Paths []string
	}
)

func New(c interface{}, opt *Option) (*Configor, error) {

	cr := &Configor{}
	cr.cr = configor.New(&configor.Config{
		AutoReload:         true,
		AutoReloadInterval: time.Minute,
		AutoReloadCallback: cr.reloadCallback,
	})

	return cr, cr.cr.Load(&c, opt.Paths...)
}

func (cr *Configor) RegisterReload(reload ReloadFn) {
	if nil != cr {
		cr.reloadBroadcasts = append(cr.reloadBroadcasts, reload)
	}
}

func (cr *Configor) reloadCallback(config interface{}) {
	if nil != cr {
		for _, reload := range cr.reloadBroadcasts {
			go reload()
		}
	}
}
