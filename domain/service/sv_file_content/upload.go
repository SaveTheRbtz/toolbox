package sv_file_content

import (
	"github.com/watermint/toolbox/domain/model/mo_file"
	"github.com/watermint/toolbox/domain/model/mo_path"
	"github.com/watermint/toolbox/infra/api/api_context"
	"github.com/watermint/toolbox/infra/api/api_util"
	"go.uber.org/zap"
	"io"
	"os"
	"path/filepath"
)

type Upload interface {
	Add(destPath mo_path.Path, filePath string) (entry mo_file.Entry, err error)
	Overwrite(destPath mo_path.Path, filePath string) (entry mo_file.Entry, err error)
	Update(destPath mo_path.Path, filePath string, revision string) (entry mo_file.Entry, err error)
}

type UploadOpt func(o *UploadOpts) *UploadOpts
type UploadOpts struct {
	ChunkSize int64
	Mute      bool
}

const (
	DefaultChunkSize = 150 * 1048576
)

func NewUpload(ctx api_context.Context, opts ...UploadOpt) Upload {
	uo := &UploadOpts{
		ChunkSize: DefaultChunkSize,
		Mute:      false,
	}
	for _, o := range opts {
		o(uo)
	}
	return &uploadImpl{
		ctx: ctx,
		uo:  uo,
	}
}
func ChunkSize(chunkSize int64) UploadOpt {
	return func(o *UploadOpts) *UploadOpts {
		o.ChunkSize = chunkSize
		return o
	}
}

type uploadImpl struct {
	ctx api_context.Context
	uo  *UploadOpts
}

func (z *uploadImpl) Add(destPath mo_path.Path, filePath string) (entry mo_file.Entry, err error) {
	return z.upload(destPath, filePath, "add", "")
}

func (z *uploadImpl) Overwrite(destPath mo_path.Path, filePath string) (entry mo_file.Entry, err error) {
	return z.upload(destPath, filePath, "overwrite", "")
}

func (z *uploadImpl) Update(destPath mo_path.Path, filePath string, revision string) (entry mo_file.Entry, err error) {
	return z.upload(destPath, filePath, "update", revision)
}

func (z *uploadImpl) upload(destPath mo_path.Path, filePath string, mode string, revision string) (entry mo_file.Entry, err error) {
	info, err := os.Lstat(filePath)
	if err != nil {
		return nil, err
	}
	if info.Size() < z.uo.ChunkSize {
		return z.uploadSingle(info, destPath, filePath, mode, revision)
	} else {
		return z.uploadChunked(info, destPath, filePath, mode, revision)
	}
}

type UploadParamMode struct {
	Tag    string `json:".tag"`
	Update string `json:"update,omitempty"`
}

type UploadParams struct {
	Path           string           `json:"path"`
	Mode           *UploadParamMode `json:"mode"`
	Mute           bool             `json:"mute"`
	ClientModified string           `json:"client_modified"`
	Autorename     bool             `json:"autorename"`
}

func UploadPath(destPath mo_path.Path, f os.FileInfo) mo_path.Path {
	return destPath.ChildPath(filepath.Base(f.Name()))
}

func (z *uploadImpl) makeParams(info os.FileInfo, destPath mo_path.Path, mode string, revision string) *UploadParams {
	upm := &UploadParamMode{
		Tag:    mode,
		Update: "",
	}
	up := &UploadParams{
		Path:           UploadPath(destPath, info).Path(),
		Mode:           upm,
		Mute:           false,
		ClientModified: api_util.RebaseAsString(info.ModTime()),
	}
	switch mode {
	case "update":
		upm.Update = revision
	case "add":
		up.Autorename = true
	}
	return up
}

func (z *uploadImpl) uploadSingle(info os.FileInfo, destPath mo_path.Path, filePath string, mode string, revision string) (entry mo_file.Entry, err error) {
	l := z.ctx.Log().With(zap.String("filePath", filePath), zap.Int64("size", info.Size()))
	l.Debug("Uploading file")

	r, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	res, err := z.ctx.Upload("files/upload", r).
		Param(z.makeParams(info, destPath, mode, revision)).Call()
	if err != nil {
		return nil, err
	}
	entry = &mo_file.Metadata{}
	if err := res.Model(entry); err != nil {
		return nil, err
	}
	return entry, nil
}

func (z *uploadImpl) uploadChunked(info os.FileInfo, destPath mo_path.Path, filePath string, mode string, revision string) (entry mo_file.Entry, err error) {
	l := z.ctx.Log().With(zap.String("filePath", filePath), zap.Int64("size", info.Size()))

	total := info.Size()
	var written int64
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	r := io.LimitReader(f, z.uo.ChunkSize)

	type SessionId struct {
		SessionId string `json:"session_id"`
	}
	type CursorInfo struct {
		SessionId string `json:"session_id"`
		Offset    int64  `json:"offset"`
	}
	type AppendInfo struct {
		Cursor *CursorInfo `json:"cursor"`
	}
	type CommitInfo struct {
		Cursor *CursorInfo   `json:"cursor"`
		Commit *UploadParams `json:"commit"`
	}

	l.Debug("Upload session start")
	res, err := z.ctx.Upload("files/upload_session/start", r).Call()
	if err != nil {
		return nil, err
	}
	sid := &SessionId{}
	if j, err := res.Json(); err != nil {
		return nil, err
	} else {
		sid.SessionId = j.Get("session_id").String()
	}
	written += z.uo.ChunkSize
	l = l.With(zap.String("sessionId", sid.SessionId))

	for (total - written) > z.uo.ChunkSize {
		l.Debug("Append chunk", zap.Int64("written", written))
		ai := &AppendInfo{
			Cursor: &CursorInfo{
				SessionId: sid.SessionId,
				Offset:    written,
			},
		}
		r = io.LimitReader(f, z.uo.ChunkSize)
		_, err := z.ctx.Upload("files/upload_session/append_v2", r).Param(ai).Call()
		if err != nil {
			return nil, err
		}
		written += z.uo.ChunkSize
	}

	l.Debug("Finish")
	ci := &CommitInfo{
		Cursor: &CursorInfo{
			SessionId: sid.SessionId,
			Offset:    written,
		},
		Commit: z.makeParams(info, destPath, mode, revision),
	}
	res, err = z.ctx.Upload("files/upload_session/finish", f).Param(ci).Call()
	if err != nil {
		return nil, err
	}
	entry = &mo_file.Metadata{}
	if err := res.Model(entry); err != nil {
		return nil, err
	}
	return entry, nil
}
