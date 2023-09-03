package biz

import (
	"cloudwegodemo/pkg"

	"github.com/google/wire"
)

var (
	BizProvider = wire.NewSet(NewBizSet)
)

type (
	BizSet struct {
		pkg.Base
	}
)

func NewBizSet() (*BizSet, error) {
	return &BizSet{}, nil
}

func (b *BizSet) Init() error {
	return nil
}
