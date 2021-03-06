package file

import (
	"errors"
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_conn"
	"github.com/watermint/toolbox/domain/dropbox/model/mo_file"
	"github.com/watermint/toolbox/domain/dropbox/model/mo_path"
	"github.com/watermint/toolbox/domain/dropbox/service/sv_file_content"
	"github.com/watermint/toolbox/essentials/log/esl"
	mo_path2 "github.com/watermint/toolbox/essentials/model/mo_path"
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/infra/recipe/rc_exec"
	"github.com/watermint/toolbox/infra/recipe/rc_recipe"
	"github.com/watermint/toolbox/infra/report/rp_model"
	"github.com/watermint/toolbox/quality/infra/qt_errors"
	"github.com/watermint/toolbox/quality/infra/qt_file"
	"github.com/watermint/toolbox/quality/recipe/qtr_endtoend"
	"os"
	"path/filepath"
	"time"
)

type Download struct {
	rc_recipe.RemarkExperimental
	Peer         dbx_conn.ConnUserFile
	DropboxPath  mo_path.DropboxPath
	LocalPath    mo_path2.FileSystemPath
	OperationLog rp_model.RowReport
}

func (z *Download) Preset() {
	z.OperationLog.SetModel(
		&mo_file.ConcreteEntry{},
		rp_model.HiddenColumns(
			"id",
			"path_lower",
			"revision",
			"content_hash",
			"shared_folder_id",
			"parent_shared_folder_id",
		),
	)
}

func (z *Download) Exec(c app_control.Control) error {
	l := c.Log()
	ctx := z.Peer.Context()

	if err := z.OperationLog.Open(); err != nil {
		return err
	}

	entry, f, err := sv_file_content.NewDownload(ctx).Download(z.DropboxPath)
	if err != nil {
		return err
	}
	if err := os.Rename(f.Path(), filepath.Join(z.LocalPath.Path(), entry.Name())); err != nil {
		l.Debug("Unable to move file to specified path",
			esl.Error(err),
			esl.String("downloaded", f.Path()),
			esl.String("destination", z.LocalPath.Path()),
		)
		return err
	}

	z.OperationLog.Row(entry.Concrete())
	return nil
}

func (z *Download) Test(c app_control.Control) error {
	// replay test
	{
		path, err := qt_file.MakeTestFolder("download", false)
		if err != nil {
			return err
		}
		defer func() {
			_ = os.RemoveAll(path)
		}()
		err = rc_exec.ExecReplay(c, &Download{}, "recipe-file-download.json.gz", func(r rc_recipe.Recipe) {
			m := r.(*Download)
			m.DropboxPath = mo_path.NewDropboxPath("/watermint-toolbox-test/watermint-toolbox.txt")
			m.LocalPath = mo_path2.NewFileSystemPath(path)
		})
		if err2, _ := qt_errors.ErrorsForTest(c.Log(), err); err2 != nil {
			return err2
		}

		// in case the test passed
		if err == nil {
			testFile := filepath.Join(path, "watermint-toolbox.txt")
			testFileInfo, err := os.Lstat(testFile)
			if err != nil {
				return err
			}

			if !testFileInfo.ModTime().Equal(time.Unix(1593502474, 0)) {
				return errors.New("invalid mod time")
			}
		}
	}

	return rc_exec.ExecMock(c, &Download{}, func(r rc_recipe.Recipe) {
		m := r.(*Download)
		m.LocalPath = qtr_endtoend.NewTestFileSystemFolderPath(c, "download")
		m.DropboxPath = qtr_endtoend.NewTestDropboxFolderPath("file-download")
	})
}
