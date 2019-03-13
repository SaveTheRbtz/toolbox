package cmd_audit

import (
	"flag"
	"github.com/tidwall/gjson"
	"github.com/watermint/toolbox/cmd"
	"github.com/watermint/toolbox/model/dbx_api"
	"github.com/watermint/toolbox/model/dbx_auth"
	"github.com/watermint/toolbox/model/dbx_group"
	"github.com/watermint/toolbox/model/dbx_member"
	"github.com/watermint/toolbox/model/dbx_namespace"
	"github.com/watermint/toolbox/model/dbx_profile"
	"github.com/watermint/toolbox/model/dbx_sharing"
	"github.com/watermint/toolbox/model/dbx_team"
	"github.com/watermint/toolbox/report"
	"go.uber.org/zap"
)

type CmdTeamAuditSharing struct {
	*cmd.SimpleCommandlet
	groupMembers   map[string][]*dbx_group.GroupMember
	report         report.Factory
	optExpandGroup bool
}

func (CmdTeamAuditSharing) Name() string {
	return "sharing"
}

func (CmdTeamAuditSharing) Desc() string {
	return "cmd.team.audit.sharing.desc"
}

func (z *CmdTeamAuditSharing) Usage() func(usage cmd.CommandUsage) {
	return nil
}

func (z *CmdTeamAuditSharing) FlagConfig(f *flag.FlagSet) {
	z.report.ExecContext = z.ExecContext
	z.report.FlagConfig(f)

	descExpandGroup := z.ExecContext.Msg("cmd.team.audit.sharing.flag.expand_group").T()
	f.BoolVar(&z.optExpandGroup, "expand-group", false, descExpandGroup)
}

func (z *CmdTeamAuditSharing) Exec(args []string) {
	au := dbx_auth.NewDefaultAuth(z.ExecContext)
	apiFile, err := au.Auth(dbx_auth.DropboxTokenBusinessFile)
	if err != nil {
		return
	}
	z.report.Init(z.ExecContext)
	defer z.report.Close()

	// Identify admin
	z.ExecContext.Msg("cmd.team.audit.sharing.progress.identify_admin").Tell()
	admin, err := dbx_profile.AuthenticatedAdmin(apiFile)
	if err != nil {
		z.DefaultErrorHandler(err)
		return
	}
	z.ExecContext.Msg("cmd.team.audit.sharing.progress.do_as_admin").WithData(struct {
		Email string
	}{
		Email: admin.Email,
	}).Tell()
	z.Log().Info("Execute scan as admin", zap.String("email", admin.Email))

	// Scan shared links
	z.ExecContext.Msg("cmd.team.audit.sharing.progress.shared_link").Tell()
	z.Log().Info("Scanning Shared links")
	if !z.reportSharedLink(apiFile) {
		return
	}

	// Scan team info
	z.ExecContext.Msg("cmd.team.audit.sharing.progress.team_info").Tell()
	z.Log().Info("Scanning Team Info")
	if !z.reportInfo(apiFile) {
		return
	}

	// Scan team feature
	z.ExecContext.Msg("cmd.team.audit.sharing.progress.team_feature").Tell()
	z.Log().Info("Scanning Team Feature")
	if !z.reportFeature(apiFile) {
		return
	}

	// Scan Team members
	z.ExecContext.Msg("cmd.team.audit.sharing.progress.members").Tell()
	z.Log().Info("Scanning Team Members")
	if !z.reportMember(apiFile) {
		return
	}

	// Scan Team group
	z.ExecContext.Msg("cmd.team.audit.sharing.progress.groups").Tell()
	z.Log().Info("Scanning Team Group")
	if !z.reportGroup(apiFile) {
		return
	}

	// Scan Team group members
	z.ExecContext.Msg("cmd.team.audit.sharing.progress.group_members").Tell()
	z.Log().Info("Scanning Team Group Member")
	if !z.reportGroupMember(apiFile) {
		return
	}

	// Scan namespaces
	z.ExecContext.Msg("cmd.team.audit.sharing.progress.namespace").Tell()
	z.Log().Info("Scanning Namespace")
	if !z.reportNamespace(apiFile) {
		return
	}

	// Expand group
	if z.optExpandGroup {
		z.ExecContext.Msg("cmd.team.audit.sharing.progress.expand_group").Tell()
		z.Log().Info("Preparing for `-expand-group`")
		z.groupMembers = dbx_group.GroupMembers(apiFile, z.Log(), z.DefaultErrorHandler)
		if z.groupMembers == nil {
			z.ExecContext.Msg("cmd.team.audit.sharing.err.fail_prepare_expand_group").TellError()
			z.Log().Warn("Unable to list group members")
		}
	}

	// Scan namespace members
	z.ExecContext.Msg("cmd.team.audit.sharing.progress.namespace_members").Tell()
	z.Log().Info("Scanning Namespace members")
	if !z.reportNamespaceMember(apiFile, admin) {
		return
	}

	// Scan namespace files
	z.ExecContext.Msg("cmd.team.audit.sharing.progress.namespace_file").Tell()
	z.Log().Info("Scanning Namespace files")
	if !z.reportNamespaceFile(apiFile, admin) {
		return
	}
}

