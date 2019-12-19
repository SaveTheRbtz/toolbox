package fd_file

import "github.com/watermint/toolbox/infra/control/app_control"

type Rows interface {
	EachRow(exec func(m interface{}, rowIndex int) error) error
}

// Row feed interface for SelfContainedRecipe
type RowFeed interface {
	Rows
	Model(m interface{})
}

// File interface for SideCarRecipe
// Deprecated: use RowFeed
type ModelFile interface {
	Rows
	Model(ctl app_control.Control, m interface{}) error
}
