package member

import (
	"github.com/watermint/toolbox/domain/usecase/uc_member_mirror"
	"github.com/watermint/toolbox/infra/api/api_util"
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/infra/quality/qt_test"
	"github.com/watermint/toolbox/infra/recpie/app_conn"
	"github.com/watermint/toolbox/infra/recpie/app_file"
	"github.com/watermint/toolbox/infra/recpie/app_kitchen"
	"github.com/watermint/toolbox/infra/recpie/app_vo"
	"github.com/watermint/toolbox/infra/report/rp_model"
	"github.com/watermint/toolbox/infra/report/rp_spec"
	"github.com/watermint/toolbox/infra/report/rp_spec_impl"
	"github.com/watermint/toolbox/infra/ui/app_msg"
)

type ReplicationRow struct {
	SrcEmail string `json:"src_email"`
	DstEmail string `json:"dst_email"`
}

type ReplicationVO struct {
	Src  app_conn.ConnBusinessFile
	Dst  app_conn.ConnBusinessFile
	File app_file.Data
}

const (
	reportReplication = "replication"
)

type Replication struct {
}

func (z *Replication) Reports() []rp_spec.ReportSpec {
	return []rp_spec.ReportSpec{
		rp_spec_impl.Spec(reportReplication, rp_model.TransactionHeader(&ReplicationRow{}, nil)),
	}
}

func (z *Replication) Hidden() {
}

func (z *Replication) Console() {
}

func (z *Replication) Requirement() app_vo.ValueObject {
	return &ReplicationVO{}
}

func (z *Replication) Exec(k app_kitchen.Kitchen) error {
	vo := k.Value().(*ReplicationVO)
	ui := k.UI()

	ui.Info("recipe.member.replication.conn_src_file")
	src, err := vo.Src.Connect(k.Control())
	if err != nil {
		return err
	}

	ui.Info("recipe.member.replication.conn_dst_file")
	dst, err := vo.Src.Connect(k.Control())
	if err != nil {
		return err
	}

	if err := vo.File.Model(k.Control(), &ReplicationRow{}); err != nil {
		return err
	}

	rep, err := rp_spec_impl.New(z, k.Control()).Open(reportReplication)
	if err != nil {
		return err
	}
	defer rep.Close()

	return vo.File.EachRow(func(m interface{}, rowIndex int) error {
		row := m.(*ReplicationRow)

		ui.Info("recipe.member.replication.progress", app_msg.P{
			"SrcEmail": row.SrcEmail,
			"DstEmail": row.DstEmail,
		})
		err = uc_member_mirror.New(src, dst).Mirror(row.SrcEmail, row.DstEmail)
		if err != nil {
			rep.Failure(api_util.MsgFromError(err), row, nil)
			return err
		}
		rep.Success(row, nil)
		return nil
	})
}

func (z *Replication) Test(c app_control.Control) error {
	return qt_test.HumanInteractionRequired()
}