func (z *CmdTeamAuditSharing) reportInfo(c *dbx_api.Context) bool {
	l := dbx_team.TeamInfoList{
		OnError: z.DefaultErrorHandler,
		OnEntry: func(info *dbx_team.TeamInfo) bool {
			z.report.Report(info)
			return true
		},
	}
	return l.List(c)
}

func (z *CmdTeamAuditSharing) reportFeature(c *dbx_api.Context) bool {
	l := dbx_team.FeatureList{
		OnError: z.DefaultErrorHandler,
		OnEntry: func(feature *dbx_team.Feature) bool {
			z.report.Report(feature)
			return true
		},
	}
	return l.List(c)
}

func (z *CmdTeamAuditSharing) reportMember(c *dbx_api.Context) bool {
	l := dbx_member.MembersList{
		OnError: z.DefaultErrorHandlerIgnoreError,
		OnEntry: func(member *dbx_profile.Member) bool {
			z.report.Report(member)
			return true
		},
	}
	return l.List(c, true)
}

func (z *CmdTeamAuditSharing) reportGroup(c *dbx_api.Context) bool {
	gl := dbx_group.GroupList{
		OnError: z.DefaultErrorHandler,
		OnEntry: func(group *dbx_group.Group) bool {
			z.report.Report(group)
			return true
		},
	}
	return gl.List(c)
}

func (z *CmdTeamAuditSharing) reportGroupMember(c *dbx_api.Context) bool {
	gl := dbx_group.GroupList{
		OnError: z.DefaultErrorHandler,
		OnEntry: func(group *dbx_group.Group) bool {

			gml := dbx_group.GroupMemberList{
				OnError: z.DefaultErrorHandler,
				OnEntry: func(gm *dbx_group.GroupMember) bool {
					z.report.Report(gm)
					return true
				},
			}
			gml.List(c, group)

			return true
		},
	}
	return gl.List(c)
}

func (z *CmdTeamAuditSharing) reportSharedLink(c *dbx_api.Context) bool {
	ml := dbx_member.MembersList{
		OnError: z.DefaultErrorHandlerIgnoreError,
		OnEntry: func(member *dbx_profile.Member) bool {
			sl := dbx_sharing.SharedLinkList{
				AsMemberId:    member.Profile.TeamMemberId,
				AsMemberEmail: member.Profile.Email,
				OnError:       z.DefaultErrorHandler,
				OnEntry: func(link *dbx_sharing.SharedLink) bool {
					z.report.Report(link)
					return true
				},
			}
			sl.List(c)
			return true
		},
	}
	return ml.List(c, false)
}

func (z *CmdTeamAuditSharing) reportNamespace(c *dbx_api.Context) bool {
	l := dbx_namespace.NamespaceList{
		OnError: z.DefaultErrorHandler,
		OnEntry: func(namespace *dbx_namespace.Namespace) bool {
			z.report.Report(namespace)
			return true
		},
	}
	return l.List(c)
}

