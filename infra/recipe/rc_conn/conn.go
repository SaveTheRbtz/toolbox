package rc_conn

import (
	"github.com/watermint/toolbox/infra/control/app_control"
)

type Connection interface {
	// Connect to api
	Connect(ctl app_control.Control) (err error)

	// Peer name
	PeerName() string

	// Update peer name
	SetPeerName(name string)

	// Scope label
	ScopeLabel() string
}
