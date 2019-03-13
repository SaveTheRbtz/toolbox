package cmd_namespace_member

import (
	"github.com/watermint/toolbox/cmd"
)

func NewCmdTeamNamespaceMember() cmd.Commandlet {
	return &cmd.CommandletGroup{
		CommandName: "member",
		CommandDesc: "cmd.team.namespace.member.desc",
		SubCommands: []cmd.Commandlet{
			&CmdTeamNamespaceMemberList{
				SimpleCommandlet: &cmd.SimpleCommandlet{},
			},
		},
	}
}