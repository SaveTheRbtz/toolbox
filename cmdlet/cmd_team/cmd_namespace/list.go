package cmd_namespace

import (
	"flag"
	"github.com/watermint/toolbox/cmdlet"
	"github.com/watermint/toolbox/dbx_api"
	"github.com/watermint/toolbox/infra"
)

type CmdTeamNamespaceList struct {
	*cmdlet.SimpleCommandlet

	apiContext *dbx_api.Context
	report     cmdlet.Report
}

func (CmdTeamNamespaceList) Name() string {
	return "list"
}

func (CmdTeamNamespaceList) Desc() string {
	return "List all namespaces of the team"
}

func (CmdTeamNamespaceList) Usage() string {
	return ""
}

func (c *CmdTeamNamespaceList) FlagConfig(f *flag.FlagSet) {
	c.report.FlagConfig(f)
}

func (c *CmdTeamNamespaceList) Exec(ec *infra.ExecContext, args []string) {
	if err := ec.Startup(); err != nil {
		return
	}
	defer ec.Shutdown()

	//apiMgmt, err := ec.LoadOrAuthBusinessFile()
	//if err != nil {
	//	return
	//}

	//rt, rs, err := c.report.ReportStages()
	//if err != nil {
	//	return
	//}
	//
	//wkNamespaceList := &team.WorkerTeamNamespaceList{
	//	Api:      apiMgmt,
	//	NextTask: rt,
	//}
	//
	//stages := []workflow.Worker{
	//	wkNamespaceList,
	//}
	//
	//stages = append(stages, rs...)
	//
	//p := workflow.Pipeline{
	//	Infra:  ec,
	//	Stages: stages,
	//}
	//
	//p.Init()
	//defer p.Close()
	//
	//p.Enqueue(
	//	workflow.MarshalTask(
	//		wkNamespaceList.Prefix(),
	//		wkNamespaceList.Prefix(),
	//		nil,
	//	),
	//)
	//p.Loop()
}
