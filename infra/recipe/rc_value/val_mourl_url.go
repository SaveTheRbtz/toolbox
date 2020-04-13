package rc_value

import (
	"github.com/watermint/toolbox/domain/dropbox/model/mo_url"
	"github.com/watermint/toolbox/infra/app"
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/infra/control/app_root"
	"github.com/watermint/toolbox/infra/recipe/rc_recipe"
	"github.com/watermint/toolbox/infra/util/ut_reflect"
	"go.uber.org/zap"
	"reflect"
)

func newValueMoUrlUrl(name string) rc_recipe.Value {
	v := &ValueMoUrlUrl{name: name}
	v.url = mo_url.NewEmptyUrl()
	return v
}

type ValueMoUrlUrl struct {
	name   string
	rawUrl string
	url    mo_url.Url
}

func (z *ValueMoUrlUrl) ValueText() string {
	return z.rawUrl
}

func (z *ValueMoUrlUrl) Spec() (typeName string, typeAttr interface{}) {
	return ut_reflect.Key(app.Pkg, z.url), nil
}

func (z *ValueMoUrlUrl) Accept(t reflect.Type, v0 interface{}, name string) rc_recipe.Value {
	if t.Implements(reflect.TypeOf((*mo_url.Url)(nil)).Elem()) {
		return newValueMoUrlUrl(name)
	}
	return nil
}

func (z *ValueMoUrlUrl) Bind() interface{} {
	return &z.rawUrl
}

func (z *ValueMoUrlUrl) Init() (v interface{}) {
	return z.url
}

func (z *ValueMoUrlUrl) ApplyPreset(v0 interface{}) {
	z.url = v0.(mo_url.Url)
	z.rawUrl = z.url.Value()
}

func (z *ValueMoUrlUrl) Apply() (v interface{}) {
	l := app_root.Log()
	u, err := mo_url.NewUrl(z.rawUrl)
	if err != nil {
		l.Debug("Unable to parse", zap.String("url", z.rawUrl), zap.Error(err))
	} else {
		z.url = u
	}
	return z.url
}

func (z *ValueMoUrlUrl) Debug() interface{} {
	return map[string]string{
		"url": z.rawUrl,
	}
}

func (z *ValueMoUrlUrl) SpinUp(ctl app_control.Control) error {
	l := app_root.Log()
	if z.rawUrl == "" {
		return ErrorMissingRequiredOption
	}

	u, err := mo_url.NewUrl(z.rawUrl)
	if err != nil {
		l.Debug("Unable to parse", zap.String("url", z.rawUrl), zap.Error(err))
		return ErrorInvalidValue
	} else {
		z.url = u
		return nil
	}
}

func (z *ValueMoUrlUrl) SpinDown(ctl app_control.Control) error {
	return nil
}
