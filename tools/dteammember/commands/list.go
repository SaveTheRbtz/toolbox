package commands

import (
	"errors"
	"flag"
	"github.com/cihub/seelog"
	"github.com/watermint/toolbox/infra"
	"github.com/watermint/toolbox/infra/util"
	"github.com/watermint/toolbox/service/members"
	"github.com/watermint/toolbox/service/report"
	"os"
)

type ListOptions struct {
	Infra  *infra.InfraOpts
	Report *report.ReportOpts
	Status string
}

func parseListOptions(args []string) (*ListOptions, error) {
	f := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	opts := &ListOptions{}
	opts.Infra = infra.PrepareInfraFlags(f)
	opts.Report = report.PrepareReportFlags(f)

	descStatus := "Status (all|active|invited|suspended|removed)"
	f.StringVar(&opts.Status, "status", "all", descStatus)

	f.SetOutput(os.Stderr)
	f.Parse(args)

	switch opts.Status {
	case "all":
	case "active":
	case "invited":
	case "suspended":
	case "removed":
	default:
		return nil, errors.New("Invalid status: " + opts.Status)
	}

	return opts, nil
}

func List(args []string) error {
	opts, err := parseListOptions(args)
	if err != nil {
		return err
	}

	defer opts.Infra.Shutdown()
	err = opts.Infra.Startup()
	if err != nil {
		seelog.Errorf("Unable to start operation: %s", err)
		return err
	}
	seelog.Tracef("options: %s", util.MarshalObjectToString(opts))

	token, err := opts.Infra.LoadOrAuthBusinessInfo()
	if err != nil || token == "" {
		seelog.Errorf("Unable to acquire token (error: %s)", err)
		return err
	}

	rows := make(chan report.ReportRow)
	opts.Report.Write(rows)

	err = members.ListMembers(token, rows, opts.Status)
	if err != nil {
		seelog.Errorf("Unable to load members: error[%s]", err)
		return err
	}

	opts.Report.Wait()

	return nil
}