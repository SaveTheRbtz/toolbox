package app_report

import (
	"errors"
	"flag"
	"fmt"
	"github.com/watermint/toolbox/app"
	"github.com/watermint/toolbox/app/app_report/app_report_csv"
	"github.com/watermint/toolbox/app/app_report/app_report_json"
	"github.com/watermint/toolbox/app/app_report/app_report_xlsx"
	"github.com/watermint/toolbox/app/app_ui"
	"github.com/watermint/toolbox/app/app_util"
	"io"
	"os"
	"path/filepath"
	"runtime"
)

type Report interface {
	Init(ec *app.ExecContext) error
	Close()
	Report(row interface{}) error
}

type Factory struct {
	ExecContext   *app.ExecContext
	reports       []Report
	DefaultWriter io.Writer
	wrapper       *app_util.LineWriter
	Path          string
	Suppress      bool
}

func (z *Factory) FlagConfig(f *flag.FlagSet) {
	descReportPath := z.ExecContext.Msg("report.common.flag.report_path").T()
	f.StringVar(&z.Path, "report-path", filepath.Join(z.ExecContext.JobsPath(), "reports"), descReportPath)
}

func (z *Factory) Init(ec *app.ExecContext) error {
	var consoleWriter io.Writer
	if runtime.GOOS == "windows" {
		consoleWriter = os.Stdout
	} else {
		z.wrapper = app_util.NewLineWriter(func(line string) error {
			if line == "" {
				return nil
			}
			app_ui.ColorPrint(os.Stdout, "REPORT\t", app_ui.ColorBlue)
			fmt.Println(line)
			return nil
		})
		consoleWriter = z.wrapper
	}
	if z.DefaultWriter == nil {
		z.DefaultWriter = os.Stdout
	}
	if z.reports == nil {
		z.reports = make([]Report, 0)
		if !z.Suppress {
			z.reports = append(z.reports, &app_report_json.JsonReport{
				DefaultWriter: consoleWriter,
				ReportPath:    "",
			})
		}
		z.reports = append(z.reports, &app_report_json.JsonReport{
			DefaultWriter: z.DefaultWriter,
			ReportPath:    z.Path,
		})
		z.reports = append(z.reports, &app_report_csv.CsvReport{
			DefaultWriter: z.DefaultWriter,
			ReportPath:    z.Path,
			ReportHeader:  true,
			ReportUseBom:  false,
		})
		z.reports = append(z.reports, &app_report_xlsx.XlsxReport{
			ReportPath: z.Path,
		})

		for _, r := range z.reports {
			if err := r.Init(ec); err != nil {
				return err
			}
		}
	}
	return nil
}

func (z *Factory) Report(row interface{}) error {
	if z.reports == nil {
		z.ExecContext.Log().Fatal("open report before write")
		return errors.New("report was not opened")
	}

	for _, r := range z.reports {
		if err := r.Report(row); err != nil {
			return err
		}
	}
	return nil
}

func (z *Factory) Close() {
	if z.reports == nil {
		z.ExecContext.Log().Debug("Report already closed")
		return
	}
	for _, r := range z.reports {
		r.Close()
	}
	if z.wrapper != nil {
		z.wrapper.Flush()
	}

	if !z.ExecContext.Quiet && !z.Suppress {
		z.ExecContext.Msg("report.common.done.tell_location").WithData(struct {
			Path string
		}{
			z.Path,
		}).TellSuccess()
	}
}