package app_job_impl

import (
	"github.com/watermint/toolbox/essentials/log/esl"
	app2 "github.com/watermint/toolbox/infra/app"
	"github.com/watermint/toolbox/infra/control/app_job"
	"github.com/watermint/toolbox/infra/recipe/rc_recipe"
	"os"
	"path/filepath"
	"time"
)

func newOrphanHistory(path string) (h app_job.History, found bool) {
	l := esl.Default().With(esl.String("path", path))

	startLogPath := filepath.Join(path, app2.LogNameStart)
	startLogStat, err := os.Lstat(startLogPath)
	if err != nil || startLogStat.IsDir() {
		l.Debug("The path is not a history directory", esl.Error(err))
		return nil, false
	}
	start := &app_job.StartLog{}
	finish := &app_job.ResultLog{}
	finishLogPath := filepath.Join(path, app2.LogNameFinish)

	if err := parse(startLogPath, start); err != nil {
		l.Debug("Unable to load start log", esl.Error(err))
		return nil, false
	}
	if err := parse(finishLogPath, finish); err != nil {
		l.Debug("Unable to load finish log", esl.Error(err))
	}
	return &OrphanHistory{
		path:   path,
		start:  start,
		finish: finish,
	}, true
}

type OrphanHistory struct {
	path   string
	start  *app_job.StartLog
	finish *app_job.ResultLog
}

func (z OrphanHistory) IsOrphaned() bool {
	return true
}

func (z OrphanHistory) JobId() string {
	if z.start.JobId != "" {
		return z.start.JobId
	} else {
		return z.start.TimeStart
	}
}

func (z OrphanHistory) JobPath() string {
	return z.path
}

func (z OrphanHistory) RecipeName() string {
	return z.start.Name
}

func (z OrphanHistory) Recipe() (r rc_recipe.Spec, found bool) {
	return getRecipe(z.start.Name)
}

func (z OrphanHistory) AppName() string {
	return z.start.AppName
}

func (z OrphanHistory) AppVersion() string {
	return z.start.AppVersion
}

func (z OrphanHistory) TimeStart() (t time.Time, found bool) {
	return app_job.TimeFromLog(z.start, z.JobId())
}

func (z OrphanHistory) TimeFinish() (t time.Time, found bool) {
	return app_job.TimeFromLog(z.finish, "")
}

func (z OrphanHistory) IsNested() bool {
	return false
}

func (z OrphanHistory) Logs() (logs []app_job.LogFile, err error) {
	return getLogs(z.path)
}
