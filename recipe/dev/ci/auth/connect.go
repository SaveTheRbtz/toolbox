package auth

import (
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/infra/recipe/rc_conn"
	"github.com/watermint/toolbox/infra/recipe/rc_exec"
	"github.com/watermint/toolbox/infra/recipe/rc_recipe"
	"github.com/watermint/toolbox/quality/infra/qt_endtoend"
)

type Connect struct {
	Full  rc_conn.ConnUserFile
	Info  rc_conn.ConnBusinessInfo
	File  rc_conn.ConnBusinessFile
	Audit rc_conn.ConnBusinessAudit
	Mgmt  rc_conn.ConnBusinessMgmt
}

func (z *Connect) Preset() {
	z.Full.SetPeerName(qt_endtoend.EndToEndPeer)
	z.Info.SetPeerName(qt_endtoend.EndToEndPeer)
	z.File.SetPeerName(qt_endtoend.EndToEndPeer)
	z.Audit.SetPeerName(qt_endtoend.EndToEndPeer)
	z.Mgmt.SetPeerName(qt_endtoend.EndToEndPeer)
}

func (z *Connect) Exec(c app_control.Control) error {
	return nil
}

func (z *Connect) Test(c app_control.Control) error {
	return rc_exec.Exec(c, &Connect{}, rc_recipe.NoCustomValues)
}