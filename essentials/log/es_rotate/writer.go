package es_rotate

import (
	"errors"
	"github.com/watermint/toolbox/essentials/concurrency/es_mutex"
	"github.com/watermint/toolbox/essentials/log/es_log"
	"io"
	"os"
)

const (
	logFileExtension = ".log"
)

var (
	ErrorLogFileNotAvailable = errors.New("log file is not available")
)

type Writer interface {
	io.WriteCloser
	Open(opts ...RotateOpt) error
	UpdateOpt(opt RotateOpt)
}

func NewWriter(basePath, baseName string) Writer {
	w := writerImpl{
		m: es_mutex.New(),
	}
	w.ro = RotateOpts{}.Apply(
		BasePath(basePath),
		BaseName(baseName),
	)
	Startup()
	return &w
}

// This implementation is fully stateful
type writerImpl struct {
	current *os.File
	ro      RotateOpts
	written int64
	m       es_mutex.Mutex
}

func (z *writerImpl) UpdateOpt(opt RotateOpt) {
	z.ro = z.ro.Apply(opt)
}

func (z *writerImpl) Write(p []byte) (n int, err error) {
	z.m.Do(func() {
		if z.current == nil {
			err = ErrorLogFileNotAvailable
			return
		}
		n, err = z.current.Write(p)
		z.written += int64(n)
		if err != nil {
			return
		}
		if z.ro.ChunkSize() < z.written {
			err = z.rotate()
		}
		return
	})
	return
}

// this function must called from caller who owns mutex lock
func (z *writerImpl) rotate() (err error) {
	// close current log
	err = z.closeCurrent()
	if err != nil {
		return
	}

	// create new current log
	err = z.createCurrent()
	if err != nil {
		return
	}

	// Enqueue rotation
	rotate(MsgRotate{
		Opts: z.ro,
	})

	return
}

// this function must called from caller who owns mutex lock
func (z *writerImpl) createCurrent() (err error) {
	l := es_log.Default()
	path := z.ro.CurrentPath()

	l.Debug("create", es_log.String("path", path))
	z.current, err = os.Create(path)
	return
}

// this function must called from caller who owns mutex lock
func (z *writerImpl) closeCurrent() (err error) {
	l := es_log.Default().With(es_log.Int64("written", z.written))

	// flush written bytes
	z.written = 0

	// return on the file already closed
	if z.current == nil {
		l.Debug("already closed")
		return nil
	}

	// close
	name := z.current.Name()
	err = z.current.Close()

	l.Debug("Close", es_log.String("name", name), es_log.Error(err))
	if err != nil {
		l.Error("Unable to close current log file", es_log.String("path", name), es_log.Error(err))
	}
	z.current = nil

	rotateOut(MsgOut{
		Path: name,
		Opts: z.ro,
	})

	return
}

func (z *writerImpl) Close() (err error) {
	z.m.Do(func() {
		err = z.closeCurrent()
	})
	return
}

func (z *writerImpl) Open(opts ...RotateOpt) (err error) {
	z.m.Do(func() {
		// Do nothing when it already opened
		if z.current != nil {
			return
		}
		z.ro = z.ro.Apply(opts...)
		err = z.createCurrent()
	})
	return
}