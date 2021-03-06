package rc_value

import (
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/infra/recipe/rc_recipe"
	"reflect"
)

func newValueString() rc_recipe.Value {
	return &ValueString{}
}

type ValueString struct {
	v string
}

func (z *ValueString) Spec() (typeName string, typeAttr interface{}) {
	return "string", nil
}

func (z *ValueString) Accept(t reflect.Type, v0 interface{}, name string) rc_recipe.Value {
	if t.Kind() == reflect.String {
		return newValueString()
	}
	return nil
}

func (z *ValueString) Bind() interface{} {
	return &z.v
}

func (z *ValueString) Init() (v interface{}) {
	return z.v
}

func (z *ValueString) ApplyPreset(v0 interface{}) {
	z.v = v0.(string)
}

func (z *ValueString) Apply() (v interface{}) {
	return z.v
}

func (z *ValueString) Debug() interface{} {
	return z.v
}

func (z *ValueString) SpinUp(ctl app_control.Control) error {
	if z.v == "" {
		return ErrorMissingRequiredOption
		//return nil
	}
	return nil
}

func (z *ValueString) SpinDown(ctl app_control.Control) error {
	return nil
}
