package file

import (
	"github.com/watermint/toolbox/infra/recpie/app_test"
	"testing"
)

func TestMove_Exec(t *testing.T) {
	app_test.TestRecipe(t, &Move{})
}