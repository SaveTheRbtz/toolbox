package eq_registry

import (
	"github.com/watermint/toolbox/essentials/queue/eq_bundle"
	"github.com/watermint/toolbox/essentials/queue/eq_mould"
	"sync"
)

type Registry interface {
	Define(mouldId string, f interface{}, ctx ...interface{}) eq_mould.Mould

	Get(mouldId string) (mould eq_mould.Mould, found bool)
}

func New(bundle eq_bundle.Bundle, ehs []eq_mould.ErrorHandler) Registry {
	return &regImpl{
		bundle:        bundle,
		moulds:        make(map[string]eq_mould.Mould),
		mouldsLock:    sync.Mutex{},
		errorHandlers: ehs,
	}
}

type regImpl struct {
	bundle        eq_bundle.Bundle
	moulds        map[string]eq_mould.Mould
	mouldsLock    sync.Mutex
	errorHandlers []eq_mould.ErrorHandler
}

func (z *regImpl) Define(mouldId string, f interface{}, ctx ...interface{}) eq_mould.Mould {
	mould := eq_mould.New(mouldId, z.bundle, z.errorHandlers, f, ctx...)
	z.moulds[mouldId] = mould
	return mould
}

func (z *regImpl) Get(mouldId string) (mould eq_mould.Mould, found bool) {
	mould, found = z.moulds[mouldId]
	return
}
