package cmd_migrate

import (
	"flag"
	"github.com/watermint/toolbox/app/app_report"
	"github.com/watermint/toolbox/cmd"
	"github.com/watermint/toolbox/domain/infra/api_auth_impl"
	"github.com/watermint/toolbox/domain/usecase/uc_team_migration"
)

type CmdTeamMigrateTransfer struct {
	*cmd.SimpleCommandlet
	report          app_report.Factory
	optSrcTeamAlias string
	optDstTeamAlias string
	optResume       string
}

func (z *CmdTeamMigrateTransfer) Name() string {
	return "transfer"
}

func (z *CmdTeamMigrateTransfer) Desc() string {
	return "cmd.team.migrate.transfer.desc"
}

func (z *CmdTeamMigrateTransfer) Usage() func(cmd.CommandUsage) {
	return nil
}

func (z *CmdTeamMigrateTransfer) FlagConfig(f *flag.FlagSet) {
	z.report.ExecContext = z.ExecContext
	z.report.FlagConfig(f)

	descFromAccount := z.ExecContext.Msg("cmd.teamfolder.mirror.flag.src_account").T()
	f.StringVar(&z.optSrcTeamAlias, "alias-src", "migration-src", descFromAccount)

	descToAccount := z.ExecContext.Msg("cmd.teamfolder.mirror.flag.dst_account").T()
	f.StringVar(&z.optDstTeamAlias, "alias-dest", "migration-dst", descToAccount)

	descResume := z.ExecContext.Msg("cmd.team.migrate.content.flag.resume").T()
	f.StringVar(&z.optResume, "resume", "", descResume)
}

func (z *CmdTeamMigrateTransfer) Exec(args []string) {
	var err error

	// Ask for SRC account authentication
	z.ExecContext.Msg("cmd.teamfolder.mirror.prompt.ask_src_file_account_auth").WithData(struct {
		Alias string
	}{
		Alias: z.optSrcTeamAlias,
	}).Tell()
	ctxFileSrc, err := api_auth_impl.Auth(z.ExecContext, api_auth_impl.PeerName(z.optSrcTeamAlias), api_auth_impl.BusinessFile())
	if err != nil {
		return
	}

	// Ask for SRC account authentication
	z.ExecContext.Msg("cmd.teamfolder.mirror.prompt.ask_src_mgmt_account_auth").WithData(struct {
		Alias string
	}{
		Alias: z.optSrcTeamAlias,
	}).Tell()
	ctxMgtSrc, err := api_auth_impl.Auth(z.ExecContext, api_auth_impl.PeerName(z.optSrcTeamAlias), api_auth_impl.BusinessManagement())
	if err != nil {
		return
	}

	// Ask for DST account authentication
	z.ExecContext.Msg("cmd.teamfolder.mirror.prompt.ask_dst_file_account_auth").WithData(struct {
		Alias string
	}{
		Alias: z.optDstTeamAlias,
	}).Tell()
	ctxFileDst, err := api_auth_impl.Auth(z.ExecContext, api_auth_impl.PeerName(z.optDstTeamAlias), api_auth_impl.BusinessFile())
	if err != nil {
		return
	}

	// Ask for DST account authentication
	z.ExecContext.Msg("cmd.teamfolder.mirror.prompt.ask_dst_mgmt_account_auth").WithData(struct {
		Alias string
	}{
		Alias: z.optDstTeamAlias,
	}).Tell()
	ctxMgtDst, err := api_auth_impl.Auth(z.ExecContext, api_auth_impl.PeerName(z.optDstTeamAlias), api_auth_impl.BusinessManagement())
	if err != nil {
		return
	}
	ucm := uc_team_migration.New(z.ExecContext, ctxFileSrc, ctxMgtSrc, ctxFileDst, ctxMgtDst)

	mc, err := ucm.Resume(uc_team_migration.ResumeExecContext(z.ExecContext), uc_team_migration.ResumeFromPath(z.optResume))
	if err != nil {
		return
	}
	if err = ucm.Migrate(mc); err != nil {
		ctxFileSrc.ErrorMsg(err).TellError()
	}
}
