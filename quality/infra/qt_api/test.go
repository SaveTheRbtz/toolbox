package qt_api

import (
	"github.com/watermint/toolbox/domain/dropbox/model/mo_path"
	"github.com/watermint/toolbox/infra/api/api_auth"
	"github.com/watermint/toolbox/infra/api/api_context"
	"github.com/watermint/toolbox/infra/api/dbx_auth"
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/quality/infra/qt_endtoend"
	"github.com/watermint/toolbox/quality/infra/qt_recipe"
	"go.uber.org/zap"
)

var (
	ToolboxTestSuiteFolder = mo_path.NewDropboxPath("/toolbox-testsuite")
	legacyTestsEnabled     = false
)

func DoTestTokenFull(test func(ctx api_context.DropboxApiContext)) {
	doTest(api_auth.DropboxTokenFull, test)
}
func DoTestBusinessInfo(test func(ctx api_context.DropboxApiContext)) {
	doTest(api_auth.DropboxTokenBusinessInfo, test)
}
func DoTestBusinessFile(test func(ctx api_context.DropboxApiContext)) {
	doTest(api_auth.DropboxTokenBusinessFile, test)
}
func DoTestBusinessManagement(test func(ctx api_context.DropboxApiContext)) {
	doTest(api_auth.DropboxTokenBusinessManagement, test)
}
func DoTestBusinessAudit(test func(ctx api_context.DropboxApiContext)) {
	doTest(api_auth.DropboxTokenBusinessAudit, test)
}

func doTest(tokenType string, test func(ctx api_context.DropboxApiContext)) {
	qt_recipe.TestWithControl(nil, func(ctl app_control.Control) {
		l := ctl.Log()
		a := dbx_auth.NewCached(ctl, dbx_auth.PeerName(qt_endtoend.EndToEndPeer))
		ctx, err := a.Auth(tokenType)
		if err != nil {
			l.Info("Skip test", zap.Error(err))
			return
		}
		dtx, ok := ctx.(api_context.DropboxApiContext)
		if !ok {
			l.Warn("Invalid context type", zap.Any("ctx", ctx))
			return
		}
		if legacyTestsEnabled {
			test(dtx)
		}
	})
}