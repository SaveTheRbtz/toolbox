package content

import (
	"github.com/watermint/toolbox/domain/common/model/mo_filter"
	"github.com/watermint/toolbox/domain/dropbox/api/dbx_context"
	"github.com/watermint/toolbox/domain/dropbox/model/mo_member"
	"github.com/watermint/toolbox/domain/dropbox/service/sv_member"
	"github.com/watermint/toolbox/domain/dropbox/service/sv_namespace"
	"github.com/watermint/toolbox/domain/dropbox/service/sv_profile"
	"github.com/watermint/toolbox/domain/dropbox/service/sv_teamfolder"
	"github.com/watermint/toolbox/essentials/log/esl"
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/infra/recipe/rc_worker"
)

type TeamScanner struct {
	ctx                 dbx_context.Context
	ctl                 app_control.Control
	teamOwnedNamespaces map[string]bool
	scanner             ScanNamespace
	queue               rc_worker.Queue
	filter              mo_filter.Filter
}

func (z *TeamScanner) namespacesOfTeam() error {
	l := z.ctx.Log()

	l.Debug("Scanning admin")
	admin, err := sv_profile.NewTeam(z.ctx).Admin()
	if err != nil {
		return err
	}
	l = l.With(esl.String("admin", admin.Email))

	l.Debug("Scanning team folders")
	teamfolders, err := sv_teamfolder.New(z.ctx).List()
	if err != nil {
		return err
	}

	l.Debug("Scanning namespaces")
	namespaces, err := sv_namespace.New(z.ctx).List()
	if err != nil {
		return err
	}

	l.Debug("Computing duplicates")
	z.teamOwnedNamespaces = make(map[string]bool)
	teamOwnedNamespaceWithName := make(map[string]string)
	for _, f := range teamfolders {
		if z.filter.Accept(f.Name) {
			z.teamOwnedNamespaces[f.TeamFolderId] = true
			teamOwnedNamespaceWithName[f.TeamFolderId] = f.Name
		}
	}
	for _, n := range namespaces {
		if !z.filter.Accept(n.Name) {
			l.Debug("Skip folder that unmatched to filter condition", esl.String("name", n.Name))
			continue
		}

		switch n.NamespaceType {
		case "app_folder", "team_member_folder":
			l.Debug("Skip non-shared namespace", esl.Any("namespace", n))

		default:
			z.teamOwnedNamespaces[n.NamespaceId] = true
			teamOwnedNamespaceWithName[n.NamespaceId] = n.Name
		}
	}

	l.Debug("Enqueue to metadata scan")
	for id, name := range teamOwnedNamespaceWithName {
		z.scanner.Scan(z.ctl, z.ctx.AsAdminId(admin.TeamMemberId), name, id)
	}

	l.Debug("Metadata of teams finished")
	return nil
}

func (z *TeamScanner) namespaceOfMember(member *mo_member.Member) error {
	z.queue.Enqueue(&MemberScannerWorker{
		Member:              member,
		Control:             z.ctl,
		Context:             z.ctx.AsMemberId(member.TeamMemberId),
		TeamOwnedNamespaces: z.teamOwnedNamespaces,
		Scanner:             z.scanner,
		Folder:              z.filter,
	})
	return nil
}

func (z *TeamScanner) iterateMembers(f func(member *mo_member.Member) error) error {
	l := z.ctl.Log()

	if z.teamOwnedNamespaces == nil {
		l.Debug("Team owned namespaces is not initialized")
		return ErrorTeamOwnedNamespaceIsNotInitialized
	}

	l.Debug("Scanning members")
	members, err := sv_member.New(z.ctx).List()
	if err != nil {
		return err
	}

	for _, member := range members {
		if err := f(member); err != nil {
			return err
		}
	}
	return nil
}

func (z *TeamScanner) namespacesOfMembers() error {
	return z.iterateMembers(z.namespaceOfMember)
}

func (z *TeamScanner) Scan() error {
	if err := z.namespacesOfTeam(); err != nil {
		return err
	}
	if err := z.namespacesOfMembers(); err != nil {
		return err
	}
	return nil
}