func (z *CmdTeamAuditSharing) reportNamespaceMember(c *dbx_api.Context, admin *dbx_profile.Profile) bool {
	l := dbx_namespace.NamespaceList{
		OnError: z.DefaultErrorHandler,
		OnEntry: func(namespace *dbx_namespace.Namespace) bool {
			if namespace.NamespaceType != "shared_folder" &&
				namespace.NamespaceType != "team_folder" {
				return true
			}

			sl := dbx_sharing.SharedFolderMembers{
				AsAdminId: admin.TeamMemberId,
				OnError:   z.DefaultErrorHandler,
				OnUser: func(user *dbx_sharing.MembershipUser) bool {
					nu := &dbx_namespace.NamespaceUser{
						Namespace: namespace,
						User:      user,
					}
					z.report.Report(nu)
					return true
				},
				OnGroup: func(group *dbx_sharing.MembershipGroup) bool {
					if z.optExpandGroup {
						if gmm, ok := z.groupMembers[group.Group.GroupId]; ok {
							for _, gm := range gmm {
								nu := &dbx_namespace.NamespaceUser{
									Namespace: namespace,
									User: &dbx_sharing.MembershipUser{
										Membership: group.Membership,
										User: &dbx_sharing.User{
											UserAccountId: gm.Profile.AccountId,
											Email:         gm.Profile.Email,
											DisplayName:   gjson.Get(string(gm.Profile.Profile), "name.display_name").String(),
											SameTeam:      true,
											TeamMemberId:  gm.TeamMemberId,
										},
									},
								}
								z.report.Report(nu)
							}
						} else {
							z.ExecContext.Msg("cmd.team.audit.sharing.err.fail_expand_group").WithData(struct {
								GroupId   string
								GroupName string
							}{
								GroupId:   group.Group.GroupId,
								GroupName: group.Group.GroupName,
							}).TellError()
							z.Log().Warn(
								"Could not expand group",
								zap.String("group_id", group.Group.GroupId),
								zap.String("group_name", group.Group.GroupName),
							)
							ng := &dbx_namespace.NamespaceGroup{
								Namespace: namespace,
								Group:     group,
							}
							z.report.Report(ng)
						}
					} else {
						ng := &dbx_namespace.NamespaceGroup{
							Namespace: namespace,
							Group:     group,
						}
						z.report.Report(ng)
					}
					return true
				},
				OnInvitee: func(invitee *dbx_sharing.MembershipInvitee) bool {
					ni := &dbx_namespace.NamespaceInvitee{
						Namespace: namespace,
						Invitee:   invitee,
					}
					z.report.Report(ni)
					return true
				},
			}
			sl.List(c, namespace.NamespaceId)
			return true
		},
	}
	return l.List(c)
}

type NamespaceMembershipFileUser struct {
	File *dbx_namespace.NamespaceFile `json:"file"`
	User *dbx_sharing.MembershipUser  `json:"user"`
}
type NamespaceMembershipFileGroup struct {
	File  *dbx_namespace.NamespaceFile `json:"file"`
	Group *dbx_sharing.MembershipGroup `json:"group"`
}
type NamespaceMembershipFileInvitee struct {
	File    *dbx_namespace.NamespaceFile   `json:"file"`
	Invitee *dbx_sharing.MembershipInvitee `json:"invitee"`
}
type NamespaceMembershipError struct {
	NamespaceId      string `json:"namespace_id"`
	NamespaceOwnerId string `json:"namespace_owner_id,omitempty"`
	AsMemberId       string `json:"as_member_id,omitempty"`
	AsAdminId        string `json:"as_admin_id,omitempty"`
	FileId           string `json:"file_id"`
	FilePath         string `json:"file_path"`
}

