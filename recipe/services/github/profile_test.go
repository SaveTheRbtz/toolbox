package github

import (
	"github.com/watermint/toolbox/quality/infra/qt_recipe"
	"testing"
)

func TestProfile_Exec(t *testing.T) {
	qt_recipe.TestRecipe(t, &Profile{})
}