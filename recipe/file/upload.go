package file

import (
	mo_path2 "github.com/watermint/toolbox/domain/common/model/mo_path"
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_conn"
	"github.com/watermint/toolbox/domain/dropbox/model/mo_path"
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/infra/recipe/rc_exec"
	"github.com/watermint/toolbox/infra/recipe/rc_recipe"
	"github.com/watermint/toolbox/ingredient/file"
	"github.com/watermint/toolbox/quality/infra/qt_errors"
	"github.com/watermint/toolbox/quality/infra/qt_recipe"
	"os"
)

type Upload struct {
	Peer        dbx_conn.ConnUserFile
	LocalPath   mo_path2.FileSystemPath
	DropboxPath mo_path.DropboxPath
	Overwrite   bool
	ChunkSizeKb int
	Upload      *file.Upload
}

func (z *Upload) Preset() {
	z.ChunkSizeKb = 150 * 1024
}

func (z *Upload) Exec(c app_control.Control) error {
	return rc_exec.Exec(c, z.Upload, func(r rc_recipe.Recipe) {
		ru := r.(*file.Upload)
		ru.EstimateOnly = false
		ru.LocalPath = z.LocalPath
		ru.DropboxPath = z.DropboxPath
		ru.Overwrite = z.Overwrite
		ru.CreateFolder = false
		ru.Context = z.Peer.Context()
		if z.ChunkSizeKb > 0 {
			ru.ChunkSizeKb = z.ChunkSizeKb
		}
	})
}

func (z *Upload) Test(c app_control.Control) error {
	l := c.Log()
	fileCandidates := []string{"README.md", "upload.go", "upload_test.go"}
	file := ""
	for _, f := range fileCandidates {
		if _, err := os.Lstat(f); err == nil {
			file = f
			break
		}
	}
	if file == "" {
		l.Warn("No file to upload")
		return qt_errors.ErrorNotEnoughResource
	}

	{
		err := rc_exec.Exec(c, &Upload{}, func(r rc_recipe.Recipe) {
			ru := r.(*Upload)
			ru.LocalPath = mo_path2.NewFileSystemPath(file)
			ru.DropboxPath = qt_recipe.NewTestDropboxFolderPath()
			ru.Overwrite = true
		})
		if err != nil {
			return err
		}
	}

	// Chunked
	{
		err := rc_exec.Exec(c, &Upload{}, func(r rc_recipe.Recipe) {
			ru := r.(*Upload)
			ru.LocalPath = mo_path2.NewFileSystemPath(file)
			ru.DropboxPath = qt_recipe.NewTestDropboxFolderPath()
			ru.Overwrite = true
			ru.ChunkSizeKb = 1
		})
		if err != nil {
			return err
		}
	}
	return nil
}
