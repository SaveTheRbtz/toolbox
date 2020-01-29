package rc_group_impl

import (
	"bytes"
	"errors"
	"flag"
	"github.com/watermint/toolbox/infra/app"
	"github.com/watermint/toolbox/infra/recipe/rc_group"
	"github.com/watermint/toolbox/infra/recipe/rc_recipe"
	"github.com/watermint/toolbox/infra/ui/app_msg"
	"github.com/watermint/toolbox/infra/ui/app_ui"
	"os"
	"sort"
	"strings"
)

type MsgGroup struct {
	RecipeHeaderUsage      app_msg.Message
	RecipeUsage            app_msg.Message
	RecipeAvailableFlags   app_msg.Message
	GroupHeaderUsage       app_msg.Message
	GroupUsage             app_msg.Message
	GroupAvailableCommands app_msg.Message
}

var (
	MGroup = app_msg.Apply(&MsgGroup{}).(*MsgGroup)
)

func NewGroup() rc_group.Group {
	return newGroupWithPath([]string{}, "")
}

func newGroupWithPath(path []string, name string) rc_group.Group {
	return &groupImpl{
		name:      name,
		path:      path,
		recipes:   make(map[string]rc_recipe.Spec),
		subGroups: make(map[string]rc_group.Group),
	}
}

type groupImpl struct {
	name      string
	basePkg   string
	path      []string
	recipes   map[string]rc_recipe.Spec
	subGroups map[string]rc_group.Group
}

func (z *groupImpl) Name() string {
	return z.name
}

func (z *groupImpl) BasePkg() string {
	return z.basePkg
}

func (z *groupImpl) Path() []string {
	return z.path
}

func (z *groupImpl) Recipes() map[string]rc_recipe.Spec {
	return z.recipes
}

func (z *groupImpl) SubGroups() map[string]rc_group.Group {
	return z.subGroups
}

func (z *groupImpl) GroupDesc() app_msg.Message {
	grpDesc := make([]string, 0)
	grpDesc = append(grpDesc, "recipe")
	grpDesc = append(grpDesc, z.path...)
	grpDesc = append(grpDesc, "title")

	return app_msg.M(strings.Join(grpDesc, "."))
}

func (z *groupImpl) AddToPath(fullPath []string, relPath []string, name string, r rc_recipe.Spec) {
	if len(relPath) > 0 {
		p0 := relPath[0]
		sg, ok := z.subGroups[p0]
		if !ok {
			sg = newGroupWithPath(fullPath, p0)
			z.subGroups[p0] = sg
		}
		sg.AddToPath(fullPath, relPath[1:], name, r)
	} else {
		z.recipes[name] = r
	}
}

func (z *groupImpl) Add(r rc_recipe.Spec) {
	path, name := r.Path()

	z.AddToPath(path, path, name, r)
}

func (z *groupImpl) usageHeader(ui app_ui.UI, desc app_msg.Message, version string) {
	rc_group.AppHeader(ui, version)
	ui.Break()
	ui.Info(desc)
	ui.Break()
}

func (z *groupImpl) PrintRecipeUsage(ui app_ui.UI, spec rc_recipe.Spec, f *flag.FlagSet) {
	z.usageHeader(ui, spec.Title(), app.Version)

	ui.Header(MGroup.RecipeHeaderUsage)
	ui.Info(MGroup.RecipeUsage.
		With("Exec", os.Args[0]).
		With("Recipe", spec.CliPath()).
		With("Args", ui.TextOrEmpty(spec.CliArgs())))

	ui.Break()
	ui.Header(MGroup.RecipeAvailableFlags)

	buf := new(bytes.Buffer)
	f.SetOutput(buf)
	f.PrintDefaults()
	ui.Info(app_msg.Raw(buf.String()))
	ui.Break()
}

func (z *groupImpl) PrintGroupUsage(ui app_ui.UI, exec, version string) {
	z.usageHeader(ui, z.GroupDesc(), version)

	ui.Header(MGroup.GroupHeaderUsage)
	ui.Info(MGroup.GroupUsage.
		With("Exec", exec).
		With("Group", strings.Join(z.path, " ")))
	ui.Break()

	ui.Header(MGroup.GroupAvailableCommands)
	cmdTable := ui.InfoTable("usage")
	cmds, ca := z.commandAnnotations(ui)
	for _, cmd := range cmds {
		ann := ca[cmd]
		cmdTable.Row(app_msg.Raw(" "), app_msg.Raw(cmd), z.CommandTitle(cmd), app_msg.Raw(ann))
	}
	cmdTable.Flush()
}

func (z *groupImpl) CommandTitle(cmd string) app_msg.Message {
	keyPath := make([]string, 0)
	keyPath = append(keyPath, "recipe")
	keyPath = append(keyPath, z.path...)
	keyPath = append(keyPath, cmd)
	keyPath = append(keyPath, "title")
	key := strings.Join(keyPath, ".")

	return app_msg.M(key)
}

func (z *groupImpl) IsSecret() bool {
	for _, r := range z.recipes {
		if !r.IsSecret() {
			return false
		}
	}
	return true
}

func (z *groupImpl) commandAnnotations(ui app_ui.UI) (cmds []string, annotation map[string]string) {
	cmds = make([]string, 0)
	annotation = make(map[string]string)
	for _, g := range z.subGroups {
		if !g.IsSecret() {
			cmds = append(cmds, g.Name())
		}
		annotation[g.Name()] = ""
	}
	for n, r := range z.recipes {
		if !r.IsSecret() {
			cmds = append(cmds, n)
		}
		annotation[n] = ui.TextOrEmpty(r.Remarks())
	}
	sort.Strings(cmds)
	return
}

func (z *groupImpl) Select(args []string) (name string, g rc_group.Group, r rc_recipe.Spec, remainder []string, err error) {
	if len(args) < 1 {
		return "", z, nil, args, nil
	}
	arg := args[0]
	for k, sg := range z.subGroups {
		if arg == k {
			return sg.Select(args[1:])
		}
	}
	for k, sr := range z.Recipes() {
		if arg == k {
			return k, z, sr, args[1:], nil
		}
	}
	return "", z, nil, args, errors.New("not found")
}