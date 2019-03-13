package sv_file

import (
	"github.com/watermint/toolbox/domain/infra/api_context"
	"github.com/watermint/toolbox/domain/infra/api_list"
	"github.com/watermint/toolbox/domain/model/mo_file"
	"github.com/watermint/toolbox/domain/model/mo_path"
	"go.uber.org/zap"
)

type Files interface {
	Resolve(path mo_path.Path) (entry mo_file.Entry, err error)
	List(path mo_path.Path) (entries []mo_file.Entry, err error)

	Delete(path mo_path.Path) (entry mo_file.Entry, err error)
	DeleteWithRevision(path mo_path.Path, revision string) (entry mo_file.Entry, err error)

	PermDelete(path mo_path.Path) (err error)
	PermDeleteWithRevision(path mo_path.Path, revision string) (err error)
}

func NewFiles(ctx api_context.Context) Files {
	return &filesImpl{
		dc: ctx,
	}
}

func newFilesTest(ctx api_context.Context) Files {
	return &filesImpl{
		dc: ctx,
		//limit: 3,
	}
}

type filesImpl struct {
	dc                              api_context.Context
	recursive                       bool
	includeMediaInfo                bool
	includeDeleted                  bool
	includeHasExplicitSharedMembers bool
	limit                           int
}

func (z *filesImpl) Resolve(path mo_path.Path) (entry mo_file.Entry, err error) {
	panic("implement me")
}

func (z *filesImpl) List(path mo_path.Path) (entries []mo_file.Entry, err error) {
	entries = make([]mo_file.Entry, 0)
	p := struct {
		Path                            string `json:"path"`
		Recursive                       bool   `json:"recursive,omitempty"`
		IncludeMediaInfo                bool   `json:"include_media_info,omitempty"`
		IncludeDeleted                  bool   `json:"include_deleted,omitempty"`
		IncludeHasExplicitSharedMembers bool   `json:"include_has_explicit_shared_members,omitempty"`
		Limit                           int    `json:"limit,omitempty"`
	}{
		Path:                            path.Path(),
		Recursive:                       z.recursive,
		IncludeMediaInfo:                z.includeMediaInfo,
		IncludeDeleted:                  z.includeDeleted,
		IncludeHasExplicitSharedMembers: z.includeHasExplicitSharedMembers,
	}

	req := z.dc.List("files/list_folder").
		Continue("files/list_folder/continue").
		Param(p).
		UseHasMore(true).
		ResultTag("entries").
		OnEntry(func(entry api_list.ListEntry) error {
			e := &mo_file.Metadata{}
			if err := entry.Model(e); err != nil {
				j, _ := entry.Json()
				z.dc.Log().Error("invalid", zap.Error(err), zap.String("entry", j.Raw))
				return err
			}
			entries = append(entries, e)
			return nil
		})
	if err := req.Call(); err != nil {
		return nil, err
	}
	return entries, nil
}

func (z *filesImpl) Delete(path mo_path.Path) (entry mo_file.Entry, err error) {
	panic("implement me")
}

func (z *filesImpl) DeleteWithRevision(path mo_path.Path, revision string) (entry mo_file.Entry, err error) {
	panic("implement me")
}

func (z *filesImpl) PermDelete(path mo_path.Path) (err error) {
	panic("implement me")
}

func (z *filesImpl) PermDeleteWithRevision(path mo_path.Path, revision string) (err error) {
	panic("implement me")
}