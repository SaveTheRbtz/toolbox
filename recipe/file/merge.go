package file

import (
	"github.com/watermint/toolbox/domain/model/mo_path"
	"github.com/watermint/toolbox/domain/usecase/uc_file_merge"
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/infra/recpie/app_conn"
	"github.com/watermint/toolbox/infra/recpie/app_kitchen"
	"github.com/watermint/toolbox/infra/recpie/app_vo"
)

type MergeVO struct {
	Peer                app_conn.ConnUserFile
	From                string
	To                  string
	DryRun              bool
	KeepEmptyFolder     bool
	WithinSameNamespace bool
}

type Merge struct {
}

func (z *Merge) Console() {
}

func (z *Merge) Requirement() app_vo.ValueObject {
	return &MergeVO{
		DryRun:              true,
		KeepEmptyFolder:     false,
		WithinSameNamespace: false,
	}
}

func (z *Merge) Exec(k app_kitchen.Kitchen) error {
	vo := k.Value().(*MergeVO)

	ctx, err := vo.Peer.Connect(k.Control())
	if err != nil {
		return err
	}

	ufm := uc_file_merge.New(ctx, k)
	opts := make([]uc_file_merge.MergeOpt, 0)
	if vo.DryRun {
		opts = append(opts, uc_file_merge.DryRun())
	}
	if !vo.KeepEmptyFolder {
		opts = append(opts, uc_file_merge.ClearEmptyFolder())
	}
	if vo.WithinSameNamespace {
		opts = append(opts, uc_file_merge.WithinSameNamespace())
	}

	return ufm.Merge(mo_path.NewPath(vo.From), mo_path.NewPath(vo.To), opts...)
}

func (z *Merge) Test(c app_control.Control) error {
	return nil
}
