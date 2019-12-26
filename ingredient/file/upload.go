package file

import (
	"github.com/watermint/toolbox/domain/model/mo_file"
	"github.com/watermint/toolbox/domain/model/mo_path"
	"github.com/watermint/toolbox/domain/service/sv_file"
	"github.com/watermint/toolbox/domain/service/sv_file_content"
	"github.com/watermint/toolbox/domain/service/sv_file_folder"
	"github.com/watermint/toolbox/infra/api/api_context"
	"github.com/watermint/toolbox/infra/api/api_util"
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/infra/recipe/rc_kitchen"
	"github.com/watermint/toolbox/infra/report/rp_model"
	"github.com/watermint/toolbox/infra/ui/app_msg"
	"github.com/watermint/toolbox/infra/util/ut_filepath"
	"github.com/watermint/toolbox/quality/infra/qt_endtoend"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	statusReportInterval = 15 * time.Second
)

type Upload struct {
	Context         api_context.Context
	EstimateOnly    bool
	Overwrite       bool
	ChunkSizeKb     int
	CreateFolder    bool
	LocalPath       mo_path.FileSystemPath
	DropboxPath     mo_path.DropboxPath
	ProgressUpload  app_msg.Message
	ProgressSummary app_msg.Message
	Uploaded        rp_model.TransactionReport
	Skipped         rp_model.TransactionReport
	Summary         rp_model.RowReport
}

func (z *Upload) Preset() {
	z.Uploaded.SetModel(&UploadRow{}, &mo_file.ConcreteEntry{}, rp_model.HiddenColumns(
		"result.id", "result.tag",
	))
	z.Skipped.SetModel(&UploadRow{}, &mo_file.ConcreteEntry{}, rp_model.HiddenColumns(
		"result.id", "result.tag",
	))
	z.Summary.SetModel(&UploadSummary{})
	z.ChunkSizeKb = 150 * 1024
}

