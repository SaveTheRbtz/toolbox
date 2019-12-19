package device

import (
	"errors"
	"github.com/watermint/toolbox/domain/model/mo_device"
	"github.com/watermint/toolbox/domain/model/mo_member"
	"github.com/watermint/toolbox/domain/service/sv_device"
	"github.com/watermint/toolbox/domain/service/sv_member"
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/infra/recipe/rc_conn"
	"github.com/watermint/toolbox/infra/recipe/rc_kitchen"
	"github.com/watermint/toolbox/infra/recipe/rc_vo"
	"github.com/watermint/toolbox/infra/report/rp_spec"
	"github.com/watermint/toolbox/infra/report/rp_spec_impl"
	"github.com/watermint/toolbox/quality/infra/qt_recipe"
)

type ListVO struct {
	Peer rc_conn.ConnBusinessFile
}

const (
	reportList = "device"
)

type List struct {
}

func (z *List) Reports() []rp_spec.ReportSpec {
	return []rp_spec.ReportSpec{
		rp_spec_impl.Spec(reportList, &mo_device.MemberSession{}),
	}
}

func (z *List) Requirement() rc_vo.ValueObject {
	return &ListVO{}
}

func (z *List) Exec(k rc_kitchen.Kitchen) error {
	lvo := k.Value().(*ListVO)
	ctx, err := lvo.Peer.Connect(k.Control())
	if err != nil {
		return err
	}

	memberList, err := sv_member.New(ctx).List()
	if err != nil {
		return err
	}
	members := mo_member.MapByTeamMemberId(memberList)

	sessions, err := sv_device.New(ctx).List()
	if err != nil {
		return err
	}

	rep, err := rp_spec_impl.New(z, k.Control()).Open(reportList)
	if err != nil {
		return err
	}
	defer rep.Close()

	for _, session := range sessions {
		if m, e := members[session.EntryTeamMemberId()]; e {
			ma := mo_device.NewMemberSession(m, session)
			rep.Row(ma)
		}
	}
	return nil
}

func (z *List) Test(c app_control.Control) error {
	lvo := &ListVO{}
	if !qt_recipe.ApplyTestPeers(c, lvo) {
		return qt_recipe.NotEnoughResource()
	}
	if err := z.Exec(rc_kitchen.NewKitchen(c, lvo)); err != nil {
		return err
	}
	return qt_recipe.TestRows(c, "device", func(cols map[string]string) error {
		if _, ok := cols["team_member_id"]; !ok {
			return errors.New("team_member_id is not found")
		}
		return nil
	})
}
