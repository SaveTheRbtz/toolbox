package qt_messages

import (
	"errors"
	"flag"
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/infra/control/app_control_launcher"
	"github.com/watermint/toolbox/infra/recipe/rc_group"
	"github.com/watermint/toolbox/infra/recipe/rc_recipe"
	"github.com/watermint/toolbox/infra/recipe/rc_spec"
	"github.com/watermint/toolbox/infra/ui/app_msg_container"
	"github.com/watermint/toolbox/infra/ui/app_ui"
	"go.uber.org/zap"
)

func VerifyMessages(ctl app_control.Control) error {
	cl := ctl.(app_control_launcher.ControlLauncher)
	root := rc_group.NewGroup([]string{}, "")
	for _, r := range cl.Catalogue() {
		root.Add(r)
	}

	qui := app_ui.NewQuiet(ctl.Messages())
	if qui, ok := qui.(*app_ui.Quiet); ok {
		qui.SetLogger(ctl.Log())
	}
	verifyGroup(root, qui)

	qm := ctl.Messages().(app_msg_container.Quality)
	missing := qm.MissingKeys()
	if len(missing) > 0 {
		for _, k := range missing {
			ctl.Log().Error("Key missing", zap.String(k, ""))
		}
		return errors.New("missing key found")
	}
	return nil
}

func verifyGroup(g *rc_group.Group, ui app_ui.UI) {
	g.PrintUsage(ui)
	for _, sg := range g.SubGroups {
		verifyGroup(sg, ui)
	}
	for _, r := range g.Recipes {
		verifyRecipe(g, r, ui)
	}
}

func verifyRecipe(g *rc_group.Group, r rc_recipe.Recipe, ui app_ui.UI) {
	f := flag.NewFlagSet("", flag.ContinueOnError)

	spec := rc_spec.New(r)
	if spec == nil {
		// skip
		return
	}
	spec.SetFlags(f, ui)
	g.PrintRecipeUsage(ui, spec, f)
}
