package cmd_migrate

import (
	"flag"
	"github.com/watermint/toolbox/app/app_io"
	"github.com/watermint/toolbox/app/app_report"
	"github.com/watermint/toolbox/cmd"
	"github.com/watermint/toolbox/domain/infra/api_auth_impl"
	"github.com/watermint/toolbox/domain/usecase/uc_team_migration"
	"go.uber.org/zap"
)

type CmdTeamMigratePermission struct {
	*cmd.SimpleCommandlet
	report          app_report.Factory
	optSrcTeamAlias string
	optDstTeamAlias string
	optResume       string
	optMappingCsv   string
}

func (z *CmdTeamMigratePermission) Name() string {
	return "permission"
}

func (z *CmdTeamMigratePermission) Desc() string {
	return "cmd.team.migrate.permission.desc"
}

func (z *CmdTeamMigratePermission) Usage() func(cmd.CommandUsage) {
	return nil
}

func (z *CmdTeamMigratePermission) FlagConfig(f *flag.FlagSet) {
	z.report.ExecContext = z.ExecContext
	z.report.FlagConfig(f)

	descFromAccount := z.ExecContext.Msg("cmd.teamfolder.mirror.flag.src_account").T()
	f.StringVar(&z.optSrcTeamAlias, "alias-src", "migration-src", descFromAccount)

	descToAccount := z.ExecContext.Msg("cmd.teamfolder.mirror.flag.dst_account").T()
	f.StringVar(&z.optDstTeamAlias, "alias-dest", "migration-dst", descToAccount)

	descResume := z.ExecContext.Msg("cmd.team.migrate.content.flag.resume").T()
	f.StringVar(&z.optResume, "resume", "", descResume)

	descMappingCsv := "Mapping of pre/post migration email addresses of an account"
	f.StringVar(&z.optMappingCsv, "mapping-csv", "", descMappingCsv)
}

func (z *CmdTeamMigratePermission) Exec(args []string) {
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

	ucm := uc_team_migration.New(z.ExecContext, ctxFileSrc, ctxMgtSrc, ctxFileDst, ctxMgtDst, &z.report)

	opts := make([]uc_team_migration.PermOpt, 0)
	if z.optMappingCsv != "" {
		emailMapping := make(map[string]string)
		err := app_io.NewCsvLoader(z.ExecContext, z.optMappingCsv).OnRow(func(cols []string) error {
			if len(cols) < 2 {
				z.Log().Warn("Not enough column in the row", zap.Strings("cols", cols))
				return nil
			}

			oldEmail := cols[0]
			newEmail := cols[1]
			emailMapping[oldEmail] = newEmail

			z.Log().Info("Mapping", zap.String("oldEmail", oldEmail), zap.String("newEmail", newEmail))
			return nil
		}).Load()
		if err != nil {
			return
		}
		opts = append(opts, uc_team_migration.PermWithEmailMapping(emailMapping))
	}

	z.report.Init(z.ExecContext)
	defer z.report.Close()

	mc, err := ucm.Resume(uc_team_migration.ResumeExecContext(z.ExecContext), uc_team_migration.ResumeFromPath(z.optResume))
	if err != nil {
		return
	}
	if err = ucm.Permissions(mc, opts...); err != nil {
		ctxFileSrc.ErrorMsg(err).TellError()
	}
}