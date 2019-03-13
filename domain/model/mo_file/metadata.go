package mo_file

import (
	"encoding/json"
	"github.com/watermint/toolbox/domain/infra/api_parser"
	"github.com/watermint/toolbox/domain/model/mo_path"
)

type Metadata struct {
	Raw              json.RawMessage
	EntryTag         string `path:"\\.tag"`
	EntryName        string `path:"name"`
	EntryPathDisplay string `path:"path_display"`
	EntryPathLower   string `path:"path_lower"`
}

func (z *Metadata) Tag() string {
	return z.EntryTag
}

func (z *Metadata) Name() string {
	return z.EntryName
}

func (z *Metadata) PathDisplay() string {
	return z.EntryPathDisplay
}

func (z *Metadata) PathLower() string {
	return z.EntryPathLower
}

func (z *Metadata) Path() mo_path.Path {
	return mo_path.NewPathDisplay(z.EntryPathDisplay)
}

func (z *Metadata) File() (*File, bool) {
	if z.EntryTag != "file" {
		return nil, false
	}
	f := &File{}
	if err := api_parser.ParseModelRaw(f, z.Raw); err != nil {
		return nil, false // Should not happen
	}
	return f, true
}

func (z *Metadata) Folder() (*Folder, bool) {
	if z.EntryTag != "folder" {
		return nil, false
	}
	f := &Folder{}
	if err := api_parser.ParseModelRaw(f, z.Raw); err != nil {
		return nil, false
	}
	return f, true
}

func (z *Metadata) Deleted() (*Deleted, bool) {
	if z.EntryTag != "deleted" {
		return nil, false
	}
	d := &Deleted{}
	if err := api_parser.ParseModelRaw(d, z.Raw); err != nil {
		return nil, false
	}
	return d, true
}