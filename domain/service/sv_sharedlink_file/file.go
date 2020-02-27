package sv_sharedlink_file

import (
	"encoding/json"
	"github.com/watermint/toolbox/domain/model/mo_file"
	"github.com/watermint/toolbox/domain/model/mo_path"
	"github.com/watermint/toolbox/domain/model/mo_url"
	"github.com/watermint/toolbox/domain/service/sv_sharedlink"
	"github.com/watermint/toolbox/infra/api/api_context"
	"github.com/watermint/toolbox/infra/api/api_list"
	"go.uber.org/zap"
	"strings"
)

type File interface {
	List(url mo_url.Url, path mo_path.DropboxPath, nEntry func(entry mo_file.Entry), opt ...ListOpt) error
	ListRecursive(url mo_url.Url, nEntry func(entry mo_file.Entry), opt ...ListOpt) error
}

type ListOpt func(opt *ListOpts) *ListOpts
type ListOpts struct {
	IncludeDeleted                  bool
	IncludeHasExplicitSharedMembers bool
	Password                        string
}

func IncludeDeleted() ListOpt {
	return func(opt *ListOpts) *ListOpts {
		opt.IncludeDeleted = true
		return opt
	}
}
func IncludeHasExplicitSharedMembers() ListOpt {
	return func(opt *ListOpts) *ListOpts {
		opt.IncludeHasExplicitSharedMembers = true
		return opt
	}
}
func Password(password string) ListOpt {
	return func(opt *ListOpts) *ListOpts {
		opt.Password = password
		return opt
	}
}

func New(ctx api_context.Context) File {
	return &fileImpl{ctx: ctx}
}

type fileImpl struct {
	ctx api_context.Context
}

func (z *fileImpl) ListRecursive(url mo_url.Url, nEntry func(entry mo_file.Entry), opts ...ListOpt) error {
	lo := &ListOpts{}
	for _, o := range opts {
		o(lo)
	}
	var ls func(entry0 mo_file.Entry, rel string) error
	ls = func(entry0 mo_file.Entry, rel string) error {
		c := entry0.Concrete()
		c.PathDisplay = rel
		if !strings.HasPrefix(c.PathDisplay, "/") {
			c.PathDisplay = "/" + c.PathDisplay
		}
		c.PathLower = strings.ToLower(c.PathDisplay) // TODO: i18n
		r := make(map[string]interface{})
		if err0 := json.Unmarshal(c.Raw, &r); err0 == nil {
			r["path_display"] = c.PathDisplay
			r["path_lower"] = c.PathLower
			if r0, err0 := json.Marshal(&r); err0 == nil {
				c.Raw = r0
			}
		}

		if f, ok := entry0.File(); ok {
			f.Raw = c.Raw
			f.EntryPathDisplay = c.PathDisplay
			f.EntryPathLower = c.PathLower
			nEntry(f)
			return nil
		}
		if f, ok := entry0.Deleted(); ok {
			f.Raw = c.Raw
			f.EntryPathDisplay = c.PathDisplay
			f.EntryPathLower = c.PathLower
			nEntry(f)
			return nil
		}
		if f, ok := entry0.Folder(); ok {
			f.Raw = c.Raw
			f.EntryPathDisplay = c.PathDisplay
			f.EntryPathLower = c.PathLower
			nEntry(f)
		}

		return z.List(url, mo_path.NewDropboxPath(rel), func(entry1 mo_file.Entry) {
			ls(entry1, rel+"/"+entry1.Name())
		}, opts...)
	}

	entry, err := sv_sharedlink.New(z.ctx).Resolve(url, lo.Password)
	if err != nil {
		return err
	}

	return ls(entry, "")
}

func (z *fileImpl) List(url mo_url.Url, path mo_path.DropboxPath, onEntry func(entry mo_file.Entry), opts ...ListOpt) error {
	lo := &ListOpts{}
	for _, o := range opts {
		o(lo)
	}

	type SL struct {
		Url      string `json:"url"`
		Password string `json:"password,omitempty"`
	}
	p := struct {
		Path                            string `json:"path"`
		SharedLink                      *SL    `json:"shared_link"`
		Recursive                       bool   `json:"recursive,omitempty"`
		IncludeDeleted                  bool   `json:"include_deleted,omitempty"`
		IncludeHasExplicitSharedMembers bool   `json:"include_has_explicit_shared_members,omitempty"`
		Limit                           int    `json:"limit,omitempty"`
	}{
		Path: path.Path(),
		SharedLink: &SL{
			Url:      url.String(),
			Password: lo.Password,
		},
		IncludeDeleted:                  lo.IncludeDeleted,
		IncludeHasExplicitSharedMembers: lo.IncludeHasExplicitSharedMembers,
	}

	req := z.ctx.List("files/list_folder").
		Continue("files/list_folder/continue").
		Param(p).
		UseHasMore(true).
		ResultTag("entries").
		OnEntry(func(entry api_list.ListEntry) error {
			e := &mo_file.Metadata{}
			if err := entry.Model(e); err != nil {
				j, _ := entry.Json()
				z.ctx.Log().Error("invalid", zap.Error(err), zap.String("entry", j.Raw))
				return err
			}
			onEntry(e)
			return nil
		})
	return req.Call()
}