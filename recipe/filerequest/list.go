package filerequest

import (
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_conn"
	"github.com/watermint/toolbox/domain/dropbox/model/mo_filerequest"
	"github.com/watermint/toolbox/domain/dropbox/service/sv_filerequest"
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/infra/recipe/rc_exec"
	"github.com/watermint/toolbox/infra/recipe/rc_recipe"
	"github.com/watermint/toolbox/infra/report/rp_model"
)

type List struct {
	Peer         dbx_conn.ConnUserFile
	FileRequests rp_model.RowReport
}

func (z *List) Preset() {
	z.FileRequests.SetModel(&mo_filerequest.FileRequest{})
}

func (z *List) Exec(c app_control.Control) error {
	if err := z.FileRequests.Open(); err != nil {
		return err
	}
	reqs, err := sv_filerequest.New(z.Peer.Context()).List()
	if err != nil {
		return err
	}
	for _, r := range reqs {
		z.FileRequests.Row(r)
	}
	return nil
}

func (z *List) Test(c app_control.Control) error {
	return rc_exec.Exec(c, z, rc_recipe.NoCustomValues)
}
