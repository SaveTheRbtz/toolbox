package gh_rest

import (
	"github.com/watermint/toolbox/domain/github/api/gh_recovery"
	"github.com/watermint/toolbox/infra/network/nw_capture"
	"github.com/watermint/toolbox/infra/network/nw_client"
	"github.com/watermint/toolbox/infra/network/nw_http"
)

var (
	defaultClient = gh_recovery.New(nw_capture.New(nw_http.NewClient()))
)

func Default() nw_client.Rest {
	return defaultClient
}