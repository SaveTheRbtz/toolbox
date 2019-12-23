package rp_model_deprecated

import (
	"github.com/watermint/toolbox/infra/report/rp_model"
)

// Deprecated:
func NewColumn(row interface{}, opts ...rp_model.ReportOpt) Column {
	ro := &rp_model.ReportOpts{}
	for _, opt := range opts {
		opt(ro)
	}
	ri := &columnJsonImpl{
		opts: ro,
	}
	_ = ri.Parse(row)

	return ri
}

// Deprecated:
type Column interface {
	Header() []string
	Values(r interface{}) []interface{}
	ValuesAsString(r interface{}) []string
}