func (z *Upload) exec(k rc_kitchen.Kitchen, localPath string, dropboxPath string, estimate bool) (summary *UploadSummary, err error) {
	// TODO: refactor localPath to mo_path.FileSystemPath, and DropboxPath to mo_path.DropboxPath
	l := k.Log().With(zap.String("localPath", localPath), zap.String("dropboxPath", dropboxPath), zap.Bool("estimate", estimate))
	l.Debug("execute")

	status := &UploadStatus{
		summary: UploadSummary{
			UploadStart: time.Now(),
		},
	}

	go func() {
		for {
			time.Sleep(statusReportInterval)

			dur := time.Now().Sub(status.summary.UploadStart) / time.Second
			if dur == 0 {
				continue
			}

			kps := status.summary.NumBytes / int64(dur) / 1024

			k.UI().InfoM(z.ProgressSummary.
				With("Time", time.Now().Truncate(time.Second).Format("15:04:05")).
				With("NumFileUpload", status.summary.NumFilesUpload).
				With("NumFileSkip", status.summary.NumFilesSkip).
				With("NumFileError", status.summary.NumFilesError).
				With("NumBytes", status.summary.NumBytes/1_048_576).
				With("Kps", kps).
				With("NumApiCall", status.summary.NumApiCall))
		}
	}()

	l.Debug("upload", zap.Int("chunkSize", z.ChunkSizeKb))
	up := sv_file_content.NewUpload(z.Context, sv_file_content.ChunkSizeKb(z.ChunkSizeKb))
	q := k.NewQueue()

	info, err := os.Lstat(localPath)
	if err != nil {
		l.Debug("Unable to fetch info", zap.Error(err))
		return nil, err
	}

	createFolder := func(path string) error {
		ll := l.With(zap.String("path", path))
		ll.Debug("Prepare create folder")
		rel, err := ut_filepath.Rel(localPath, path)
		if err != nil {
			l.Debug("unable to calculate rel path", zap.Error(err))
			z.Uploaded.Failure(err, &UploadRow{File: path})
			status.error()
			return err
		}
		if rel == "." {
			ll.Debug("Skip")
			return nil
		}

		folderPath := mo_path.NewDropboxPath(dropboxPath).ChildPath(rel)
		ll = ll.With(zap.String("folderPath", folderPath.Path()), zap.String("rel", rel))
		ll.Debug("Create folder")

		entry, err := sv_file_folder.New(z.Context).Create(folderPath)
		if err != nil {
			if api_util.ErrorSummaryPrefix(err, "path/conflict/folder") {
				ll.Debug("The folder already exist, ignore it", zap.Error(err))
				return nil
			} else {
				ll.Debug("Unable to create folder", zap.Error(err))
				z.Uploaded.Failure(err, &UploadRow{File: path})
				return err
			}
		}
		z.Uploaded.Success(&UploadRow{File: path}, entry.Concrete())

		return nil
	}

	var scanFolder func(path string) error
	scanFolder = func(path string) error {
		ll := l.With(zap.String("path", path))

		ll.Debug("Scanning folder")
		localEntries, err := ioutil.ReadDir(path)
		if err != nil {
			ll.Debug("Unable to read dir", zap.Error(err))
			return err
		}
		localPathRel, err := ut_filepath.Rel(localPath, path)
		if err != nil {
			ll.Debug("Unable to calc rel path", zap.Error(err))
			return err
		}

		dbxPath := mo_path.NewDropboxPath(dropboxPath)
		if localPathRel != "." {
			dbxPath = dbxPath.ChildPath(localPathRel)
		}

		dbxEntries, err := sv_file.NewFiles(z.Context).List(dbxPath)
		if err != nil {
			if api_util.ErrorSummaryPrefix(err, "path/not_found") {
				ll.Debug("Dropbox entry not found", zap.String("dbxPath", dbxPath.Path()), zap.Error(err))
				dbxEntries = make([]mo_file.Entry, 0)
			} else {
				ll.Debug("Unable to read Dropbox entries", zap.String("dbxPath", dbxPath.Path()), zap.Error(err))
				return err
			}
		}
		dbxEntryByName := mo_file.MapByNameLower(dbxEntries)

		numEntriesProceed := 0
		var lastErr error
		for _, e := range localEntries {
			p := filepath.Join(path, e.Name())
			if api_util.IsFileNameIgnored(p) {
				ll.Debug("Ignore file", zap.String("p", p))
				var ps int64 = 0
				pi, err := os.Lstat(p)
				if err == nil {
					ps = pi.Size()
				}
				status.skip()
				z.Skipped.Skip(
					app_msg.M("usecase.uc_file_upload.skip.dont_sync"),
					UploadRow{
						File: p,
						Size: ps,
					})
				continue
			}
			numEntriesProceed++
			if e.IsDir() {
				lastErr = scanFolder(filepath.Join(path, e.Name()))
			} else {
				dbxEntry := dbxEntryByName[strings.ToLower(e.Name())]
				ll.Debug("Enqueue", zap.String("p", p))
				q.Enqueue(&UploadWorker{
					dropboxBasePath: dropboxPath,
					localBasePath:   localPath,
					localFilePath:   p,
					dbxEntry:        dbxEntry,
					ctx:             z.Context,
					ctl:             k.Control(),
					up:              up,
					estimateOnly:    estimate,
					status:          status,
					upload:          z,
				})
			}
		}
		l.Debug("folder scan finished", zap.Int("numEntriesProceed", numEntriesProceed), zap.Error(lastErr))
		if numEntriesProceed == 0 && z.CreateFolder {
			l.Debug("Create folder for empty folder")
			return createFolder(path)
		}
		return lastErr
	}

	var lastErr error
	if info.IsDir() {
		lastErr = scanFolder(localPath)
	} else {
		q.Enqueue(&UploadWorker{
			dropboxBasePath: dropboxPath,
			localBasePath:   localPath,
			localFilePath:   localPath,
			ctx:             z.Context,
			ctl:             k.Control(),
			up:              up,
			estimateOnly:    estimate,
			upload:          z,
			status:          status,
		})
	}

	q.Wait()

	status.summary.UploadEnd = time.Now()
	z.Summary.Row(&status.summary)
	return &status.summary, lastErr
}

func (z *Upload) Exec(k rc_kitchen.Kitchen) error {
	if err := z.Uploaded.Open(); err != nil {
		return err
	}
	if err := z.Skipped.Open(); err != nil {
		return err
	}
	if err := z.Summary.Open(); err != nil {
		return err
	}
	_, err := z.exec(k, z.LocalPath.Path(), z.DropboxPath.Path(), z.EstimateOnly)
	return err
}

func (z *Upload) Test(c app_control.Control) error {
	return qt_endtoend.ScenarioTest()
}