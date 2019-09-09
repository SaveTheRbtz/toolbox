package app_ui

import (
	"github.com/watermint/toolbox/infra/ui/app_msg"
)

type UI interface {
	Header(key string, p ...app_msg.Param)
	Info(key string, p ...app_msg.Param)
	InfoTable(name string) Table
	Error(key string, p ...app_msg.Param)
	Break()
	Text(key string, p ...app_msg.Param) string

	AskCont(key string, p ...app_msg.Param) (cont bool, cancel bool)
	AskText(key string, p ...app_msg.Param) (text string, cancel bool)
	AskSecure(key string, p ...app_msg.Param) (secure string, cancel bool)

	OpenArtifact(path string)
	Success(key string, p ...app_msg.Param)
	Failure(key string, p ...app_msg.Param)

	IsConsole() bool
	IsWeb() bool
}

type Table interface {
	Header(h ...app_msg.Message)
	HeaderRaw(h ...string)
	Row(m ...app_msg.Message)
	RowRaw(m ...string)
	Flush()
}