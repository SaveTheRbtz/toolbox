package app_test

import (
	"encoding/csv"
	rice "github.com/GeertJohan/go.rice"
	"github.com/watermint/toolbox/infra/api/api_auth"
	"github.com/watermint/toolbox/infra/api/api_auth_impl"
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/infra/control/app_control_impl"
	"github.com/watermint/toolbox/infra/control/app_run_impl"
	"github.com/watermint/toolbox/infra/quality/qt_control_impl"
	"github.com/watermint/toolbox/infra/recpie/app_conn"
	"github.com/watermint/toolbox/infra/recpie/app_conn_impl"
	"github.com/watermint/toolbox/infra/recpie/app_recipe"
	"github.com/watermint/toolbox/infra/recpie/app_vo"
	"github.com/watermint/toolbox/infra/recpie/app_vo_impl"
	"github.com/watermint/toolbox/infra/ui/app_msg_container"
	"github.com/watermint/toolbox/infra/ui/app_ui"
	"go.uber.org/zap"
	"io"
	"os"
	"path/filepath"
	"testing"
)

const (
	EndToEndPeer = "end_to_end_test"
)

func ApplyTestPeers(ctl app_control.Control, vo app_vo.ValueObject) bool {
	l := ctl.Log()
	l.Debug("Prepare for applying test peers")
	a := api_auth_impl.NewCached(ctl, api_auth_impl.PeerName(EndToEndPeer))

	vc := app_vo_impl.NewValueContainer(vo)
	for k, v := range vc.Values {
		if _, ok := v.(app_conn.ConnBusinessInfo); ok {
			if _, err := a.Auth(api_auth.DropboxTokenBusinessInfo); err != nil {
				l.Info("BusinessInfo: Skip end to end test", zap.String("k", k))
				return false
			}
			vc.Values[k] = &app_conn_impl.ConnBusinessInfo{
				PeerName: EndToEndPeer,
			}
		} else if _, ok := v.(app_conn.ConnBusinessFile); ok {
			if _, err := a.Auth(api_auth.DropboxTokenBusinessFile); err != nil {
				l.Info("BusinessFile: Skip end to end test", zap.String("k", k))
				return false
			}
			vc.Values[k] = &app_conn_impl.ConnBusinessFile{
				PeerName: EndToEndPeer,
			}
		} else if _, ok := v.(app_conn.ConnBusinessAudit); ok {
			if _, err := a.Auth(api_auth.DropboxTokenBusinessAudit); err != nil {
				l.Info("BusinessAudit: Skip end to end test", zap.String("k", k))
				return false
			}
			vc.Values[k] = &app_conn_impl.ConnBusinessAudit{
				PeerName: EndToEndPeer,
			}
		} else if _, ok := v.(app_conn.ConnBusinessMgmt); ok {
			if _, err := a.Auth(api_auth.DropboxTokenBusinessManagement); err != nil {
				l.Info("BusinessManagement: Skip end to end test", zap.String("k", k))
				return false
			}
			vc.Values[k] = &app_conn_impl.ConnBusinessMgmt{
				PeerName: EndToEndPeer,
			}
		} else if _, ok := v.(app_conn.ConnUserFile); ok {
			if _, err := a.Auth(api_auth.DropboxTokenFull); err != nil {
				l.Info("UserFull: Skip end to end test", zap.String("k", k))
				return false
			}
			vc.Values[k] = &app_conn_impl.ConnUserFile{
				PeerName: EndToEndPeer,
			}
		}
	}

	l.Debug("Applying for debug")
	vc.Apply(vo)

	return true
}

func TestResources(t *testing.T) (bx, web *rice.Box, mc app_msg_container.Container, ui app_ui.UI) {
	bx = rice.MustFindBox("../../../resources")
	web = rice.MustFindBox("../../../web")

	mc = app_run_impl.NewContainer(bx)
	ui = app_ui.NewConsole(mc, qt_control_impl.NewMessageTest(t), true)
	return
}

func TestRecipe(t *testing.T, re app_recipe.Recipe) {
	bx, web, mc, ui := TestResources(t)

	ctl := app_control_impl.NewSingle(ui, bx, web, mc, false, make([]app_recipe.Recipe, 0))
	err := ctl.Up(app_control.Test())
	if err != nil {
		os.Exit(app_control.FatalStartup)
	}
	defer ctl.Down()

	if err := re.Test(ctl); err != nil {
		t.Error("test failed", err)
	}
}

type RowTester func(cols map[string]string) error

func TestRows(ctl app_control.Control, reportName string, tester RowTester) error {
	l := ctl.Log().With(zap.String("reportName", reportName))
	job := ctl.Workspace().Job()
	rep := filepath.Join(job, "reports")
	csvFile := filepath.Join(rep, reportName+".csv")

	l.Debug("Start loading report", zap.String("csvFile", csvFile))

	cf, err := os.Open(csvFile)
	if err != nil {
		l.Warn("Unable to open report CSV", zap.Error(err))
		return err
	}
	defer cf.Close()
	csf := csv.NewReader(cf)
	var header []string
	isFirstLine := true

	for {
		cols, err := csf.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			l.Warn("An error occurred during read report file", zap.Error(err))
			return err
		}
		if isFirstLine {
			header = cols
			isFirstLine = false
		} else {
			colMap := make(map[string]string)
			for i, h := range header {
				colMap[h] = cols[i]
			}
			if err := tester(colMap); err != nil {
				l.Warn("Tester returned an error", zap.Error(err), zap.Any("cols", colMap))
				return err
			}
		}
	}

	return nil
}