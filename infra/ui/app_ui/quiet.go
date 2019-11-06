package app_ui

import (
	"github.com/watermint/toolbox/infra/ui/app_msg"
	"github.com/watermint/toolbox/infra/ui/app_msg_container"
	"go.uber.org/zap"
)

func NewQuiet(container app_msg_container.Container) UI {
	q := container.(app_msg_container.Quality)
	return &Quiet{
		mc: container,
		mq: q,
	}
}

type Quiet struct {
	mc  app_msg_container.Container
	mq  app_msg_container.Quality
	log *zap.Logger
}

func (z *Quiet) Success(key string, p ...app_msg.P) {
	z.mq.Verify(key)
	z.log.Debug(key, zap.Any("params", p))
}

func (z *Quiet) Failure(key string, p ...app_msg.P) {
	z.mq.Verify(key)
	z.log.Debug(key, zap.Any("params", p))
}

func (z *Quiet) IsConsole() bool {
	return true
}

func (z *Quiet) IsWeb() bool {
	return false
}

func (z *Quiet) OpenArtifact(path string) {
	z.log.Debug("Open artifact", zap.String("path", path))
}

func (z *Quiet) Text(key string, p ...app_msg.P) string {
	z.mq.Verify(key)
	return z.mc.Compile(app_msg.M(key, p...))
}

func (z *Quiet) TextOrEmpty(key string, p ...app_msg.P) string {
	if z.mc.Exists(key) {
		return z.mc.Compile(app_msg.M(key, p...))
	} else {
		return ""
	}
}
func (z *Quiet) SetLogger(log *zap.Logger) {
	z.log = log
}

func (z *Quiet) Break() {
	z.log.Debug("Break")
}

func (z *Quiet) Header(key string, p ...app_msg.P) {
	z.mq.Verify(key)
	z.log.Debug(key, zap.Any("params", p))
}

func (z *Quiet) InfoTable(name string) Table {
	return &QuietTable{
		log: z.log,
		mq:  z.mq,
	}
}

func (z *Quiet) Info(key string, p ...app_msg.P) {
	z.mq.Verify(key)
	z.log.Debug(key, zap.Any("params", p))
}

func (z *Quiet) Error(key string, p ...app_msg.P) {
	z.mq.Verify(key)
	z.log.Debug(key, zap.Any("params", p))
	z.log.Error(z.mc.Compile(app_msg.M(key, p...)))
}

// always cancel process
func (z *Quiet) AskCont(key string, p ...app_msg.P) (cont bool, cancel bool) {
	z.mq.Verify(key)
	z.log.Debug(key, zap.Any("params", p))
	return false, true
}

// always cancel
func (z *Quiet) AskText(key string, p ...app_msg.P) (text string, cancel bool) {
	z.mq.Verify(key)
	z.log.Debug(key, zap.Any("params", p))
	return "", true
}

// always cancel
func (z *Quiet) AskSecure(key string, p ...app_msg.P) (secure string, cancel bool) {
	z.mq.Verify(key)
	z.log.Debug(key, zap.Any("params", p))
	return "", true
}

type QuietTable struct {
	log *zap.Logger
	mq  app_msg_container.Quality
}

func (z *QuietTable) HeaderRaw(h ...string) {
	z.log.Debug("header", zap.Any("h", h))
}

func (z *QuietTable) RowRaw(m ...string) {
	z.log.Debug("row", zap.Any("m", m))
}

func (z *QuietTable) Header(h ...app_msg.Message) {
	z.log.Debug("header", zap.Any("h", h))
	for _, m := range h {
		z.mq.Verify(m.Key())
	}
}

func (z *QuietTable) Row(m ...app_msg.Message) {
	z.log.Debug("row", zap.Any("m", m))
	for _, r := range m {
		z.mq.Verify(r.Key())
	}
}

func (z *QuietTable) Flush() {
	z.log.Debug("Flush")
}
