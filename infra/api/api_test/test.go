package api_test

import (
	"github.com/watermint/toolbox/domain/model/mo_path"
	"github.com/watermint/toolbox/infra/api/api_auth"
	"github.com/watermint/toolbox/infra/api/api_auth_impl"
	"github.com/watermint/toolbox/infra/api/api_context"
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/infra/recpie/app_test"
	"go.uber.org/zap"
)

var (
	ToolboxTestSuiteFolder = mo_path.NewPath("/toolbox-testsuite")
	legacyTestsEnabled     = false
)

func DoTestTokenFull(test func(ctx api_context.Context)) {
	doTest(api_auth.DropboxTokenFull, test)
}
func DoTestBusinessInfo(test func(ctx api_context.Context)) {
	doTest(api_auth.DropboxTokenBusinessInfo, test)
}
func DoTestBusinessFile(test func(ctx api_context.Context)) {
	doTest(api_auth.DropboxTokenBusinessFile, test)
}
func DoTestBusinessManagement(test func(ctx api_context.Context)) {
	doTest(api_auth.DropboxTokenBusinessManagement, test)
}
func DoTestBusinessAudit(test func(ctx api_context.Context)) {
	doTest(api_auth.DropboxTokenBusinessAudit, test)
}

func doTest(tokenType string, test func(ctx api_context.Context)) {
	app_test.TestWithControl(nil, func(ctl app_control.Control) {
		l := ctl.Log()
		a := api_auth_impl.NewCached(ctl, api_auth_impl.PeerName(app_test.EndToEndPeer))
		ctx, err := a.Auth(tokenType)
		if err != nil {
			l.Info("Skip test", zap.Error(err))
			return
		}

		if legacyTestsEnabled {
			test(ctx)
		}
	})
}
