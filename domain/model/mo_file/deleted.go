package mo_file

import (
	"encoding/json"
	"github.com/watermint/toolbox/domain/model/mo_path"
)

type Deleted struct {
	Raw              json.RawMessage
	EntryTag         string `path:"\\.tag"`
	EntryName        string `path:"name"`
	EntryPathLower   string `path:"path_lower"`
	EntryPathDisplay string `path:"path_display"`
}

func (z *Deleted) Tag() string {
	return z.EntryTag
}

func (z *Deleted) Name() string {
	return z.EntryName
}

func (z *Deleted) PathDisplay() string {
	return z.EntryPathDisplay
}

func (z *Deleted) PathLower() string {
	return z.EntryPathLower
}

func (z *Deleted) Path() mo_path.Path {
	return mo_path.NewPathDisplay(z.EntryPathDisplay)
}

func (z *Deleted) File() (*File, bool) {
	return nil, false
}

func (z *Deleted) Folder() (*Folder, bool) {
	return nil, false
}

func (z *Deleted) Deleted() (*Deleted, bool) {
	return z, true
}