package benchmark

import (
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_conn"
	"github.com/watermint/toolbox/domain/dropbox/filesystem"
	"github.com/watermint/toolbox/domain/dropbox/model/mo_path"
	"github.com/watermint/toolbox/domain/dropbox/service/sv_file_content"
	"github.com/watermint/toolbox/essentials/file/es_filesystem_model"
	"github.com/watermint/toolbox/essentials/file/es_sync"
	"github.com/watermint/toolbox/essentials/model/em_file"
	"github.com/watermint/toolbox/essentials/model/mo_int"
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/infra/recipe/rc_exec"
	"github.com/watermint/toolbox/infra/recipe/rc_recipe"
	"github.com/watermint/toolbox/quality/recipe/qtr_endtoend"
)

type Upload struct {
	Peer        dbx_conn.ConnUserFile
	Path        mo_path.DropboxPath
	Lambda      int
	MinNodes    int
	MaxNodes    int
	MinFileSize int64
	MaxFileSize int64
	ChunkSizeKb mo_int.RangeInt
}

func (z *Upload) Preset() {
	z.Lambda = 100
	z.MinNodes = 10
	z.MaxNodes = 1000
	z.MinFileSize = 100 * 1024 * 1024
	z.MaxFileSize = 1 * 1024 * 1024 * 1024
	z.ChunkSizeKb.SetRange(1, 150*1024, 64*1024)
}

func (z *Upload) Exec(c app_control.Control) error {
	model := em_file.NewGenerator().Generate(
		em_file.NumNodes(z.Lambda, z.MinNodes, z.MaxNodes),
		em_file.FileSize(z.MinFileSize, z.MaxFileSize),
	)
	copier := filesystem.NewModelToDropbox(model, z.Peer.Context(), sv_file_content.ChunkSizeKb(z.ChunkSizeKb.Value()))
	syncer := es_sync.New(
		c.Log(),
		c.Sequence(),
		es_filesystem_model.NewFileSystem(model),
		filesystem.NewFileSystem(z.Peer.Context()),
		copier,
	)

	return syncer.Sync(es_filesystem_model.NewPath("/"),
		filesystem.NewPath("", z.Path))

}

func (z *Upload) Test(c app_control.Control) error {
	return rc_exec.ExecMock(c, &Upload{}, func(r rc_recipe.Recipe) {
		m := r.(*Upload)
		m.Path = qtr_endtoend.NewTestDropboxFolderPath("benchmark")
	})
}
