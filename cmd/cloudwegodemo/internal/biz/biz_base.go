package biz

import (
	"cloudwegodemo/cmd/cloudwegodemo/internal/repo"
	"cloudwegodemo/pkg"

	"github.com/google/wire"
)

var (
	BizProvider = wire.NewSet(NewBaseBiz)
)

type (
	BaseBiz struct {
		pkg.Base

		baseRepo repo.BaseRepo
	}
)

func NewBaseBiz(baseRepo repo.BaseRepo) (*BaseBiz, error) {

	b := &BaseBiz{
		baseRepo: baseRepo,
	}
	if err := b.Init(); nil != err {
		return nil, err
	}
	return b, nil
}

func (b *BaseBiz) Init() error {
	return nil
}
