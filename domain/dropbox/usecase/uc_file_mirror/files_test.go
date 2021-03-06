package uc_file_mirror

import (
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_context_impl"
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/quality/infra/qt_errors"
	"github.com/watermint/toolbox/quality/recipe/qtr_endtoend"
	"testing"
)

func TestFilesImpl_Mirror(t *testing.T) {
	qtr_endtoend.TestWithControl(t, func(ctl app_control.Control) {
		ctx := dbx_context_impl.NewMock(ctl)
		sv := New(ctx, ctx)
		err := sv.Mirror(qtr_endtoend.NewTestDropboxFolderPath("from"), qtr_endtoend.NewTestDropboxFolderPath("to"))
		if err != nil && err != qt_errors.ErrorMock {
			t.Error(err)
		}
	})
}
