package uc_file_relocation

import (
	"github.com/watermint/toolbox/infra/api/api_context"
	"github.com/watermint/toolbox/quality/infra/qt_errors"
	"github.com/watermint/toolbox/quality/infra/qt_recipe"
	"testing"
)

func TestRelocationImpl_Copy(t *testing.T) {
	qt_recipe.TestWithApiContext(t, func(ctx api_context.DropboxApiContext) {
		sv := New(ctx)
		err := sv.Copy(qt_recipe.NewTestDropboxFolderPath("from"), qt_recipe.NewTestDropboxFolderPath("to"))
		if err != nil && err != qt_errors.ErrorMock {
			t.Error(err)
		}
	})
}

func TestRelocationImpl_Move(t *testing.T) {
	qt_recipe.TestWithApiContext(t, func(ctx api_context.DropboxApiContext) {
		sv := New(ctx)
		err := sv.Move(qt_recipe.NewTestDropboxFolderPath("from"), qt_recipe.NewTestDropboxFolderPath("to"))
		if err != nil && err != qt_errors.ErrorMock {
			t.Error(err)
		}
	})
}