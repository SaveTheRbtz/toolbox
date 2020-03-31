package dbx_auth

import (
	"github.com/watermint/toolbox/infra/api/api_auth"
	"github.com/watermint/toolbox/infra/control/app_control"
)

func New(control app_control.Control, opts ...AuthOpt) api_auth.Console {
	ao := &authOpts{
		tokenType: api_auth.DropboxTokenNoAuth,
		peerName:  "default",
	}
	for _, o := range opts {
		o(ao)
	}
	ua := &CcAuth{
		control: control,
	}
	ua.init()
	ca := &CcCachedAuth{
		peerName: ao.peerName,
		control:  control,
		auth:     ua,
	}
	ca.init()
	return ca
}

func NewCached(control app_control.Control, opts ...AuthOpt) api_auth.Console {
	ao := &authOpts{
		tokenType: api_auth.DropboxTokenNoAuth,
		peerName:  "default",
	}
	for _, o := range opts {
		o(ao)
	}
	ca := &CcCachedAuth{
		peerName: ao.peerName,
		control:  control,
	}
	ca.init()
	return ca
}

type AuthOpt func(opt *authOpts) *authOpts
type authOpts struct {
	peerName  string
	tokenType string
}

func PeerName(name string) AuthOpt {
	return func(opt *authOpts) *authOpts {
		opt.peerName = name
		return opt
	}
}
func Full() AuthOpt {
	return func(opt *authOpts) *authOpts {
		opt.tokenType = api_auth.DropboxTokenFull
		return opt
	}
}
func BusinessFile() AuthOpt {
	return func(opt *authOpts) *authOpts {
		opt.tokenType = api_auth.DropboxTokenBusinessFile
		return opt
	}
}
func BusinessManagement() AuthOpt {
	return func(opt *authOpts) *authOpts {
		opt.tokenType = api_auth.DropboxTokenBusinessManagement
		return opt
	}
}
func BusinessInfo() AuthOpt {
	return func(opt *authOpts) *authOpts {
		opt.tokenType = api_auth.DropboxTokenBusinessInfo
		return opt
	}
}
func BusinessAudit() AuthOpt {
	return func(opt *authOpts) *authOpts {
		opt.tokenType = api_auth.DropboxTokenBusinessAudit
		return opt
	}
}