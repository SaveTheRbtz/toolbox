package es_open

import (
	"errors"
	"github.com/skratchdot/open-golang/open"
	"github.com/watermint/toolbox/essentials/log/esl"
)

var (
	ErrorAlwaysError = errors.New("always fail")
)

type Open interface {
	// Open a file, folder, or URI using the OS's default app.
	// Wait for the open command complete when `blocking` is true.
	Open(input string, blocking bool) error
}

func New() Open {
	return &openWrapper{}
}

func NewTestDummy() Open {
	return &dummyOpener{}
}
func NewTestError() Open {
	return &errorOpener{}
}

type openWrapper struct {
}

func (z *openWrapper) Open(input string, blocking bool) error {
	if blocking {
		return open.Run(input)
	} else {
		return open.Start(input)
	}
}

type dummyOpener struct {
}

func (z *dummyOpener) Open(input string, blocking bool) error {
	l := esl.Default()
	l.Debug("Open", esl.String("input", input), esl.Bool("blocking", blocking))
	return nil
}

type errorOpener struct {
}

func (z *errorOpener) Open(input string, blocking bool) error {
	l := esl.Default()
	l.Debug("Open", esl.String("input", input), esl.Bool("blocking", blocking))
	return ErrorAlwaysError
}
