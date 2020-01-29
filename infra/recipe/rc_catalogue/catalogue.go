package rc_catalogue

import (
	"github.com/watermint/toolbox/infra/recipe/rc_group"
	"github.com/watermint/toolbox/infra/recipe/rc_group_impl"
	"github.com/watermint/toolbox/infra/recipe/rc_recipe"
	"github.com/watermint/toolbox/infra/recipe/rc_spec"
)

type Catalogue interface {
	Recipes() []rc_recipe.Recipe
	Ingredients() []rc_recipe.Recipe
	Messages() []interface{}
	RootGroup() rc_group.Group
}

type catalogueImpl struct {
	recipes     []rc_recipe.Recipe
	ingredients []rc_recipe.Recipe
	messages    []interface{}
	root        rc_group.Group
}

func (z *catalogueImpl) Recipes() []rc_recipe.Recipe {
	return z.recipes
}

func (z *catalogueImpl) Ingredients() []rc_recipe.Recipe {
	return z.ingredients
}

func (z *catalogueImpl) Messages() []interface{} {
	return z.messages
}

func (z *catalogueImpl) RootGroup() rc_group.Group {
	return z.root
}

func NewCatalogue(recipes, ingredients []rc_recipe.Recipe, messages []interface{}) Catalogue {
	root := rc_group_impl.NewGroup()
	for _, r := range recipes {
		s := rc_spec.New(r)
		root.Add(s)
	}

	return &catalogueImpl{
		recipes:     recipes,
		ingredients: ingredients,
		messages:    messages,
		root:        root,
	}
}

func NewEmptyCatalogue() Catalogue {
	return NewCatalogue([]rc_recipe.Recipe{}, []rc_recipe.Recipe{}, []interface{}{})
}