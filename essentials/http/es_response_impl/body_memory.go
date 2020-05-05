package es_response_impl

import (
	"github.com/tidwall/gjson"
	"github.com/watermint/toolbox/essentials/encoding/es_json"
	"github.com/watermint/toolbox/essentials/http/es_context"
	"github.com/watermint/toolbox/essentials/http/es_response"
	"github.com/watermint/toolbox/essentials/log/es_encode"
	"github.com/watermint/toolbox/essentials/log/es_log"
	"io/ioutil"
	"os"
)

func newMemoryBody(ctx es_context.Context, content []byte) es_response.Body {
	return &bodyMemoryImpl{
		ctx:     ctx,
		content: content,
	}
}

type bodyMemoryImpl struct {
	ctx     es_context.Context
	content []byte
}

func (z bodyMemoryImpl) Json() es_json.Json {
	if j, err := z.AsJson(); err != nil {
		return es_json.Null()
	} else {
		return j
	}
}

func (z bodyMemoryImpl) Error() error {
	return nil
}

func (z bodyMemoryImpl) BodyString() string {
	return string(z.content)
}

func (z bodyMemoryImpl) AsJson() (es_json.Json, error) {
	l := z.ctx.Log()
	if !gjson.ValidBytes(z.content) {
		l.Debug("Invalid bytes", es_log.Any("bytes", es_encode.ByteDigest(z.content)))
		return nil, es_response.ErrorContentIsNotAJSON
	}
	return es_json.Parse(z.content)
}

func toFile(ctx es_context.Context, content []byte) (string, error) {
	l := ctx.Log()
	p, err := ioutil.TempFile("", ctx.ClientHash())
	if err != nil {
		l.Debug("Unable to create temp file", es_log.Error(err))
		return "", err
	}
	cleanupOnError := func() {
		if err := p.Close(); err != nil {
			l.Debug("unable to close", es_log.Error(err))
		}
		if err := os.Remove(p.Name()); err != nil {
			l.Debug("unable to remove", es_log.Error(err))
		}
	}
	if err := ioutil.WriteFile(p.Name(), content, 0600); err != nil {
		l.Debug("Unable to write", es_log.Error(err))
		cleanupOnError()
		return "", err
	}
	if err := p.Close(); err != nil {
		l.Debug("unable to close", es_log.Error(err))
		cleanupOnError()
		return "", err
	}
	return p.Name(), nil
}

func (z bodyMemoryImpl) AsFile() (string, error) {
	return toFile(z.ctx, z.content)
}

func (z bodyMemoryImpl) ContentLength() int64 {
	return int64(len(z.content))
}

func (z bodyMemoryImpl) Body() []byte {
	return z.content
}

func (z bodyMemoryImpl) File() string {
	return ""
}

func (z bodyMemoryImpl) IsFile() bool {
	return false
}