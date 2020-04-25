package sv_teamfolder

import (
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_context"
	"github.com/watermint/toolbox/domain/dropbox/model/mo_teamfolder"
	"github.com/watermint/toolbox/essentials/format/tjson"
)

type TeamFolder interface {
	List() (teamfolders []*mo_teamfolder.TeamFolder, err error)
	Resolve(teamFolderId string) (teamfolder *mo_teamfolder.TeamFolder, err error)
	Create(name string, opts ...CreateOption) (teamfolder *mo_teamfolder.TeamFolder, err error)
	Activate(tf *mo_teamfolder.TeamFolder) (teamfolder *mo_teamfolder.TeamFolder, err error)
	Archive(tf *mo_teamfolder.TeamFolder) (teamfolder *mo_teamfolder.TeamFolder, err error)
	Rename(tf *mo_teamfolder.TeamFolder, newName string) (updated *mo_teamfolder.TeamFolder, err error)
	PermDelete(tf *mo_teamfolder.TeamFolder) (err error)
}

type createOptions struct {
	syncSetting string
}

type CreateOption func(opt *createOptions) *createOptions

func SyncDefault() CreateOption {
	return func(opt *createOptions) *createOptions {
		opt.syncSetting = "default"
		return opt
	}
}
func SyncNoSync() CreateOption {
	return func(opt *createOptions) *createOptions {
		opt.syncSetting = "not_synced"
		return opt
	}
}

func New(ctx dbx_context.Context) TeamFolder {
	return &teamFolderImpl{
		ctx: ctx,
	}
}

type teamFolderImpl struct {
	ctx dbx_context.Context
}

func (z *teamFolderImpl) List() (teamfolders []*mo_teamfolder.TeamFolder, err error) {
	teamfolders = make([]*mo_teamfolder.TeamFolder, 0)
	err = z.ctx.List("team/team_folder/list").
		Continue("team/team_folder/list/continue").
		UseHasMore(true).
		ResultTag("team_folders").
		OnEntry(func(entry tjson.Json) error {
			tf := &mo_teamfolder.TeamFolder{}
			if _, err := entry.Model(tf); err != nil {
				return err
			}
			teamfolders = append(teamfolders, tf)
			return nil
		}).Call()
	if err != nil {
		return nil, err
	}
	return teamfolders, nil
}

func (z *teamFolderImpl) Resolve(teamFolderId string) (teamfolder *mo_teamfolder.TeamFolder, err error) {
	teamfolder = &mo_teamfolder.TeamFolder{}
	p := struct {
		TeamFolderIds []string `json:"team_folder_ids"`
	}{
		TeamFolderIds: []string{teamFolderId},
	}
	res, err := z.ctx.Post("team/team_folder/get_info").Param(p).Call()
	if err != nil {
		return nil, err
	}
	if _, err = res.Success().Json().FindModel(tjson.PathArrayFirst, teamfolder); err != nil {
		return nil, err
	}
	return teamfolder, nil
}

func (z *teamFolderImpl) Create(name string, opts ...CreateOption) (teamfolder *mo_teamfolder.TeamFolder, err error) {
	co := &createOptions{}
	for _, o := range opts {
		o(co)
	}
	p := struct {
		Name        string `json:"name"`
		SyncSetting string `json:"sync_setting,omitempty"`
	}{
		Name:        name,
		SyncSetting: co.syncSetting,
	}

	teamfolder = &mo_teamfolder.TeamFolder{}
	res, err := z.ctx.Post("team/team_folder/create").Param(p).Call()
	if err != nil {
		return nil, err
	}
	if _, err = res.Success().Json().Model(teamfolder); err != nil {
		return nil, err
	}
	return teamfolder, nil
}

func (z *teamFolderImpl) Activate(tf *mo_teamfolder.TeamFolder) (teamfolder *mo_teamfolder.TeamFolder, err error) {
	p := struct {
		TeamFolderId string `json:"team_folder_id"`
	}{
		TeamFolderId: tf.TeamFolderId,
	}
	teamfolder = &mo_teamfolder.TeamFolder{}
	res, err := z.ctx.Post("team/team_folder/activate").Param(p).Call()
	if err != nil {
		return nil, err
	}
	if _, err = res.Success().Json().Model(teamfolder); err != nil {
		return nil, err
	}
	return teamfolder, nil
}

func (z *teamFolderImpl) Archive(tf *mo_teamfolder.TeamFolder) (teamfolder *mo_teamfolder.TeamFolder, err error) {
	p := struct {
		TeamFolderId string `json:"team_folder_id"`
	}{
		TeamFolderId: tf.TeamFolderId,
	}
	teamfolder = &mo_teamfolder.TeamFolder{}
	res, err := z.ctx.Async("team/team_folder/archive").
		Status("team/team_folder/archive/check").
		Param(p).
		Call()
	if err != nil {
		return nil, err
	}
	if _, err = res.Success().Json().Model(teamfolder); err != nil {
		return nil, err
	}
	return teamfolder, nil
}

func (z *teamFolderImpl) Rename(tf *mo_teamfolder.TeamFolder, newName string) (updated *mo_teamfolder.TeamFolder, err error) {
	p := struct {
		TeamFolderId string `json:"team_folder_id"`
		Name         string `json:"name"`
	}{
		TeamFolderId: tf.TeamFolderId,
		Name:         newName,
	}
	updated = &mo_teamfolder.TeamFolder{}
	res, err := z.ctx.Post("team/team_folder/rename").Param(p).Call()
	if err != nil {
		return nil, err
	}
	if _, err = res.Success().Json().Model(updated); err != nil {
		return nil, err
	}
	return updated, nil
}

func (z *teamFolderImpl) PermDelete(tf *mo_teamfolder.TeamFolder) (err error) {
	p := struct {
		TeamFolderId string `json:"team_folder_id"`
	}{
		TeamFolderId: tf.TeamFolderId,
	}
	_, err = z.ctx.Post("team/team_folder/permanently_delete").Param(p).Call()
	return err
}
