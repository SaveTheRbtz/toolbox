package auth

import (
	"github.com/watermint/toolbox/quality/infra/qt_recipe"
	"testing"
)

func TestAuth_Exec(t *testing.T) {
	qt_recipe.TestRecipe(t, &Connect{})
}