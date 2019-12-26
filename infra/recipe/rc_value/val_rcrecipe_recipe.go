package rc_value

import (
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/infra/recipe/rc_kitchen"
	"github.com/watermint/toolbox/infra/recipe/rc_recipe"
	"github.com/watermint/toolbox/infra/report/rp_model"
	"reflect"
)

func newValueRcRecipeRecipe(name string, t reflect.Type) Value {
	v := &ValueRcRecipeRecipe{name: name}
	if t == nil {
		v.recipe = &EmptyRecipe{}
		v.recipeType = reflect.TypeOf(v.recipe).Elem()
	} else {
		v.recipeType = t
		v.recipe = reflect.New(t.Elem()).Interface().(rc_recipe.Recipe)
	}
	return v
}

type ValueRcRecipeRecipe struct {
	name       string
	recipe     rc_recipe.Recipe
	recipeType reflect.Type
}

func (z *ValueRcRecipeRecipe) Reports() map[string]rp_model.Report {
	spec := NewRepository(z.recipe)
	return spec.Reports()
}

func (z *ValueRcRecipeRecipe) Accept(t reflect.Type, r rc_recipe.Recipe, name string) Value {
	if t.Implements(reflect.TypeOf((*rc_recipe.Recipe)(nil)).Elem()) {
		return newValueRcRecipeRecipe(name, t)
	}
	return nil
}

func (z *ValueRcRecipeRecipe) Bind() interface{} {
	return nil
}

func (z *ValueRcRecipeRecipe) Init() (v interface{}) {
	return z.recipe
}

func (z *ValueRcRecipeRecipe) Apply(v0 interface{}) (v interface{}) {
	return z.recipe
}

func (z *ValueRcRecipeRecipe) Debug() interface{} {
	return map[string]string{
		"typePkg":  z.recipeType.PkgPath(),
		"typeName": z.recipeType.Name(),
	}
}

func (z *ValueRcRecipeRecipe) SpinUp(ctl app_control.Control) error {
	return nil
}

func (z *ValueRcRecipeRecipe) SpinDown(ctl app_control.Control) error {
	return nil
}

type EmptyRecipe struct {
}

func (z *EmptyRecipe) Exec(k rc_kitchen.Kitchen) error {
	return nil
}

func (z *EmptyRecipe) Test(c app_control.Control) error {
	return nil
}

func (z *EmptyRecipe) Preset() {
}