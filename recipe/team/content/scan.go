package content

import (
	"errors"
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_context"
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/infra/ui/app_msg"
)

type MsgScan struct {
	ProgressScanNamespaceMetadata app_msg.Message
	ProgressScanNamespaceMember   app_msg.Message
	ProgressScanMember            app_msg.Message
}

var (
	ErrorTeamOwnedNamespaceIsNotInitialized = errors.New("team owned namespace is not initialized")
	MScanMetadata                           = app_msg.Apply(&MsgScan{}).(*MsgScan)
)

type ScanNamespace interface {
	Scan(ctl app_control.Control, ctx dbx_context.Context, namespaceName string, namespaceId string)
}

type ScanNamespaceMetadataAndMembership struct {
	metadata   ScanNamespace
	membership ScanNamespace
}

func (z *ScanNamespaceMetadataAndMembership) Scan(ctl app_control.Control, ctx dbx_context.Context, namespaceName string, namespaceId string) {
	z.membership.Scan(ctl, ctx, namespaceName, namespaceId)
	z.metadata.Scan(ctl, ctx, namespaceName, namespaceId)
}
