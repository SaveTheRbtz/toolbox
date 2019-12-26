package rc_value

import (
	"github.com/watermint/toolbox/infra/api/api_auth_impl"
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/infra/recipe/rc_conn"
	"github.com/watermint/toolbox/infra/recipe/rc_conn_impl"
	"github.com/watermint/toolbox/infra/recipe/rc_recipe"
	"github.com/watermint/toolbox/quality/infra/qt_recipe"
	"reflect"
)

func newValueRcConnBusinessInfo(peerName string) Value {
	v := &ValueRcConnBusinessInfo{peerName: peerName}
	v.conn = rc_conn_impl.NewConnBusinessInfo(peerName)
	return v
}

type ValueRcConnBusinessInfo struct {
	conn     rc_conn.ConnBusinessInfo
	peerName string
}

func (z *ValueRcConnBusinessInfo) ValueText() string {
	return z.peerName
}

func (z *ValueRcConnBusinessInfo) Accept(t reflect.Type, r rc_recipe.Recipe, name string) Value {
	if t.Implements(reflect.TypeOf((*rc_conn.ConnBusinessInfo)(nil)).Elem()) {
		return newValueRcConnBusinessInfo(z.peerName)
	}
	return nil
}

func (z *ValueRcConnBusinessInfo) Bind() interface{} {
	return &z.peerName
}

func (z *ValueRcConnBusinessInfo) Init() (v interface{}) {
	return z.conn
}

func (z *ValueRcConnBusinessInfo) Apply() (v interface{}) {
	z.conn.SetName(z.peerName)
	return z.conn
}

func (z *ValueRcConnBusinessInfo) SpinUp(ctl app_control.Control) error {
	if ctl.IsTest() {
		if qt_recipe.IsSkipEndToEndTest() {
			return qt_recipe.ErrorSkipEndToEndTest
		}
		a := api_auth_impl.NewCached(ctl, api_auth_impl.PeerName(z.peerName))
		if _, err := a.Auth(z.conn.ScopeLabel()); err != nil {
			return err
		}
	}
	return z.conn.Connect(ctl)
}

func (z *ValueRcConnBusinessInfo) SpinDown(ctl app_control.Control) error {
	return nil
}

func (z *ValueRcConnBusinessInfo) Conn() (conn rc_conn.ConnDropboxApi, valid bool) {
	return z.conn, true
}

func (z *ValueRcConnBusinessInfo) Debug() interface{} {
	return map[string]string{
		"peerName": z.peerName,
		"scope":    z.conn.ScopeLabel(),
	}
}
