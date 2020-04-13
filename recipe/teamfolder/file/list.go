package file

import (
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_conn"
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/infra/recipe/rc_exec"
	"github.com/watermint/toolbox/infra/recipe/rc_recipe"
	namespacefile "github.com/watermint/toolbox/ingredient/team/namespace/file"
	"github.com/watermint/toolbox/quality/infra/qt_errors"
)

type List struct {
	Peer     dbx_conn.ConnBusinessFile
	FileList *namespacefile.List
}

func (z *List) Preset() {
}

func (z *List) Exec(c app_control.Control) error {
	return rc_exec.Exec(c, z.FileList, func(r rc_recipe.Recipe) {
		rc := r.(*namespacefile.List)
		rc.IncludeTeamFolder = true
		rc.IncludeSharedFolder = false
	})
}

func (z *List) Test(c app_control.Control) error {
	return qt_errors.ErrorNoTestRequired
}
