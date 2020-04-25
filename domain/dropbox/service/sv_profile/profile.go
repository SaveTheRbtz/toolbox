package sv_profile

import (
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_context"
	"github.com/watermint/toolbox/domain/dropbox/model/mo_profile"
)

type Profile interface {
	Current() (profile *mo_profile.Profile, err error)
}

func NewProfile(ctx dbx_context.Context) Profile {
	return &profileImpl{
		ctx: ctx,
	}
}

type Team interface {
	Admin() (profile *mo_profile.Profile, err error)
}

func NewTeam(ctx dbx_context.Context) Team {
	return &teamImpl{
		ctx: ctx,
	}
}

type profileImpl struct {
	ctx dbx_context.Context
}

func (z *profileImpl) Current() (profile *mo_profile.Profile, err error) {
	profile = &mo_profile.Profile{}
	res, err := z.ctx.Post("users/get_current_account").Call()
	if err != nil {
		return nil, err
	}
	if _, err = res.Success().Json().Model(profile); err != nil {
		return nil, err
	}
	return profile, nil
}

type teamImpl struct {
	ctx dbx_context.Context
}

func (z *teamImpl) Admin() (profile *mo_profile.Profile, err error) {
	profile = &mo_profile.Profile{}
	res, err := z.ctx.Post("team/token/get_authenticated_admin").Call()
	if err != nil {
		return nil, err
	}
	if _, err = res.Success().Json().FindModel("admin_profile", profile); err != nil {
		return nil, err
	}
	return profile, nil
}
