package benchmark

import (
	"github.com/watermint/toolbox/essentials/file/es_filesystem_copier"
	"github.com/watermint/toolbox/essentials/file/es_filesystem_local"
	"github.com/watermint/toolbox/essentials/file/es_filesystem_model"
	"github.com/watermint/toolbox/essentials/file/es_sync"
	"github.com/watermint/toolbox/essentials/model/em_file"
	"github.com/watermint/toolbox/essentials/model/mo_path"
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/infra/recipe/rc_exec"
	"github.com/watermint/toolbox/infra/recipe/rc_recipe"
	"github.com/watermint/toolbox/quality/infra/qt_file"
	"os"
)

type Local struct {
	Path       mo_path.FileSystemPath
	NodeMin    int
	NodeMax    int
	NodeLambda int
	SizeMin    int
	SizeMax    int
}

func (z *Local) Preset() {
	z.NodeMin = 100
	z.NodeLambda = 100
	z.NodeMax = 1000
	z.SizeMin = 0
	z.SizeMax = 2 * 1_048_576 // 2 MiB
}

func (z *Local) Exec(c app_control.Control) error {
	model := em_file.NewGenerator().Generate(
		em_file.NumNodes(z.NodeLambda, z.NodeMin, z.NodeMax),
		em_file.FileSize(int64(z.SizeMin), int64(z.SizeMax)),
	)

	copier := es_filesystem_copier.NewModelToLocal(c.Log(), model)
	syncer := es_sync.New(
		c.Log(),
		c.Sequence(),
		es_filesystem_model.NewFileSystem(model),
		es_filesystem_local.NewFileSystem(),
		copier,
	)

	return syncer.Sync(es_filesystem_model.NewPath("/"),
		es_filesystem_local.NewPath(z.Path.Path()))
}

func (z *Local) Test(c app_control.Control) error {
	workPath, err := qt_file.MakeTestFolder("local", false)
	if err != nil {
		return err
	}
	defer func() {
		_ = os.RemoveAll(workPath)
	}()

	return rc_exec.ExecMock(c, &Local{}, func(r rc_recipe.Recipe) {
		m := r.(*Local)
		m.Path = mo_path.NewFileSystemPath(workPath)
	})
}
