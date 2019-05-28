package app_recipe

import (
	"github.com/watermint/toolbox/app86/app_control"
	"github.com/watermint/toolbox/app86/app_report"
	"github.com/watermint/toolbox/app86/app_ui"
	"github.com/watermint/toolbox/app86/app_vo"
	"github.com/watermint/toolbox/domain/infra/api_context"
	"go.uber.org/zap"
)

type Recipe interface {
	Requirement() app_vo.ValueObject
	Exec(k Kitchen) error
}

// SecretRecipe will not be listed in available commands.
type SecretRecipe interface {
	Hidden()
}

type Kitchen interface {
	Value() app_vo.ValueObject
	Control() app_control.Control
	UI() app_ui.UI
	Log() *zap.Logger
	Report() app_report.Report
}

type ApiKitchen interface {
	Kitchen
	Context() api_context.Context
}

func WithBusinessFile(exec func(k ApiKitchen) error) error {
	panic("implement me")
}

func WithBusinessManagement(exec func(k ApiKitchen) error) error {
	panic("implement me")
}

type kitchenImpl struct {
	vo  app_vo.ValueObject
	ctl app_control.Control
}

func (z *kitchenImpl) Value() app_vo.ValueObject {
	return z.vo
}

func (z *kitchenImpl) Control() app_control.Control {
	return z.ctl
}

func (z *kitchenImpl) UI() app_ui.UI {
	return z.UI()
}

func (z *kitchenImpl) Log() *zap.Logger {
	return z.Log()
}

func (z *kitchenImpl) Report() app_report.Report {
	panic("implement me")
}

func NewKitchen(ctl app_control.Control, vo app_vo.ValueObject) Kitchen {
	return &kitchenImpl{
		ctl: ctl,
		vo:  vo,
	}
}
