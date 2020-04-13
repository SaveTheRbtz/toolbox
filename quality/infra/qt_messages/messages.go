package qt_messages

import (
	"errors"
	"flag"
	"fmt"
	"github.com/watermint/toolbox/infra/app"
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/infra/control/app_control_launcher"
	"github.com/watermint/toolbox/infra/recipe/rc_group"
	"github.com/watermint/toolbox/infra/recipe/rc_group_impl"
	"github.com/watermint/toolbox/infra/recipe/rc_recipe"
	"github.com/watermint/toolbox/infra/recipe/rc_spec"
	"github.com/watermint/toolbox/infra/ui/app_msg_container"
	"github.com/watermint/toolbox/infra/ui/app_ui"
	"github.com/watermint/toolbox/infra/util/ut_io"
	"go.uber.org/zap"
	"io"
	"os"
	"sort"
	"strings"
)

func SuggestMessages(ctl app_control.Control, suggest func(out io.Writer)) {
	out := ut_io.NewDefaultOut(ctl.Feature().IsTest())
	fmt.Fprintln(out, "Please add those messages to message resource files:")
	fmt.Fprintln(out, "====================================================")
	suggest(out)
	fmt.Fprintln(out, "====================================================")
}

func VerifyMessages(ctl app_control.Control) error {
	cl := ctl.(app_control_launcher.ControlLauncher)
	cat := cl.Catalogue()
	recipes := cat.Recipes()
	root := rc_group_impl.NewGroup()
	for _, r := range recipes {
		s := rc_spec.New(r)
		root.Add(s)
	}

	qui := app_ui.NewQuiet(ctl.Messages())
	if qui, ok := qui.(*app_ui.Quiet); ok {
		qui.SetLogger(ctl.Log())
	}
	verifyGroup(root, qui)

	qm := ctl.Messages().(app_msg_container.Quality)
	missing := qm.MissingKeys()
	if len(missing) > 0 {
		sort.Strings(missing)
		for _, k := range missing {
			ctl.Log().Error("Key missing", zap.String(k, ""))
		}

		SuggestMessages(ctl, func(out io.Writer) {
			for _, m := range missing {
				switch {
				case strings.HasSuffix(m, ".flag.peer"):
					fmt.Fprintf(out, `"%s":"Account alias",`, m)
				case strings.HasSuffix(m, ".flag.file"):
					fmt.Fprintf(out, `"%s":"Path to data file",`, m)
				case strings.HasSuffix(m, ".agreement"):
					fmt.Fprintf(out, `"%s":"This feature is in an early stage of development. This is not well tested. Please proceed by typing 'yes' to agree & enable this feature.",`, m)
				case strings.HasSuffix(m, ".disclaimer"):
					fmt.Fprintf(out, `"%s":"WARN: The early access feature is enabled.",`, m)
				default:
					fmt.Fprintf(out, `"%s":"",`, m)
				}
				fmt.Fprintln(out)
			}
		})

		return errors.New("missing key found")
	}
	return nil
}

func verifyGroup(g rc_group.Group, ui app_ui.UI) {
	g.PrintUsage(ui, os.Args[0], app.Version)
	for _, sg := range g.SubGroups() {
		verifyGroup(sg, ui)
	}
	for _, r := range g.Recipes() {
		verifyRecipe(g, r, ui)
	}
}

func verifyRecipe(g rc_group.Group, r rc_recipe.Spec, ui app_ui.UI) {
	f := flag.NewFlagSet("", flag.ContinueOnError)

	r.SetFlags(f, ui)
	r.PrintUsage(ui)
}