func (z *CmdTeamAuditSharing) reportNamespaceFile(c *dbx_api.Context, admin *dbx_profile.Profile) bool {
	fileSharing := func(file *dbx_namespace.NamespaceFile) bool {
		lfm := dbx_sharing.SharedFileMembers{
			OnUser: func(user *dbx_sharing.MembershipUser) bool {
				r := NamespaceMembershipFileUser{
					File: file,
					User: user,
				}
				z.report.Report(r)
				return true
			},
			OnInvitee: func(invitee *dbx_sharing.MembershipInvitee) bool {
				r := NamespaceMembershipFileInvitee{
					File:    file,
					Invitee: invitee,
				}
				z.report.Report(r)
				return true
			},
			OnGroup: func(group *dbx_sharing.MembershipGroup) bool {
				gr := NamespaceMembershipFileGroup{
					File:  file,
					Group: group,
				}

				if !z.optExpandGroup {
					z.report.Report(gr)
					return true
				}

				if gmm, ok := z.groupMembers[group.Group.GroupId]; ok {
					for _, gm := range gmm {
						nu := &NamespaceMembershipFileUser{
							File: file,
							User: &dbx_sharing.MembershipUser{
								Membership: group.Membership,
								User: &dbx_sharing.User{
									UserAccountId: gm.Profile.AccountId,
									Email:         gm.Profile.Email,
									DisplayName:   gjson.Get(string(gm.Profile.Profile), "name.display_name").String(),
									SameTeam:      true,
									TeamMemberId:  gm.TeamMemberId,
								},
							},
						}
						z.report.Report(nu)
					}
				} else {
					z.ExecContext.Msg("cmd.team.audit.sharing.err.fail_expand_group").WithData(struct {
						GroupId   string
						GroupName string
					}{
						GroupId:   group.Group.GroupId,
						GroupName: group.Group.GroupName,
					}).TellError()

					z.Log().Warn(
						"Could not expand group",
						zap.String("group_id", group.Group.GroupId),
						zap.String("group_name", group.Group.GroupName),
					)
					z.report.Report(gr)
				}

				return true
			},
		}
		lfm.OnError = func(err error) bool {
			nme := NamespaceMembershipError{
				NamespaceId:      file.Namespace.TeamMemberId,
				NamespaceOwnerId: file.Namespace.TeamMemberId,
				AsMemberId:       lfm.AsMemberId,
				AsAdminId:        lfm.AsAdminId,
				FileId:           file.File.FileId,
				FilePath:         file.File.PathLower,
			}
			// Out error report
			z.report.Report(nme)
			z.ExecContext.Msg("cmd.team.audit.sharing.err.unable_to_acquire_sharing_info").WithData(struct {
				NamespaceId      string
				NamespaceOwnerId string
				FileId           string
				FilePath         string
			}{
				NamespaceId:      file.Namespace.NamespaceId,
				NamespaceOwnerId: file.Namespace.TeamMemberId,
				FileId:           file.File.FileId,
				FilePath:         file.File.PathDisplay,
			}).TellError()
			z.Log().Warn(
				"Unable to acquire sharing information",
				zap.String("namespace_id", file.Namespace.NamespaceId),
				zap.String("namespace_owner_id", file.Namespace.TeamMemberId),
				zap.String("as_member_id", lfm.AsMemberId),
				zap.String("as_admin_id", lfm.AsAdminId),
				zap.String("file_id", file.File.FileId),
				zap.String("file_path", file.File.PathLower),
			)
			return z.DefaultErrorHandler(err)
		}

		if file.Namespace.TeamMemberId != "" {
			lfm.AsMemberId = file.Namespace.TeamMemberId
		} else {
			lfm.AsAdminId = admin.TeamMemberId
		}

		return lfm.List(c, file.File.FileId)
	}

	lns := dbx_namespace.ListNamespaceFile{}
	lns.OptIncludeDeleted = false
	lns.OptIncludeMediaInfo = false
	lns.OptIncludeAppFolder = true
	lns.OptIncludeMemberFolder = true
	lns.OptIncludeSharedFolder = true
	lns.OptIncludeTeamFolder = true
	lns.AsAdminId = admin.TeamMemberId
	lns.OnError = z.DefaultErrorHandler
	lns.OnNamespace = func(namespace *dbx_namespace.Namespace) bool {
		z.ExecContext.Msg("cmd.team.audit.sharing.progress.folder").WithData(struct {
			Type        string
			NamespaceId string
			Name        string
		}{
			Type:        namespace.NamespaceType,
			NamespaceId: namespace.NamespaceId,
			Name:        namespace.Name,
		}).Tell()
		z.Log().Info("Scanning folder",
			zap.String("namespace_type", namespace.NamespaceType),
			zap.String("namespace_id", namespace.NamespaceId),
			zap.String("name", namespace.Name),
		)
		return true
	}
	lns.OnFolder = func(folder *dbx_namespace.NamespaceFolder) bool {
		z.report.Report(folder)
		return true
	}
	lns.OnFile = func(file *dbx_namespace.NamespaceFile) bool {
		z.report.Report(file)

		if file.File.HasExplicitSharedMembers {
			fileSharing(file) // ignore error
			return true
		}
		return true
	}
	lns.OnDelete = func(deleted *dbx_namespace.NamespaceDeleted) bool {
		z.report.Report(deleted)
		return true
	}
	return lns.List(c)
}