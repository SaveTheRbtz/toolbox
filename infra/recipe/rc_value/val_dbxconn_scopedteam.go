package rc_value

import (
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_conn"
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_conn_impl"
	"github.com/watermint/toolbox/essentials/go/es_reflect"
	"github.com/watermint/toolbox/infra/api/api_conn"
	"github.com/watermint/toolbox/infra/app"
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/infra/recipe/rc_recipe"
	"reflect"
	"strings"
)

func newValueDbxConnScopedTeam(peerName string) rc_recipe.Value {
	return &ValueDbxConnScopedTeam{
		conn:     dbx_conn_impl.NewConnScopedTeam(peerName),
		peerName: peerName,
	}
}

type ValueDbxConnScopedTeam struct {
	conn     dbx_conn.ConnScopedTeam
	peerName string
}

func (z *ValueDbxConnScopedTeam) Accept(t reflect.Type, v0 interface{}, name string) rc_recipe.Value {
	if t.Implements(reflect.TypeOf((*dbx_conn.ConnScopedTeam)(nil)).Elem()) {
		return newValueDbxConnScopedTeam(name)
	}
	return nil
}

func (z *ValueDbxConnScopedTeam) Bind() interface{} {
	return &z.peerName
}

func (z *ValueDbxConnScopedTeam) Init() (v interface{}) {
	return z.conn
}

func (z *ValueDbxConnScopedTeam) ApplyPreset(v0 interface{}) {
	z.conn = v0.(dbx_conn.ConnScopedTeam)
	z.peerName = z.conn.PeerName()
}

func (z *ValueDbxConnScopedTeam) Apply() (v interface{}) {
	z.conn.SetPeerName(z.peerName)
	return z.conn
}

func (z *ValueDbxConnScopedTeam) Debug() interface{} {
	return map[string]string{
		"peerName": z.peerName,
		"scopes":   strings.Join(z.conn.Scopes(), ","),
		"appType":  z.conn.ScopeLabel(),
	}
}

func (z *ValueDbxConnScopedTeam) SpinUp(ctl app_control.Control) error {
	return z.conn.Connect(ctl)
}

func (z *ValueDbxConnScopedTeam) SpinDown(ctl app_control.Control) error {
	return nil
}

func (z *ValueDbxConnScopedTeam) Conn() (conn api_conn.Connection, valid bool) {
	return z.conn, true
}

func (z *ValueDbxConnScopedTeam) Spec() (typeName string, typeAttr interface{}) {
	return es_reflect.Key(app.Pkg, z.conn), z.conn.Scopes()
}
